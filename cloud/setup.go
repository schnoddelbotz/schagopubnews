package cloud

import (
	"log"
)

// Setup creates necessary GoogleCloud infra for schagopubnews operations --- NOT USABLE YET
func Setup(projectID string) {
	log.Print("SETTING UP infrastructure for schagopubnews ...")
	// todo: deploy cfn
	log.Print("SUCCESS setting up schagopubnews cloud infra. Have fun!")
}

// Destroy removes infra created by setup routine
func Destroy(projectID string) {
	log.Print("DESTROYING infrastructure for schagopubnews ...")
	// todo: delete cfn
	//       clear datastore
	//       clear VMs?
	//       clear SD logs?
	log.Print("SUCCESS destroying schagopubnews cloud infra. Have fun!")
}

func bailOnError(err error, message string) {
	if err == nil {
		log.Printf("Success: %s", message)
	} else {
		log.Fatalf("ERROR: %s: %s", message, err.Error())
	}
}
