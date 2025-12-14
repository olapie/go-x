package main

import (
	"fmt"
)

func main() {
	exampleCamelToSnake()
	examplePascalToSnake()
	examplePascalToSnakeWithAcronym()
	exampleSnakeToCamel()
	exampleSnakeToCamelWithAcronym()
	exampleSnakeToPascal()
	exampleSnakeToPascalWithAcronym()
	examplePascalToKebab()
	examplePascalToKebabWithAcronym()
	exampleCustomAcronym()
}

func exampleCamelToSnake() {
	camel := "httpRequestId"
	fmt.Println(xname.ToSnake(camel)) // output: http_request_id
}

func examplePascalToSnake() {
	pascal := "HttpRequestId"
	fmt.Println(xname.ToSnake(pascal)) // output: http_request_id
}

func examplePascalToSnakeWithAcronym() {
	pascal := "HTTPRequestID"
	fmt.Println(xname.ToSnake(pascal, xname.WithAcronym())) // output: http_request_id
}

func exampleSnakeToCamel() {
	snake := "http_request_id"
	fmt.Println(xname.ToCamel(snake)) // output: httpRequestId
}

func exampleSnakeToCamelWithAcronym() {
	snake := "http_request_id"
	fmt.Println(xname.ToCamel(snake, xname.WithAcronym())) // output: httpRequestID
}

func exampleSnakeToPascal() {
	snake := "http_request_id"
	fmt.Println(xname.ToPascal(snake)) // output: HttpRequestId
}

func exampleSnakeToPascalWithAcronym() {
	snake := "http_request_id"
	fmt.Println(xname.ToPascal(snake, xname.WithAcronym())) // output: HTTPRequestID
}

func examplePascalToKebab() {
	pascal := "HttpRequestId"
	fmt.Println(xname.ToKebab(pascal)) // output: http-request-id
}

func examplePascalToKebabWithAcronym() {
	pascal := "HTTPRequestID"
	fmt.Println(xname.ToKebab(pascal)) // output: http-request-id
}

func exampleCustomAcronym() {
	pascal := "CUSTOMAcronym"
	xname.AddAcronym("CUSTOM")
	fmt.Println(xname.ToSnake(pascal, xname.WithAcronym())) // output: custom_acronym
}
