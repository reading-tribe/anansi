package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/reading-tribe/anansi/pkg/dbmodel"
	"github.com/reading-tribe/anansi/pkg/dynamodbx"
)

const UserTableName = "zula_users"

type UserRepository interface {
	GetUser(ctx context.Context, emailAddress string) (dbmodel.User, error)
	CreateUser(ctx context.Context, newUser dbmodel.User) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return userRepository{}
}

func (u userRepository) GetUser(ctx context.Context, emailAddress string) (dbmodel.User, error) {
	item := dbmodel.User{}

	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return item, fmt.Errorf("GetUser > GetClient: %\n", getClientErr)
	}

	data, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(UserTableName),
		Key: map[string]types.AttributeValue{
			"email_address": &types.AttributeValueMemberS{Value: emailAddress},
		},
	})

	if err != nil {
		return item, fmt.Errorf("GetUser: %v\n", err)
	}

	if data.Item == nil {
		return item, fmt.Errorf("GetUser: User not found.\n")
	}

	err = attributevalue.UnmarshalMap(data.Item, &item)
	if err != nil {
		return item, fmt.Errorf("GetUser > UnmarshalMap: %v\n", err)
	}

	return item, nil
}

func (u userRepository) CreateUser(ctx context.Context, newUser dbmodel.User) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("CreateUser > GetClient: %\n", getClientErr)
	}

	data, err := attributevalue.MarshalMap(newUser)
	if err != nil {
		return fmt.Errorf("CreateUser > MarshalMap: %v\n", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(UserTableName),
		Item:      data,
	})

	if err != nil {
		return fmt.Errorf("CreateUser > PutItem: %v\n", err)
	}

	return nil
}
