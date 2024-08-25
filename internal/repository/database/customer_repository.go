package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/icrxz/crm-api-core/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type customerRepository struct {
	client *mongo.Client
}

func NewCustomerRepository(client *mongo.Client) domain.CustomerRepository {
	return &customerRepository{
		client: client,
	}
}

func (db *customerRepository) customerCollection(ctx context.Context) *mongo.Collection {
	leadCollection := GetCollection(db.client, "leads")
	return leadCollection
}

func (db *customerRepository) Create(ctx context.Context, customer domain.Customer) (string, error) {
	customerDTO := mapCustomerToCustomerDTO(customer)

	_, err := db.customerCollection(ctx).InsertOne(ctx, customerDTO)

	if err != nil {
		return "", err
	}

	return customer.CustomerID, nil
}

func (db *customerRepository) GetByID(ctx context.Context, customerID string) (*domain.Customer, error) {
	var customerDTO CustomerDTO
	err := db.customerCollection(ctx).FindOne(context.TODO(), bson.D{{"_id", customerID}}).Decode(&customerDTO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no customer found with this id", map[string]any{"customer_id": customerID})
		}
		return nil, err
	}

	customer := mapCustomerDTOToCustomer(customerDTO)

	return &customer, nil
}

func (db *customerRepository) Search(ctx context.Context, filters domain.CustomerFilters) (domain.PagingResult[domain.Customer], error) {

	filter := bson.M{}
	if len(filters.CustomerID) > 0 {
		filter["_id"] = filters.CustomerID
	}
	if len(filters.CustomerType) > 0 {
		filter["document"] = filters.CustomerType
	}
	if len(filters.Document) > 0 {
		filter["lead_type"] = filters.Document
	}

	cursor, err := db.customerCollection(ctx).Find(ctx, filter)

	var results []CustomerDTO
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	customers := mapCustomerDTOsToCustomers(results)

	result := domain.PagingResult[domain.Customer]{
		Result: customers,
		Paging: domain.Paging{
			Total:  len(results),
			Limit:  filters.Limit,
			Offset: filters.Offset,
		},
	}

	return result, nil
}

func (db *customerRepository) Update(ctx context.Context, customer domain.Customer) error {
	customerDTO := mapCustomerToCustomerDTO(customer)

	filter := bson.D{{"_id", customer.CustomerID}}

	result := db.customerCollection(ctx).FindOneAndUpdate(ctx, filter, customerDTO)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (db *customerRepository) Delete(ctx context.Context, customerID string) error {
	if customerID == "" {
		return domain.NewValidationError("customer id is required", map[string]any{"customer_id": customerID})
	}

	_, err := db.customerCollection(ctx).DeleteOne(ctx, bson.D{{"_id", customerID}})
	if err != nil {
		return err
	}

	return nil
}
