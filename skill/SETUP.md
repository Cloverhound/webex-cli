---
name: webex-cli-setup
description: "Set up the Webex CLI, then sign in with device login. Use at the start of any Webex task, or when the user asks to install or log in to Webex."
---

# Webex CLI setup

Run these steps in order. Never store, print or write Webex tokens or credentials anywhere, including in this skill.

## 1. Install the CLI

Run `command -v webex`. If nothing is found, run:

```bash
npm install -g @cloverhound/webex-cli
```

## 2. Load the Webex skill

If `~/.claude/skills/webex-cli/SKILL.md` does not exist, run:

```bash
mkdir -p ~/.claude && webex post-install -y
```

Read `~/.claude/skills/webex-cli/SKILL.md` and follow it for all Webex work. Open the area skill files beside it (admin, calling, cc, device, meetings, messaging) as needed.

## 3. Check for a login

Run `webex auth status`. If it shows `Token: valid` or `Refresh: available`, skip to step 5.

## 4. Sign in with device login

Start the login in the background, because it blocks until the user approves. It writes the instructions to stderr, so capture both streams:

```bash
nohup webex login --device --output json > /tmp/webex-login.log 2>&1 &
```

Read `/tmp/webex-login.log`. Relay the line `Open <url> and enter <code>` to the user word for word, and include the expiry from the `Waiting for approval` line.

End the turn and wait. When the user replies, read the log again. A JSON line with `"status":"logged_in"` means the login worked. If the code expired or the login failed, start a new device login and relay the new code.

## 5. Confirm

Run `webex auth status` and report the user, org and token expiry.

Tell the user:

- If `Mode: read-only` appears, the login is read-only. Otherwise it has full access, so ask before any command that creates, changes or deletes Webex resources.
- If `Store` shows `file, plain text`, there is no keyring here, so the token is unencrypted on disk. They can run `webex logout` when they're done.

## 6. Offer to save this skill

Ask whether to save these steps for future sessions. If the user agrees, save the output of `webex agent-setup` as a skill named `webex-cli-setup`. In Claude Code, write it to `~/.claude/skills/webex-cli-setup/SKILL.md`. Where that folder does not outlast the session, give the user the file to add to their skills.
