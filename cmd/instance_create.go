package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var (
	instanceName         string
	instancePlan         string
	instanceRegion       string
	instanceTags         []string
	instanceVPCSubnet    string
	instanceVPCID        string
	instanceCopyFromID   string
	instanceCopySettings []string
	instanceWait         bool
	instanceWaitTimeout  string
)

var instanceCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new CloudAMQP instance",
	Long: `Create a new CloudAMQP instance with the specified configuration.

Required flags:
  --name: Name of the instance
  --plan: Subscription plan (e.g., lemming, bunny-1, rabbit-1)
  --region: Region identifier (e.g., amazon-web-services::us-east-1)

Optional flags:
  --tags: Instance tags (can be specified multiple times)
  --vpc-subnet: VPC subnet for dedicated VPC
  --vpc-id: ID of existing VPC to add instance to
  --copy-from-id: Instance ID to copy settings from (dedicated instances only)
  --copy-settings: Settings to copy (alarms, metrics, logs, firewall, config)
  --wait: Wait for instance to be ready before returning
  --wait-timeout: Timeout for waiting (default: 15m)`,
	Example: `  cloudamqp instance create --name=my-instance --plan=bunny-1 --region=amazon-web-services::us-east-1
  cloudamqp instance create --name=my-instance --plan=bunny-1 --region=amazon-web-services::us-east-1 --tags=production --tags=web-app
  cloudamqp instance create --name=my-copy --plan=bunny-1 --region=amazon-web-services::us-east-1 --copy-from-id=12345 --copy-settings=metrics,firewall
  cloudamqp instance create --name=my-instance --plan=bunny-1 --region=amazon-web-services::us-east-1 --wait`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		req := &client.InstanceCreateRequest{
			Name:   instanceName,
			Plan:   instancePlan,
			Region: instanceRegion,
			Tags:   instanceTags,
		}

		if instanceVPCSubnet != "" {
			req.VPCSubnet = instanceVPCSubnet
		}

		if instanceVPCID != "" {
			vpcID, err := strconv.Atoi(instanceVPCID)
			if err != nil {
				return fmt.Errorf("invalid VPC ID: %v", err)
			}
			req.VPCID = &vpcID
		}

		if instanceCopyFromID != "" {
			copyFromID, err := strconv.Atoi(instanceCopyFromID)
			if err != nil {
				return fmt.Errorf("invalid copy-from-id: %v", err)
			}
			req.CopySettings = &client.CopySettings{
				SubscriptionID: copyFromID,
				Settings:       instanceCopySettings,
			}
		}

		resp, err := c.CreateInstance(req)
		if err != nil {
			fmt.Printf("Error creating instance: %v\n", err)
			return err
		}

		if instanceWait {
			timeout, err := time.ParseDuration(instanceWaitTimeout)
			if err != nil {
				return fmt.Errorf("invalid wait-timeout value: %v", err)
			}

			if err := waitForInstanceReady(c, resp.ID, timeout); err != nil {
				// Instance was created but failed to become ready
				output, _ := json.MarshalIndent(resp, "", "  ")
				fmt.Printf("Instance created but not ready:\n%s\n", string(output))
				return fmt.Errorf("wait failed: %w", err)
			}
		}

		output, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format response: %v", err)
		}

		fmt.Printf("Instance created successfully:\n%s\n", string(output))
		return nil
	},
}

func init() {
	instanceCreateCmd.Flags().StringVar(&instanceName, "name", "", "Name of the instance (required)")
	instanceCreateCmd.Flags().StringVar(&instancePlan, "plan", "", "Subscription plan (required)")
	instanceCreateCmd.Flags().StringVar(&instanceRegion, "region", "", "Region identifier (required)")
	instanceCreateCmd.Flags().StringSliceVar(&instanceTags, "tags", []string{}, "Instance tags")
	instanceCreateCmd.Flags().StringVar(&instanceVPCSubnet, "vpc-subnet", "", "VPC subnet")
	instanceCreateCmd.Flags().StringVar(&instanceVPCID, "vpc-id", "", "VPC ID")
	instanceCreateCmd.Flags().StringVar(&instanceCopyFromID, "copy-from-id", "", "Instance ID to copy settings from")
	instanceCreateCmd.Flags().StringSliceVar(&instanceCopySettings, "copy-settings", []string{}, "Settings to copy (alarms, metrics, logs, firewall, config)")
	instanceCreateCmd.Flags().BoolVar(&instanceWait, "wait", false, "Wait for instance to be ready")
	instanceCreateCmd.Flags().StringVar(&instanceWaitTimeout, "wait-timeout", "15m", "Timeout for waiting (e.g., 15m, 30m)")

	instanceCreateCmd.MarkFlagRequired("name")
	instanceCreateCmd.MarkFlagRequired("plan")
	instanceCreateCmd.MarkFlagRequired("region")

	instanceCreateCmd.RegisterFlagCompletionFunc("plan", completePlans)
	instanceCreateCmd.RegisterFlagCompletionFunc("region", completeRegions)
	instanceCreateCmd.RegisterFlagCompletionFunc("vpc-id", completeVPCIDFlag)
	instanceCreateCmd.RegisterFlagCompletionFunc("copy-from-id", completeCopyFromIDFlag)
	instanceCreateCmd.RegisterFlagCompletionFunc("copy-settings", completeCopySettings)
}
