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
)

const TranslationTableName = "zula_translations"

type TranslationRepository interface {
	GetTranslation(ctx context.Context, id idx.TranslationID) (dbmodel.Translation, error)
	CreateTranslation(ctx context.Context, newTranslation dbmodel.Translation) (dbmodel.Translation, error)
	ListTranslations(ctx context.Context) ([]dbmodel.Translation, error)
	ListTranslationsByBookID(ctx context.Context, bookID idx.BookID) ([]dbmodel.Translation, error)
	UpdateTranslation(ctx context.Context, updatedTranslation dbmodel.Translation) error
	DeleteTranslation(ctx context.Context, id idx.TranslationID) error
}

type translationRepository struct{}

func NewTranslationRepository() TranslationRepository {
	return translationRepository{}
}

func (t translationRepository) GetTranslation(ctx context.Context, id idx.TranslationID) (dbmodel.Translation, error) {
	item := dbmodel.Translation{}

	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return item, fmt.Errorf("GetTranslation > GetClient: %v\n", getClientErr)
	}

	data, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TranslationTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})

	if err != nil {
		return item, fmt.Errorf("GetTranslation: %v\n", err)
	}

	if data.Item == nil {
		return item, fmt.Errorf("GetTranslation: Translation not found.\n")
	}

	err = attributevalue.UnmarshalMap(data.Item, &item)
	if err != nil {
		return item, fmt.Errorf("GetTranslation > UnmarshalMap: %v\n", err)
	}

	return item, nil
}

func (t translationRepository) CreateTranslation(ctx context.Context, newTranslation dbmodel.Translation) (dbmodel.Translation, error) {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return dbmodel.Translation{}, fmt.Errorf("CreateTranslation > GetClient: %v\n", getClientErr)
	}

	id, idxErr := idx.NewTranslationID()
	if idxErr != nil {
		return dbmodel.Translation{}, fmt.Errorf("CreateTranslation > NewPageID: %\n", getClientErr)
	}

	newTranslation.ID = id

	data, err := attributevalue.MarshalMap(newTranslation)
	if err != nil {
		return dbmodel.Translation{}, fmt.Errorf("CreateTranslation > MarshalMap: %v\n", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(TranslationTableName),
		Item:      data,
	})

	if err != nil {
		return dbmodel.Translation{}, fmt.Errorf("CreateTranslation > PutItem: %v\n", err)
	}

	return newTranslation, nil
}

func (t translationRepository) ListTranslations(ctx context.Context) ([]dbmodel.Translation, error) {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return nil, fmt.Errorf("ListTranslations > GetClient: %v\n", getClientErr)
	}

	items := []dbmodel.Translation{}

	data, err := client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(TranslationTableName),
	})
	if err != nil {
		return items, fmt.Errorf("ListTranslations > Scan: %v\n", err)
	}

	err = attributevalue.UnmarshalListOfMaps(data.Items, &items)
	if err != nil {
		return items, fmt.Errorf("ListTranslations > UnmarshalListOfMaps: %v\n", err)
	}

	return items, nil
}

func (t translationRepository) UpdateTranslation(ctx context.Context, updatedTranslation dbmodel.Translation) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("UpdateTranslation > GetClient: %v\n", getClientErr)
	}

	_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(TranslationTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: updatedTranslation.ID.String()},
		},
		UpdateExpression: aws.String("set book_id = :book_id, localised_title = :localised_title, lang = :lang"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":book_id":         &types.AttributeValueMemberS{Value: updatedTranslation.BookID.String()},
			":localised_title": &types.AttributeValueMemberS{Value: updatedTranslation.LocalisedTitle},
			":lang":            &types.AttributeValueMemberS{Value: string(updatedTranslation.Language)},
		},
	})

	if err != nil {
		return fmt.Errorf("UpdateTranslation > UpdateItem: %v\n", err)
	}

	return nil
}

func (t translationRepository) ListTranslationsByBookID(ctx context.Context, bookID idx.BookID) ([]dbmodel.Translation, error) {
	items := []dbmodel.Translation{}

	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return items, fmt.Errorf("ListTranslationsByBookID > GetClient: %v\n", getClientErr)
	}

	data, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(TranslationTableName),
		IndexName: aws.String("book_id-index"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":book_id": &types.AttributeValueMemberS{
				Value: bookID.String(),
			},
		},
		KeyConditionExpression: aws.String("book_id = :book_id"),
	})

	if err != nil {
		return items, fmt.Errorf("ListTranslationsByBookID > Query: %v\n", err)
	}

	err = attributevalue.UnmarshalListOfMaps(data.Items, &items)
	if err != nil {
		return items, fmt.Errorf("ListTranslationsByBookID > UnmarshalListOfMaps: %v\n", err)
	}

	return items, nil
}

func (t translationRepository) DeleteTranslation(ctx context.Context, id idx.TranslationID) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("DeleteTranslation > GetClient: %v\n", getClientErr)
	}

	_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(TranslationTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		return fmt.Errorf("DeleteTranslation > DeleteItem: %v\n", err)
	}

	return nil
}
