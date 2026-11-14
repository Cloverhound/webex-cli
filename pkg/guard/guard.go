// Package guard is the write protocol. Webex has no conditional writes, can accept a
// write and ignore it, and can commit a write it reports as failed, so every write
// here is checked before and after it is sent.
package guard

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Cloverhound/webex-cli/pkg/resource"
	"github.com/Cloverhound/webex-cli/pkg/webex"
)

// Kind classifies a failure. The app shows each kind differently.
type Kind string

const (
	Refused        Kind = "refused"         // writes are not permitted in this run
	Stopped        Kind = "stopped"         // an earlier failure or cancellation ended the run
	Conflict       Kind = "conflict"        // the object changed after the plan
	Rejected       Kind = "rejected"        // Webex refused the write; nothing changed
	OutcomeUnknown Kind = "outcome_unknown" // the write may or may not have landed
	Divergent      Kind = "divergent"       // Webex accepted the write but the object differs
	Unverified     Kind = "unverified"      // the write was accepted but could not be read back
)

// Failure is what the admin sees.
type Failure struct {
	Kind    Kind
	Summary string
	Detail  string
}

func (f *Failure) Error() string { return f.Summary + ": " + f.Detail }

// API is the Webex access the guard needs.
type API interface {
	Read(ctx context.Context, t resource.Type, key resource.Key) (webex.Object, error)
	Create(ctx context.Context, t resource.Type, desired resource.Values) (resource.Key, error)
	Update(ctx context.Context, t resource.Type, key resource.Key, raw map[string]any, delta resource.Values) error
	Delete(ctx context.Context, t resource.Type, key resource.Key) error
}

const verifyAttempts = 3

// Guard wraps writes for one apply. It lives for one provider process, which is what
// makes its stop flag effective: OpenTofu keeps scheduling independent resources after
// a failure, and every one of them comes through here.
type Guard struct {
	API API

	// Writable must be set by the reviewed apply. Sync and plan runs leave it false,
	// so they cannot change Webex even by mistake.
	Writable bool

	// VerifyDelay separates read-back attempts, to allow for read-after-write lag.
	VerifyDelay time.Duration

	mu      sync.Mutex
	stopped string
}

func (g *Guard) stop(reason string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.stopped == "" {
		g.stopped = reason
	}
}

func (g *Guard) fail(f *Failure) *Failure {
	g.stop(f.Summary)
	return f
}

func (g *Guard) admit(ctx context.Context) error {
	if !g.Writable {
		return &Failure{Refused, "Refused: this run may not write to Webex", "writes happen only during a reviewed apply"}
	}
	if err := g.halted(); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		g.stop("the run was cancelled")
		return &Failure{Stopped, "Stopped: the run was cancelled", err.Error()}
	}
	return nil
}

// halted reports a stop. Writes check it again just before sending, because another
// write can fail while this one reads.
func (g *Guard) halted() error {
	g.mu.Lock()
	reason := g.stopped
	g.mu.Unlock()
	if reason != "" {
		return &Failure{Stopped, "Stopped: an earlier failure blocks further writes", reason}
	}
	return nil
}

// recoveryContext deliberately ignores the caller's cancellation: the outcome of a
// write that was already sent must be learned even when the run is stopping.
func recoveryContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 8*time.Second)
}

// plannedChanges lists fields whose planned value is known and differs. A nil planned
// value means Terraform left it for the server, so it is not a change.
func plannedChanges(t resource.Type, have, planned resource.Values) []string {
	return t.Unmet(have, planned)
}

// Update writes the fields that differ between prior and planned. It returns the
// object's real values, also on failure, so Terraform state never records a value
// Webex does not hold.
func (g *Guard) Update(ctx context.Context, t resource.Type, key resource.Key, prior, planned resource.Values) (resource.Values, error) {
	if err := g.admit(ctx); err != nil {
		return prior, err
	}
	fresh, err := g.API.Read(ctx, t, key)
	if errors.Is(err, webex.ErrNotFound) {
		return prior, g.fail(&Failure{Conflict, "Conflict: this object was deleted after the plan was made", t.Title})
	}
	if err != nil {
		return prior, g.fail(&Failure{Stopped, "Stopped before writing: could not read the current object", err.Error()})
	}
	// Webex offers no conditional write, so a change landing between this read and
	// the PUT cannot be detected.
	if d := resource.Changed(fresh.Values, prior, t.FieldNames()); len(d) > 0 {
		return fresh.Values, g.fail(&Failure{Conflict, "Conflict: this object changed after the plan was made", strings.Join(d, ", ")})
	}
	fields := plannedChanges(t, prior, planned)
	if len(fields) == 0 {
		return fresh.Values, nil
	}
	delta := resource.Values{}
	for _, f := range fields {
		delta[f] = planned[f]
	}
	return g.write(ctx, t, key, fresh.Raw, fresh.Values, planned, delta)
}

func (g *Guard) write(ctx context.Context, t resource.Type, key resource.Key, raw map[string]any, baseline, planned, delta resource.Values) (resource.Values, error) {
	if err := g.halted(); err != nil {
		return baseline, err
	}
	// Each attempt checks again, as a 429 retry can come after another write failed.
	if err := g.API.Update(webex.WithBeforeWrite(ctx, g.halted), t, key, raw, delta); err != nil {
		var stopped *Failure
		if errors.As(err, &stopped) {
			return baseline, err
		}
		// Other writes stop now, not after the read below, which can take seconds.
		g.stop(writeFailure(err, "The write", "").Summary)
		actual := g.observe(t, key, baseline)
		return actual, g.fail(writeFailure(err, "The write", observed(actual, baseline, planned, t)))
	}
	actual, diff, err := g.verify(t, key, planned)
	if err != nil {
		return baseline, g.fail(&Failure{Unverified, "Webex accepted the write but it could not be verified", err.Error()})
	}
	if len(diff) > 0 {
		return actual, g.fail(&Failure{Divergent, "The observed state differs from the approved plan", strings.Join(diff, ", ")})
	}
	return actual, nil
}

