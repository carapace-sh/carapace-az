package cmd

import (
	"fmt"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-az/cmd/carapace-az/cmd/azcli"
	"github.com/carapace-sh/carapace-az/cmd/carapace-az/cmd/common"
	_ "github.com/carapace-sh/carapace-az/pkg/actions"
	"github.com/carapace-sh/carapace-az/pkg/actions/az"
	spec "github.com/carapace-sh/carapace-spec"
	"github.com/carapace-sh/carapace/pkg/style"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "az",
	Short: "An enriched az completer",
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetUsageFunc(func(c *cobra.Command) error { return nil })
	carapace.Gen(rootCmd).Standalone()

	rootCmd.PersistentFlags().String("output", "", "Output format.")
	rootCmd.PersistentFlags().StringP("query", "q", "", "JMESPath query string.")
	rootCmd.PersistentFlags().String("subscription", "", "Name or ID of subscription.")
	rootCmd.PersistentFlags().Bool("verbose", false, "Increase logging verbosity.")
	rootCmd.PersistentFlags().Bool("debug", false, "Increase logging verbosity to show all debug logs.")
	rootCmd.PersistentFlags().Bool("only-show-errors", false, "Only show errors, suppressing all warnings.")
	rootCmd.PersistentFlags().Bool("help", false, "Show help message.")
	rootCmd.PersistentFlags().Bool("version", false, "Show version information.")

	carapace.Gen(rootCmd).FlagCompletion(carapace.ActionMap{
		"output":       carapace.ActionValues("json", "jsonc", "none", "table", "tsv", "yaml", "yamlc").StyleF(style.ForExtension),
		"subscription": az.ActionSubscriptions(),
	})

	for name, description := range azcli.Services() {
		serviceCmd := &cobra.Command{
			Use:   name,
			Short: description,
			Run:   func(cmd *cobra.Command, args []string) {},
		}
		carapace.Gen(serviceCmd).Standalone()
		rootCmd.AddCommand(serviceCmd)
		carapace.Gen(serviceCmd).PreRun(func(cmd *cobra.Command, args []string) {
			azCommand, err := azcli.Get(fmt.Sprintf("az.%s.yaml", serviceCmd.Use))
			if err != nil {
				carapace.LOG.Println(err.Error())
				return
			}

			for _, subCmd := range azCommand.Commands {
				operationCmd := spec.Command(subCmd).ToCobra()
				serviceCmd.AddCommand(operationCmd)

				carapace.Gen(operationCmd).PreInvoke(func(cmd *cobra.Command, flag *pflag.Flag, action carapace.Action) carapace.Action {
					if flag != nil && flag.Value.Type() != "bool" {
						if _, ok := subCmd.Completion.Flag[flag.Name]; !ok {
							return common.ActionBridgeAzCompleter()
						}
					}
					return action
				})
			}
		})
	}

	for _, extension := range []string{
		"configure",
		"feedback",
		"find",
		"init",
		"interactive",
		"login",
		"logout",
		"rest",
		"survey",
		"upgrade",
		"version",
	} {
		subCmd := &cobra.Command{
			Use:                extension,
			Run:                func(cmd *cobra.Command, args []string) {},
			DisableFlagParsing: true,
		}
		rootCmd.AddCommand(subCmd)

		carapace.Gen(subCmd).PositionalAnyCompletion(
			common.ActionBridgeAzCompleter(),
		)
	}

	spec.Register(rootCmd)
}
