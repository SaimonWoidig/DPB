package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/opensearch-project/opensearch-go/v4/opensearchutil"
)

const OPENSEARCH_ADMIN_PASSWORD string = "ThisIsASecret123" // used only for local development
const INDEX_NAME string = "person"

type Person struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func main() {
	fmt.Println("Connecting to Opensearch")
	client, err := opensearchapi.NewClient(opensearchapi.Config{
		Client: opensearch.Config{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Addresses: []string{"https://localhost:9200"},
			Username:  "admin",
			Password:  OPENSEARCH_ADMIN_PASSWORD,
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
	createIndexReq := opensearchapi.IndicesCreateReq{Index: INDEX_NAME}
	createIndexRes, err := client.Indices.Create(context.Background(), createIndexReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Index created - ack: %v, idx: %v\n", createIndexRes.Acknowledged, createIndexRes.Index)

	// create new person
	newPersonId := "123"
	newPerson := Person{Name: "John", Surname: "Doe"}
	newPersonReq := opensearchapi.IndexReq{
		Index:      INDEX_NAME,
		Body:       opensearchutil.NewJSONReader(&newPerson),
		DocumentID: newPersonId,
	}
	newPersonRes, err := client.Index(context.Background(), newPersonReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created person - id: %v\n", newPersonRes.ID)

	// print created person
	getPersonReq := opensearchapi.DocumentGetReq{Index: INDEX_NAME, DocumentID: newPersonId}
	getPersonRes, err := client.Document.Get(context.Background(), getPersonReq)
	if err != nil {
		panic(err.Error())
	}
	data, err := getPersonRes.Source.MarshalJSON()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Person - id: %v, data: %v\n", getPersonRes.ID, string(data))

	// delete index
	deleteIndexReq := opensearchapi.IndicesDeleteReq{Indices: []string{INDEX_NAME}}
	deleteIndexRes, err := client.Indices.Delete(context.Background(), deleteIndexReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Index deleted - ack: %v\n", deleteIndexRes.Acknowledged)
}
