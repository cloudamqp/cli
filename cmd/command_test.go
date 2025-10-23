package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	cmd := rootCmd

	assert.Equal(t, "cloudamqp", cmd.Use)
	assert.Contains(t, cmd.Long, "CloudAMQP API")
	assert.Contains(t, cmd.Long, "CLOUDAMQP_APIKEY environment variable")
	assert.Contains(t, cmd.Long, "~/.cloudamqprc file")
}

func TestInstanceCommand(t *testing.T) {
	cmd := instanceCmd

	assert.Equal(t, "instance", cmd.Use)
	assert.Equal(t, "Manage CloudAMQP instances", cmd.Short)

	// Check subcommands are present
	subcommands := cmd.Commands()
	commandNames := make([]string, len(subcommands))
	for i, subcmd := range subcommands {
		commandNames[i] = subcmd.Use
	}

	assert.Contains(t, commandNames, "create")
	assert.Contains(t, commandNames, "list")
	assert.Contains(t, commandNames, "get <id>")
	assert.Contains(t, commandNames, "update <id>")
	assert.Contains(t, commandNames, "delete <id>")
	assert.Contains(t, commandNames, "resize <id>")
	assert.Contains(t, commandNames, "manage <instance_id>")
}

func TestVPCCommand(t *testing.T) {
	cmd := vpcCmd

	assert.Equal(t, "vpc", cmd.Use)
	assert.Equal(t, "Manage CloudAMQP VPCs", cmd.Short)

	// Check subcommands are present
	subcommands := cmd.Commands()
	commandNames := make([]string, len(subcommands))
	for i, subcmd := range subcommands {
		commandNames[i] = subcmd.Use
	}

	assert.Contains(t, commandNames, "create")
	assert.Contains(t, commandNames, "list")
	assert.Contains(t, commandNames, "get <id>")
	assert.Contains(t, commandNames, "update <id>")
	assert.Contains(t, commandNames, "delete <id>")
}

func TestInstanceCreateCommand_Validation(t *testing.T) {
	cmd := instanceCreateCmd

	// Test required flags
	requiredFlags := []string{"name", "plan", "region"}
	for _, flagName := range requiredFlags {
		flag := cmd.Flag(flagName)
		assert.NotNil(t, flag, "Flag %s should exist", flagName)
	}

	// Test optional flags
	optionalFlags := []string{"tags", "vpc-subnet", "vpc-id"}
	for _, flagName := range optionalFlags {
		flag := cmd.Flag(flagName)
		assert.NotNil(t, flag, "Flag %s should exist", flagName)
	}
}

func TestInstanceResizeCommand_Validation(t *testing.T) {
	cmd := instanceResizeCmd

	// Test required flags
	diskSizeFlag := cmd.Flag("disk-size")
	assert.NotNil(t, diskSizeFlag)

	// Test optional flags
	downtimeFlag := cmd.Flag("allow-downtime")
	assert.NotNil(t, downtimeFlag)
}

func TestVPCCreateCommand_Validation(t *testing.T) {
	cmd := vpcCreateCmd

	// Test required flags
	requiredFlags := []string{"name", "region", "subnet"}
	for _, flagName := range requiredFlags {
		flag := cmd.Flag(flagName)
		assert.NotNil(t, flag, "Flag %s should exist", flagName)
	}

	// Test optional flags
	tagsFlag := cmd.Flag("tags")
	assert.NotNil(t, tagsFlag)
}

func TestTeamInviteCommand_Validation(t *testing.T) {
	cmd := teamInviteCmd

	// Test required flags
	emailFlag := cmd.Flag("email")
	assert.NotNil(t, emailFlag)

	// Test optional flags
	roleFlag := cmd.Flag("role")
	assert.NotNil(t, roleFlag)

	tagsFlag := cmd.Flag("tags")
	assert.NotNil(t, tagsFlag)
}

