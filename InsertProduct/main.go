package main

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

type Product struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Price int `json:"price"`
}

func InsertProduct(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
	var product Product
	err:= json.Unmarshal([]byte(request.Body), &product)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	product.ID = uuid.New().String()
	
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	input := &dynamodb.PutItemInput{
		TableName: aws.String("ProductsVideo"),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(product.ID),
			},
			"name": {
				S: aws.String((product.Name)),
			},
			"price": {
				N: aws.String(strconv.Itoa(product.Price)),
			},
		},
	}

	_, err = svc.PutItem((input))
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	body, err := json.Marshal(product)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(body),
	}, nil

}

func main() {
	lambda.Start(InsertProduct)
}