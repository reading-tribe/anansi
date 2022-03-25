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

const BookDiversityTagsTableName = "zula_book_diversity_tags"

type BookDiversityTagsRepository interface {
	GetBookDiversityTag(ctx context.Context, id string) (dbmodel.BookDiversityTags, error)
	CreateBookDiversityTag(ctx context.Context, newBookDiversityTags dbmodel.BookDiversityTags) error
	ListBookDiversityTags(ctx context.Context) ([]dbmodel.BookDiversityTags, error)
	DeleteBookDiversityTag(ctx context.Context, id string) error
}

type bookDiversityTagsRepository struct{}

func NewBookDiversityTagsRepository() BookDiversityTagsRepository {
	return bookDiversityTagsRepository{}
}

func (b bookDiversityTagsRepository) GetBookDiversityTag(ctx context.Context, id string) (dbmodel.BookDiversityTags, error) {
	item := dbmodel.BookDiversityTags{}

	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return item, fmt.Errorf("GetBookDiversityTag > GetClient: %\n", getClientErr)
	}

	data, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(BookDiversityTagsTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return item, fmt.Errorf("GetBookDiversityTag: %v\n", err)
	}

	if data.Item == nil {
		return item, fmt.Errorf("GetBookDiversityTag: BookDiversityTags not found.\n")
	}

	err = attributevalue.UnmarshalMap(data.Item, &item)
	if err != nil {
		return item, fmt.Errorf("GetBookDiversityTag > UnmarshalMap: %v\n", err)
	}

	return item, nil
}

func (b bookDiversityTagsRepository) CreateBookDiversityTag(ctx context.Context, newBookDiversityTags dbmodel.BookDiversityTags) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("CreateBookDiversityTag > GetClient: %\n", getClientErr)
	}

	data, err := attributevalue.MarshalMap(newBookDiversityTags)
	if err != nil {
		return fmt.Errorf("CreateBookDiversityTag > MarshalMap: %v\n", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(BookDiversityTagsTableName),
		Item:      data,
	})

	if err != nil {
		return fmt.Errorf("CreateBookDiversityTag > PutItem: %v\n", err)
	}

	return nil
}

func (b bookDiversityTagsRepository) ListBookDiversityTags(ctx context.Context) ([]dbmodel.BookDiversityTags, error) {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return nil, fmt.Errorf("ListBookDiversityTags > GetClient: %\n", getClientErr)
	}

	items := []dbmodel.BookDiversityTags{}

	data, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(BookDiversityTagsTableName),
		IndexName:                 aws.String("id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{},
	})
	if err != nil {
		return items, fmt.Errorf("ListBookDiversityTags > Query: %v\n", err)
	}

	err = attributevalue.UnmarshalListOfMaps(data.Items, &items)
	if err != nil {
		return items, fmt.Errorf("ListBookDiversityTags > UnmarshalListOfMaps: %v\n", err)
	}

	return items, nil
}

func (b bookDiversityTagsRepository) DeleteBookDiversityTag(ctx context.Context, id string) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("DeleteBookDiversityTag > GetClient: %\n", getClientErr)
	}

	_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(BookDiversityTagsTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return fmt.Errorf("DeleteBookDiversityTag > DeleteItem: %v\n", err)
	}

	return nil
}
