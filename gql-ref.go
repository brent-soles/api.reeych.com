package lol

import (
	// "encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	//	"github.com/graphql-go/handler"
	"log"
	//	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
	shuuid "github.com/lithammer/shortuuid"
	"reeychapi"
	"reflect"
)

type Foo struct {
	Name string
}

var FieldFooType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Foo",
	Fields: graphql.Fields{
		"name": &graphql.Field{Type: graphql.String},
	},
})

type Bar struct {
	Name string
}

var FieldBarType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Bar",
	Fields: graphql.Fields{
		"name": &graphql.Field{Type: graphql.String},
	},
})

func resolveBarType(p graphql.ResolveParams) (interface{}, error) {

	fmt.Println(reflect.TypeOf(p))
	type result struct {
		data interface{}
		err  error
	}
	ch := make(chan *result, 1)
	go func() {
		defer close(ch)
		bar := &Bar{Name: "Bar's name"}
		ch <- &result{data: bar, err: nil}
	}()
	return func() (interface{}, error) {
		r := <-ch
		return r.data, r.err
	}, nil
}

// QueryType fields: `concurrentFieldFoo` and `concurrentFieldBar` are resolved
// concurrently because they belong to the same field-level and their `Resolve`
// function returns a function (thunk).
var QueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"concurrentFieldFoo": &graphql.Field{
			Type: FieldFooType,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var foo = Foo{Name: "Foo's name"}
				return func() (interface{}, error) {
					return &foo, nil
				}, nil
			},
		},
		"concurrentFieldBar": &graphql.Field{
			Type:    FieldBarType,
			Resolve: resolveBarType,
		},
	},
})

func main() {

	/*schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: QueryType,
	})
	*/
	/* gqlHandler := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})
	*/

	connStr := "user=brent dbname=reeych_backend_dev sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(db)

	rows, err := db.Query("SELECT uc.id, first_name, last_name, title, owner FROM (users JOIN cards_users_rel ON users.id = user_id) as uc JOIN cards ON uc.card_id = cards.id")

	defer rows.Close()

	for rows.Next() {
		var (
			userId    string
			firstName string
			lastName  string
			title     string
			owner     string
		)

		if err := rows.Scan(&userId, &firstName, &lastName, &title, &owner); err != nil {
			log.Fatal(err)
		}

		fmt.Println(userId, firstName, lastName, title, owner)
	}

	fmt.Println(shuuid.New())
	// http.Handle("/gql", gqlHandler)
	// http.ListenAndServe(":7080", nil)

	/*
		schema, err := graphql.NewSchema(graphql.SchemaConfig{
			Query: QueryType,
		})
		if err != nil {
			log.Fatal(err)
		}

		query := `
			query {
				concurrentFieldFoo {
					name
				}
				concurrentFieldBar {
					name
				}
			}
		`

		query := `
			query {
				concurrentFieldFoo {
					name
				}
			}
		`

		result := graphql.Do(graphql.Params{
			RequestString: query,
			Schema:        schema,
		})
		b, err := json.Marshal(result)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", b)
	*/
	/*
		{
		  "data": {
		    "concurrentFieldBar": {
		      "name": "Bar's name"
		    },
		    "concurrentFieldFoo": {
		      "name": "Foo's name"
		    }
		  }
		}
	*/
}
