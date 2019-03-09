package main

import (
	"fmt"

	ginprom "github.com/Depado/ginprom"
	_ "github.com/dgraph-io/badger"
	gin "github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Interval struct should reflect database record
type Interval struct {
	ID      string `json:"id"`
	TsStart int64  `json:"tsStart"`
	TsEnd   int64  `json:"tsEnd"`
	G       string `json:"g"`
	U       string `json:"u"`
	D       string `json:"d"`
}

// IntervalList is and array of Interval
// global variable?
// var IntervalList []Interval

// IntervalType - custom GraphQL ObjectType for our Golang struct `Interval`
// Note that
// - the fields in our IntervalType maps with the json tags for the fields in our struct
// - the field type matches the field type in our struct
var IntervalType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Interval",
	Fields: graphql.Fields{
		"id":      &graphql.Field{Type: graphql.String},
		"tsStart": &graphql.Field{Type: graphql.Int},
		"tsEnd":   &graphql.Field{Type: graphql.Int},
		"g":       &graphql.Field{Type: graphql.String},
		"u":       &graphql.Field{Type: graphql.String},
		"d":       &graphql.Field{Type: graphql.String},
	},
})

// intervalsMutation should be empty, it is readonly for GraphQL users
var intervalsMutation = graphql.Fields{}

// intervalsQuery
// we just define a trivial example here, since root query is required.
var intervalsQuery = graphql.Fields{
	"interval": &graphql.Field{
		Type:        IntervalType,
		Description: "Get single interval",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			idQuery, isOK := params.Args["id"].(string)
			var IntervalList []Interval
			if isOK {
				// Search for el with id
				for _, interval := range IntervalList {
					if interval.ID == idQuery {
						return interval, nil
					}
				}
			}
			return Interval{}, nil
		},
	},

	/*
	   curl -g 'http://localhost:8092/graphql?query={intervalList{id,tsStart,tsEnd,u,g,d}}'
	*/
	"intervalList": &graphql.Field{
		Type:        graphql.NewList(IntervalType),
		Description: "List of intervals",
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			var IntervalList []Interval
			return IntervalList, nil
		},
	},
}

// root mutation
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields(
		intervalsMutation,
	),
})

// root query
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields(
		intervalsQuery,
	),
})

// Schema is an export
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func main() {
	fmt.Println("Starting")

	r := gin.Default()

	p := ginprom.New(
		ginprom.Engine(r),
		ginprom.Subsystem("gin"),
		ginprom.Path("/metrics"),
	)
	r.Use(p.Instrument())

	// mode #1 - heartbeats tracking, using badger as temp storage
	// (should not be replicated)
	r.POST("/heartbeat", func(c *gin.Context) {})
	r.GET("/active", func(c *gin.Context) {})
	r.GET("/all", func(c *gin.Context) {})

	// mode #2 - graphql server
	// (stateless, can be replicated)
	// Creates a GraphQL-go HTTP handler with the defined schema
	hGraphql := handler.New(&handler.Config{
		Schema:   &Schema,
		Pretty:   true,
		GraphiQL: true,
	})
	graphQL := func(c *gin.Context) {
		hGraphql.ServeHTTP(c.Writer, c.Request)
	}
	r.POST("/graphql", graphQL)
	r.GET("/graphql", graphQL)

	// TODO: port configuration using mow.cli
	r.Run("0.0.0.0:8092")
}
