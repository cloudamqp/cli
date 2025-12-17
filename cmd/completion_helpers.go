package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

// completionAPIKey retrieves the API key without prompting the user
func completionAPIKey() (string, error) {
	// First, check environment variable
	if apiKey := os.Getenv("CLOUDAMQP_APIKEY"); apiKey != "" {
		return apiKey, nil
	}

	// Second, check config file
	apiKey, err := loadAPIKey()
	if apiKey != "" {
		return apiKey, nil
	}
	return "", fmt.Errorf("API key not configured: %w", err)
}

// completeInstances returns a list of instance IDs and names for completion
func completeInstances(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	apiKey, err := completionAPIKey()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	c := client.New(apiKey, Version)

	// Try to get from cache
	var instances []client.Instance
	if cachedData, ok := getCachedData("instances", instancesCacheTTL); ok {
		if err := json.Unmarshal(cachedData, &instances); err == nil {
			goto formatOutput
		}
	}

	// Cache miss or error, fetch from API
	instances, err = c.ListInstances()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Store in cache
	setCachedData("instances", instancesCacheTTL, instances)

formatOutput:
	var suggestions []string
	for _, instance := range instances {
		suggestions = append(suggestions, fmt.Sprintf("%d\t%s", instance.ID, instance.Name))
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

// completePlans returns a list of plan names for completion
func completePlans(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	apiKey, err := completionAPIKey()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	c := client.New(apiKey, Version)

	// Try to get from cache
	var plans []client.Plan
	if cachedData, ok := getCachedData("plans", plansCacheTTL); ok {
		if err := json.Unmarshal(cachedData, &plans); err == nil {
			goto formatOutput
		}
	}

	// Cache miss or error, fetch from API
	plans, err = c.ListPlans("")
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Store in cache
	setCachedData("plans", plansCacheTTL, plans)

formatOutput:
	var suggestions []string
	for _, plan := range plans {
		suggestions = append(suggestions, fmt.Sprintf("%s\t%s", plan.Name, plan.Backend))
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

// completeRegions returns a list of region identifiers for completion
func completeRegions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	apiKey, err := completionAPIKey()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	c := client.New(apiKey, Version)

	// Try to get from cache
	var regions []client.Region
	if cachedData, ok := getCachedData("regions", regionsCacheTTL); ok {
		if err := json.Unmarshal(cachedData, &regions); err == nil {
			goto formatOutput
		}
	}

	// Cache miss or error, fetch from API
	regions, err = c.ListRegions("")
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Store in cache
	setCachedData("regions", regionsCacheTTL, regions)

formatOutput:
	var suggestions []string
	for _, region := range regions {
		fullRegion := fmt.Sprintf("%s::%s", region.Provider, region.Region)
		suggestions = append(suggestions, fmt.Sprintf("%s\t%s", fullRegion, region.Name))
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

// completeVPCs returns a list of VPC IDs and names for completion
func completeVPCs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	apiKey, err := completionAPIKey()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	c := client.New(apiKey, Version)

	// Try to get from cache
	var vpcs []client.VPC
	if cachedData, ok := getCachedData("vpcs", vpcsCacheTTL); ok {
		if err := json.Unmarshal(cachedData, &vpcs); err == nil {
			goto formatOutput
		}
	}

	// Cache miss or error, fetch from API
	vpcs, err = c.ListVPCs()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Store in cache
	setCachedData("vpcs", vpcsCacheTTL, vpcs)

formatOutput:
	var suggestions []string
	for _, vpc := range vpcs {
		suggestions = append(suggestions, fmt.Sprintf("%d\t%s (%s)", vpc.ID, vpc.Name, vpc.Region))
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

// completeCopySettings returns the valid copy-settings options
func completeCopySettings(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	settings := []string{
		"alarms\tCopy alarm settings",
		"metrics\tCopy metrics settings",
		"logs\tCopy logs settings",
		"firewall\tCopy firewall settings",
		"config\tCopy configuration settings",
	}
	return settings, cobra.ShellCompDirectiveNoFileComp
}

// completeInstanceIDFlag is a wrapper for instance ID flag completion
func completeInstanceIDFlag(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return completeInstances(cmd, args, toComplete)
}

// completeVPCIDFlag is a wrapper for VPC ID flag completion
func completeVPCIDFlag(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	apiKey, err := completionAPIKey()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	c := client.New(apiKey, Version)

	// Try to get from cache
	var vpcs []client.VPC
	if cachedData, ok := getCachedData("vpcs", vpcsCacheTTL); ok {
		if err := json.Unmarshal(cachedData, &vpcs); err == nil {
			goto formatOutput
		}
	}

	// Cache miss or error, fetch from API
	vpcs, err = c.ListVPCs()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Store in cache
	setCachedData("vpcs", vpcsCacheTTL, vpcs)

formatOutput:
	var suggestions []string
	for _, vpc := range vpcs {
		suggestions = append(suggestions, fmt.Sprintf("%d\t%s", vpc.ID, vpc.Name))
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

// completeCopyFromIDFlag completes instance IDs for the --copy-from-id flag
func completeCopyFromIDFlag(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	apiKey, err := completionAPIKey()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	c := client.New(apiKey, Version)

	// Try to get from cache
	var instances []client.Instance
	if cachedData, ok := getCachedData("instances", instancesCacheTTL); ok {
		if err := json.Unmarshal(cachedData, &instances); err == nil {
			goto formatOutput
		}
	}

	// Cache miss or error, fetch from API
	instances, err = c.ListInstances()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Store in cache
	setCachedData("instances", instancesCacheTTL, instances)

formatOutput:
	var suggestions []string
	for _, instance := range instances {
		suggestions = append(suggestions, fmt.Sprintf("%s\t%s", strconv.Itoa(instance.ID), instance.Name))
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}

// completeVPCArgs returns a list of VPC IDs for positional argument completion
func completeVPCArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	apiKey, err := completionAPIKey()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	c := client.New(apiKey, Version)

	// Try to get from cache
	var vpcs []client.VPC
	if cachedData, ok := getCachedData("vpcs", vpcsCacheTTL); ok {
		if err := json.Unmarshal(cachedData, &vpcs); err == nil {
			goto formatOutput
		}
	}

	// Cache miss or error, fetch from API
	vpcs, err = c.ListVPCs()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Store in cache
	setCachedData("vpcs", vpcsCacheTTL, vpcs)

formatOutput:
	var suggestions []string
	for _, vpc := range vpcs {
		suggestions = append(suggestions, fmt.Sprintf("%d\t%s (%s)", vpc.ID, vpc.Name, vpc.Region))
	}

	return suggestions, cobra.ShellCompDirectiveNoFileComp
}
