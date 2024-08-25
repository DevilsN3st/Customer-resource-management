package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// mockgen -source ticket.go -destination mocks/productMocks.go
type TicketRepository interface {
	Create(ctx context.Context, crmTicket Ticket) (string, error)
	GetByID(ctx context.Context, ticketID string) (*Ticket, error)
	Search(ctx context.Context, filters TicketFilters) (PagingResult[Ticket], error)
	Update(ctx context.Context, crmTicket Ticket) error
}

type CreateTicket struct {
	Ticket  Ticket
	Product Product
}

type Ticket struct {
	TicketID          string
	TenantID          string
	CustomerID        string
	LeadID            string
	OwnerID           string
	OriginChannel     string
	Type              string
	Subject           string
	Priority          TicketPriority
	Transactions      []Transaction
	Comments          []Comment
	Status            TicketStatus
	DueDate           time.Time
	CreatedBy         string
	CreatedAt         time.Time
	UpdatedBy         string
	UpdatedAt         time.Time
	Region            int
	ProductID         string
	ClosedAt          *time.Time
	ExternalReference string
	TargetDate        *time.Time
}

type TicketFilters struct {
	OwnerID    []string
	LeadID     []string
	TenantID   []string
	CustomerID []string
	Status     []string
	Region     []string
	PagingFilter
}

type TicketUpdate struct {
	Status     *TicketStatus
	LeadID     *string
	OwnerID    *string
	TargetDate *time.Time
	ClosedAt   *time.Time
	UpdatedBy  string
}

type TicketStatus string

const (
	NEW           TicketStatus = "New"
	CUSTOMER_INFO TicketStatus = "CustomerInfo"
	WAITING_LEAD  TicketStatus = "WaitingLead"
	ONGOING       TicketStatus = "Ongoing"
	REPORT        TicketStatus = "Report"
	PAYMENT       TicketStatus = "Payment"
	RECEIPT       TicketStatus = "Receipt"
	CLOSED        TicketStatus = "Closed"
	CANCELED      TicketStatus = "Canceled"
)

type TicketPriority string

const (
	LOW    TicketPriority = "Low"
	MEDIUM TicketPriority = "Medium"
	HIGH   TicketPriority = "High"
)

func NewTicket(
	tenantID string,
	customerID string,
	originChannel string,
	ticketType string,
	subject string,
	dueDate time.Time,
	author string,
	externalReference string,
) (Ticket, error) {
	now := time.Now().UTC()

	ticketID, err := uuid.NewUUID()
	if err != nil {
		return Ticket{}, err
	}

	return Ticket{
		TicketID:          ticketID.String(),
		TenantID:          tenantID,
		CustomerID:        customerID,
		OriginChannel:     originChannel,
		Type:              ticketType,
		Subject:           subject,
		Priority:          MEDIUM,
		Status:            NEW,
		DueDate:           dueDate,
		CreatedAt:         now,
		CreatedBy:         author,
		UpdatedAt:         now,
		UpdatedBy:         author,
		ExternalReference: externalReference,
	}, nil
}

func (c *Ticket) MergeUpdate(updateTicket TicketUpdate) {
	c.UpdatedAt = time.Now().UTC()
	c.UpdatedBy = updateTicket.UpdatedBy

	if updateTicket.Status != nil {
		c.Status = *updateTicket.Status
	}

	if updateTicket.OwnerID != nil {
		c.OwnerID = *updateTicket.OwnerID
	}

	if updateTicket.LeadID != nil {
		c.LeadID = *updateTicket.LeadID
	}

	if updateTicket.TargetDate != nil {
		c.TargetDate = updateTicket.TargetDate
	}

	if updateTicket.ClosedAt != nil {
		c.ClosedAt = updateTicket.ClosedAt
	}
}
