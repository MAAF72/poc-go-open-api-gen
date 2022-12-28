package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/iancoleman/strcase"
)

var (
	httpStatus           = []int{200, 400, 401, 500}
	mapRefWithStructName = map[string]string{}
)

// GeneratePath :nodoc
func GeneratePath(path string, data openapi3.PathItem) {
	fmt.Println("Generating Path", path)
	GeneratePathOperation(path, data.Get)
	GeneratePathOperation(path, data.Post)
	GeneratePathOperation(path, data.Put)
	GeneratePathOperation(path, data.Patch)
	GeneratePathOperation(path, data.Delete)
	fmt.Println("======================================")
}

// GeneratePathOperation :nodoc
func GeneratePathOperation(path string, data *openapi3.Operation) {
	if data == nil {
		return
	}

	if data.OperationID == "" {
		panic(fmt.Sprintf("must have OperationID in %s", path))
	}

	GenerateParameters(path, data.OperationID, data.Parameters)

	if data.RequestBody != nil {
		GenerateRequestBody(path, data.OperationID, *data.RequestBody)
	}

	for _, code := range httpStatus {
		if dataResponse := data.Responses.Get(code); dataResponse != nil {
			GenerateResponseBody(path, data.OperationID, code, *dataResponse)
		}
	}
}

// GenerateRequestBody :nodoc
// example
//
//	type ApplyPromoRequest {
//		... ....
//	}
func GenerateRequestBody(path string, operationID string, data openapi3.RequestBodyRef) {
	if data.Value == nil {
		return
	}

	fmt.Println("Generating Request Body", path, operationID)

	if mediaType := data.Value.Content.Get("application/json"); mediaType != nil {
		if mediaType.Schema == nil {
			return
		}

		GenerateSchema(fmt.Sprintf("%sRequest", strcase.ToCamel(operationID)), "body", operationID, *mediaType.Schema)
	}
}

// GenerateParameter :nodoc
func GenerateParameter(path string, operationID string, data openapi3.ParameterRef) {
	if data.Value == nil {
		return
	}

	// fmt.Println("Generating Parameter", path, operationID, data.Value.Name, data.Value.In, data.Value.Required)

	if data.Value.Schema != nil {
		// fmt.Println("Validate", data.Value.Schema.Validate(context.TODO()))
		GenerateSchema(data.Value.Name, data.Value.In, operationID, *data.Value.Schema)
	}

}

// GenerateParameters :nodoc
// example
//
//	type ApplyPromoParameters {
//		Path ApplyPromoPathParameter
//		Query ApplyPromoQueryParameter
//		Header ApplyPromoHeaderParameter
//		Cookie ApplyPromoCookieParameter
//	}
func GenerateParameters(path string, operationID string, data openapi3.Parameters) {
	fmt.Println("Generating Parameters", path, data.Validate(context.TODO()))

	listParameterInPath := make([]openapi3.ParameterRef, 0)
	listParameterInQuery := make([]openapi3.ParameterRef, 0)
	listParameterInHeader := make([]openapi3.ParameterRef, 0)
	listParameterInCookie := make([]openapi3.ParameterRef, 0)

	for _, val := range data {
		if val == nil {
			continue
		}

		switch val.Value.In {
		case openapi3.ParameterInPath:
			listParameterInPath = append(listParameterInPath, *val)
		case openapi3.ParameterInQuery:
			listParameterInQuery = append(listParameterInQuery, *val)
		case openapi3.ParameterInHeader:
			listParameterInHeader = append(listParameterInHeader, *val)
		case openapi3.ParameterInCookie:
			listParameterInCookie = append(listParameterInCookie, *val)

		}
	}

	fmt.Printf("type %sParameters struct {\n", strcase.ToCamel(operationID))

	if len(listParameterInPath) > 0 {
		fmt.Printf("\tPath\t%sPathParameter\n", strcase.ToCamel(operationID))
	}

	if len(listParameterInQuery) > 0 {
		fmt.Printf("\tQuery\t%sQueryParameter\n", strcase.ToCamel(operationID))
	}

	if len(listParameterInHeader) > 0 {
		fmt.Printf("\tHeader\t%sHeaderParameter\n", strcase.ToCamel(operationID))
	}

	if len(listParameterInCookie) > 0 {
		fmt.Printf("\tCookie\t%sCookieParameter\n", strcase.ToCamel(operationID))
	}

	fmt.Printf("}\n")

	GenerateParameterByLocation(path, operationID, "path", listParameterInPath)
	GenerateParameterByLocation(path, operationID, "query", listParameterInQuery)
	GenerateParameterByLocation(path, operationID, "header", listParameterInHeader)
	GenerateParameterByLocation(path, operationID, "cookie", listParameterInCookie)

}

// GenerateParameterByLocation :nodoc
func GenerateParameterByLocation(path string, operationID string, location string, parameters []openapi3.ParameterRef) {
	fmt.Printf("type %s%sParameter struct {\n", strcase.ToCamel(operationID), strcase.ToCamel(location))
	for _, parameter := range parameters {
		GenerateParameter(path, operationID, parameter)
	}
	fmt.Printf("}\n")
}

// GenerateResponseBody :nodoc
// example
//
//	type ApplyPromoOKResponse {
//		.... .....
//	}
//
//	type ApplyPromoInternalServerErrorResponse {
//		.... .....
//	}
//
//	type GeneralNotAuthorizedResponse {
//		.... .....
//	}
//
//	type GeneralForbiddenResponse {
//		.... .....
//	}
func GenerateResponseBody(path string, operationID string, code int, data openapi3.ResponseRef) {
	if data.Value == nil {
		return
	}

	// fmt.Println("Generating Response Body", path, operationID, code)

	if mediaType := data.Value.Content.Get("application/json"); mediaType != nil {
		if mediaType.Schema == nil {
			return
		}

		// fmt.Println("Validate", mediaType.Validate(context.TODO()))

		fmt.Printf("type %s%sResponse {\n", strcase.ToCamel(operationID), http.StatusText(code))
		GenerateSchema(fmt.Sprintf("%sResponse", strcase.ToCamel(operationID)), "response", operationID, *mediaType.Schema)
		fmt.Printf("}\n")
	}
}

// GenerateSchema :nodoc
func GenerateSchema(primaryVariableName string, in string, operationID string, data openapi3.SchemaRef) {
	var variableName, dataType string

	// fmt.Println("Validate", data.Validate(context.TODO()))
	// fmt.Print(operationID, "->", in, " ")
	if data.Ref != "" {
		ref := strings.Split(data.Ref, "/")
		variableName = ref[len(ref)-1]
	} else if primaryVariableName != "" {
		variableName = primaryVariableName
	} else if data.Value.Title != "" {
		variableName = data.Value.Title
	} else {
		variableName = "NO_NAME"
	}

	switch strings.ToLower(data.Value.Type) {
	case "string":
		dataType = "string"
	case "integer":
		dataType = "int32"
	case "number":
		dataType = "float64"
	case "array":
		if data.Value.Items != nil {
			dataType = fmt.Sprintf("[]%s", data.Value.Items.Value.Type)
		}
	case "object":
		fmt.Println("Schema type : object")
		fmt.Println("Properties")
		// Give nested object variable name

		for key, val := range data.Value.Properties {
			fmt.Println(key, val)
			// if val.Value.Type == "object" {
			// 	GenerateSchema(variableName+strcase.ToCamel(key), in, operationID, name+"_"+key, *val)
			// }
		}
	}

	fmt.Printf("\t%s\t%s\t%s\n", variableName, dataType, strings.ToLower(data.Value.Type))
}
