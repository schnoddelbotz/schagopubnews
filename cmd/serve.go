package cmd

import (
	"fmt"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"

	"github.com/schnoddelbotz/schagopubnews/handlers"
	"github.com/schnoddelbotz/schagopubnews/settings"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves the SPN WebUI like the CloudFunction, but locally",
	Long: `Serving locally, in contrast to CloudFunction runtime environment, 
expects GOOGLE_APPLICATION_CREDENTIALS for FireStore access.`,
	Run: func(cmd *cobra.Command, args []string) {
		runtimeSettings := settings.ViperToRuntimeSettings()
		runtimeSettings.SPNVersion = AppVersion
		fmt.Printf("SPN starting service on port : %s\n", runtimeSettings.Port)
		handlers.Serve(runtimeSettings)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP(settings.FlagPort, "p", "8080", "HTTP port to serve on")
	serveCmd.Flags().StringP(settings.FlagRootURL, "r", "/SPN/", "Makes little sense to tweak in serve?")
	serveCmd.Flags().StringP(settings.FlagAPIURL, "a", "/SPN/", "(Opt) Any remote API that allows us via CORS")
	// FIXME TODO
	// serveCmd.Flags().StringP(settings.FlagUIBaseURI, "p", "8080", "Ember UI base URL (/SPN)")
	// serveCmd.Flags().StringP(settings.FlagAPIBaseURI, "p", "8080", "CFn URL")
	// use go template engine (?) to replace ember's index.html config stuff -- ember build per customer would suck big

	_ = viper.BindPFlag(settings.FlagPort, serveCmd.Flags().Lookup(settings.FlagPort))
	_ = viper.BindPFlag(settings.FlagRootURL, serveCmd.Flags().Lookup(settings.FlagRootURL))
	_ = viper.BindPFlag(settings.FlagAPIURL, serveCmd.Flags().Lookup(settings.FlagAPIURL))
}

// offer base image including gcloud sdk --> ease gsutil ... stuff. how to give console access to bucket?
