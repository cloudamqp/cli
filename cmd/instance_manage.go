package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var currentInstanceID string

var instanceManageCmd = &cobra.Command{
	Use:   "manage <instance_id>",
	Short: "Manage a specific CloudAMQP instance",
	Long: `Use instance-specific API to manage nodes, plugins, actions, and more.

This command uses the instance API key, not your main API key.
Instance API keys are automatically saved when you run 'cloudamqp instance get <id>'.`,
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			cmd.SilenceUsage = true
			return fmt.Errorf("instance ID is required")
		}

		if len(args) == 1 {
			cmd.Help()
			cmd.SilenceUsage = true
			return fmt.Errorf("subcommand is required")
		}

		// Parse custom: manage <instance_id> <subcommand> [args...]
		instanceID := args[0]
		subcommandName := args[1]
		subArgs := args[2:]

		// Set the global instance ID
		currentInstanceID = instanceID

		// Find and execute the appropriate subcommand
		for _, subCmd := range cmd.Commands() {
			if subCmd.Name() == subcommandName {
				// Handle subcommands with their own subcommands (like config)
				if len(subArgs) > 0 && len(subCmd.Commands()) > 0 {
					// Look for nested subcommand
					nestedSubcommandName := subArgs[0]
					nestedArgs := subArgs[1:]

					for _, nestedSubCmd := range subCmd.Commands() {
						if nestedSubCmd.Name() == nestedSubcommandName {
							if nestedSubCmd.RunE != nil {
								return nestedSubCmd.RunE(nestedSubCmd, nestedArgs)
							}
							return fmt.Errorf("subcommand %s %s has no implementation", subcommandName, nestedSubcommandName)
						}
					}
					return fmt.Errorf("unknown subcommand: %s %s", subcommandName, nestedSubcommandName)
				}

				// Execute direct subcommand
				if subCmd.RunE != nil {
					return subCmd.RunE(subCmd, subArgs)
				}

				// If no RunE, show help for the subcommand
				if len(subCmd.Commands()) > 0 {
					subCmd.Help()
					cmd.SilenceUsage = true
					return fmt.Errorf("subcommand required for %s", subcommandName)
				}

				return fmt.Errorf("subcommand %s has no implementation", subcommandName)
			}
		}

		cmd.Help()
		cmd.SilenceUsage = true
		return fmt.Errorf("unknown subcommand: %s", subcommandName)
	},
}

func init() {
	// All commands have been moved to direct instance subcommands
	// Keeping manage command for backward compatibility but it will be empty
}
