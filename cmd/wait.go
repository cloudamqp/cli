package cmd

import (
	"context"
	"fmt"
	"time"

	"cloudamqp-cli/client"
)

func waitForInstanceReady(c *client.Client, instanceID int, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()

	// Check immediately first
	instance, err := c.GetInstance(instanceID)
	if err != nil {
		return fmt.Errorf("failed to check instance status: %w", err)
	}
	if instance.Ready {
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			elapsed := time.Since(startTime)
			return fmt.Errorf("timeout after %s waiting for instance to be ready", elapsed.Round(time.Second))
		case <-ticker.C:
			instance, err := c.GetInstance(instanceID)
			if err != nil {
				return fmt.Errorf("failed to check instance status: %w", err)
			}
			if instance.Ready {
				return nil
			}
		}
	}
}
