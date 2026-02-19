package cmd

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

func maskPassword(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}

	password, _ := parsedURL.User.Password()
	return strings.Replace(urlStr, password, "****", 1)
}

var instanceGetCmd = &cobra.Command{
	Use:     "get --id <id>",
	Short:   "Get details of a specific CloudAMQP instance",
	Long:    `Retrieves and displays detailed information about a specific CloudAMQP instance.`,
	Example: `  cloudamqp instance get --id 1234`,
	RunE: func(cmd *cobra.Command, args []string) error {
		idFlag, _ := cmd.Flags().GetString("id")
		if idFlag == "" {
			return fmt.Errorf("instance ID is required. Use --id flag")
		}

		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		instanceID, err := strconv.Atoi(idFlag)
		if err != nil {
			return fmt.Errorf("invalid instance ID: %v", err)
		}

		c := client.New(apiKey, Version)

		instance, err := c.GetInstance(instanceID)
		if err != nil {
			fmt.Printf("Error getting instance: %v\n", err)
			return err
		}

		p, err := getPrinter(cmd)
		if err != nil {
			return err
		}

		showURL, _ := cmd.Flags().GetBool("show-url")
		ready := "No"
		if instance.Ready {
			ready = "Yes"
		}

		urlVal := maskPassword(instance.URL)
		if showURL {
			urlVal = instance.URL
		}

		p.PrintRecord(
			[]string{"ID", "NAME", "PLAN", "REGION", "TAGS", "URL", "HOSTNAME", "READY"},
			[]string{
				strconv.Itoa(instance.ID),
				instance.Name,
				instance.Plan,
				instance.Region,
				strings.Join(instance.Tags, ","),
				urlVal,
				instance.HostnameExternal,
				ready,
			},
		)

		return nil
	},
}

func init() {
	instanceGetCmd.Flags().StringP("id", "", "", "Instance ID (required)")
	instanceGetCmd.MarkFlagRequired("id")
	instanceGetCmd.Flags().BoolP("show-url", "", false, "Show full connection URL with credentials")
	instanceGetCmd.RegisterFlagCompletionFunc("id", completeInstanceIDFlag)
}
