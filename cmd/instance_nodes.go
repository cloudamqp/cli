package cmd

import (
	"encoding/json"
	"fmt"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var instanceNodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Manage instance nodes",
	Long:  `List nodes and get available versions for the instance.`,
}

var instanceNodesListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List nodes in the instance",
	Long:    `Retrieves all nodes in the instance.`,
	Example: `  cloudamqp instance manage 1234 nodes list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceID := currentInstanceID

		instanceAPIKey, err := getInstanceAPIKey(instanceID)
		if err != nil {
			return fmt.Errorf("failed to get instance API key: %w", err)
		}

		c := client.NewInstanceAPI(instanceAPIKey)

		nodes, err := c.ListNodes()
		if err != nil {
			fmt.Printf("Error listing nodes: %v\n", err)
			return err
		}

		if len(nodes) == 0 {
			fmt.Println("No nodes found.")
			return nil
		}

		output, err := json.MarshalIndent(nodes, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("Nodes:\n%s\n", string(output))
		return nil
	},
}

var instanceNodesVersionsCmd = &cobra.Command{
	Use:     "versions",
	Short:   "Get available versions",
	Long:    `Lists RabbitMQ and Erlang versions to which the instance can be upgraded.`,
	Example: `  cloudamqp instance manage 1234 nodes versions`,
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceID := currentInstanceID

		instanceAPIKey, err := getInstanceAPIKey(instanceID)
		if err != nil {
			return fmt.Errorf("failed to get instance API key: %w", err)
		}

		c := client.NewInstanceAPI(instanceAPIKey)

		versions, err := c.GetAvailableVersions()
		if err != nil {
			fmt.Printf("Error getting available versions: %v\n", err)
			return err
		}

		output, err := json.MarshalIndent(versions, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("Available versions:\n%s\n", string(output))
		return nil
	},
}

func init() {
	instanceNodesCmd.AddCommand(instanceNodesListCmd)
	instanceNodesCmd.AddCommand(instanceNodesVersionsCmd)
}
