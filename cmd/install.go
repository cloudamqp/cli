package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

//go:embed all:skills
var skillsFS embed.FS

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install integrations",
}

var installSkillsCmd = &cobra.Command{
	Use:   "skills",
	Short: "Install Claude Code skills to ~/.claude/skills/",
	Long: `Install the CloudAMQP CLI skills for Claude Code.

Skills teach Claude how to use the cloudamqp CLI. After installation,
Claude Code will automatically discover and use them.

Skills are installed to: ~/.claude/skills/cloudamqp-cli/`,
	RunE: func(cmd *cobra.Command, args []string) error {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("could not determine home directory: %w", err)
		}
		dest := filepath.Join(home, ".claude", "skills", "cloudamqp-cli")

		err = fs.WalkDir(skillsFS, "skills/cloudamqp-cli", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			// path relative to dest: strip "skills/cloudamqp-cli" prefix
			rel, err := filepath.Rel("skills/cloudamqp-cli", path)
			if err != nil {
				return err
			}
			target := filepath.Join(dest, rel)
			if d.IsDir() {
				return os.MkdirAll(target, 0755)
			}
			data, err := skillsFS.ReadFile(path)
			if err != nil {
				return err
			}
			return os.WriteFile(target, data, 0644)
		})
		if err != nil {
			return fmt.Errorf("failed to install skills: %w", err)
		}
		fmt.Printf("Skills installed to %s\n", dest)
		return nil
	},
}

func init() {
	installCmd.AddCommand(installSkillsCmd)
	rootCmd.AddCommand(installCmd)
}
