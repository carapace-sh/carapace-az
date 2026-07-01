package main

import "github.com/carapace-sh/carapace-az/cmd/carapace-az/cmd"

//go:generate sh -c "go run -C ./generate ."
func main() {
	cmd.Execute()
}
