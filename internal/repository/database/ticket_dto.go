package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type TicketDTO struct {
	TicketID          string     `db:"ticket_id"`
	TenantID          string     `db:"tenant_id"`
	CustomerID        string     `db:"customer_id"`
	LeadID            *string    `db:"lead_id"`
	OwnerID           *string    `db:"owner_id"`
	OriginChannel     string     `db:"origin"`
	Type              string     `db:"type"`
	Subject           string     `db:"subject"`
	Priority          string     `db:"priority"`
	Status            string     `db:"status"`
	DueDate           time.Time  `db:"due_date"`
	CreatedBy         string     `db:"created_by"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedBy         string     `db:"updated_by"`
	UpdatedAt         time.Time  `db:"updated_at"`
	ExternalReference string     `db:"external_reference"`
	ProductID         string     `db:"product_id"`
	Region            int        `db:"region"`
	ClosedAt          *time.Time `db:"closed_at"`
	TargetDate        *time.Time `db:"target_date"`
}

func mapTicketToTicketDTO(crmTicket domain.Ticket) TicketDTO {
	var leadID *string
	if crmTicket.LeadID != "" {
		leadID = &crmTicket.LeadID
	}

	var ownerID *string
	if crmTicket.OwnerID != "" {
		ownerID = &crmTicket.OwnerID
	}

	return TicketDTO{
		TicketID:          crmTicket.TicketID,
		TenantID:          crmTicket.TenantID,
		CustomerID:        crmTicket.CustomerID,
		LeadID:            leadID,
		OwnerID:           ownerID,
		OriginChannel:     crmTicket.OriginChannel,
		Type:              crmTicket.Type,
		Subject:           crmTicket.Subject,
		Priority:          string(crmTicket.Priority),
		Status:            string(crmTicket.Status),
		DueDate:           crmTicket.DueDate,
		CreatedBy:         crmTicket.CreatedBy,
		CreatedAt:         crmTicket.CreatedAt,
		UpdatedBy:         crmTicket.UpdatedBy,
		UpdatedAt:         crmTicket.UpdatedAt,
		ExternalReference: crmTicket.ExternalReference,
		Region:            crmTicket.Region,
		ProductID:         crmTicket.ProductID,
		TargetDate:        crmTicket.TargetDate,
		ClosedAt:          crmTicket.ClosedAt,
	}
}

func mapTicketDTOToTicket(crmTicketDTO TicketDTO) domain.Ticket {
	var leadID string
	if crmTicketDTO.LeadID != nil {
		leadID = *crmTicketDTO.LeadID
	}

	var ownerID string
	if crmTicketDTO.OwnerID != nil {
		ownerID = *crmTicketDTO.OwnerID
	}

	return domain.Ticket{
		TicketID:          crmTicketDTO.TicketID,
		TenantID:          crmTicketDTO.TenantID,
		CustomerID:        crmTicketDTO.CustomerID,
		LeadID:            leadID,
		OwnerID:           ownerID,
		OriginChannel:     crmTicketDTO.OriginChannel,
		Type:              crmTicketDTO.Type,
		Subject:           crmTicketDTO.Subject,
		Priority:          domain.TicketPriority(crmTicketDTO.Priority),
		Status:            domain.TicketStatus(crmTicketDTO.Status),
		DueDate:           crmTicketDTO.DueDate,
		CreatedBy:         crmTicketDTO.CreatedBy,
		CreatedAt:         crmTicketDTO.CreatedAt,
		UpdatedBy:         crmTicketDTO.UpdatedBy,
		UpdatedAt:         crmTicketDTO.UpdatedAt,
		ExternalReference: crmTicketDTO.ExternalReference,
		Region:            crmTicketDTO.Region,
		ProductID:         crmTicketDTO.ProductID,
		ClosedAt:          crmTicketDTO.ClosedAt,
		TargetDate:        crmTicketDTO.TargetDate,
	}
}

func mapTicketDTOsToTickets(crmTicketDTOs []TicketDTO) []domain.Ticket {
	crmTickets := make([]domain.Ticket, 0, len(crmTicketDTOs))
	for _, crmTicketDTO := range crmTicketDTOs {
		crmTicket := mapTicketDTOToTicket(crmTicketDTO)
		crmTickets = append(crmTickets, crmTicket)
	}

	return crmTickets
}
