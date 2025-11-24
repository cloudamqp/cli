package cmd

import (
	"fmt"
	"strconv"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var (
	resizeInstanceID string
	diskSize         int
	allowDowntime    bool
)

var instanceResizeCmd = &cobra.Command{
	Use:   "resize-disk --id <id>",
	Short: "Resize instance disk",
	Long: `Resize the disk size of an instance. Default behavior is to expand the disk without any downtime.
Currently limited to instances in Amazon Web Services (AWS) and Google Compute Engine (GCE).

Note: This action is asynchronous. The request will return almost immediately. The disk resize runs in the background.

Note: Due to restrictions from cloud providers, it's only possible to resize the disk every 8 hours unless --allow-downtime is set.

Available disk sizes: 0, 25, 50, 100, 250, 500, 1000, 2000 GB`,
	Example: `  cloudamqp instance resize-disk --id 1234 --disk-size=100
  cloudamqp instance resize-disk --id 1234 --disk-size=250 --allow-downtime`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		if resizeInstanceID == "" {
			return fmt.Errorf("--id is required")
		}

		instanceID, err := strconv.Atoi(resizeInstanceID)
		if err != nil {
			return fmt.Errorf("invalid instance ID: %v", err)
		}

		// Validate disk size
		validSizes := []int{0, 25, 50, 100, 250, 500, 1000, 2000}
		isValid := false
		for _, size := range validSizes {
			if diskSize == size {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("invalid disk size. Valid sizes are: 0, 25, 50, 100, 250, 500, 1000, 2000 GB")
		}

		c := client.New(apiKey)

		req := &client.DiskResizeRequest{
			ExtraDiskSize: diskSize,
			AllowDowntime: allowDowntime,
		}

		err = c.ResizeInstanceDisk(instanceID, req)
		if err != nil {
			fmt.Printf("Error resizing instance disk: %v\n", err)
			return err
		}

		fmt.Printf("Disk resize initiated for instance %d. Additional disk size: %d GB\n", instanceID, diskSize)
		if allowDowntime {
			fmt.Println("Note: Downtime is allowed for this resize operation.")
		}
		return nil
	},
}

func init() {
	instanceResizeCmd.Flags().StringVar(&resizeInstanceID, "id", "", "Instance ID (required)")
	instanceResizeCmd.Flags().IntVar(&diskSize, "disk-size", 0, "Disk size to add in gigabytes (0, 25, 50, 100, 250, 500, 1000, 2000)")
	instanceResizeCmd.Flags().BoolVar(&allowDowntime, "allow-downtime", false, "Allow cluster downtime if needed when resizing disk")
	instanceResizeCmd.MarkFlagRequired("id")
	instanceResizeCmd.MarkFlagRequired("disk-size")
	instanceResizeCmd.RegisterFlagCompletionFunc("id", completeInstances)
}
