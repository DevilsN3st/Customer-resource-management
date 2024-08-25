package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/icrxz/crm-api-core/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ticketRepository struct {
	client *mongo.Client
}

func NewTicketRepository(client *mongo.Client) domain.TicketRepository {
	return &ticketRepository{
		client: client,
	}
}

func (r *ticketRepository) ticketCollection(ctx context.Context) *mongo.Collection {
	userCollection := GetCollection(r.client, "user")
	return userCollection
}

func (r *ticketRepository) Create(ctx context.Context, crmTicket domain.Ticket) (string, error) {
	crmTicketDTO := mapTicketToTicketDTO(crmTicket)

	_, err := r.ticketCollection(ctx).InsertOne(ctx, crmTicketDTO)

	if err != nil {
		return "", err
	}

	return crmTicket.TicketID, nil
}

func (r *ticketRepository) GetByID(ctx context.Context, ticketID string) (*domain.Ticket, error) {
	if ticketID == "" {
		return nil, domain.NewValidationError("ticketID is required", nil)
	}

	var crmTicketDTO TicketDTO
	err := r.ticketCollection(ctx).FindOne(context.TODO(), bson.D{{"_id", ticketID}}).Decode(&crmTicketDTO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no ticket found with this id", map[string]any{"ticket_id": ticketID})
		}
		return nil, err
	}

	crmTicket := mapTicketDTOToTicket(crmTicketDTO)
	return &crmTicket, nil
}

func (r *ticketRepository) Search(ctx context.Context, filters domain.TicketFilters) (domain.PagingResult[domain.Ticket], error) {

	filter := bson.M{}

	if len(filters.TenantID) > 0 {
		filter["tenant_id"] = filters.TenantID
	}
	if len(filters.OwnerID) > 0 {
		filter["owner_id"] = filters.OwnerID
	}
	if len(filters.CustomerID) > 0 {
		filter["customer_id"] = filters.CustomerID
	}
	if len(filters.LeadID) > 0 {
		filter["lead_id"] = filters.LeadID
	}
	if len(filters.Status) > 0 {
		filter["status"] = filters.Status
	}
	if len(filters.Region) > 0 {
		filter["region"] = filters.Region
	}

	var crmTicketsDTO []TicketDTO

	err := r.ticketCollection(ctx).FindOne(ctx, filter).Decode(&crmTicketsDTO)

	if err != nil {
		return domain.PagingResult[domain.Ticket]{}, err
	}
	crmTickets := mapTicketDTOsToTickets(crmTicketsDTO)

	result := domain.PagingResult[domain.Ticket]{
		Result: crmTickets,
		Paging: domain.Paging{
			Total:  len(crmTicketsDTO),
			Limit:  filters.Limit,
			Offset: filters.Offset,
		},
	}

	return result, nil
}

func (r *ticketRepository) Update(ctx context.Context, crmTicket domain.Ticket) error {
	crmTicketDTO := mapTicketToTicketDTO(crmTicket)

	filter := bson.M{
		"_id": crmTicket.TicketID,
	}

	result := r.ticketCollection(ctx).FindOneAndUpdate(ctx, filter, crmTicketDTO)

	if result.Err() != nil {
		return result.Err()
	}

	return nil
}
