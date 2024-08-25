package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/icrxz/crm-api-core/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type leadRepository struct {
	client *mongo.Client
}

func NewLeadRepository(client *mongo.Client) domain.LeadRepository {
	return &leadRepository{
		client: client,
	}
}

func (db *leadRepository) leadCollection(ctx context.Context) *mongo.Collection {
	leadCollection := GetCollection(db.client, "leads")
	return leadCollection
}

func (db *leadRepository) Create(ctx context.Context, lead domain.Lead) (string, error) {
	leadDTO := mapLeadToLeadDTO(lead)

	_, err := db.leadCollection(ctx).InsertOne(ctx, leadDTO)
	if err != nil {
		return "", err
	}

	return lead.LeadID, nil
}

func (db *leadRepository) Delete(ctx context.Context, leadID string) error {
	if leadID == "" {
		return domain.NewValidationError("leadID is required", map[string]any{"lead_id": leadID})
	}

	_, err := db.leadCollection(ctx).DeleteOne(ctx, bson.D{{"_id", leadID}})
	if err != nil {
		return err
	}

	return nil
}

func (db *leadRepository) GetByID(ctx context.Context, leadID string) (*domain.Lead, error) {
	var leadDTO LeadDTO
	err := db.leadCollection(ctx).FindOne(context.TODO(), bson.D{{"_id", leadID}}).Decode(&leadDTO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no lead found with this id", map[string]any{"lead_id": leadID})
		}
		return nil, err
	}

	lead := mapLeadDTOToLead(leadDTO)

	return &lead, nil
}

func (db *leadRepository) Search(ctx context.Context, filters domain.LeadFilters) (domain.PagingResult[domain.Lead], error) {

	filter := bson.M{}
	if len(filters.LeadID) > 0 {
		filter["_id"] = filters.LeadID
	}
	if len(filters.Document) > 0 {
		filter["document"] = filters.Document
	}
	if len(filters.LeadType) > 0 {
		filter["lead_type"] = filters.LeadType
	}
	if len(filters.State) > 0 {
		filter["shipping_state"] = filters.State
	}

	cursor, err := db.leadCollection(ctx).Find(ctx, filter)

	var results []LeadDTO
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	leads := mapLeadDTOsToLeads(results)

	result := domain.PagingResult[domain.Lead]{
		Result: leads,
		Paging: domain.Paging{
			Total:  len(results),
			Limit:  filters.Limit,
			Offset: filters.Offset,
		},
	}

	return result, nil
}

func (db *leadRepository) Update(ctx context.Context, lead domain.Lead) error {
	leadDTO := mapLeadToLeadDTO(lead)

	filter := bson.M{
		"_id": leadDTO.LeadID,
	}

	result := db.leadCollection(ctx).FindOneAndUpdate(ctx, filter, leadDTO)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func (db *leadRepository) CreateBatch(ctx context.Context, leads []domain.Lead) ([]string, error) {
	chunks := db.createChunks(leads, 100)

	insertedIDs := make([]string, 0, len(leads))
	for _, chunk := range chunks {
		leadDTOs := mapLeadsToLeadDTOs(chunk)

		_, err := db.leadCollection(ctx).InsertMany(context.TODO(), leadDTOs)

		if err != nil {
			return nil, err
		}

		for _, lead := range leadDTOs {
			insertedIDs = append(insertedIDs, lead.(*LeadDTO).LeadID)
		}
	}

	return insertedIDs, nil
}

func (db *leadRepository) createChunks(slice []domain.Lead, size int) [][]domain.Lead {
	var chunks [][]domain.Lead
	for i := 0; i < len(slice); i += size {
		end := i + size

		if end > len(slice) {
			end = len(slice)
		}

		chunks = append(chunks, slice[i:end])
	}

	return chunks
}
