package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var ddb *dynamodb.DynamoDB

func init() {
	region := os.Getenv("AWS_REGION")
	if session, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		fmt.Println(fmt.Sprintf("Failed to connect to AWS: %s", err.Error()))
	} else {
		ddb = dynamodb.New(session)
	}
}

// CompleteTodo takes in the id of a Todo Item, this will change the Completed bool value
func CompleteTodo(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("CompleteTodo")

	var (
		id        = req.PathParameters["id"]
		tableName = aws.String(os.Getenv("TODOS_TABLE_NAME"))
		done      = "done"
	)

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
		UpdateExpression: aws.String("set #d = :id"),
		ExpressionAttributeNames: map[string]*string{
			"#d": &done,
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":d": {
				BOOL: aws.Bool(true),
			},
		},
		ReturnValues: aws.String("UPDATED_NEW"),
		TableName:    tableName,
	}

	_, err := ddb.UpdateItem(input)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       req.Body,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(CompleteTodo)
}
