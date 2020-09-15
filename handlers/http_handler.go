package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	firebase "firebase.google.com/go/v4"

	"github.com/schnoddelbotz/schagopubnews/settings"
)

type server struct {
	RuntimeSettings settings.RuntimeSettings
}

// Schagopubnews is the CloudFunction entry point for SPN
func Schagopubnews(w http.ResponseWriter, r *http.Request, env *Environment) {
	log.Printf("Request for... %s", r.URL.Path)
	CFNMux := contentNegotiationHandler(accessLogHandler(serveMux("", env.RuntimeSettings))) // todo: move out of here, reuse.
	CFNMux.ServeHTTP(w, r)
}

// Serve runs the SPN app on a local TCP port
func Serve(runtimeSettings settings.RuntimeSettings) {
	log.Printf("Starting SPN server on port %s\n", runtimeSettings.Port)
	LocalServerMux := contentNegotiationHandler(accessLogHandler(serveMux("/SPN", runtimeSettings))) // FIXME viper setting -- no trail?!!
	err := http.ListenAndServe(":"+runtimeSettings.Port, LocalServerMux)
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}

func serveMux(muxPrefix string, runtimeSettings settings.RuntimeSettings) *http.ServeMux {
	srv := &server{}
	srv.RuntimeSettings = runtimeSettings
	mux := http.NewServeMux()
	if muxPrefix == "" {
		mux.Handle(muxPrefix+"/assets/", http.FileServer(_escFS(false)))
	} else {
		mux.Handle(muxPrefix+"/assets/", http.StripPrefix(muxPrefix, http.FileServer(_escFS(false))))
	}
	mux.HandleFunc(muxPrefix+"/token", srv.tokenHandler)
	mux.HandleFunc(muxPrefix+"/", srv.indexHandler)
	return mux
}

func contentNegotiationHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// log.Printf("cNeg NO ZIP support on client FOR: %s", r.URL.Path)
			h.ServeHTTP(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, ".js") {
			r.URL.Path = r.URL.Path + ".gz"
			w.Header().Set("Content-Type", "text/javascript;charset=UTF-8")
			w.Header().Set("Content-Encoding", "gzip")
			h.ServeHTTP(w, r)
			return
		}
		if strings.HasSuffix(r.URL.Path, ".css") {
			r.URL.Path = r.URL.Path + ".gz"
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Content-Type", "text/css;charset=UTF-8")
			h.ServeHTTP(w, r)
			return
		}
		// log.Printf("NO GZIP FOR: %s", r.URL.Path)
		// DEFAULT:
		h.ServeHTTP(w, r)
	})
}

func (s *server) indexHandler(w http.ResponseWriter, r *http.Request) {
	// use default index.html template, which has /SPN/ in config for both URLs
	if s.RuntimeSettings.RootURL == "/SPN/" && s.RuntimeSettings.ApiURL == "/SPN/" {
		w.Write(_escFSMustByte(false, "/index.html"))
		return
	}
	// otherwise, adjust URLs by using template; TODO: only fixes apiURL; index.html itself contains relative CSS/JS
	w.Write(renderIndexTemplate(s.RuntimeSettings))
}

func renderIndexTemplate(runtimeSettings settings.RuntimeSettings) []byte {
	buf := &bytes.Buffer{}
	templateBinary := _escFSMustByte(false, "/index.html.tpl")
	tpl, err := template.New("index").Parse(string(templateBinary))
	runtimeSettings.ApiURL = url.QueryEscape(fmt.Sprintf(`"apiURL":"%s"`, runtimeSettings.ApiURL))
	runtimeSettings.RootURL = url.QueryEscape(fmt.Sprintf(`"rootURL":"%s"`, runtimeSettings.RootURL))
	if err != nil {
		log.Fatalf("Template parsing error: %v\n", err)
	}
	err = tpl.Execute(buf, runtimeSettings)
	if err != nil {
		log.Printf("Template execution error: %v\n", err)
	}
	return buf.Bytes()
}

type tokenRequest struct {
	Username string
	Password string
}
type tokenResponse struct {
	AccessToken string `json:"access_token,omitempty"`
	Error       string `json:"error,omitempty"`
}

func (s *server) tokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	request := tokenRequest{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	if request.Username == "undefined" || request.Password == "undefined" ||
		request.Username == "" || request.Password == "" {
		http.Error(w, `{"error":"empty credentials"}`, 401)
		return
	}

	// verify user/pass, should yield firestore-, bucket-, and spn session tokens
	if request.Username != "popeye" || request.Password != "spin@t" {
		time.Sleep(1 * time.Second)
		http.Error(w, `{"error":"bad credentials"}`, 401)
		return
	}
	/// OUTSOURCE client usage/creation!
	//  https://firebase.google.com/docs/auth/admin/create-custom-tokens#go
	//  https://godoc.org/cloud.google.com/go/firestore
	//  https://godoc.org/firebase.google.com/go
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := client.CustomToken(context.Background(), "some-uid")
	if err != nil {
		log.Fatalf("error minting custom token: %v\n", err)
	}

	log.Printf("Got custom token: %v\n", token)

	response := tokenResponse{
		AccessToken: token,
	}

	responseBody, _ := json.Marshal(response)
	numBytes, err := w.Write(responseBody)
	if err != nil {
		log.Fatalf("Waaaah! Failed sending off %d bytes to client, who will be unhappy for sure: %s", numBytes, err)
	}
}

