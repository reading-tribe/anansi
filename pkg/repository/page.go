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

const PageTableName = "zula_pages"

type PageRepository interface {
	GetPage(ctx context.Context, id idx.PageID) (dbmodel.Page, error)
	CreatePage(ctx context.Context, newPage dbmodel.Page) (dbmodel.Page, error)
	ListPages(ctx context.Context) ([]dbmodel.Page, error)
	ListPagesByTranslationID(ctx context.Context, translationID idx.TranslationID) ([]dbmodel.Page, error)
	UpdatePage(ctx context.Context, updatedPage dbmodel.Page) error
	DeletePage(ctx context.Context, id idx.PageID) error
}

type pageRepository struct{}

func NewPageRepository() PageRepository {
	return pageRepository{}
}

func (p pageRepository) GetPage(ctx context.Context, id idx.PageID) (dbmodel.Page, error) {
	item := dbmodel.Page{}

	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return item, fmt.Errorf("GetPage > GetClient: %v\n", getClientErr)
	}

	data, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(PageTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})

	if err != nil {
		return item, fmt.Errorf("GetPage: %v\n", err)
	}

	if data.Item == nil {
		return item, fmt.Errorf("GetPage: Page not found.\n")
	}

	err = attributevalue.UnmarshalMap(data.Item, &item)
	if err != nil {
		return item, fmt.Errorf("GetPage > UnmarshalMap: %v\n", err)
	}

	return item, nil
}

func (p pageRepository) CreatePage(ctx context.Context, newPage dbmodel.Page) (dbmodel.Page, error) {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return dbmodel.Page{}, fmt.Errorf("CreatePage > GetClient: %v\n", getClientErr)
	}

	id, idxErr := idx.NewPageID()
	if idxErr != nil {
		return dbmodel.Page{}, fmt.Errorf("CreatePage > NewPageID: %\n", getClientErr)
	}

	newPage.ID = id

	data, err := attributevalue.MarshalMap(newPage)
	if err != nil {
		return dbmodel.Page{}, fmt.Errorf("CreatePage > MarshalMap: %v\n", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(PageTableName),
		Item:      data,
	})

	if err != nil {
		return dbmodel.Page{}, fmt.Errorf("CreatePage > PutItem: %v\n", err)
	}

	return newPage, nil
}

func (p pageRepository) ListPages(ctx context.Context) ([]dbmodel.Page, error) {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return nil, fmt.Errorf("ListPages > GetClient: %v\n", getClientErr)
	}

	items := []dbmodel.Page{}

	data, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(PageTableName),
		IndexName:                 aws.String("id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{},
	})
	if err != nil {
		return items, fmt.Errorf("ListPages > Query: %v\n", err)
	}

	err = attributevalue.UnmarshalListOfMaps(data.Items, &items)
	if err != nil {
		return items, fmt.Errorf("ListPages > UnmarshalListOfMaps: %v\n", err)
	}

	return items, nil
}

func (p pageRepository) ListPagesByTranslationID(ctx context.Context, translationID idx.TranslationID) ([]dbmodel.Page, error) {
	items := []dbmodel.Page{}

	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return items, fmt.Errorf("ListPagesByTranslationID > GetClient: %v\n", getClientErr)
	}

	data, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(PageTableName),
		IndexName: aws.String("translation_id-index"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":translation_id": &types.AttributeValueMemberS{
				Value: translationID.String(),
			},
		},
		KeyConditionExpression: aws.String("translation_id = :translation_id"),
	})

	if err != nil {
		return items, fmt.Errorf("ListPagesByTranslationID > Query: %v\n", err)
	}

	err = attributevalue.UnmarshalListOfMaps(data.Items, &items)
	if err != nil {
		return items, fmt.Errorf("ListPagesByTranslationID > UnmarshalListOfMaps: %v\n", err)
	}

	return items, nil
}

func (p pageRepository) UpdatePage(ctx context.Context, updatedPage dbmodel.Page) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("UpdatePage > GetClient: %v\n", getClientErr)
	}

	_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(PageTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: updatedPage.ID.String()},
		},
		UpdateExpression: aws.String("set image_url = :image_url, number = :number, translation_id = :translation_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":image_url":      &types.AttributeValueMemberS{Value: updatedPage.ImageURL},
			":number":         &types.AttributeValueMemberS{Value: strconv.Itoa(updatedPage.Number)},
			":translation_id": &types.AttributeValueMemberS{Value: updatedPage.TranslationID.String()},
		},
	})

	if err != nil {
		return fmt.Errorf("UpdatePage > UpdateItem: %v\n", err)
	}

	return nil
}

func (p pageRepository) DeletePage(ctx context.Context, id idx.PageID) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("DeletePage > GetClient: %\n", getClientErr)
	}

	_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(PageTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		return fmt.Errorf("DeletePage > DeleteItem: %v\n", err)
	}

	return nil
}
