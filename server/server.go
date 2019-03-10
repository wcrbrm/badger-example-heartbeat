package main

import (
	"log"
	"os"

	ginprom "github.com/Depado/ginprom"
	badger "github.com/dgraph-io/badger"
	gin "github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	cli "github.com/jawher/mow.cli"
	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// see more at liu0fanyi/web_for_fun

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

var app = cli.App("heartbeat", "Heartbeat server")

func main() {

	var (
		enableGraphQL = app.Bool(cli.BoolOpt{
			Name:  "graphql",
			Desc:  "Serve API via GraphQL",
			Value: true,
		})
		enableHeartBeat = app.Bool(cli.BoolOpt{
			Name:  "heartbeat",
			Desc:  "Accept heartbeats",
			Value: true,
		})
		enableMetrics = app.Bool(cli.BoolOpt{
			Name:   "metrics",
			Desc:   "Enable prometheus metrics",
			EnvVar: "ENABLE_PROMETHEUS",
			Value:  true,
		})
		webListenAddr = app.String(cli.StringOpt{
			Name:   "L listen",
			Desc:   "Sets listen address for web server",
			EnvVar: "LISTEN_ADDR",
			Value:  "0.0.0.0:8092",
		})
		badgerPath = app.String(cli.StringOpt{
			Name:   "b path",
			Desc:   "Sets path to the badger database",
			EnvVar: "BADGER_PATH",
			Value:  "../db/badger",
		})
	)
	// Specify the action to execute when the app is invoked correctly
	app.Action = func() {
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()
	if *enableMetrics {
		p := ginprom.New(
			ginprom.Engine(r),
			ginprom.Subsystem("gin"),
			ginprom.Path("/metrics"),
		)
		r.Use(p.Instrument())
	}

	if *enableHeartBeat {

		// initializing badger database
		opts := badger.DefaultOptions
		opts.Dir = *badgerPath
		opts.ValueDir = *badgerPath
		db, err := badger.Open(opts)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// mode #1 - heartbeats tracking, using badger as temp storage
		// (this API should not be replicated)
		r.POST("/heartbeat", func(c *gin.Context) {
			// Start a writable transaction.
			txn := db.NewTransaction(true)
			defer txn.Discard()

			// Use the transaction...
			err := txn.Set([]byte("exampleKey"), []byte("exampleValue"))
			if err != nil {
				c.JSON(426, gin.H{"error": err})
				return
			}
			// Commit the transaction and check for error.
			if err := txn.Commit(); err != nil {
				c.JSON(426, gin.H{"error": err})
				return
			}
			c.JSON(200, gin.H{"message": "ok"})
		})

		r.GET("/active", func(c *gin.Context) {

		})
	}

	if *enableGraphQL {
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
	}

	if *enableGraphQL || *enableHeartBeat {
		r.Run(*webListenAddr)
	} else {
		log.Fatal("Either GraphQL or Heartbeat should be enabled")
	}
}
