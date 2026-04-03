package plugin

import (
	_ "embed"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

//go:embed skill.md
var skillContent string

// SkillBody returns the SKILL.md content with YAML frontmatter stripped.
func SkillBody() string {
	content := skillContent
	if strings.HasPrefix(content, "---") {
		if idx := strings.Index(content[3:], "---"); idx != -1 {
			content = strings.TrimLeft(content[3+idx+3:], "\n")
		}
	}
	return content
}

// ExtractClaudePlugin writes a complete .claude-plugin/ directory to dest,
// including plugin.json, hooks, commands, and a copy of SKILL.md.
func ExtractClaudePlugin(dest string) error {
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return err
	}

	// plugin.json
	manifest := map[string]interface{}{
		"name":        "neetocal",
		"description": "NeetoCal integration for Claude Code. Manage meetings, bookings, and scheduling.",
		"author": map[string]string{
			"name":  "BigBinary",
			"email": "support@bigbinary.com",
		},
		"homepage":   "https://github.com/neetozone/neeto-cal-cli",
		"repository": "https://github.com/neetozone/neeto-cal-cli",
		"license":    "MIT",
	}
	manifestJSON, _ := json.MarshalIndent(manifest, "", "  ")
	if err := writeFile(filepath.Join(dest, "plugin.json"), manifestJSON, 0o644); err != nil {
		return err
	}

	// hooks/hooks.json
	hooks := map[string]interface{}{
		"hooks": map[string]interface{}{
			"SessionStart": []interface{}{
				map[string]interface{}{
					"hooks": []interface{}{
						map[string]interface{}{
							"type":    "command",
							"command": "${CLAUDE_PLUGIN_ROOT}/hooks/session-start.sh",
							"timeout": 5,
						},
					},
				},
			},
		},
	}
	hooksJSON, _ := json.MarshalIndent(hooks, "", "  ")
	if err := writeFile(filepath.Join(dest, "hooks", "hooks.json"), hooksJSON, 0o644); err != nil {
		return err
	}

	// hooks/session-start.sh
	sessionStartSh := `#!/bin/sh
# NeetoCal CLI — session-start hook for Claude Code
# Lightweight auth liveness check. Always exits 0 (informational).

if ! command -v neetocal >/dev/null 2>&1; then
  echo "NeetoCal CLI is not installed or not on PATH."
  exit 0
fi

if neetocal whoami >/dev/null 2>&1; then
  echo "NeetoCal plugin active."
else
  echo "NeetoCal CLI installed but not authenticated. Run 'neetocal login' to authenticate."
fi

exit 0
`
	if err := writeFile(filepath.Join(dest, "hooks", "session-start.sh"), []byte(sessionStartSh), 0o755); err != nil {
		return err
	}

	// commands/doctor.md
	doctorMd := `---
name: neetocal-doctor
description: Check NeetoCal CLI health — auth, API connectivity.
invocable: true
---

Run ` + "`neetocal doctor`" + ` and report the results to the user.
`
	if err := writeFile(filepath.Join(dest, "commands", "doctor.md"), []byte(doctorMd), 0o644); err != nil {
		return err
	}

	// skills/neetocal/SKILL.md
	if err := writeFile(filepath.Join(dest, "skills", "neetocal", "SKILL.md"), []byte(skillContent), 0o644); err != nil {
		return err
	}

	return nil
}

func writeFile(path string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}
