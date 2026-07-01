package az

import (
	"encoding/json"
	"os/exec"

	"github.com/carapace-sh/carapace"
)

func ActionResourceGroups() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		cmd := exec.Command("az", "group", "list", "--output", "json")
		output, err := cmd.Output()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		var groups []struct {
			Name     string `json:"name"`
			Location string `json:"location"`
		}
		if err := json.Unmarshal(output, &groups); err != nil {
			return carapace.ActionMessage(err.Error())
		}

		vals := make([]string, 0)
		for _, rg := range groups {
			vals = append(vals, rg.Name, rg.Location)
		}
		return carapace.ActionValuesDescribed(vals...)
	})
}
