package main

import (
	"log"

	"github.com/graphql-go/graphql"
)

func main() {
	playerType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Player",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"number": &graphql.Field{
				Type: graphql.Int,
			},
			"club": &graphql.Field{
				Type: graphql.String,
			},
			"position": &graphql.Field{
				Type: graphql.NewList(graphql.String),
			},
		},
	})

	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"player": &graphql.Field{
				Type:        playerType,
				Description: "Get single player",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return nil, nil
				},
			},
		},
	})

	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
}
