package get

import (
	"fmt"
	"os"
	"strings"

	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/supergloo/cli/pkg/clients"
	"github.com/solo-io/supergloo/cli/pkg/cmd/options"
	"github.com/solo-io/supergloo/cli/pkg/cmd/printers"
	"github.com/solo-io/supergloo/cli/pkg/constants"
	"github.com/solo-io/supergloo/cli/pkg/util"
	"github.com/spf13/cobra"
)

var supportedOutputFormats = []string{"wide"}

func Cmd(opts *options.Options) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: `Display one or many supergloo resources`,
		Long:  `Display one or many supergloo resources`,
		Args:  cobra.RangeArgs(1, 2),
		Run: func(c *cobra.Command, args []string) {
			err := get(args, opts)
			if err != nil {
				fmt.Println(err)
			}
		},
	}
	getOpts := &opts.Get
	pFlags := cmd.Flags()
	pFlags.StringVarP(&getOpts.Output, "output", "o", "",
		"Output format. Must be one of: \n"+strings.Join(supportedOutputFormats, "|"))
	return cmd
}

func get(args []string, opts *options.Options) error {

	output := opts.Get.Output
	if output != "" && !util.Contains(supportedOutputFormats, output) {
		return errors.Errorf(constants.UnknownOutputFormat, output, strings.Join(supportedOutputFormats, "|"))
	}

	if argNumber := len(args); argNumber == 1 {
		return getResource(args[0], "", opts.Get)
	} else {
		// Show the resource of the given type with the given name
		return getResource(args[0], args[1], opts.Get)
	}
}

func getResource(resourceType, resourceName string, opts options.Get) error {
	sgClient, err := clients.NewClient()
	if err != nil {
		return err
	}

	// Get available resource types
	resourceTypes, err := sgClient.ListResourceTypes()
	if err != nil {
		return err
	}

	// Validate input resource type
	if !util.Contains(resourceTypes, resourceType) {
		return errors.Errorf(constants.UnknownResourceTypeMsg, resourceType)
	}

	// Fetch the resource information
	info, err := sgClient.ListResources(resourceType, resourceName)
	if err != nil {
		return err
	}

	// Write the resource information to stdout
	writer := printers.NewTableWriter(os.Stdout)
	if err = writer.WriteLine(info.Headers(opts)); err != nil {
		return err
	}

	for _, line := range info.Resources(opts) {
		if err = writer.WriteLine(line); err != nil {
			return err
		}
	}

	return writer.Flush()
}
