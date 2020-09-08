package settings

import (
	"strings"

	"github.com/spf13/viper"
)

// RuntimeSettings define anything Google related (project, service account, ...)
// Also passed to go index.html tpl. May cause leaks...? clean up.
type RuntimeSettings struct {
	ProjectID           string
	Region              string
	Token               string // access protects the HTTP CFN
	FireStoreCollection string
	Port                string
	Verbose             bool
	RootURL             string
	ApiURL              string
	SPNVersion          string
}

const (
	// FlagProject ...
	FlagProject = "project"
	// FlagRegion defines region of CFNs
	FlagRegion = "region"
	// FlagToken is the access token for CloudVMDocker HTTP CloudFunction for run/ps/...
	FlagToken = "token"
	// FlagVerbose is just a bool here
	FlagVerbose = "verbose"
	// FlagPort is
	FlagPort = "port"

	FlagRootURL = "rootURL"
	FlagAPIURL = "apiURL"

	// FireStoreCollection is the name of our firestore collection (static for now)
	FireStoreCollection = "schagopubnews"
)

// ViperToRuntimeSettings translates environment variables into a RuntimeSettings struct.
func ViperToRuntimeSettings() RuntimeSettings {
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("SPN")
	s := RuntimeSettings{
		ProjectID:           viper.GetString(FlagProject),
		Region:              viper.GetString(FlagRegion),
		Token:               viper.GetString(FlagToken),
		Port:                viper.GetString(FlagPort),
		Verbose:             viper.GetBool(FlagVerbose),
		RootURL: viper.GetString(FlagRootURL),
		ApiURL: viper.GetString(FlagAPIURL),
		FireStoreCollection: FireStoreCollection,
	}
	return s
}
