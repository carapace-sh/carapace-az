package common

import (
	"os"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-bridge/pkg/actions/bridge"
)

func ActionBridgeAzCompleter() carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		c.Args = carapace.NewContext(os.Args[4:]...).Args
		return bridge.ActionArgcomplete("az").Invoke(c).ToA()
	})
}
