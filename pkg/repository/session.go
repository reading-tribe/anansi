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

const SessionTableName = "zula_sessions"

type SessionRepository interface {
	GetSession(ctx context.Context, key string) (dbmodel.Session, error)
	CreateSession(ctx context.Context, newSession dbmodel.Session) error
}

type sessionRepository struct{}

func NewSessionRepository() SessionRepository {
	return sessionRepository{}
}

func (s sessionRepository) GetSession(ctx context.Context, key string) (dbmodel.Session, error) {
	item := dbmodel.Session{}

	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return item, fmt.Errorf("GetSession > GetClient: %\n", getClientErr)
	}

	data, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(UserTableName),
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: key},
		},
	})

	if err != nil {
		return item, fmt.Errorf("GetSession: %v\n", err)
	}

	if data.Item == nil {
		return item, fmt.Errorf("GetSession: Session not found.\n")
	}

	err = attributevalue.UnmarshalMap(data.Item, &item)
	if err != nil {
		return item, fmt.Errorf("GetSession > UnmarshalMap: %v\n", err)
	}

	return item, nil
}

func (s sessionRepository) CreateSession(ctx context.Context, newSession dbmodel.Session) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("CreateSession > GetClient: %\n", getClientErr)
	}

	data, err := attributevalue.MarshalMap(newSession)
	if err != nil {
		return fmt.Errorf("CreateSession > MarshalMap: %v\n", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(SessionTableName),
		Item:      data,
	})

	if err != nil {
		return fmt.Errorf("CreateSession > PutItem: %v\n", err)
	}

	return nil
}
