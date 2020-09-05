package cmd

import (
	"fmt"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"

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
		fmt.Printf("serve called, using port : %s", runtimeSettings.Port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP(settings.FlagPort, "p", "8080", "HTTP port to serve on")
	_ = viper.BindPFlag(settings.FlagPort, serveCmd.Flags().Lookup("port"))
}