// observe reads what is there after a failed write. The result is only a first look:
// a delayed write can still land afterwards, so the caller reports the outcome as
// unknown whatever this says.
func (g *Guard) observe(t resource.Type, key resource.Key, fallback resource.Values) resource.Values {
	ctx, cancel := recoveryContext()
	defer cancel()
	obj, err := g.API.Read(ctx, t, key)
	if err != nil {
		return fallback
	}
	return obj.Values
}

func observed(actual, baseline, planned resource.Values, t resource.Type) string {
	switch {
	case len(plannedChanges(t, actual, planned)) == 0:
		return "it appears to have landed"
	case len(resource.Changed(actual, baseline, t.FieldNames())) == 0:
		return "it does not appear to have landed yet"
	default:
		return "the object holds other values"
	}
}

func writeFailure(err error, what, seen string) *Failure {
	if webex.Classify(err) == webex.Rejected {
		return &Failure{Rejected, "Webex rejected " + strings.ToLower(what), err.Error()}
	}
	detail := err.Error()
	if seen != "" {
		detail = fmt.Sprintf("%s; %s", detail, seen)
	}
	return &Failure{OutcomeUnknown, "The outcome of " + strings.ToLower(what) + " is unknown", detail}
}

func (g *Guard) verify(t resource.Type, key resource.Key, planned resource.Values) (resource.Values, []string, error) {
	ctx, cancel := recoveryContext()
	defer cancel()
	var last error
	var actual resource.Values
	for i := 0; i < verifyAttempts; i++ {
		if i > 0 {
			time.Sleep(g.VerifyDelay)
		}
		obj, err := g.API.Read(ctx, t, key)
		if err != nil {
			last = err
			continue
		}
		last, actual = nil, obj.Values
		if len(plannedChanges(t, actual, planned)) == 0 {
			return actual, nil, nil
		}
	}
	if last != nil {
		return nil, nil, last
	}
	return actual, plannedChanges(t, actual, planned), nil
}

// Create makes the object, then writes any field Webex ignored on create. A returned
// key with an error means the object exists but is not as planned, and Terraform
// state must hold it.
func (g *Guard) Create(ctx context.Context, t resource.Type, desired resource.Values) (resource.Key, resource.Values, error) {
	if err := g.admit(ctx); err != nil {
		return resource.Key{}, nil, err
	}
	// A create that copies from an existing object reads first, and a failure
	// elsewhere during those reads must still stop its POST.
	key, err := g.API.Create(webex.WithBeforeWrite(ctx, g.halted), t, desired)
	var stopped *Failure
	if errors.As(err, &stopped) {
		return resource.Key{}, nil, err
	}
	if err != nil {
		return resource.Key{}, nil, g.fail(writeFailure(err, "The create", ""))
	}
	obj, err := g.readBack(t, key, func(err error) bool { return err == nil })
	if err != nil {
		return key, desired, g.fail(&Failure{Unverified, "The object was created but could not be verified", err.Error()})
	}
	missed := plannedChanges(t, obj.Values, desired)
	if len(missed) == 0 {
		return key, obj.Values, nil
	}
	delta := resource.Values{}
	for _, f := range missed {
		delta[f] = desired[f]
	}
	actual, err := g.write(ctx, t, key, obj.Raw, obj.Values, desired, delta)
	return key, actual, err
}

// Delete removes the object if it still matches prior. An object already gone is
// success, since that is the intended end state.
func (g *Guard) Delete(ctx context.Context, t resource.Type, key resource.Key, prior resource.Values) error {
	if err := g.admit(ctx); err != nil {
		return err
	}
	fresh, err := g.API.Read(ctx, t, key)
	if errors.Is(err, webex.ErrNotFound) {
		return nil
	}
	if err != nil {
		return g.fail(&Failure{Stopped, "Stopped before writing: could not read the current object", err.Error()})
	}
	if d := resource.Changed(fresh.Values, prior, t.FieldNames()); len(d) > 0 {
		return g.fail(&Failure{Conflict, "Conflict: this object changed after the plan was made", strings.Join(d, ", ")})
	}
	if err := g.halted(); err != nil {
		return err
	}
	if err := g.API.Delete(webex.WithBeforeWrite(ctx, g.halted), t, key); err != nil {
		var stopped *Failure
		if errors.As(err, &stopped) {
			return err
		}
		return g.fail(writeFailure(err, "The delete", ""))
	}
	_, err = g.readBack(t, key, func(err error) bool { return errors.Is(err, webex.ErrNotFound) })
	switch {
	case errors.Is(err, webex.ErrNotFound):
		return nil
	case err != nil:
		return g.fail(&Failure{Unverified, "The delete could not be verified", err.Error()})
	default:
		return g.fail(&Failure{Divergent, "The observed state differs from the approved plan", "object still exists"})
	}
}

// readBack reads the object after a write until done accepts the result, up to
// verifyAttempts reads apart, because Webex can answer a read with the state before
// the write. It returns the last read.
func (g *Guard) readBack(t resource.Type, key resource.Key, done func(error) bool) (webex.Object, error) {
	ctx, cancel := recoveryContext()
	defer cancel()
	var obj webex.Object
	var err error
	for i := 0; i < verifyAttempts; i++ {
		if i > 0 {
			time.Sleep(g.VerifyDelay)
		}
		if obj, err = g.API.Read(ctx, t, key); done(err) {
			break
		}
	}
	return obj, err
}
