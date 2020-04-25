package main

import (
	"fmt"
	"time"

	"github.com/Knetic/govaluate"
)

func test1() {
	expression, _ := govaluate.NewEvaluableExpression("date > '2014-01-01 23:59:59'")
	parameters := make(map[string]interface{}, 8)
	date, _ := time.ParseInLocation("2006-01-02", "2016-12-15", time.Local)
	fmt.Println(date)
	parameters["date"] = date.Unix()
	result, _ := expression.Evaluate(parameters)

	// para := make(map[string]interface{}, 8)
	// para["foo"] = -1

	fmt.Println(result)
}

func test2() {
	functions := map[string]govaluate.ExpressionFunction{
		"strlen": func(args ...interface{}) (interface{}, error) {
			length := len(args[0].(string))
			return (float64)(length), nil
		},
	}

	expString := "strlen(str) <= 16"
	expression, _ := govaluate.NewEvaluableExpressionWithFunctions(expString, functions)

	parameters := make(map[string]interface{}, 8)
	parameters["str"] = "somereallylongstring"
	result, _ := expression.Evaluate(parameters)
	fmt.Println(result)
}
func main() {
	test1()
}
