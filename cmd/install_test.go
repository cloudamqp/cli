package cmd

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallSkillsCmd(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	err := installSkillsCmd.RunE(installSkillsCmd, []string{})
	require.NoError(t, err)

	base := filepath.Join(home, ".claude", "skills", "cloudamqp-cli")
	assert.FileExists(t, filepath.Join(base, "SKILL.md"))
	assert.FileExists(t, filepath.Join(base, "references", "scripting.md"))
	assert.FileExists(t, filepath.Join(base, "references", "upgrades.md"))
	assert.FileExists(t, filepath.Join(base, "references", "vpc-setup.md"))
}

func TestInstallSkillsCmd_NoArgs(t *testing.T) {
	err := installSkillsCmd.Args(installSkillsCmd, []string{"unexpected"})
	assert.Error(t, err)
}
