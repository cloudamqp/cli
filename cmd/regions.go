package cmd

import (
	"fmt"
	"os"

	"cloudamqp-cli/client"
	"cloudamqp-cli/internal/table"
	"github.com/spf13/cobra"
)

var providerFilter string

var regionsCmd = &cobra.Command{
	Use:   "regions",
	Short: "List available regions",
	Long:  `Retrieves all available regions, optionally filtered by provider.`,
	Example: `  cloudamqp regions
  cloudamqp regions --provider=amazon-web-services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey)

		regions, err := c.ListRegions(providerFilter)
		if err != nil {
			fmt.Printf("Error listing regions: %v\n", err)
			return err
		}

		if len(regions) == 0 {
			fmt.Println("No regions found.")
			return nil
		}

		t := table.New(os.Stdout, "PROVIDER", "REGION", "NAME")
		for _, region := range regions {
			t.AddRow(region.Provider, region.Region, region.Name)
		}
		t.Print()

		return nil
	},
}

func init() {
	regionsCmd.Flags().StringVar(&providerFilter, "provider", "", "Filter by specific provider (e.g., amazon-web-services)")
}
