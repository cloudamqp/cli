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
	Use:   "restart-rabbitmq --id <instance_id>",
	Short: "Restart RabbitMQ",
	Long:  `Restart RabbitMQ on specified nodes or all nodes.`,
	Example: `  cloudamqp instance restart-rabbitmq --id 1234
  cloudamqp instance restart-rabbitmq --id 1234 --nodes=node1,node2`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, "restart-rabbitmq")
	},
}

var restartClusterCmd = &cobra.Command{
	Use:     "restart-cluster --id <instance_id>",
	Short:   "Restart cluster",
	Long:    `Restart the entire cluster.`,
	Example: `  cloudamqp instance restart-cluster --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performClusterAction(cmd, "restart-cluster")
	},
}

var restartManagementCmd = &cobra.Command{
	Use:     "restart-management --id <instance_id>",
	Short:   "Restart management interface",
	Long:    `Restart the RabbitMQ management interface.`,
	Example: `  cloudamqp instance restart-management --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, "restart-management")
	},
}

// Stop/Start commands
var stopCmd = &cobra.Command{
	Use:     "stop --id <instance_id>",
	Short:   "Stop instance",
	Long:    `Stop specified nodes or all nodes.`,
	Example: `  cloudamqp instance stop --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, "stop")
	},
}

var startCmd = &cobra.Command{
	Use:     "start --id <instance_id>",
	Short:   "Start instance",
	Long:    `Start specified nodes or all nodes.`,
	Example: `  cloudamqp instance start --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, "start")
	},
}

var rebootCmd = &cobra.Command{
	Use:     "reboot --id <instance_id>",
	Short:   "Reboot instance",
	Long:    `Reboot specified nodes or all nodes.`,
	Example: `  cloudamqp instance reboot --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performNodeAction(cmd, "reboot")
	},
}

// Cluster commands
var stopClusterCmd = &cobra.Command{
	Use:     "stop-cluster --id <instance_id>",
	Short:   "Stop cluster",
	Long:    `Stop the entire cluster.`,
	Example: `  cloudamqp instance stop-cluster --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performClusterAction(cmd, "stop-cluster")
	},
}

var startClusterCmd = &cobra.Command{
	Use:     "start-cluster --id <instance_id>",
	Short:   "Start cluster",
	Long:    `Start the entire cluster.`,
	Example: `  cloudamqp instance start-cluster --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performClusterAction(cmd, "start-cluster")
	},
}

// Upgrade commands
var upgradeErlangCmd = &cobra.Command{
	Use:   "upgrade-erlang --id <instance_id>",
	Short: "Upgrade Erlang",
	Long: `Always updates to latest compatible version.

Note: This action is asynchronous. The request will return immediately, the process runs in the background.`,
	Example: `  cloudamqp instance upgrade-erlang --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performUpgradeAction(cmd, "upgrade-erlang", "")
	},
}

var upgradeRabbitMQCmd = &cobra.Command{
	Use:   "upgrade-rabbitmq --id <instance_id>",
	Short: "Upgrade RabbitMQ",
	Long: `Upgrade RabbitMQ to specified version.

Note: This action is asynchronous. The request will return immediately, the process runs in the background.`,
	Example: `  cloudamqp instance upgrade-rabbitmq --id 1234 --version=3.10.7`,
	RunE: func(cmd *cobra.Command, args []string) error {
		version, _ := cmd.Flags().GetString("version")
		if version == "" {
			return fmt.Errorf("version flag is required")
		}
		return performUpgradeAction(cmd, "upgrade-rabbitmq", version)
	},
}

