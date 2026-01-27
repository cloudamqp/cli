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
		return urlStr // Return original if parsing fails
	}

	if parsedURL.User == nil {
		return urlStr // No user info to mask
	}

	username := parsedURL.User.Username()

	// Manually construct the URL to avoid encoding issues with asterisks
	result := parsedURL.Scheme + "://" + username + ":****@" + parsedURL.Host
	if parsedURL.Path != "" {
		result += parsedURL.Path
	}
	if parsedURL.RawQuery != "" {
		result += "?" + parsedURL.RawQuery
	}
	if parsedURL.Fragment != "" {
		result += "#" + parsedURL.Fragment
	}

	return result
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

		// Format output as "Name = Value"
		fmt.Printf("Name = %s\n", instance.Name)
		fmt.Printf("Plan = %s\n", instance.Plan)
		fmt.Printf("Region = %s\n", instance.Region)
		fmt.Printf("Tags = %s\n", strings.Join(instance.Tags, ","))

		showURL, _ := cmd.Flags().GetBool("show-url")
		if showURL {
			fmt.Printf("URL = %s\n", instance.URL)
		} else {
			fmt.Printf("URL = %s\n", maskPassword(instance.URL))
		}

		fmt.Printf("Hostname = %s\n", instance.HostnameExternal)
		ready := "No"
		if instance.Ready {
			ready = "Yes"
		}
		fmt.Printf("Ready = %s\n", ready)

		return nil
	},
}

func init() {
	instanceGetCmd.Flags().StringP("id", "", "", "Instance ID (required)")
	instanceGetCmd.MarkFlagRequired("id")
	instanceGetCmd.Flags().BoolP("show-url", "", false, "Show full connection URL with credentials")
	instanceGetCmd.RegisterFlagCompletionFunc("id", completeInstanceIDFlag)
}
