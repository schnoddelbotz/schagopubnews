package cloud

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"github.com/schnoddelbotz/schagopubnews/article"
	"github.com/schnoddelbotz/schagopubnews/settings"
)

// Thing below should be split into NewArticle(proj, args).Store()

// StoreNewArticle saves a new task in FireStore DB.
func StoreNewArticle(projectID string /*, taskArguments article.ArticleArguments*/) article.Article {
	ctx := context.Background()
	client := NewFireStoreClient(ctx, projectID)
	doc := article.Article{
		Status: article.ArticleStatusCreated,
	}
	// Saves the new entity.
	if _, err := client.Collection(settings.FireStoreCollection).Doc("my-first-doc").Set(ctx, doc); err != nil {
		log.Fatalf("Failed to save doc: %v", err)
	}
	log.Printf("Saved %v: %v", "my-first-doc", doc.Status)
	return doc
}

// ListArticles provides 'docker ps' functionality by querying FireStore
func ListArticles(projectID string) {
	ctx := context.Background() // fixme pass in
	client := NewFireStoreClient(ctx, projectID)
	docList := make([]article.Article, 0)

	// todo: add filter (-a arg), sorting, FIX CREATED OUTPUT
	iter := client.Collection(settings.FireStoreCollection).
		OrderBy("CreatedAt", firestore.Asc).
		//Limit(25). // FIXME: no-constant-magic ... make user flag
		Documents(ctx)
	defer iter.Stop()
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("FireStore iterator boom: %s", err)
		}

		var task article.Article
		if err := doc.DataTo(&task); err != nil {
			log.Printf("FireStore data error: %s", err)
			continue
		}
		docList = append(docList, task)
	}

	for _, doc := range docList {
		fmt.Printf("%+v", doc)
	}

	// log.Printf("Got %d tasks as response, showed X, 3 running, 2 deleted.", len(something = client.GetAll retval))
}

// NewFireStoreClient returns a dataStore client and its context, exits fatally on error
func NewFireStoreClient(ctx context.Context, projectID string) *firestore.Client {
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

/*
// UpdateArticleStatus sets a Article's 'status' field to given string value
func UpdateArticleStatus(projectID, vmID, status string, exitCode ...int) error {
	log.Printf("UpdateArticleStatus: %s -> %s", vmID, status)
	ctx := context.Background()
	client := NewFireStoreClient(ctx, projectID)
	update := map[string]interface{}{"Status": status}
	if len(exitCode) > 0 {
		update["DockerExitCode"] = exitCode[0]
	}
	_, err := client.Collection(settings.FireStoreCollection).Doc(vmID).
		Set(ctx, update, firestore.MergeAll)
	if err != nil {
		return err
	}
	return nil
}

// SetArticleContainerID sets a Article's 'status' field to given string value
func SetArticleContainerID(projectID, vmID, containerID string) error {
	log.Printf("SetArticleContainerID: %s -> %s", vmID, containerID)
	ctx := context.Background()
	client := NewFireStoreClient(ctx, projectID)
	_, err := client.Collection(settings.FireStoreCollection).Doc(vmID).
		Set(ctx, map[string]interface{}{"DockerContainerId": containerID}, firestore.MergeAll)
	if err != nil {
		return err
	}
	return nil
}

// UpdateArticleStatus sets a Article's 'status' field to given string value
func SetArticleInstanceId(projectID, vmID string, instanceID uint64) error {
	log.Printf("SetArticleInstanceId: VM_ID %s -> InstanceID %d", vmID, instanceID)
	ctx := context.Background()
	client := NewFireStoreClient(ctx, projectID)
	_, err := client.Collection(settings.FireStoreCollection).Doc(vmID).
		Set(ctx, map[string]interface{}{"InstanceID": strconv.FormatUint(instanceID, 10)}, firestore.MergeAll)
	if err != nil {
		return err
	}
	return nil
}

// GetArticle tries to fetch given record from FireStore
func GetArticle(projectID, vmID string) (article.Article, error) {
	ctx := context.Background() // fixme pass in
	client := NewFireStoreClient(ctx, projectID)
	var myArticle article.Article
	d, err := client.Collection(settings.FireStoreCollection).Doc(vmID).Get(ctx)
	if err != nil {
		return article.Article{}, err
	}
	if err := d.DataTo(&myArticle); err != nil {
		log.Fatalf("Ooops. Cannot convert FireStore data to task %s: %s", vmID, err)
	}
	return myArticle, err
}
*/
