package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"cloudamqp-cli/client"
	"github.com/spf13/cobra"
)

var instanceListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all CloudAMQP instances",
	Long:    `Retrieves and displays all CloudAMQP instances in your account.`,
	Example: `  cloudamqp instance list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		apiKey, err = getAPIKey()
		if err != nil {
			return fmt.Errorf("failed to get API key: %w", err)
		}

		c := client.New(apiKey, Version)

		instances, err := c.ListInstances()
		if err != nil {
			fmt.Printf("Error listing instances: %v\n", err)
			return err
		}

		if len(instances) == 0 {
			fmt.Println("No instances found.")
			return nil
		}

		p, err := getPrinter(cmd)
		if err != nil {
			return err
		}

		details, _ := cmd.Flags().GetBool("details")

		if details {
			showURL, _ := cmd.Flags().GetBool("show-url")
			detailed := make([]*client.Instance, len(instances))
			headers := []string{"ID", "NAME", "PLAN", "REGION", "TAGS", "URL", "HOSTNAME", "READY"}
			rows := make([][]string, len(instances))
			var (
				mu       sync.Mutex
				firstErr error
				wg       sync.WaitGroup
			)
			for i, instance := range instances {
				wg.Add(1)
				go func(idx, id int) {
					defer wg.Done()
					det, err := c.GetInstance(id)
					mu.Lock()
					defer mu.Unlock()
					if err != nil {
						if firstErr == nil {
							firstErr = fmt.Errorf("error fetching instance %d: %w", id, err)
						}
						return
					}
					detailed[idx] = det
				}(i, instance.ID)
			}
			wg.Wait()
			if firstErr != nil {
				return firstErr
			}

			for i, inst := range detailed {
				ready := "No"
				if inst.Ready {
					ready = "Yes"
				}
				urlVal := maskPassword(inst.URL)
				if showURL {
					urlVal = inst.URL
				}
				rows[i] = []string{
					strconv.Itoa(inst.ID),
					inst.Name,
					inst.Plan,
					inst.Region,
					strings.Join(inst.Tags, ","),
					urlVal,
					inst.HostnameExternal,
					ready,
				}
			}
			p.PrintRecords(headers, rows)
			return nil
		}

		headers := []string{"ID", "NAME", "PLAN", "REGION"}
		rows := make([][]string, len(instances))
		for i, instance := range instances {
			rows[i] = []string{
				strconv.Itoa(instance.ID),
				instance.Name,
				instance.Plan,
				instance.Region,
			}
		}
		p.PrintRecords(headers, rows)

		return nil
	},
}

func init() {
	instanceListCmd.Flags().BoolP("details", "", false, "Fetch full details for each instance (one GET request per instance)")
	instanceListCmd.Flags().BoolP("show-url", "", false, "Show full connection URL with credentials (requires --details)")
}
