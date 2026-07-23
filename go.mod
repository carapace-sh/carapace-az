module github.com/carapace-sh/carapace-az

go 1.25.0

require (
	github.com/carapace-sh/carapace v1.13.3
	github.com/carapace-sh/carapace-bridge v1.6.2
	github.com/carapace-sh/carapace-spec v1.8.0
	github.com/neurosnap/sentences v1.1.2
	github.com/spf13/cobra v1.10.2
	github.com/spf13/pflag v1.0.10
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/carapace-sh/carapace-shlex v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
)

replace github.com/spf13/pflag => github.com/carapace-sh/carapace-pflag v1.1.0
