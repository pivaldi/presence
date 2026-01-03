package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/pivaldi/presence/examples/gqlgen/graph"
	"github.com/pivaldi/presence/examples/gqlgen/graph/generated"
)

func main() {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{
			Resolvers: graph.NewResolver(),
		}),
	)

	http.Handle("/", playground.Handler("GraphQL Playground", "/query"))
	http.Handle("/query", srv)

	log.Println("GraphQL server running at http://localhost:8181/")
	log.Println("Open http://localhost:8181/ for GraphQL Playground")
	log.Fatal(http.ListenAndServe(":8181", nil)) //nolint:gosec // G114: example server, timeouts not needed
}
