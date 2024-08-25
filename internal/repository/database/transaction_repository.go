package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/icrxz/crm-api-core/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type transactionRepository struct {
	client *mongo.Client
}

func NewTransactionRepository(client *mongo.Client) domain.TransactionRepository {
	return &transactionRepository{
		client: client,
	}
}

func (r *transactionRepository) transactionCollection(ctx context.Context) *mongo.Collection {
	userCollection := GetCollection(r.client, "user")
	return userCollection
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, transaction domain.Transaction) (string, error) {
	transactionDTO := mapTransactionToTransactionDTO(transaction)

	_, err := r.transactionCollection(ctx).InsertOne(ctx, transactionDTO)

	if err != nil {
		return "", err
	}

	return transaction.TransactionID, nil
}

func (r *transactionRepository) GetTransaction(ctx context.Context, transactionID string) (domain.Transaction, error) {
	if transactionID == "" {
		return domain.Transaction{}, domain.NewValidationError("transaction_id is required", nil)
	}

	var transactionDTO TransactionDTO
	err := r.transactionCollection(ctx).FindOne(context.TODO(), bson.D{{"_id", transactionID}}).Decode(&transactionDTO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Transaction{}, domain.NewNotFoundError("no transaction found with this id", map[string]any{"transaction_id": transactionID})
		}
		return domain.Transaction{}, err
	}

	return mapTransactionDTOToTransaction(transactionDTO), nil
}

func (r *transactionRepository) UpdateTransaction(ctx context.Context, transaction domain.Transaction) error {
	transactionDTO := mapTransactionToTransactionDTO(transaction)
	filter := bson.M{
		"_id": transaction.TransactionID,
	}
	result := r.transactionCollection(ctx).FindOneAndUpdate(ctx, filter, transactionDTO)

	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (r *transactionRepository) SearchTransactions(ctx context.Context, filters domain.TransactionFilters) ([]domain.Transaction, error) {

	filter := bson.M{}
	if len(filters.Status) > 0 {
		filter["status"] = filters.Status
	}
	if len(filters.Status) > 0 {
		filter["ticket_ids"] = filters.TicketIDs
	}
	if len(filters.Status) > 0 {
		filter["type"] = filters.Types
	}
	var transactionResults []TransactionDTO
	err := r.transactionCollection(ctx).FindOne(ctx, filter).Decode(&transactionResults)
	if err != nil {
		return nil, err
	}

	transactions := mapTransactionDTOsToTransactions(transactionResults)

	return transactions, nil
}
