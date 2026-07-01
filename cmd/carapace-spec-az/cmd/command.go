package cmd

type CliData struct {
	Cli      CliMeta                `json:"cli"`
	Commands map[string]CommandData `json:"commands"`
	Groups   map[string]GroupData   `json:"groups"`
}

type CliMeta struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type CommandData struct {
	Description string          `json:"description"`
	Arguments   []*ArgumentData `json:"arguments"`
	Group       string          `json:"group"`
}

type ArgumentData struct {
	Name     string   `json:"name"`
	Options  []string `json:"options"`
	Help     string   `json:"help"`
	Required bool     `json:"required"`
	Choices  []any    `json:"choices"`
	Type     string   `json:"type"`
	Nargs    string   `json:"nargs"`
	Default  any      `json:"default"`
	Metavar  any      `json:"metavar"`
}

type GroupData struct {
	Help   string         `json:"help"`
	Groups map[string]any `json:"groups"`
}
