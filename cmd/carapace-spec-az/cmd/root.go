package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/carapace-sh/carapace"
	command "github.com/carapace-sh/carapace-spec/pkg/command"
	"github.com/neurosnap/sentences"
	"github.com/neurosnap/sentences/english"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:   "carapace-spec-az",
	Short: "",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}

		var data CliData
		if err := json.Unmarshal(content, &data); err != nil {
			return err
		}

		specCommand := buildSpec(&data)

		if cmd.Flag("no-doc").Changed {
			stripDoc(&specCommand)
		}

		if cmd.Flag("stdout").Changed {
			m, err := yaml.Marshal(specCommand)
			if err != nil {
				return err
			}
			fmt.Println("# yaml-language-server: $schema=https://carapace.sh/schemas/command.json")
			fmt.Println(string(m))
			return nil
		}

		dir := cmd.Flag("target").Value.String()
		if dir == "" {
			dir, err = os.MkdirTemp("", "carapace-spec-az-*")
			if err != nil {
				return err
			}
		}
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}

		for _, subCommand := range specCommand.Commands {
			m, err := yaml.Marshal(subCommand)
			if err != nil {
				return err
			}
			m = append([]byte("# yaml-language-server: $schema=https://carapace.sh/schemas/command.json\n"), m...)
			p := path.Join(dir, fmt.Sprintf("az.%s.yaml", subCommand.Name))
			println(p)
			if err := os.WriteFile(p, m, os.ModePerm); err != nil {
				return err
			}
		}

		specCommand.Commands = nil
		m, err := yaml.Marshal(specCommand)
		if err != nil {
			return err
		}
		p := path.Join(dir, "az.yaml")
		println(p)
		if err := os.WriteFile(p, m, os.ModePerm); err != nil {
			return err
		}
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	carapace.Gen(rootCmd).Standalone()
	rootCmd.Flags().Bool("no-doc", false, "strip documentation")
	rootCmd.Flags().Bool("stdout", false, "print to stdout")
	rootCmd.Flags().String("target", "", "target directory")
	rootCmd.MarkFlagsMutuallyExclusive("stdout", "target")

	carapace.Gen(rootCmd).PositionalCompletion(
		carapace.ActionFiles(),
	)
}

func stripDoc(command *command.Command) {
	command.Documentation.Command = ""
	command.Documentation.Flag = nil
	for index := range command.Commands {
		stripDoc(&command.Commands[index])
	}
}

func buildSpec(data *CliData) command.Command {
	root := command.Command{
		Name:        data.Cli.Name,
		Description: "Manage Azure resources and services.",
	}
	root.Completion.Flag = make(map[string][]string)
	root.Documentation.Flag = make(map[string]string)

	groupCommands := make(map[string][]string)
	for cmdName := range data.Commands {
		parts := strings.SplitN(cmdName, " ", 2)
		topGroup := parts[0]
		groupCommands[topGroup] = append(groupCommands[topGroup], cmdName)
	}

	for topGroup, cmdNames := range groupCommands {
		if len(cmdNames) == 1 && cmdNames[0] == topGroup {
			continue
		}
		groupSpec := buildGroupSpec(topGroup, cmdNames, data)
		root.Commands = append(root.Commands, groupSpec)
	}

	slices.SortFunc(root.Commands, func(a, b command.Command) int { return strings.Compare(a.Name, b.Name) })
	return root
}

func buildGroupSpec(topGroup string, cmdNames []string, data *CliData) command.Command {
	groupSpec := command.Command{
		Name: topGroup,
	}

	if group, ok := data.Groups[topGroup]; ok && group.Help != "" {
		groupSpec.Description = group.Help
	}
	if groupSpec.Description == "" {
		groupSpec.Description = topGroup
	}

	groupSpec.Completion.Flag = make(map[string][]string)
	groupSpec.Documentation.Flag = make(map[string]string)

	root := &command.Command{}
	for _, cmdName := range cmdNames {
		cmdData := data.Commands[cmdName]
		specCmd := convertCommand(cmdName, &cmdData)
		insertCommand(root, cmdName, specCmd, data)
	}

	groupSpec.Commands = root.Commands
	slices.SortFunc(groupSpec.Commands, func(a, b command.Command) int { return strings.Compare(a.Name, b.Name) })

	return groupSpec
}

func convertCommand(fullName string, cmdData *CommandData) command.Command {
	parts := strings.Split(fullName, " ")
	leafName := parts[len(parts)-1]

	specCmd := command.Command{
		Name:        leafName,
		Description: firstSentence(cmdData.Description),
	}
	specCmd.Completion.Flag = make(map[string][]string)
	specCmd.Documentation.Flag = make(map[string]string)

	for _, arg := range cmdData.Arguments {
		if len(arg.Options) == 0 {
			continue
		}
		f := convertArgument(arg)
		specCmd.AddFlag(f)

		if len(arg.Choices) > 0 {
			choices := make([]string, 0, len(arg.Choices))
			for _, c := range arg.Choices {
				choices = append(choices, fmt.Sprintf("%v", c))
			}
			specCmd.Completion.Flag[f.Name()] = choices
		}
		specCmd.Documentation.Flag[f.Name()] = arg.Help
	}

	return specCmd
}

func convertArgument(arg *ArgumentData) command.Flag {
	f := command.Flag{
		Description: firstSentence(arg.Help),
		Required:    arg.Required,
	}

	for _, rawOpt := range arg.Options {
		for _, opt := range strings.Split(rawOpt, ", ") {
			opt = strings.TrimSpace(opt)
			if strings.HasPrefix(opt, "--") {
				if f.Longhand == "" {
					f.Longhand = strings.TrimPrefix(opt, "--")
				}
			} else if strings.HasPrefix(opt, "-") && len(opt) > 1 {
				if f.Shorthand == "" {
					f.Shorthand = strings.TrimPrefix(opt, "-")
				}
			}
		}
	}

	if arg.Type == "bool" || arg.Type == "<class 'bool'>" {
		f.Value = false
	} else {
		f.Value = true
	}

	if arg.Nargs != "" && arg.Nargs != "1" && arg.Nargs != "0" && arg.Nargs != "None" {
		f.Nargs = -1
	}

	return f
}

func insertCommand(root *command.Command, fullName string, specCmd command.Command, data *CliData) {
	parts := strings.Split(fullName, " ")
	if len(parts) <= 1 {
		root.Commands = append(root.Commands, specCmd)
		return
	}

	current := root
	for i := 1; i < len(parts)-1; i++ {
		groupName := parts[i]
		found := false
		for idx := range current.Commands {
			if current.Commands[idx].Name == groupName {
				current = &current.Commands[idx]
				found = true
				break
			}
		}
		if !found {
			newGroup := command.Command{Name: groupName}
			groupPath := strings.Join(parts[:i+1], " ")
			if group, ok := data.Groups[groupPath]; ok && group.Help != "" {
				newGroup.Description = group.Help
			}
			current.Commands = append(current.Commands, newGroup)
			current = &current.Commands[len(current.Commands)-1]
		}
	}
	current.Commands = append(current.Commands, specCmd)
}

var tokenizer *sentences.DefaultSentenceTokenizer

func init() {
	var err error
	tokenizer, err = english.NewSentenceTokenizer(nil)
	if err != nil {
		panic(err.Error())
	}
}

func firstSentence(s string) string {
	if s == "" {
		return ""
	}
	tokens := tokenizer.Tokenize(s)
	if len(tokens) > 0 {
		return strings.TrimSpace(tokens[0].Text)
	}
	return strings.TrimSpace(s)
}
