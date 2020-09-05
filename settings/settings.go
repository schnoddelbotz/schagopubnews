package settings

import (
	"strings"

	"github.com/spf13/viper"
)

// RuntimeSettings define anything Google related (project, service account, ...)
type RuntimeSettings struct {
	ProjectID           string
	Region              string
	Token               string // access protects the HTTP CFN
	FireStoreCollection string
	Verbose             bool
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

	// FireStoreCollection is the name of our firestore collection (static for now)
	FireStoreCollection = "schagopubnews"
)

// ViperToRuntimeSettings translates environment variables into a RuntimeSettings struct.
func ViperToRuntimeSettings(permitEmptyToken bool) RuntimeSettings {
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("SPN")
	s := RuntimeSettings{
		ProjectID:           viper.GetString(FlagProject),
		Region:              viper.GetString(FlagRegion),
		Token:               viper.GetString(FlagToken),
		Verbose:             viper.GetBool(FlagVerbose),
		FireStoreCollection: FireStoreCollection,
	}
	return s
}
