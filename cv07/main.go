package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/opensearch-project/opensearch-go/v4/opensearchutil"
)

const (
	EnvOpensearchUrl      string = "OPENSEARCH_URL"
	EnvOpensearchUser     string = "OPENSEARCH_USER"
	EnvOpensearchPassword string = "OPENSEARCH_PASSWORD"
)

const IndexName string = "person"

type Person struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err.Error())
	}

	// load configuration parameters from dotenv
	url := os.Getenv(EnvOpensearchUrl)
	if len(url) == 0 {
		panic("OPENSEARCH_URL is unset")
	}
	username := os.Getenv(EnvOpensearchUser)
	if len(username) == 0 {
		panic("OPENSEARCH_USER is unset")
	}
	password := os.Getenv(EnvOpensearchPassword)
	if len(password) == 0 {
		panic("OPENSEARCH_PASSWORD is unset")
	}

	fmt.Printf("Connecting to Opensearch at %q\n", url)
	client, err := opensearchapi.NewClient(opensearchapi.Config{
		Client: opensearch.Config{
			Addresses: []string{url},
			Username:  username,
			Password:  password,
		},
	})
	if err != nil {
		panic(err.Error())
	}
	infoRes, err := client.Info(context.Background(), nil)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Cluster INFO:\n  Cluster Name: %s\n  Cluster UUID: %s\n  Version Number: %s\n", infoRes.ClusterName, infoRes.ClusterUUID, infoRes.Version.Number)

	// create index
	createIndexReq := opensearchapi.IndicesCreateReq{
		Index: IndexName,
	}
	createIndexRes, err := client.Indices.Create(context.Background(), createIndexReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Index created - ack: %v, idx: %v\n", createIndexRes.Acknowledged, createIndexRes.Index)

	// create new person
	newPersonId := "123"
	newPerson := Person{Name: "John", Surname: "Doe"}
	newPersonReq := opensearchapi.IndexReq{
		Index:      IndexName,
		Body:       opensearchutil.NewJSONReader(&newPerson),
		DocumentID: newPersonId,
		Params: opensearchapi.IndexParams{
			Refresh: "true",
		},
	}
	newPersonRes, err := client.Index(context.Background(), newPersonReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created person - id: %v\n", newPersonRes.ID)

	// print created person
	getPersonReq := opensearchapi.DocumentGetReq{
		Index:      IndexName,
		DocumentID: newPersonId,
	}
	getPersonRes, err := client.Document.Get(context.Background(), getPersonReq)
	if err != nil {
		panic(err.Error())
	}
	data, err := getPersonRes.Source.MarshalJSON()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Person - id: %v, data: %v\n", getPersonRes.ID, string(data))

	// update person
	updatePersonReq := opensearchapi.UpdateReq{
		Index:      IndexName,
		DocumentID: newPersonId,
		Body:       strings.NewReader(`{"doc": {"name": "Bob"}}`),
		Params: opensearchapi.UpdateParams{
			Refresh: "true",
		},
	}
	updatePersonRes, err := client.Update(context.Background(), updatePersonReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Updated person - id: %v, status %v\n", updatePersonRes.ID, updatePersonRes.Result)

	// print all in index person
	getAllPersonReq := opensearchapi.SearchReq{
		Indices: []string{IndexName},
		Body:    strings.NewReader(`{"query": {"match_all": {}}}`),
	}
	getAllPersonRes, err := client.Search(context.Background(), &getAllPersonReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Found:")
	for _, hit := range getAllPersonRes.Hits.Hits {
		data, err := hit.Source.MarshalJSON()
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("  - id: %v, data: %v\n", hit.ID, string(data))
	}

	// delete person
	deletePersonReq := opensearchapi.DocumentDeleteReq{
		Index:      IndexName,
		DocumentID: newPersonId,
	}
	deletePersonRes, err := client.Document.Delete(context.Background(), deletePersonReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Deleted person - id: %v, status %v\n", deletePersonRes.ID, deletePersonRes.Result)

	// delete index
	deleteIndexReq := opensearchapi.IndicesDeleteReq{
		Indices: []string{IndexName},
	}
	deleteIndexRes, err := client.Indices.Delete(context.Background(), deleteIndexReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Index deleted - ack: %v\n", deleteIndexRes.Acknowledged)
}
