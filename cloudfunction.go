package cloudfunction

import (
	"log"
	"net/http"
	"os"

	"github.com/schnoddelbotz/schagopubnews/handlers"
	"github.com/schnoddelbotz/schagopubnews/settings"
	// to load viper defaults for our flags...
	_ "github.com/schnoddelbotz/schagopubnews/cmd"
)

var runtimeEnvironment *handlers.Environment

func init() {
	// dunno how to ldflag on `gcloud functions deploy` ... so we pass version at deploy time. m(
	version := os.Getenv("SPN_VERSION")

	// gcloud / stackdriver logs have own timestamps, so drop Go's
	log.SetFlags(0)

	// import environment vars, using same defaults as CLI
	googleSettings := settings.ViperToRuntimeSettings(true)
	log.Printf(`schagopubnews version %s starting in "cloudfunction" mode with env proj=%s/cfn-region=%s`,
		version, googleSettings.ProjectID, googleSettings.Region)

	// we initialize all clients here, albeit different needs of CFNs. Solve.
	runtimeEnvironment = handlers.NewEnvironment(googleSettings, true)
}

// CloudVMDocker handles VMCreate, TaskStatus and TaskProgress requests
func Schagopubnews(w http.ResponseWriter, r *http.Request) {
	handlers.Schagopubnews(w, r, runtimeEnvironment)
}
