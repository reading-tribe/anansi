package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/reading-tribe/anansi/pkg/dbmodel"
	"github.com/reading-tribe/anansi/pkg/dynamodbx"
	"github.com/reading-tribe/anansi/pkg/idx"
)

const SessionTableName = "zula_sessions"

type SessionRepository interface {
	GetSession(ctx context.Context, key idx.SessionID) (dbmodel.Session, error)
	CreateSession(ctx context.Context, newSession dbmodel.Session) error
	ListSessions(ctx context.Context) ([]dbmodel.Session, error)
	UpdateSession(ctx context.Context, updatedSession dbmodel.Session) error
	DeleteSession(ctx context.Context, key idx.SessionID) error
}

type sessionRepository struct{}

func NewSessionRepository() SessionRepository {
	return sessionRepository{}
}

func (s sessionRepository) GetSession(ctx context.Context, key idx.SessionID) (dbmodel.Session, error) {
	item := dbmodel.Session{}

	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return item, fmt.Errorf("GetSession > GetClient: %v\n", getClientErr)
	}

	data, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(UserTableName),
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: key.String()},
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
		return fmt.Errorf("CreateSession > GetClient: %v\n", getClientErr)
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

func (s sessionRepository) ListSessions(ctx context.Context) ([]dbmodel.Session, error) {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return nil, fmt.Errorf("ListSessions > GetClient: %v\n", getClientErr)
	}

	items := []dbmodel.Session{}

	data, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(SessionTableName),
		IndexName:                 aws.String("key"),
		ExpressionAttributeValues: map[string]types.AttributeValue{},
	})
	if err != nil {
		return items, fmt.Errorf("ListSessions > Query: %v\n", err)
	}

	err = attributevalue.UnmarshalListOfMaps(data.Items, &items)
	if err != nil {
		return items, fmt.Errorf("ListSessions > UnmarshalListOfMaps: %v\n", err)
	}

	return items, nil
}

func (s sessionRepository) UpdateSession(ctx context.Context, updatedSession dbmodel.Session) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("UpdateSession > GetClient: %v\n", getClientErr)
	}

	_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(SessionTableName),
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: updatedSession.Key.String()},
		},
		UpdateExpression: aws.String("set expires_at = :expires_at"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":expires_at": &types.AttributeValueMemberS{Value: strconv.FormatInt(updatedSession.ExpiresAtUnix, 10)},
		},
	})

	if err != nil {
		return fmt.Errorf("UpdateSession > UpdateItem: %v\n", err)
	}

	return nil
}

func (s sessionRepository) DeleteSession(ctx context.Context, key idx.SessionID) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("DeleteSession > GetClient: %v\n", getClientErr)
	}

	_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(SessionTableName),
		Key: map[string]types.AttributeValue{
			"key": &types.AttributeValueMemberS{Value: key.String()},
		},
	})
	if err != nil {
		return fmt.Errorf("DeleteSession > DeleteItem: %v\n", err)
	}

	return nil
}
