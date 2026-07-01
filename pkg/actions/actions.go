package actions

import (
	"github.com/carapace-sh/carapace-az/pkg/actions/az"
	spec "github.com/carapace-sh/carapace-spec"
)

func init() {
	spec.AddMacro("az.subscriptions", spec.MacroN(az.ActionSubscriptions))
	spec.AddMacro("az.resourcegroups", spec.MacroN(az.ActionResourceGroups))
}
