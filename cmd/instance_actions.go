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
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		cmd.SilenceUsage = true
		return fmt.Errorf("subcommand required")
	},
}

// Restart commands
var restartRabbitMQCmd = &cobra.Command{
	Use:   "restart-rabbitmq <instance_id>",
	Short: "Restart RabbitMQ",
	Long:  `Restart RabbitMQ on specified nodes or all nodes.`,
	Example: `  cloudamqp instance restart-rabbitmq 1234
  cloudamqp instance restart-rabbitmq 1234 --nodes=node1,node2`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, args[0], "restart-rabbitmq")
	},
}

var restartClusterCmd = &cobra.Command{
	Use:     "restart-cluster <instance_id>",
	Short:   "Restart cluster",
	Long:    `Restart the entire cluster.`,
	Example: `  cloudamqp instance restart-cluster 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return performClusterAction(cmd, args[0], "restart-cluster")
	},
}

var restartManagementCmd = &cobra.Command{
	Use:     "restart-management <instance_id>",
	Short:   "Restart management interface",
	Long:    `Restart the RabbitMQ management interface.`,
	Example: `  cloudamqp instance restart-management 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, args[0], "restart-management")
	},
}

// Stop/Start commands
var stopCmd = &cobra.Command{
	Use:     "stop <instance_id>",
	Short:   "Stop instance",
	Long:    `Stop specified nodes or all nodes.`,
	Example: `  cloudamqp instance stop 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, args[0], "stop")
	},
}

var startCmd = &cobra.Command{
	Use:     "start <instance_id>",
	Short:   "Start instance",
	Long:    `Start specified nodes or all nodes.`,
	Example: `  cloudamqp instance start 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, args[0], "start")
	},
}

var rebootCmd = &cobra.Command{
	Use:     "reboot <instance_id>",
	Short:   "Reboot instance",
	Long:    `Reboot specified nodes or all nodes.`,
	Example: `  cloudamqp instance reboot 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, args[0], "reboot")
	},
}

// Cluster commands
var stopClusterCmd = &cobra.Command{
	Use:     "stop-cluster <instance_id>",
	Short:   "Stop cluster",
	Long:    `Stop the entire cluster.`,
	Example: `  cloudamqp instance stop-cluster 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return performClusterAction(cmd, args[0], "stop-cluster")
	},
}

var startClusterCmd = &cobra.Command{
	Use:     "start-cluster <instance_id>",
	Short:   "Start cluster",
	Long:    `Start the entire cluster.`,
	Example: `  cloudamqp instance start-cluster 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return performClusterAction(cmd, args[0], "start-cluster")
	},
}

// Upgrade commands
var upgradeErlangCmd = &cobra.Command{
	Use:   "upgrade-erlang <instance_id>",
	Short: "Upgrade Erlang",
	Long: `Always updates to latest compatible version.

Note: This action is asynchronous. The request will return immediately, the process runs in the background.`,
	Example: `  cloudamqp instance upgrade-erlang 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return performUpgradeAction(cmd, args[0], "upgrade-erlang", "")
	},
}

var upgradeRabbitMQCmd = &cobra.Command{
	Use:   "upgrade-rabbitmq <instance_id>",
	Short: "Upgrade RabbitMQ",
	Long: `Upgrade RabbitMQ to specified version.

Note: This action is asynchronous. The request will return immediately, the process runs in the background.`,
	Example: `  cloudamqp instance upgrade-rabbitmq 1234 --version=3.10.7`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version, _ := cmd.Flags().GetString("version")
		if version == "" {
			return fmt.Errorf("version flag is required")
		}
		return performUpgradeAction(cmd, args[0], "upgrade-rabbitmq", version)
	},
}

