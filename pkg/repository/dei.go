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

const DiversityAndInclusionCatalogueTableName = "zula_diversity_and_inclusion_catalogue"

type DiversityAndInclusionCatalogueRepository interface {
	GetDiversityAndInclusionCatalogueItem(ctx context.Context, id idx.DiversityAndInclusionID) (dbmodel.DiversityAndInclusionCatalogue, error)
	CreateDiversityAndInclusionCatalogueItem(ctx context.Context, newDiversityAndInclusionCatalogue dbmodel.DiversityAndInclusionCatalogue) error
	ListDiversityAndInclusionCatalogueItems(ctx context.Context) ([]dbmodel.DiversityAndInclusionCatalogue, error)
	UpdateDiversityAndInclusionCatalogueItem(ctx context.Context, updatedDiversityAndInclusionCatalogue dbmodel.DiversityAndInclusionCatalogue) error
	DeleteDiversityAndInclusionCatalogueItem(ctx context.Context, id idx.DiversityAndInclusionID) error
}

type diversityAndInclusionCatalogueRepository struct{}

func NewDiversityAndInclusionCatalogueRepository() DiversityAndInclusionCatalogueRepository {
	return diversityAndInclusionCatalogueRepository{}
}

func (d diversityAndInclusionCatalogueRepository) GetDiversityAndInclusionCatalogueItem(ctx context.Context, id idx.DiversityAndInclusionID) (dbmodel.DiversityAndInclusionCatalogue, error) {
	item := dbmodel.DiversityAndInclusionCatalogue{}

	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return item, fmt.Errorf("GetDiversityAndInclusionCatalogueItem > GetClient: %v\n", getClientErr)
	}

	data, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(DiversityAndInclusionCatalogueTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})

	if err != nil {
		return item, fmt.Errorf("GetDiversityAndInclusionCatalogueItem: %v\n", err)
	}

	if data.Item == nil {
		return item, fmt.Errorf("GetDiversityAndInclusionCatalogueItem: DiversityAndInclusionCatalogue not found.\n")
	}

	err = attributevalue.UnmarshalMap(data.Item, &item)
	if err != nil {
		return item, fmt.Errorf("GetDiversityAndInclusionCatalogueItem > UnmarshalMap: %v\n", err)
	}

	return item, nil
}

func (d diversityAndInclusionCatalogueRepository) CreateDiversityAndInclusionCatalogueItem(ctx context.Context, newDiversityAndInclusionCatalogue dbmodel.DiversityAndInclusionCatalogue) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("CreateDiversityAndInclusionCatalogueItem > GetClient: %v\n", getClientErr)
	}

	id, idxErr := idx.NewDiversityAndInclusionCatalogueID()
	if idxErr != nil {
		return fmt.Errorf("CreateDiversityAndInclusionCatalogueItem > NewDiversityAndInclusionCatalogueID: %\n", getClientErr)
	}

	newDiversityAndInclusionCatalogue.ID = id

	data, err := attributevalue.MarshalMap(newDiversityAndInclusionCatalogue)
	if err != nil {
		return fmt.Errorf("CreateDiversityAndInclusionCatalogueItem > MarshalMap: %v\n", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(DiversityAndInclusionCatalogueTableName),
		Item:      data,
	})

	if err != nil {
		return fmt.Errorf("CreateDiversityAndInclusionCatalogueItem > PutItem: %v\n", err)
	}

	return nil
}

func (d diversityAndInclusionCatalogueRepository) ListDiversityAndInclusionCatalogueItems(ctx context.Context) ([]dbmodel.DiversityAndInclusionCatalogue, error) {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return nil, fmt.Errorf("ListDiversityAndInclusionCatalogueItems > GetClient: %v\n", getClientErr)
	}

	items := []dbmodel.DiversityAndInclusionCatalogue{}

	data, err := client.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(DiversityAndInclusionCatalogueTableName),
		IndexName:                 aws.String("id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{},
	})
	if err != nil {
		return items, fmt.Errorf("ListDiversityAndInclusionCatalogueItems > Query: %v\n", err)
	}

	err = attributevalue.UnmarshalListOfMaps(data.Items, &items)
	if err != nil {
		return items, fmt.Errorf("ListDiversityAndInclusionCatalogueItems > UnmarshalListOfMaps: %v\n", err)
	}

	return items, nil
}

func (d diversityAndInclusionCatalogueRepository) UpdateDiversityAndInclusionCatalogueItem(ctx context.Context, updatedDiversityAndInclusionCatalogue dbmodel.DiversityAndInclusionCatalogue) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("UpdateDiversityAndInclusionCatalogueItem > GetClient: %v\n", getClientErr)
	}

	_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(DiversityAndInclusionCatalogueTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: updatedDiversityAndInclusionCatalogue.ID.String()},
		},
		UpdateExpression: aws.String("set key = :key, value = :value"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":key":   &types.AttributeValueMemberS{Value: updatedDiversityAndInclusionCatalogue.Key},
			":value": &types.AttributeValueMemberS{Value: updatedDiversityAndInclusionCatalogue.Value},
		},
	})

	if err != nil {
		return fmt.Errorf("UpdateDiversityAndInclusionCatalogueItem > UpdateItem: %v\n", err)
	}

	return nil
}

func (d diversityAndInclusionCatalogueRepository) DeleteDiversityAndInclusionCatalogueItem(ctx context.Context, id idx.DiversityAndInclusionID) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("DeleteDiversityAndInclusionCatalogueItem > GetClient: %v\n", getClientErr)
	}

	_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(DiversityAndInclusionCatalogueTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		return fmt.Errorf("DeleteDiversityAndInclusionCatalogueItem > DeleteItem: %v\n", err)
	}

	return nil
}
