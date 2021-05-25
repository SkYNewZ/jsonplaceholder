package jsonplaceholder_test

import (
	"context"
	"fmt"
	"log"

	"github.com/SkYNewZ/jsonplaceholder"
)

func Example() {
	// Use the default HTTP client
	client := jsonplaceholder.New()

	todo, err := client.Todo.Get(context.Background(), 1)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(todo)
	// Output: &{1 1 delectus aut autem false}
}
