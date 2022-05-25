package monitor

import (
	"fmt"
	"github.com/containrrr/shoutrrr/pkg/router"
	"github.com/containrrr/shoutrrr/pkg/types"
)

type Notifier struct {
	router *router.ServiceRouter
}

func (n *Notifier) Send(title, message string) (errs []error) {
	if n.router == nil {
		return []error{fmt.Errorf("notifier not initialized")}
	}
	params := types.Params{}
	params.SetTitle(title)
	return n.router.Send(message, &params)
}
