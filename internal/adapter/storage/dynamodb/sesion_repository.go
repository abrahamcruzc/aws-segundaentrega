package dynamodb

import (
	"context"
	"fmt"

	"github.com/abrahamcruzc/aws-segundaentrega/internal/domain"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type SesionRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewSesionRepository(client *dynamodb.Client, tableName string) *SesionRepository {
	return &SesionRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *SesionRepository) Create(ctx context.Context, sesion *domain.Sesion) error {
	item, err := attributevalue.MarshalMap(sesion)
	if err != nil {
		return fmt.Errorf("error al serializar sesión: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("error al crear sesión en DynamoDB: %w", err)
	}

	return nil
}

func (r *SesionRepository) GetBySessionString(ctx context.Context, sessionString string) (*domain.Sesion, error) {
	output, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
		FilterExpression: aws.String("sessionString = :ss"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":ss": &types.AttributeValueMemberS{Value: sessionString},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error al buscar sesión: %w", err)
	}

	if len(output.Items) == 0 {
		return nil, nil
	}

	var sesion domain.Sesion
	if err := attributevalue.UnmarshalMap(output.Items[0], &sesion); err != nil {
		return nil, fmt.Errorf("error al deserializar sesión: %w", err)
	}

	return &sesion, nil
}


func (r *SesionRepository) Deactivate(ctx context.Context, sessionString string) error {
	sesion, err := r.GetBySessionString(ctx, sessionString)
	if err != nil {
		return err
	}
	if sesion == nil {
		return fmt.Errorf("sesión no encontrada")
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(r.tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: sesion.ID},
		},
		UpdateExpression: aws.String("SET active = :a"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":a": &types.AttributeValueMemberBOOL{Value: false},
		},
	})
	if err != nil {
		return fmt.Errorf("error al desactivar sesión: %w", err)
	}

	return nil
}

func (r *SesionRepository) CreateTable(ctx context.Context) error {
	_, err := r.client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(r.tableName),
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType: types.KeyTypeHash,
			},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		return nil
	}
	return nil
}