/*
func authenticateAdminRequest(w http.ResponseWriter, clientToken string, env *Environment) bool {
	if clientToken != env.RuntimeSettings.Token {
		log.Printf("Permission denied for admin request - bad token: %s", clientToken)
		http.Error(w, "FIXME This should be a JSON 401", 401)
		return false
	}
	return true
}

func authenticateVMRequest(w http.ResponseWriter, clientToken string, env *Environment, vmID string) bool {
	taskData, err := cloud.GetArticle(env.RuntimeSettings.ProjectID, vmID)
	if err != nil {
		log.Printf("Error loading task: %s", err)
		http.Error(w, err.Error(), 500)
		return false
	}

	if taskData.ManagementToken != clientToken {
		log.Printf("DENIED: Invalid token: %s", clientToken)
		http.Error(w, err.Error(), http.StatusForbidden)
		return false
	}
	return true
}

func parseRequestURI(uri string) (action, vmid, data string, err error) {
	parts := strings.Split(uri, "/")
	if len(parts) != 4 {
		err = errors.New("invalid URI, expected /CloudVMDocker/ACTION/VM_ID/DATA")
		return
	}
	action = parts[1]
	if action != "delete" && action != "status" && action != "run" && action != "container" {
		err = errors.New("invalid action %s, expected on of: delete, status, run, container")
		return
	}
	vmid = parts[2]
	// todo: validate VMID: starts with t, len=...12?
	data = parts[3]
	return
}

func handleRun(w http.ResponseWriter, r *http.Request, env *Environment, vmID string) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error loading request body: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	taskArguments := task.NewArticleArgumentsFromBytes(requestBody)
	log.Printf("Writing task to FireStore: %+v", taskArguments)
	task := cloud.StoreNewArticle(env.RuntimeSettings.ProjectID, *taskArguments)
	createOp, err := cloud.CreateVM(env.RuntimeSettings, task)
	if err != nil {
		log.Printf("ERROR running ArticleArguments: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}

	log.Println("VM creation requested successfully, waiting for op...")
	err = cloud.WaitForOperation(env.RuntimeSettings.ProjectID, taskArguments.Zone, createOp.Name)
	if err != nil {
		log.Printf("ERROR waiting for VM creation done status: %s", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("Saving GCE InstanceID to FireStore: %s => %d", vmID, createOp.TargetId)
	err = cloud.SetArticleInstanceId(env.RuntimeSettings.ProjectID, vmID, createOp.TargetId)
	if err != nil {
		log.Printf("ARGH!!! Could not update instanceID in FireStore: %s", err)
	}
	task.InstanceID = strconv.FormatUint(createOp.TargetId, 10)

	responseBody, _ := json.Marshal(task)
	w.Header().Set("content-type", "application/json")
	numBytes, err := w.Write(responseBody)
	if err != nil {
		log.Fatalf("Waaaah! Failed sending off %d bytes to client, who will be unhappy for sure: %s", numBytes, err)
	}
}

func handleStatusGet(w http.ResponseWriter, env *Environment, vmID string) {
	log.Printf("Serving status for vm %s", vmID)
	task, err := cloud.GetArticle(env.RuntimeSettings.ProjectID, vmID)
	if err != nil {
		log.Printf("Error handling task status get requests vmid=%s: %s", vmID, err)
		http.Error(w, err.Error(), 500)
		return
	}
	responseBody, _ := json.Marshal(task)
	w.Header().Set("content-type", "application/json")
	w.Write(responseBody)
}

func handleDelete(w http.ResponseWriter, env *Environment, vmID string, exitCodeString string, zone string) {
	log.Printf("Handling DELETE request form VMID %s (in %s) with exitCode %s", vmID, zone, exitCodeString)
	err := cloud.DeleteInstanceByName(env.RuntimeSettings, vmID, zone)
	if err != nil {
		log.Printf("Error on DeleteInstanceByName(..., %s): %s", vmID, err)
		http.Error(w, err.Error(), 500)
		return
	}
	exitCode, err := strconv.Atoi(exitCodeString)
	if err != nil {
		log.Printf("Error on DeleteInstanceByName(..., %s): Cannot convert exit code '%s': %s", vmID, exitCodeString, err)
		http.Error(w, err.Error(), 500)
		return
	}
	err = cloud.UpdateArticleStatus(env.RuntimeSettings.ProjectID, vmID, task.ArticleStatusDone, exitCode)
	if err != nil {
		log.Printf("Error on DeleteInstanceByName(..., %s): Unable to update FireStore after successful VM deletion %s", vmID, err)
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, `Thanks for your DELETE request -- processed successfully`)
}

func handleContainer(w http.ResponseWriter, r *http.Request, env *Environment, vmID string) {
	log.Printf("Handling CONTAINER(ID) request form VMID %s", vmID)
	containerID, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: Cannot read body of container update request: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	if len(containerID) != 64 {
		log.Printf("ERROR: Bad container update request, expected 64-char container ID, got: `%s`", containerID)
		http.Error(w, err.Error(), 500)
		return
	}
	err = cloud.SetArticleContainerID(env.RuntimeSettings.ProjectID, vmID, string(containerID))
	if err != nil {
		log.Printf("ERROR: Failed to update containerID of vm %s to `%s`", vmID, containerID)
		http.Error(w, err.Error(), 500)
		return
	}
	fmt.Fprintf(w, `Thanks for your CONTAINER request -- processed successfully ....`)
}
*/
