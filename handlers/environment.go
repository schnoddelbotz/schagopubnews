package handlers

import (
	"context"

	"cloud.google.com/go/firestore"

	"github.com/schnoddelbotz/schagopubnews/cloud"
	"github.com/schnoddelbotz/schagopubnews/settings"
)

// Environment enables resource sharing between CFN requests, holds env settings + svc conns
type Environment struct {
	RuntimeSettings settings.RuntimeSettings
	Context         context.Context
	FireStoreClient *firestore.Client
}

// NewEnvironment creates google service clients as requested
func NewEnvironment(runtimeSettings settings.RuntimeSettings, withFireStoreClient bool) *Environment {
	var firestoreClient *firestore.Client

	if withFireStoreClient {
		firestoreClient = cloud.NewFireStoreClient(context.Background(), runtimeSettings.ProjectID)
	}

	return &Environment{
		RuntimeSettings: runtimeSettings,
		Context:         context.Background(),
		FireStoreClient: firestoreClient,
	}
}
