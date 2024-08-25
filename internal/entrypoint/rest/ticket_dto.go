package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateTicketDTO struct {
	TenantID           string    `json:"tenant_id" validate:"required"`
	CustomerID         string    `json:"customer_id" validate:"required"`
	OriginChannel      string    `json:"origin_channel" validate:"required"`
	TicketType         string    `json:"ticket_type" validate:"required"`
	Subject            string    `json:"subject" validate:"required"`
	DueDate            time.Time `json:"due_date" validate:"required"`
	CreatedBy          string    `json:"created_by" validate:"required"`
	ExternalReference  string    `json:"external_reference"`
	ProductName        string    `json:"product_name"`
	Brand              string    `json:"brand" validate:"required"`
	Model              string    `json:"model" validate:"required"`
	ProductDescription string    `json:"product_description"`
	Value              float64   `json:"value"`
	SerialNumber       string    `json:"serial_number"`
}

type TicketDTO struct {
	TicketID          string                `json:"ticket_id"`
	TenantID          string                `json:"tenant_id"`
	CustomerID        string                `json:"customer_id"`
	LeadID            string                `json:"lead_id"`
	OwnerID           string                `json:"owner_id"`
	OriginChannel     string                `json:"origin_channel"`
	Type              string                `json:"type"`
	Subject           string                `json:"subject"`
	Priority          domain.TicketPriority `json:"priority"`
	Status            domain.TicketStatus   `json:"status"`
	DueDate           time.Time             `json:"due_date"`
	CreatedBy         string                `json:"created_by"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedBy         string                `json:"updated_by"`
	UpdatedAt         time.Time             `json:"updated_at"`
	Region            int                   `json:"region"`
	ExternalReference string                `json:"external_reference"`
	ProductID         string                `json:"product_id"`
	ClosedAt          *time.Time            `json:"closed_at"`
	TargetDate        *time.Time            `json:"target_date"`
}

type UpdateTicketDTO struct {
	TargetDate *time.Time `json:"target_date"`
	UpdatedBy  string     `json:"updated_by" validate:"required"`
}

func mapCreateTicketDTOToCreateTicket(createTicketDTO CreateTicketDTO) (domain.CreateTicket, error) {
	crmTicket, err := domain.NewTicket(
		createTicketDTO.TenantID,
		createTicketDTO.CustomerID,
		createTicketDTO.OriginChannel,
		createTicketDTO.TicketType,
		createTicketDTO.Subject,
		createTicketDTO.DueDate,
		createTicketDTO.CreatedBy,
		createTicketDTO.ExternalReference,
	)
	if err != nil {
		return domain.CreateTicket{}, err
	}

	product, err := domain.NewProduct(
		createTicketDTO.ProductName,
		createTicketDTO.ProductDescription,
		createTicketDTO.Value,
		createTicketDTO.Brand,
		createTicketDTO.Model,
		createTicketDTO.SerialNumber,
		createTicketDTO.CreatedBy,
	)
	if err != nil {
		return domain.CreateTicket{}, err
	}

	return domain.CreateTicket{
		Ticket:  crmTicket,
		Product: product,
	}, nil
}

func mapTicketToTicketDTO(crmTicket domain.Ticket) TicketDTO {
	return TicketDTO{
		TicketID:          crmTicket.TicketID,
		TenantID:          crmTicket.TenantID,
		CustomerID:        crmTicket.CustomerID,
		LeadID:            crmTicket.LeadID,
		OwnerID:           crmTicket.OwnerID,
		OriginChannel:     crmTicket.OriginChannel,
		Type:              crmTicket.Type,
		Subject:           crmTicket.Subject,
		Priority:          crmTicket.Priority,
		Status:            crmTicket.Status,
		DueDate:           crmTicket.DueDate,
		CreatedBy:         crmTicket.CreatedBy,
		CreatedAt:         crmTicket.CreatedAt,
		UpdatedBy:         crmTicket.UpdatedBy,
		UpdatedAt:         crmTicket.UpdatedAt,
		Region:            crmTicket.Region,
		ExternalReference: crmTicket.ExternalReference,
		ProductID:         crmTicket.ProductID,
		ClosedAt:          crmTicket.ClosedAt,
		TargetDate:        crmTicket.TargetDate,
	}
}

func mapTicketsToTicketDTOs(crmTickets []domain.Ticket) []TicketDTO {
	crmTicketsDTO := make([]TicketDTO, len(crmTickets))
	for i, crmTicket := range crmTickets {
		crmTicketsDTO[i] = mapTicketToTicketDTO(crmTicket)
	}
	return crmTicketsDTO
}

func mapUpdateTicketDTOToUpdateTicket(dto UpdateTicketDTO) domain.TicketUpdate {
	return domain.TicketUpdate{
		TargetDate: dto.TargetDate,
		UpdatedBy:  dto.UpdatedBy,
	}
}
