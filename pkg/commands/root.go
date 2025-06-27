package commands

import (
	"os"

	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/oam-dev/kubevela/pkg/cli/vela"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                "vela",
		DisableFlagParsing: true,
		Run: func(cmd *cobra.Command, args []string) {
			argsLen := len(args)
			if argsLen == 0 {
				cmd.Help()
				return
			}
			// Find the first argument that look like a flag
			firstArg := 0
			for ; firstArg < argsLen; firstArg++ {
				if args[firstArg] == "--" {
					firstArg = firstArg + 1
					break
				}
				if args[firstArg] == "-h" || args[firstArg] == "--help" {
					cmd.Help()
					return
				}
				if args[firstArg] == "--version" {
					// Print version here or implement as needed
					cmd.Println("KubeVela version: unknown")
					return
				}
				if args[firstArg][0] == '-' {
					cmd.Help()
					return
				}
			}
			// No-op for flags.ParseInto(cmd, args)
		},
		SilenceUsage: true,
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
	}
	// Remove references to util.IOStreams, flags, etc. if you don't have those files
	// Remove f := flags.NewDefaultCliFlags(cmd)
	// Remove f.AddFlags(cmd)
	args := genericclioptions.NewConfigFlags(true)
	args.AddFlags(cmd.PersistentFlags())

	// Only add commands you have implemented
	cmd.AddCommand(
		vela.NewWorkspaceCommand(), // Adjust signature if needed
		// ...add other commands as implemented...
	)
	return cmd
}