func TestCommandHelp(t *testing.T) {
	tests := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"root", rootCmd},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that help text is accessible
			assert.NotEmpty(t, tt.cmd.Use)
			assert.NotEmpty(t, tt.cmd.Short)

			// For root command, test the long description
			if tt.name == "root" {
				assert.Contains(t, tt.cmd.Long, "CloudAMQP API")
			}
		})
	}
}

func TestEnvironmentVariablePrecedence(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "cloudamqp-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Override home directory for test
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	// Save config file with one key
	err = saveMainAPIKey("file-key")
	assert.NoError(t, err)

	// Set environment variable with different key
	os.Setenv("CLOUDAMQP_APIKEY", "env-key")
	defer os.Unsetenv("CLOUDAMQP_APIKEY")

	// Environment variable should take precedence
	apiKey, err := getAPIKey()
	assert.NoError(t, err)
	assert.Equal(t, "env-key", apiKey)
}

func TestInstanceManageCommand(t *testing.T) {
	cmd := instanceManageCmd

	assert.Equal(t, "manage <instance_id>", cmd.Use)
	assert.Contains(t, cmd.Long, "instance-specific API")
	assert.Contains(t, cmd.Long, "Instance API keys are automatically saved")

	// Check that it has args requirement (we can't directly compare functions)
	assert.NotNil(t, cmd.Args)

	// Check subcommands
	subcommands := cmd.Commands()
	commandNames := make([]string, len(subcommands))
	for i, subcmd := range subcommands {
		commandNames[i] = subcmd.Use
	}

	assert.Contains(t, commandNames, "nodes")
	assert.Contains(t, commandNames, "plugins")
	assert.Contains(t, commandNames, "actions")
	assert.Contains(t, commandNames, "account")
}

func TestInstanceActionsCommand(t *testing.T) {
	cmd := instanceActionsCmd

	// Check that actions command has all expected subcommands
	subcommands := cmd.Commands()
	commandNames := make([]string, len(subcommands))
	for i, subcmd := range subcommands {
		commandNames[i] = subcmd.Use
	}

	expectedActions := []string{
		"restart-rabbitmq",
		"restart-cluster",
		"restart-management",
		"stop",
		"start",
		"reboot",
		"stop-cluster",
		"start-cluster",
		"upgrade-erlang",
		"upgrade-rabbitmq",
		"upgrade-all",
		"toggle-hipe",
		"toggle-firehose",
		"upgrade-versions",
	}

	for _, action := range expectedActions {
		assert.Contains(t, commandNames, action, "Action %s should be available", action)
	}
}

func TestUpgradeRabbitMQCommand_RequiredFlag(t *testing.T) {
	cmd := upgradeRabbitMQCmd

	// Check that version flag is required
	versionFlag := cmd.Flag("version")
	assert.NotNil(t, versionFlag)

	// This would normally be tested with actual command execution,
	// but that requires complex mocking of the API client
}

func TestToggleCommands_RequiredFlags(t *testing.T) {
	// Test HiPE toggle command
	hipeCmd := toggleHiPECmd
	enableFlag := hipeCmd.Flag("enable")
	assert.NotNil(t, enableFlag)

	nodesFlag := hipeCmd.Flag("nodes")
	assert.NotNil(t, nodesFlag)

	// Test Firehose toggle command
	firehoseCmd := toggleFirehoseCmd
	enableFlag = firehoseCmd.Flag("enable")
	assert.NotNil(t, enableFlag)

	vhostFlag := firehoseCmd.Flag("vhost")
	assert.NotNil(t, vhostFlag)
}

func TestCommandExamples(t *testing.T) {
	tests := []struct {
		name string
		cmd  *cobra.Command
	}{
		{"instance-create", instanceCreateCmd},
		{"vpc-create", vpcCreateCmd},
		{"team-invite", teamInviteCmd},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.cmd.Example, "Command %s should have examples", tt.name)
			assert.Contains(t, tt.cmd.Example, "cloudamqp", "Example should contain cloudamqp command")
		})
	}
}
