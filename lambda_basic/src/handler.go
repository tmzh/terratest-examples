package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

type Event struct {
	Name string `json:"Name"`
}

func HandleRequest(ctx context.Context, evnt Event) (string, error) {
	return fmt.Sprintf("Hello %s!", evnt.Name), nil
}

func main() {
	lambda.Start(HandleRequest)
}
