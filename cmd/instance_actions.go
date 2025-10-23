package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var instanceActionsCmd = &cobra.Command{
	Use:   "actions",
	Short: "Perform instance actions",
	Long:  `Restart, stop, start, reboot, and upgrade instance components.`,
}

// Restart commands
var restartRabbitMQCmd = &cobra.Command{
	Use:   "restart-rabbitmq",
	Short: "Restart RabbitMQ",
	Long:  `Restart RabbitMQ on specified nodes or all nodes.`,
	Example: `  cloudamqp instance manage 1234 actions restart-rabbitmq
  cloudamqp instance manage 1234 actions restart-rabbitmq --nodes=node1,node2`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, "restart-rabbitmq")
	},
}

var restartClusterCmd = &cobra.Command{
	Use:     "restart-cluster",
	Short:   "Restart cluster",
	Long:    `Restart the entire cluster.`,
	Example: `  cloudamqp instance manage 1234 actions restart-cluster`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performClusterAction(cmd, "restart-cluster")
	},
}

var restartManagementCmd = &cobra.Command{
	Use:     "restart-management",
	Short:   "Restart management interface",
	Long:    `Restart the RabbitMQ management interface.`,
	Example: `  cloudamqp instance manage 1234 actions restart-management`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, "restart-management")
	},
}

// Stop/Start commands
var stopCmd = &cobra.Command{
	Use:     "stop",
	Short:   "Stop instance",
	Long:    `Stop specified nodes or all nodes.`,
	Example: `  cloudamqp instance manage 1234 actions stop`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, "stop")
	},
}

var startCmd = &cobra.Command{
	Use:     "start",
	Short:   "Start instance",
	Long:    `Start specified nodes or all nodes.`,
	Example: `  cloudamqp instance manage 1234 actions start`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, "start")
	},
}

var rebootCmd = &cobra.Command{
	Use:     "reboot",
	Short:   "Reboot instance",
	Long:    `Reboot specified nodes or all nodes.`,
	Example: `  cloudamqp instance manage 1234 actions reboot`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, "reboot")
	},
}

// Cluster commands
var stopClusterCmd = &cobra.Command{
	Use:     "stop-cluster",
	Short:   "Stop cluster",
	Long:    `Stop the entire cluster.`,
	Example: `  cloudamqp instance manage 1234 actions stop-cluster`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performClusterAction(cmd, "stop-cluster")
	},
}

var startClusterCmd = &cobra.Command{
	Use:     "start-cluster",
	Short:   "Start cluster",
	Long:    `Start the entire cluster.`,
	Example: `  cloudamqp instance manage 1234 actions start-cluster`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performClusterAction(cmd, "start-cluster")
	},
}

// Upgrade commands
var upgradeErlangCmd = &cobra.Command{
	Use:   "upgrade-erlang",
	Short: "Upgrade Erlang",
	Long: `Always updates to latest compatible version.

Note: This action is asynchronous. The request will return immediately, the process runs in the background.`,
	Example: `  cloudamqp instance manage 1234 actions upgrade-erlang`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performUpgradeAction(cmd, "upgrade-erlang", "")
	},
}

var upgradeRabbitMQCmd = &cobra.Command{
	Use:   "upgrade-rabbitmq",
	Short: "Upgrade RabbitMQ",
	Long: `Upgrade RabbitMQ to specified version.

Note: This action is asynchronous. The request will return immediately, the process runs in the background.`,
	Example: `  cloudamqp instance manage 1234 actions upgrade-rabbitmq --version=3.10.7`,
	RunE: func(cmd *cobra.Command, args []string) error {
		version, _ := cmd.Flags().GetString("version")
		if version == "" {
			return fmt.Errorf("version flag is required")
		}
		return performUpgradeAction(cmd, "upgrade-rabbitmq", version)
	},
}

var upgradeRabbitMQErlangCmd = &cobra.Command{
	Use:   "upgrade-all",
	Short: "Upgrade RabbitMQ and Erlang",
	Long: `Always updates to latest possible version of both RabbitMQ and Erlang.

Note: This action is asynchronous. The request will return immediately, the process runs in the background.`,
	Example: `  cloudamqp instance manage 1234 actions upgrade-all`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performUpgradeAction(cmd, "upgrade-all", "")
	},
}

