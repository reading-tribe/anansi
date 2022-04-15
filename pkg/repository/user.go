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
	"github.com/reading-tribe/anansi/pkg/idx"
	"github.com/segmentio/ksuid"
)

const UserTableName = "zula_users"

type UserRepository interface {
	GetUser(ctx context.Context, id idx.UserID) (dbmodel.User, error)
	GetUserByEmailAddress(ctx context.Context, emailAddress string) (dbmodel.User, error)
	ListUsersByCode(ctx context.Context, code string) ([]dbmodel.User, error)
	CreateUser(ctx context.Context, newUser dbmodel.User) error
	ListUsers(ctx context.Context) ([]dbmodel.User, error)
	UpdateUser(ctx context.Context, updatedUser dbmodel.User) error
	DeleteUser(ctx context.Context, id idx.UserID) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return userRepository{}
}

func (u userRepository) GetUser(ctx context.Context, id idx.UserID) (dbmodel.User, error) {
	item := dbmodel.User{}

	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return item, fmt.Errorf("GetUser > GetClient: %\n", getClientErr)
	}

	data, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(UserTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
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

func (u userRepository) GetUserByEmailAddress(ctx context.Context, emailAddress string) (dbmodel.User, error) {
	item := dbmodel.User{}
	items := []dbmodel.User{}

	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return item, fmt.Errorf("GetUserByEmailAddress > GetClient: %v\n", getClientErr)
	}

	data, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(UserTableName),
		IndexName: aws.String("email_address-index"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{
				Value: emailAddress,
			},
		},
		KeyConditionExpression: aws.String("email_address = :email"),
	})

	if err != nil {
		return item, fmt.Errorf("GetUserByEmailAddress > Query: %v\n", err)
	}

	err = attributevalue.UnmarshalListOfMaps(data.Items, &items)
	if err != nil {
		return item, fmt.Errorf("GetUserByEmailAddress > UnmarshalListOfMaps: %v\n", err)
	}

	if len(items) != 1 {
		return item, fmt.Errorf("GetUserByEmailAddress: unexpected number of items %d\n", len(items))
	}

	return items[0], nil
}

func (u userRepository) ListUsersByCode(ctx context.Context, code string) ([]dbmodel.User, error) {
	items := []dbmodel.User{}

	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return items, fmt.Errorf("ListUsersByCode > GetClient: %v\n", getClientErr)
	}

	data, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(UserTableName),
		IndexName: aws.String("invite_code-index"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":code": &types.AttributeValueMemberS{
				Value: code,
			},
		},
		KeyConditionExpression: aws.String("invite_code = :code"),
	})

	if err != nil {
		return items, fmt.Errorf("ListUsersByCode > Query: %v\n", err)
	}

	err = attributevalue.UnmarshalListOfMaps(data.Items, &items)
	if err != nil {
		return items, fmt.Errorf("ListUsersByCode > UnmarshalListOfMaps: %v\n", err)
	}

	return items, nil
}

func (u userRepository) CreateUser(ctx context.Context, newUser dbmodel.User) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("CreateUser > GetClient: %v\n", getClientErr)
	}

	random, randomErr := ksuid.NewRandom()
	if randomErr != nil {
		return fmt.Errorf("CreateUser > NewRandom: %v\n", randomErr)
	}
	code := random.String()

	newUser.InviteCode = code

	id, getIDErr := idx.NewUserID()
	if getIDErr != nil {
		return fmt.Errorf("CreateUser > NewUserID: %v\n", getIDErr)
	}

	newUser.ID = id

	data, err := attributevalue.MarshalMap(newUser.AsMap())
	if err != nil {
		return fmt.Errorf("CreateUser > MarshalMap: %v\n", err)
	}

	fmt.Printf("CreateUser: %v\n", data)

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(UserTableName),
		Item:      data,
	})

	if err != nil {
		return fmt.Errorf("CreateUser > PutItem: %v\n", err)
	}

	return nil
}

func (u userRepository) ListUsers(ctx context.Context) ([]dbmodel.User, error) {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return nil, fmt.Errorf("ListUsers > GetClient: %v\n", getClientErr)
	}

	items := []dbmodel.User{}

	data, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(UserTableName),
		IndexName:                 aws.String("email_address"),
		ExpressionAttributeValues: map[string]types.AttributeValue{},
	})
	if err != nil {
		return items, fmt.Errorf("ListUsers > Query: %v\n", err)
	}

	err = attributevalue.UnmarshalListOfMaps(data.Items, &items)
	if err != nil {
		return items, fmt.Errorf("ListUsers > UnmarshalListOfMaps: %v\n", err)
	}

	return items, nil
}

func (u userRepository) UpdateUser(ctx context.Context, updatedUser dbmodel.User) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("UpdateUser > GetClient: %\n", getClientErr)
	}

	_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(UserTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: updatedUser.ID.String()},
		},
		UpdateExpression: aws.String("set password_hash = :password_hash, invite_code = :code, confirmed = :confirmed"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":password_hash": &types.AttributeValueMemberS{Value: updatedUser.PasswordHash},
			":code":          &types.AttributeValueMemberS{Value: updatedUser.InviteCode},
			":confirmed":     &types.AttributeValueMemberBOOL{Value: updatedUser.Confirmed},
		},
	})

	if err != nil {
		return fmt.Errorf("UpdateUser > UpdateItem: %v\n", err)
	}

	return nil
}

func (u userRepository) DeleteUser(ctx context.Context, id idx.UserID) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("DeleteUser > GetClient: %\n", getClientErr)
	}

	_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(UserTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		return fmt.Errorf("DeleteUser > DeleteItem: %v\n", err)
	}

	return nil
}
