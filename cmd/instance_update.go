package cmd

import (
	"fmt"
	"strconv"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var (
	updateInstanceID   string
	updateInstanceName string
	updateInstancePlan string
	updateInstanceTags []string
)

var instanceUpdateCmd = &cobra.Command{
	Use:   "update --id <id>",
	Short: "Update a CloudAMQP instance",
	Long: `Update an existing CloudAMQP instance with new configuration.

You can update the following fields:
  --name: Instance name
  --plan: Subscription plan
  --tags: Instance tags (replaces existing tags)`,
	Example: `  cloudamqp instance update --id 1234 --name=new-name
  cloudamqp instance update --id 1234 --plan=rabbit-1
  cloudamqp instance update --id 1234 --tags=production --tags=updated`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		if updateInstanceID == "" {
			return fmt.Errorf("--id is required")
		}

		instanceID, err := strconv.Atoi(updateInstanceID)
		if err != nil {
			return fmt.Errorf("invalid instance ID: %v", err)
		}

		c := client.New(apiKey, Version)

		req := &client.InstanceUpdateRequest{
			Name: updateInstanceName,
			Plan: updateInstancePlan,
			Tags: updateInstanceTags,
		}

		if req.Name == "" && req.Plan == "" && len(req.Tags) == 0 {
			return fmt.Errorf("at least one field must be specified for update")
		}

		err = c.UpdateInstance(instanceID, req)
		if err != nil {
			fmt.Printf("Error updating instance: %v\n", err)
			return err
		}

		fmt.Printf("Instance %d updated successfully.\n", instanceID)
		return nil
	},
}

func init() {
	instanceUpdateCmd.Flags().StringVar(&updateInstanceID, "id", "", "Instance ID (required)")
	instanceUpdateCmd.Flags().StringVar(&updateInstanceName, "name", "", "New instance name")
	instanceUpdateCmd.Flags().StringVar(&updateInstancePlan, "plan", "", "New subscription plan")
	instanceUpdateCmd.Flags().StringSliceVar(&updateInstanceTags, "tags", []string{}, "New instance tags")
	instanceUpdateCmd.MarkFlagRequired("id")
	instanceUpdateCmd.RegisterFlagCompletionFunc("id", completeInstances)
	instanceUpdateCmd.RegisterFlagCompletionFunc("plan", completePlans)
}
