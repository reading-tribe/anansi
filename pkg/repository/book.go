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

const BookTableName = "zula_books"

type BookRepository interface {
	GetBook(ctx context.Context, id idx.BookID) (dbmodel.Book, error)
	CreateBook(ctx context.Context, newBook dbmodel.Book) (dbmodel.Book, error)
	ListBooks(ctx context.Context) ([]dbmodel.Book, error)
	UpdateBook(ctx context.Context, updatedBook dbmodel.Book) error
	DeleteBook(ctx context.Context, id idx.BookID) error
}

type bookRepository struct{}

func NewBookRepository() BookRepository {
	return bookRepository{}
}

func (b bookRepository) GetBook(ctx context.Context, id idx.BookID) (dbmodel.Book, error) {
	item := dbmodel.Book{}

	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return item, fmt.Errorf("GetBook > GetClient: %v\n", getClientErr)
	}

	data, err := client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(BookTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})

	if err != nil {
		return item, fmt.Errorf("GetBook: %v\n", err)
	}

	if data.Item == nil {
		return item, fmt.Errorf("GetBook: Book not found.\n")
	}

	err = attributevalue.UnmarshalMap(data.Item, &item)
	if err != nil {
		return item, fmt.Errorf("GetBook > UnmarshalMap: %v\n", err)
	}

	return item, nil
}

func (b bookRepository) CreateBook(ctx context.Context, newBook dbmodel.Book) (dbmodel.Book, error) {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return dbmodel.Book{}, fmt.Errorf("CreateBook > GetClient: %v\n", getClientErr)
	}

	id, idxErr := idx.NewBookID()
	if idxErr != nil {
		return dbmodel.Book{}, fmt.Errorf("CreateBook > NewBookID: %\n", getClientErr)
	}

	newBook.ID = id

	data, err := attributevalue.MarshalMap(newBook)
	if err != nil {
		return dbmodel.Book{}, fmt.Errorf("CreateBook > MarshalMap: %v\n", err)
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(BookTableName),
		Item:      data,
	})

	if err != nil {
		return dbmodel.Book{}, fmt.Errorf("CreateBook > PutItem: %v\n", err)
	}

	return newBook, nil
}

func (b bookRepository) ListBooks(ctx context.Context) ([]dbmodel.Book, error) {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return nil, fmt.Errorf("ListBooks > GetClient: %v\n", getClientErr)
	}

	items := []dbmodel.Book{}

	data, err := client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(BookTableName),
	})
	if err != nil {
		return items, fmt.Errorf("ListBooks > Scan: %v\n", err)
	}

	err = attributevalue.UnmarshalListOfMaps(data.Items, &items)
	if err != nil {
		return items, fmt.Errorf("ListBooks > UnmarshalListOfMaps: %v\n", err)
	}

	return items, nil
}

func (b bookRepository) UpdateBook(ctx context.Context, updatedBook dbmodel.Book) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("UpdateBook > GetClient: %v\n", getClientErr)
	}

	_, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(BookTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: updatedBook.ID.String()},
		},
		UpdateExpression: aws.String("set internal_title = :internal_title, authors = :authors"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":internal_title": &types.AttributeValueMemberS{Value: updatedBook.InternalTitle},
			":authors":        &types.AttributeValueMemberS{Value: updatedBook.Authors},
		},
	})

	if err != nil {
		return fmt.Errorf("UpdateBook > UpdateItem: %v\n", err)
	}

	return nil
}

func (b bookRepository) DeleteBook(ctx context.Context, id idx.BookID) error {
	client, getClientErr := dynamodbx.GetClient(ctx)
	if getClientErr != nil {
		return fmt.Errorf("DeleteBook > GetClient: %v\n", getClientErr)
	}

	_, err := client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(BookTableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id.String()},
		},
	})
	if err != nil {
		return fmt.Errorf("DeleteBook > DeleteItem: %v\n", err)
	}

	return nil
}
