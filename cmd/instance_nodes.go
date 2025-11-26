package cmd

import (
	"fmt"
	"os"

	"cloudamqp-cli/client"
	"cloudamqp-cli/internal/table"
	"github.com/spf13/cobra"
)

var instanceNodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Manage instance nodes",
	Long:  `List nodes and get available versions for the instance.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		cmd.SilenceUsage = true
		return fmt.Errorf("subcommand required")
	},
}

var instanceNodesListCmd = &cobra.Command{
	Use:     "list --id <instance_id>",
	Short:   "List nodes in the instance",
	Long:    `Retrieves all nodes in the instance.`,
	Example: `  cloudamqp instance nodes list --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		idFlag, _ := cmd.Flags().GetString("id")
		if idFlag == "" {
			return fmt.Errorf("instance ID is required. Use --id flag")
		}

		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey)

		nodes, err := c.ListNodes(idFlag)
		if err != nil {
			fmt.Printf("Error listing nodes: %v\n", err)
			return err
		}

		if len(nodes) == 0 {
			fmt.Println("No nodes found.")
			return nil
		}

		// Create table and populate data
		t := table.New(os.Stdout, "NAME", "CONFIGURED", "RUNNING", "DISK_SIZE", "RABBITMQ_VERSION")
		for _, node := range nodes {
			configured := "No"
			if node.Configured {
				configured = "Yes"
			}
			running := "No"
			if node.Running {
				running = "Yes"
			}
			totalDisk := node.DiskSize + node.AdditionalDiskSize
			t.AddRow(
				node.Name,
				configured,
				running,
				fmt.Sprintf("%d GB", totalDisk),
				node.RabbitMQVersion,
			)
		}
		t.Print()

		return nil
	},
}

var instanceNodesVersionsCmd = &cobra.Command{
	Use:     "versions --id <instance_id>",
	Short:   "Get available versions",
	Long:    `Lists RabbitMQ and Erlang versions to which the instance can be upgraded.`,
	Example: `  cloudamqp instance nodes versions --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		idFlag, _ := cmd.Flags().GetString("id")
		if idFlag == "" {
			return fmt.Errorf("instance ID is required. Use --id flag")
		}

		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey)

		versions, err := c.GetAvailableVersions(idFlag)
		if err != nil {
			fmt.Printf("Error getting available versions: %v\n", err)
			return err
		}

		fmt.Printf("Available versions:\n")
		fmt.Printf("RabbitMQ versions: %v\n", versions.RabbitMQVersions)
		fmt.Printf("Erlang versions: %v\n", versions.ErlangVersions)
		return nil
	},
}

func init() {
	// Add --id flag to all subcommands
	instanceNodesListCmd.Flags().StringP("id", "", "", "Instance ID (required)")
	instanceNodesListCmd.MarkFlagRequired("id")

	instanceNodesVersionsCmd.Flags().StringP("id", "", "", "Instance ID (required)")
	instanceNodesVersionsCmd.MarkFlagRequired("id")

	instanceNodesCmd.AddCommand(instanceNodesListCmd)
	instanceNodesCmd.AddCommand(instanceNodesVersionsCmd)
}
