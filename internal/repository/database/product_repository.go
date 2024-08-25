package database

import (
	"context"
	"database/sql"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type productRepository struct {
	client *mongo.Client
}

func NewProductRepository(client *mongo.Client) domain.ProductRepository {
	return &productRepository{
		client: client,
	}
}

func (db *productRepository) productCollection(ctx context.Context) *mongo.Collection {
	productCollection := GetCollection(db.client, "user")
	return productCollection
}

func (r *productRepository) CreateProduct(ctx context.Context, product domain.Product) (string, error) {
	productDTO := mapProductToProductDTO(product)

	_, err := r.productCollection(ctx).InsertOne(ctx, productDTO)

	if err != nil {
		return "", err
	}

	return product.ProductID, nil
}

func (r *productRepository) GetProductByID(ctx context.Context, productID string) (*domain.Product, error) {
	if productID == "" {
		return nil, domain.NewValidationError("product_id is required", nil)
	}

	var productDTO ProductDTO
	err := r.productCollection(ctx).FindOne(context.TODO(), bson.D{{"_id", productID}}).Decode(&productDTO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no product found with this id", map[string]any{"product_id": productID})
		}
		return nil, err
	}

	product := mapProductDTOToProduct(productDTO)
	return &product, nil
}