var upgradeRabbitMQErlangCmd = &cobra.Command{
	Use:   "upgrade-all --id <instance_id>",
	Short: "Upgrade RabbitMQ and Erlang",
	Long: `Always updates to latest possible version of both RabbitMQ and Erlang.

Note: This action is asynchronous. The request will return immediately, the process runs in the background.`,
	Example: `  cloudamqp instance upgrade-all --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performUpgradeAction(cmd, "upgrade-all", "")
	},
}

// HiPE and Firehose commands
var toggleHiPECmd = &cobra.Command{
	Use:     "toggle-hipe --id <instance_id>",
	Short:   "Enable/disable HiPE",
	Long:    `Enable or disable HiPE (High Performance Erlang) compilation.`,
	Example: `  cloudamqp instance toggle-hipe --id 1234 --enable=true`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performToggleAction(cmd, "hipe")
	},
}

var toggleFirehoseCmd = &cobra.Command{
	Use:     "toggle-firehose --id <instance_id>",
	Short:   "Enable/disable Firehose",
	Long:    `Enable or disable RabbitMQ Firehose tracing (not recommended in production).`,
	Example: `  cloudamqp instance toggle-firehose --id 1234 --enable=true --vhost=/`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return performToggleAction(cmd, "firehose")
	},
}

var upgradeVersionsCmd = &cobra.Command{
	Use:     "upgrade-versions --id <instance_id>",
	Short:   "Fetch upgrade versions",
	Long:    `Returns what version of Erlang and RabbitMQ the cluster will update to.`,
	Example: `  cloudamqp instance upgrade-versions --id 1234`,
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

		c := client.New(apiKey, Version)

		versions, err := c.GetUpgradeVersions(idFlag)
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
	idFlag, _ := cmd.Flags().GetString("id")
	if idFlag == "" {
		return fmt.Errorf("instance ID is required. Use --id flag")
	}

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
		err = c.RestartRabbitMQ(idFlag, nodes)
	case "restart-management":
		err = c.RestartManagement(idFlag, nodes)
	case "stop":
		err = c.StopInstance(idFlag, nodes)
	case "start":
		err = c.StartInstance(idFlag, nodes)
	case "reboot":
		err = c.RebootInstance(idFlag, nodes)
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
	idFlag, _ := cmd.Flags().GetString("id")
	if idFlag == "" {
		return fmt.Errorf("instance ID is required. Use --id flag")
	}

	var err error
	apiKey, err := getAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	c := client.New(apiKey, Version)

	switch action {
	case "restart-cluster":
		err = c.RestartCluster(idFlag)
	case "stop-cluster":
		err = c.StopCluster(idFlag)
	case "start-cluster":
		err = c.StartCluster(idFlag)
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
	idFlag, _ := cmd.Flags().GetString("id")
	if idFlag == "" {
		return fmt.Errorf("instance ID is required. Use --id flag")
	}

	var err error
	apiKey, err := getAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	c := client.New(apiKey, Version)

	switch action {
	case "upgrade-erlang":
		err = c.UpgradeErlang(idFlag)
	case "upgrade-rabbitmq":
		err = c.UpgradeRabbitMQ(idFlag, version)
	case "upgrade-all":
		err = c.UpgradeRabbitMQErlang(idFlag)
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
	idFlag, _ := cmd.Flags().GetString("id")
	if idFlag == "" {
		return fmt.Errorf("instance ID is required. Use --id flag")
	}

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
		err = c.ToggleHiPE(idFlag, req)

	case "firehose":
		vhost, _ := cmd.Flags().GetString("vhost")
		if vhost == "" {
			return fmt.Errorf("vhost flag is required for firehose")
		}

		req := &client.FirehoseRequest{
			Enable: enable,
			VHost:  vhost,
		}
		err = c.ToggleFirehose(idFlag, req)

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
	// Add --id flag to all action commands
	commands := []*cobra.Command{
		restartRabbitMQCmd, restartClusterCmd, restartManagementCmd,
		stopCmd, startCmd, rebootCmd,
		stopClusterCmd, startClusterCmd,
		upgradeErlangCmd, upgradeRabbitMQCmd, upgradeRabbitMQErlangCmd,
		toggleHiPECmd, toggleFirehoseCmd, upgradeVersionsCmd,
	}

	for _, cmd := range commands {
		cmd.Flags().StringP("id", "", "", "Instance ID (required)")
		cmd.MarkFlagRequired("id")
		cmd.RegisterFlagCompletionFunc("id", completeInstanceIDFlag)
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
