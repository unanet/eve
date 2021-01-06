package queue

// API Queue Commands
const (
	CommandScheduleDeployment string = "api-schedule-deployment"
	CommandUpdateDeployment   string = "api-update-deployment"
	CommandCallbackMessage    string = "api-callback-message"
)

// Scheduler Queue Commands
const (
	CommandDeployNamespace string = "sch-deploy-namespace"
	// CommandRestartNamespace is the command used for scheduling a service restart in a namespace
	CommandRestartNamespace string = "sch-restart-namespace"
)
