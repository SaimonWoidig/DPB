package main

import (
	"bytes"
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
const DatafileOrders string = "data/orders-bulk.json"
const IndexProducts string = "products"
const DatafileProducts string = "data/products-bulk.json"
const IndexRecipes string = "recipes"
const DatafileRecipes string = "data/recipes-bulk.json"

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

	// delete old indices
	deleteIndicesReq := opensearchapi.IndicesDeleteReq{Indices: []string{IndexProducts, IndexOrders, IndexRecipes}}
	deleteIndicesResp, err := client.Indices.Delete(context.Background(), deleteIndicesReq)
	if err != nil {
		fmt.Printf("Error deleting indices: %v\n", err.Error())
	} else {
		fmt.Printf("Indices deleted - ack: %v\n", deleteIndicesResp.Acknowledged)
	}
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

	// import data
	// products
	productsData, err := os.ReadFile(DatafileProducts)
	if err != nil {
		panic(err.Error())
	}
	productsDataBuf := bytes.NewBuffer(productsData)
	bulkProductsReq := opensearchapi.BulkReq{Index: IndexProducts, Body: productsDataBuf, Params: opensearchapi.BulkParams{Refresh: "true"}}
	bulkProductsResp, err := client.Bulk(context.Background(), bulkProductsReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Bulk: %v products imported (took: %vms)\n", len(bulkProductsResp.Items), bulkProductsResp.Took)
	// orders
	ordersData, err := os.ReadFile(DatafileOrders)
	if err != nil {
		panic(err.Error())
	}
	ordersDataBuf := bytes.NewBuffer(ordersData)
	bulkOrdersReq := opensearchapi.BulkReq{Index: IndexOrders, Body: ordersDataBuf, Params: opensearchapi.BulkParams{Refresh: "true"}}
	bulkOrdersResp, err := client.Bulk(context.Background(), bulkOrdersReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Bulk: %v orders imported (took: %vms)\n", len(bulkOrdersResp.Items), bulkOrdersResp.Took)
	// recipes
	recipesData, err := os.ReadFile(DatafileRecipes)
	if err != nil {
		panic(err.Error())
	}
	recipesDataBuf := bytes.NewBuffer(recipesData)
	bulkRecipesReq := opensearchapi.BulkReq{Index: IndexRecipes, Body: recipesDataBuf, Params: opensearchapi.BulkParams{Refresh: "true"}}
	bulkRecipesResp, err := client.Bulk(context.Background(), bulkRecipesReq)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Bulk: %v recipes imported (took: %vms)\n", len(bulkRecipesResp.Items), bulkRecipesResp.Took)

	fmt.Println("Done. Indexes created and data imported.")
}
