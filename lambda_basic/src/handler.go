package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context, name string) (string, error) {
	return fmt.Sprintf("Hello %s", name), fmt.Errorf("Failed to handle %#v", evnt)
}

func main() {
	lambda.Start(HandleRequest)
}
