package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

func debug(t string, v interface{}) {
	fmt.Printf("%s %#v\n", t, v)
}

func debugJSON(t string, v interface{}) {
	fmt.Println(t)
	tmp, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(tmp))
}

func main() {
	// pikirin lagi tentang generated and non generated code
	// gimana kalau usernya mau ubah generated code?
	ctx := context.Background()
	doc, err := openapi3.NewLoader().LoadFromFile("docs/api/swagger_minimal.yaml")
	// doc, err := openapi3.NewLoader().LoadFromFile("docs/api/swagger.yaml")
	if err != nil {
		fmt.Println(err)
	}

	if validationErr := doc.Validate(ctx); validationErr != nil {
		fmt.Println(validationErr)
	}

	for path, data := range doc.Paths {
		if data == nil {
			continue
		}
		GeneratePath(path, *data)
	}
}
