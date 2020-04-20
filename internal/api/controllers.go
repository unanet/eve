package api

import (
	"gitlab.unanet.io/devops/eve/internal/api/ping"
	"gitlab.unanet.io/devops/eve/pkg/mux"
)

var Controllers = []mux.EveController{
	ping.New(),
}