// HiPE and Firehose commands
var toggleHiPECmd = &cobra.Command{
	Use:     "toggle-hipe",
	Short:   "Enable/disable HiPE",
	Long:    `Enable or disable HiPE (High Performance Erlang) compilation.`,
	Example: `  cloudamqp instance manage 1234 actions toggle-hipe --enable=true`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performToggleAction(cmd, "hipe")
	},
}

var toggleFirehoseCmd = &cobra.Command{
	Use:     "toggle-firehose",
	Short:   "Enable/disable Firehose",
	Long:    `Enable or disable RabbitMQ Firehose tracing (not recommended in production).`,
	Example: `  cloudamqp instance manage 1234 actions toggle-firehose --enable=true --vhost=/`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performToggleAction(cmd, "firehose")
	},
}

var upgradeVersionsCmd = &cobra.Command{
	Use:     "upgrade-versions",
	Short:   "Fetch upgrade versions",
	Long:    `Returns what version of Erlang and RabbitMQ the cluster will update to.`,
	Example: `  cloudamqp instance manage 1234 actions upgrade-versions`,
	RunE: func(cmd *cobra.Command, args []string) error {
		instanceID := currentInstanceID

		instanceAPIKey, err := getInstanceAPIKey(instanceID)
		if err != nil {
			return fmt.Errorf("failed to get instance API key: %w", err)
		}

		c := client.NewInstanceAPI(instanceAPIKey)

		versions, err := c.GetUpgradeVersions()
		if err != nil {
			fmt.Printf("Error getting upgrade versions: %v\n", err)
			return err
		}

		output, err := json.MarshalIndent(versions, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("Upgrade versions:\n%s\n", string(output))
		return nil
	},
}

// Helper functions
func performNodeAction(cmd *cobra.Command, action string) error {
	instanceID := currentInstanceID

	instanceAPIKey, err := getInstanceAPIKey(instanceID)
	if err != nil {
		return fmt.Errorf("failed to get instance API key: %w", err)
	}

	c := client.NewInstanceAPI(instanceAPIKey)

	nodesStr, _ := cmd.Flags().GetString("nodes")
	var nodes []string
	if nodesStr != "" {
		nodes = strings.Split(nodesStr, ",")
	}

	switch action {
	case "restart-rabbitmq":
		err = c.RestartRabbitMQ(nodes)
	case "restart-management":
		err = c.RestartManagement(nodes)
	case "stop":
		err = c.StopInstance(nodes)
	case "start":
		err = c.StartInstance(nodes)
	case "reboot":
		err = c.RebootInstance(nodes)
	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	if err != nil {
		fmt.Printf("Error performing %s: %v\n", action, err)
		return err
	}

	fmt.Printf("%s initiated successfully.\n", strings.Title(strings.ReplaceAll(action, "-", " ")))
	return nil
}

func performClusterAction(cmd *cobra.Command, action string) error {
	instanceID := currentInstanceID

	instanceAPIKey, err := getInstanceAPIKey(instanceID)
	if err != nil {
		return fmt.Errorf("failed to get instance API key: %w", err)
	}

	c := client.NewInstanceAPI(instanceAPIKey)

	switch action {
	case "restart-cluster":
		err = c.RestartCluster()
	case "stop-cluster":
		err = c.StopCluster()
	case "start-cluster":
		err = c.StartCluster()
	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	if err != nil {
		fmt.Printf("Error performing %s: %v\n", action, err)
		return err
	}

	fmt.Printf("%s initiated successfully.\n", strings.Title(strings.ReplaceAll(action, "-", " ")))
	return nil
}

func performUpgradeAction(cmd *cobra.Command, action, version string) error {
	instanceID := currentInstanceID

	instanceAPIKey, err := getInstanceAPIKey(instanceID)
	if err != nil {
		return fmt.Errorf("failed to get instance API key: %w", err)
	}

	c := client.NewInstanceAPI(instanceAPIKey)

	switch action {
	case "upgrade-erlang":
		err = c.UpgradeErlang()
	case "upgrade-rabbitmq":
		err = c.UpgradeRabbitMQ(version)
	case "upgrade-all":
		err = c.UpgradeRabbitMQErlang()
	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	if err != nil {
		fmt.Printf("Error performing %s: %v\n", action, err)
		return err
	}

	fmt.Printf("%s initiated successfully.\n", strings.Title(strings.ReplaceAll(action, "-", " ")))
	return nil
}

func performToggleAction(cmd *cobra.Command, action string) error {
	instanceID := currentInstanceID

	instanceAPIKey, err := getInstanceAPIKey(instanceID)
	if err != nil {
		return fmt.Errorf("failed to get instance API key: %w", err)
	}

	c := client.NewInstanceAPI(instanceAPIKey)

	enable, _ := cmd.Flags().GetBool("enable")

	switch action {
	case "hipe":
		nodesStr, _ := cmd.Flags().GetString("nodes")
		var nodes []string
		if nodesStr != "" {
			nodes = strings.Split(nodesStr, ",")
		}

		req := &client.HiPERequest{
			Enable: enable,
			Nodes:  nodes,
		}
		err = c.ToggleHiPE(req)

	case "firehose":
		vhost, _ := cmd.Flags().GetString("vhost")
		if vhost == "" {
			return fmt.Errorf("vhost flag is required for firehose")
		}

		req := &client.FirehoseRequest{
			Enable: enable,
			VHost:  vhost,
		}
		err = c.ToggleFirehose(req)

	default:
		return fmt.Errorf("unknown action: %s", action)
	}

	if err != nil {
		fmt.Printf("Error toggling %s: %v\n", action, err)
		return err
	}

	status := "disabled"
	if enable {
		status = "enabled"
	}
	fmt.Printf("%s %s successfully.\n", strings.Title(action), status)
	return nil
}

func init() {
	// Add node flags where applicable
	restartRabbitMQCmd.Flags().String("nodes", "", "Comma-separated list of node names")
	restartManagementCmd.Flags().String("nodes", "", "Comma-separated list of node names")
	stopCmd.Flags().String("nodes", "", "Comma-separated list of node names")
	startCmd.Flags().String("nodes", "", "Comma-separated list of node names")
	rebootCmd.Flags().String("nodes", "", "Comma-separated list of node names")

	// Add version flag for RabbitMQ upgrade
	upgradeRabbitMQCmd.Flags().String("version", "", "RabbitMQ version (required)")
	upgradeRabbitMQCmd.MarkFlagRequired("version")

	// Add flags for toggle commands
	toggleHiPECmd.Flags().Bool("enable", false, "Enable or disable HiPE")
	toggleHiPECmd.Flags().String("nodes", "", "Comma-separated list of node names")
	toggleHiPECmd.MarkFlagRequired("enable")

	toggleFirehoseCmd.Flags().Bool("enable", false, "Enable or disable Firehose")
	toggleFirehoseCmd.Flags().String("vhost", "", "Virtual host to enable tracing on (required)")
	toggleFirehoseCmd.MarkFlagRequired("enable")
	toggleFirehoseCmd.MarkFlagRequired("vhost")

	// Add all commands to actions
	instanceActionsCmd.AddCommand(restartRabbitMQCmd)
	instanceActionsCmd.AddCommand(restartClusterCmd)
	instanceActionsCmd.AddCommand(restartManagementCmd)
	instanceActionsCmd.AddCommand(stopCmd)
	instanceActionsCmd.AddCommand(startCmd)
	instanceActionsCmd.AddCommand(rebootCmd)
	instanceActionsCmd.AddCommand(stopClusterCmd)
	instanceActionsCmd.AddCommand(startClusterCmd)
	instanceActionsCmd.AddCommand(upgradeErlangCmd)
	instanceActionsCmd.AddCommand(upgradeRabbitMQCmd)
	instanceActionsCmd.AddCommand(upgradeRabbitMQErlangCmd)
	instanceActionsCmd.AddCommand(toggleHiPECmd)
	instanceActionsCmd.AddCommand(toggleFirehoseCmd)
	instanceActionsCmd.AddCommand(upgradeVersionsCmd)
}
