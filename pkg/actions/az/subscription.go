package az

import (
	"encoding/json"
	"os/exec"

	"github.com/carapace-sh/carapace"
)

func ActionSubscriptions() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		cmd := exec.Command("az", "account", "list", "--output", "json")
		output, err := cmd.Output()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		var subscriptions []struct {
			Name      string `json:"name"`
			ID        string `json:"id"`
			IsDefault bool   `json:"isDefault"`
			State     string `json:"state"`
		}
		if err := json.Unmarshal(output, &subscriptions); err != nil {
			return carapace.ActionMessage(err.Error())
		}

		vals := make([]string, 0)
		for _, sub := range subscriptions {
			vals = append(vals, sub.Name, sub.ID)
		}
		return carapace.ActionValuesDescribed(vals...)
	})
}