var upgradeRabbitMQErlangCmd = &cobra.Command{
	Use:   "upgrade-all <instance_id>",
	Short: "Upgrade RabbitMQ and Erlang",
	Long: `Always updates to latest possible version of both RabbitMQ and Erlang.

Note: This action is asynchronous. The request will return immediately, the process runs in the background.`,
	Example: `  cloudamqp instance upgrade-all 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return performUpgradeAction(cmd, args[0], "upgrade-all", "")
	},
}

// HiPE and Firehose commands
var toggleHiPECmd = &cobra.Command{
	Use:     "toggle-hipe <instance_id>",
	Short:   "Enable/disable HiPE",
	Long:    `Enable or disable HiPE (High Performance Erlang) compilation.`,
	Example: `  cloudamqp instance toggle-hipe 1234 --enable=true`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return performToggleAction(cmd, args[0], "hipe")
	},
}

var toggleFirehoseCmd = &cobra.Command{
	Use:     "toggle-firehose <instance_id>",
	Short:   "Enable/disable Firehose",
	Long:    `Enable or disable RabbitMQ Firehose tracing (not recommended in production).`,
	Example: `  cloudamqp instance toggle-firehose 1234 --enable=true --vhost=/`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return performToggleAction(cmd, args[0], "firehose")
	},
}

var upgradeVersionsCmd = &cobra.Command{
	Use:     "upgrade-versions <instance_id>",
	Short:   "Fetch upgrade versions",
	Long:    `Returns what version of Erlang and RabbitMQ the cluster will update to.`,
	Example: `  cloudamqp instance upgrade-versions 1234`,
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err := getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		versions, err := c.GetUpgradeVersions(args[0])
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
func performNodeAction(cmd *cobra.Command, instanceID, action string) error {
	var err error
	apiKey, err := getAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	c := client.New(apiKey, Version)

	nodesStr, _ := cmd.Flags().GetString("nodes")
	var nodes []string
	if nodesStr != "" {
		nodes = strings.Split(nodesStr, ",")
	}

	switch action {
	case "restart-rabbitmq":
		err = c.RestartRabbitMQ(instanceID, nodes)
	case "restart-management":
		err = c.RestartManagement(instanceID, nodes)
	case "stop":
		err = c.StopInstance(instanceID, nodes)
	case "start":
		err = c.StartInstance(instanceID, nodes)
	case "reboot":
		err = c.RebootInstance(instanceID, nodes)
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

func performClusterAction(cmd *cobra.Command, instanceID, action string) error {
	var err error
	apiKey, err := getAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	c := client.New(apiKey, Version)

	switch action {
	case "restart-cluster":
		err = c.RestartCluster(instanceID)
	case "stop-cluster":
		err = c.StopCluster(instanceID)
	case "start-cluster":
		err = c.StartCluster(instanceID)
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

func performUpgradeAction(cmd *cobra.Command, instanceID, action, version string) error {
	var err error
	apiKey, err := getAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	c := client.New(apiKey, Version)

	switch action {
	case "upgrade-erlang":
		err = c.UpgradeErlang(instanceID)
	case "upgrade-rabbitmq":
		err = c.UpgradeRabbitMQ(instanceID, version)
	case "upgrade-all":
		err = c.UpgradeRabbitMQErlang(instanceID)
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

func performToggleAction(cmd *cobra.Command, instanceID, action string) error {
	var err error
	apiKey, err := getAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	c := client.New(apiKey, Version)

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
		err = c.ToggleHiPE(instanceID, req)

	case "firehose":
		vhost, _ := cmd.Flags().GetString("vhost")
		if vhost == "" {
			return fmt.Errorf("vhost flag is required for firehose")
		}

		req := &client.FirehoseRequest{
			Enable: enable,
			VHost:  vhost,
		}
		err = c.ToggleFirehose(instanceID, req)

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
	// Add completion for all action commands
	for _, cmd := range []*cobra.Command{
		restartRabbitMQCmd, restartClusterCmd, restartManagementCmd,
		stopCmd, startCmd, rebootCmd,
		stopClusterCmd, startClusterCmd,
		upgradeErlangCmd, upgradeRabbitMQCmd, upgradeRabbitMQErlangCmd,
		toggleHiPECmd, toggleFirehoseCmd, upgradeVersionsCmd,
	} {
		cmd.ValidArgsFunction = completeInstances
	}

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
