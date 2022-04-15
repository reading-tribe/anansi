package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/reading-tribe/anansi/pkg/authorizer"
)

func main() {
	fn := authorizer.GetAuthorizer()
	lambda.Start(fn)
}
