package compute

import (
	"context"

	"fmt"

	"time"

	"github.com/joyent/triton-go/cmd/internal/api"
	"github.com/joyent/triton-go/cmd/internal/config"
	"github.com/joyent/triton-go/cmd/internal/console_writer"
	"github.com/joyent/triton-go/compute"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var InstancesCommand = &cobra.Command{
	Use:        "instance",
	Aliases:    []string{"instances"},
	SuggestFor: []string{"machines"},
	Short:      "instance interaction with triton",
}

var ListInstancesCommand = &cobra.Command{
	Args:         cobra.NoArgs,
	Use:          "list instances",
	Short:        "lists instances associated with triton account",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cons := console_writer.GetTerminal()

		tritonClientConfig, err := api.InitConfig()
		if err != nil {
			return err
		}

		client, err := tritonClientConfig.GetComputeClient()
		if err != nil {
			return err
		}

		instances, err := client.Instances().List(context.Background(), &compute.ListInstancesInput{})
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(cons)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetHeaderLine(false)
		table.SetAutoFormatHeaders(true)

		table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")

		table.SetHeader([]string{"SHORTID", "NAME", "IMAGE", "STATE"})

		var numInstances uint
		for _, instance := range instances {
			table.Append([]string{string(instance.ID[:8]), instance.Name, instance.Image, instance.State})
			numInstances++
		}

		table.Render()

		return nil
	},
}

var GetInstanceCommand = &cobra.Command{
	Args:         cobra.NoArgs,
	Use:          "get",
	Short:        "get a triton instance",
	SilenceUsage: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !viper.IsSet(config.KeyImageId) && !viper.IsSet(config.KeyImageName) {
			return errors.New("Either `id` or `name` must be specified")
		}

		if viper.IsSet(config.KeyImageId) && viper.IsSet(config.KeyImageName) {
			return errors.New("Only 1 of `id` or `name` must be specified")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cons := console_writer.GetTerminal()

		var machine *compute.Instance

		tritonClientConfig, err := api.InitConfig()
		if err != nil {
			return err
		}

		client, err := tritonClientConfig.GetComputeClient()
		if err != nil {
			return err
		}

		machineName := getMachineName()
		if machineName != "" {
			machines, err := client.Instances().List(context.Background(), &compute.ListInstancesInput{
				Name: machineName,
			})
			if err != nil {
				return err
			}

			if machines != nil {
				//Name has to be unique with an account and a region
				machine = machines[0]
			}
		}

		machineID := getMachineID()
		if machineID != "" {
			machine, err = client.Instances().Get(context.Background(), &compute.GetInstanceInput{
				ID: machineID,
			})
			if err != nil {
				return err
			}
		}

		table := tablewriter.NewWriter(cons)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetHeaderLine(false)
		table.SetAutoFormatHeaders(true)

		table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_RIGHT})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("")
		table.SetColumnSeparator("")
		table.SetRowSeparator("")

		table.SetHeader([]string{"------", "------"})

		table.Append([]string{"id", machine.ID})
		table.Append([]string{"name", machine.Name})
		table.Append([]string{"package", machine.Package})
		table.Append([]string{"image", machine.Image})
		table.Append([]string{"brand", machine.Brand})
		table.Append([]string{"firewall enabled", fmt.Sprintf("%t", machine.FirewallEnabled)})

		table.Render()

		return nil
	},
}

var CreateInstanceCommand = &cobra.Command{
	Use:          "create",
	Short:        "create a triton instance",
	SilenceUsage: true,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if !viper.IsSet(config.KeyPackageId) && !viper.IsSet(config.KeyPackageName) {
			return errors.New("Either `pkg-name` or `pkg-id` must be specified")
		}

		if !viper.IsSet(config.KeyImageId) && !viper.IsSet(config.KeyImageName) {
			return errors.New("Either `img-name` or `img-id` must be specified")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cons := console_writer.GetTerminal()

		params := &compute.CreateInstanceInput{
			Name:            viper.GetString(config.KeyInstanceName),
			FirewallEnabled: viper.GetBool(config.KeyInstanceFirewall),
		}

		tritonClientConfig, err := api.InitConfig()
		if err != nil {
			return err
		}

		client, err := tritonClientConfig.GetComputeClient()
		if err != nil {
			return err
		}

		pkgID := getPkgID()
		if pkgID != "" {
			params.Package = pkgID
		} else {
			pkgName := getPkgName()
			packages, err := client.Packages().List(context.Background(), &compute.ListPackagesInput{})
			if err != nil {
				return err
			}

			for _, pkg := range packages {
				if pkg.Name == pkgName {
					params.Package = pkg.ID
					cons.Write([]byte(fmt.Sprintf("Found Package: %q\n", pkg.Name)))
					break
				}
			}
		}

		imgID := getImgID()
		if imgID != "" {
			params.Image = imgID
		} else {
			imgName := getImgName()
			images, err := client.Images().List(context.Background(), &compute.ListImagesInput{})
			if err != nil {
				return err
			}

			for _, img := range images {
				if img.Name == imgName {
					params.Image = img.ID
					cons.Write([]byte(fmt.Sprintf("Found Image: %q\n", img.Name)))
					break
				}
			}
		}

		startTime := time.Now()
		machine, err := client.Instances().Create(context.Background(), params)
		if err != nil {
			return err
		}

		cons.Write([]byte(fmt.Sprintf("Creating instance %q (%s)", machine.Name, machine.ID)))

		if viper.GetBool(config.KeyInstanceWait) {
			state := make(chan *compute.Instance, 1)
			go func(machineID string, c *compute.ComputeClient) {
				for {
					time.Sleep(1 * time.Second)
					instance, err := c.Instances().Get(context.Background(), &compute.GetInstanceInput{
						ID: machineID,
					})
					if err != nil {
						panic(err)
					}
					if instance.State == "running" {
						state <- instance
					} else {
						fmt.Print(".")
					}
				}
			}(machine.ID, client)

			select {
			case instance := <-state:
				cons.Write([]byte("\n"))
				cons.Write([]byte(fmt.Sprintf("Created instance %q (%s) in %d", instance.Name, instance.ID, time.Since(startTime))))
			case <-time.After(5 * time.Minute):
				cons.Write([]byte("Create instance operation timed out"))
			}
		}

		return nil
	},
}

func getPkgID() string {
	if viper.IsSet(config.KeyPackageId) {
		return viper.GetString(config.KeyPackageId)
	}
	return ""
}

func getPkgName() string {
	if viper.IsSet(config.KeyPackageName) {
		return viper.GetString(config.KeyPackageName)
	}
	return ""
}

func getImgID() string {
	if viper.IsSet(config.KeyImageId) {
		return viper.GetString(config.KeyImageId)
	}
	return ""
}

func getImgName() string {
	if viper.IsSet(config.KeyImageName) {
		return viper.GetString(config.KeyImageName)
	}
	return ""
}

func getMachineID() string {
	if viper.IsSet(config.KeyInstanceId) {
		return viper.GetString(config.KeyInstanceId)
	}
	return ""
}

func getMachineName() string {
	if viper.IsSet(config.KeyInstanceName) {
		return viper.GetString(config.KeyInstanceName)
	}
	return ""
}
