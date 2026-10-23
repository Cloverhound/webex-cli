// Package skill embeds the agent skill files so an installed binary always
// carries the skill that matches its own commands.
package skill

import (
	"embed"
	"io/fs"
)

//go:embed SKILL.md */SKILL.md
var embedded embed.FS

// Files returns every embedded skill file keyed by its slash-separated path
// relative to the skill directory, such as "SKILL.md" or "admin/SKILL.md".
func Files() (map[string][]byte, error) {
	files := make(map[string][]byte)
	err := fs.WalkDir(embedded, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := embedded.ReadFile(path)
		if err != nil {
			return err
		}
		files[path] = data
		return nil
	})
	return files, err
}
