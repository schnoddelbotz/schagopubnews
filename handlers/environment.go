package handlers

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/schnoddelbotz/schagopubnews/cloud"
	"github.com/schnoddelbotz/schagopubnews/settings"
)

// Environment enables resource sharing between CFN requests, holds env settings + svc conns
type Environment struct {
	GoogleSettings  settings.RuntimeSettings
	Context         context.Context
	FireStoreClient *firestore.Client
}

// NewEnvironment creates google service clients as requested
func NewEnvironment(googleSettings settings.RuntimeSettings, withFireStoreClient bool) *Environment {
	var dataStoreClient *firestore.Client

	if withFireStoreClient {
		dataStoreClient = cloud.NewFireStoreClient(context.Background(), googleSettings.ProjectID)
	}

	return &Environment{
		GoogleSettings:  googleSettings,
		Context:         context.Background(),
		FireStoreClient: dataStoreClient,
	}
}
