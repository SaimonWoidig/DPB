package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

const (
	EnvOpensearchUrl      string = "OPENSEARCH_URL"
	EnvOpensearchUser     string = "OPENSEARCH_USER"
	EnvOpensearchPassword string = "OPENSEARCH_PASSWORD"
)

const IndexOrders string = "orders"
const IndexProducts string = "products"
const IndexRecipes string = "recipes"

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

	// create indices
	var createIndexReq opensearchapi.IndicesCreateReq
	var createIndexResp *opensearchapi.IndicesCreateResp
	// products
	createIndexReq = opensearchapi.IndicesCreateReq{
		Index: IndexProducts,
	}
	createIndexResp, err = client.Indices.Create(context.Background(), createIndexReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Product index created - ack: %v, idx: %v\n", createIndexResp.Acknowledged, createIndexResp.Index)
	// orders
	createIndexReq = opensearchapi.IndicesCreateReq{
		Index: IndexOrders,
	}
	createIndexResp, err = client.Indices.Create(context.Background(), createIndexReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Orders index created - ack: %v, idx: %v\n", createIndexResp.Acknowledged, createIndexResp.Index)
	// recipes
	createIndexReq = opensearchapi.IndicesCreateReq{
		Index: IndexRecipes,
	}
	createIndexResp, err = client.Indices.Create(context.Background(), createIndexReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Recipes index created - ack: %v, idx: %v\n", createIndexResp.Acknowledged, createIndexResp.Index)
}
