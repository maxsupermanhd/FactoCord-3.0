---
name: Bug report
about: Create a report
title: "[BUG]"
labels: bug
assignees: ''

---

**BUG-FILTER-CHECKLIST.** Anyone who opens up a bug report must fill these checkboxes (add an "x" to `[ ]` to be like `[x]`) to assume that he acknowledged that he red README.md and followed checklist below.
- [ ] Bot have rights to write, read and send embeds to desired channel.
- [ ] Executable not receiving SIGTERM, SIGKILL or other terminating signals during work.
- [ ] Scenario/mods are working correctly and not crashing game.
- [ ] I'm using FactoCord-compatible scenario that writes up all desired text to console with `[DISCORD]` tag.
OR (remove this and unnecessary line above or below)
- [ ] I'm using vanilla game and acknowledged about bugs and its workarounds.
- [ ] I have direct access to shell that runs FactoCord. (via tmux, screen or whatever that pipes `stdin` **and** `stderr` to output)

**Describe the bug**
A clear and concise description of what the bug is. (Not sending message, message spam, not responding on commands, malformed output, or any other undefined behavior)

**To Reproduce**
Steps to reproduce the behavior.

**Expected behavior**
A clear and concise description of what you expected to happen.

**My instance information**
- OS: ex PepeTheFrogLinux lmao version 1.2.3 (kernel kek version 9.8.7)
- golang version: output of `go version`
- FactoCord version: standing commit hash
- Factorio version: Factorio server version
- Server config: uploaded `.env` file (don't forget to remove token!)
- Factorio and bot logs: files `error.log`, `factorio.log` and `.exit` if they exist. (`.exit` contains the latest exit code I believe)
- Run conditions: systemd, manual or explain process of starting up server. (if you are using systemd, include service file here)

**Screenshots**
Add it if you are seeing malformed/unexpected behavior in the game using FactoCord-compatible scenario, and it is related to FactoCord.
