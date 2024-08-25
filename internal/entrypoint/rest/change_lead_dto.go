package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type ChangeLeadDTO struct {
	LeadID     string              `json:"lead_id"`
	TargetDate time.Time           `json:"target_date"`
	Status     domain.TicketStatus `json:"status"`
	UpdatedBy  string              `json:"updated_by"`
}

func mapChangeLeadDTOToChangeLead(c ChangeLeadDTO) domain.ChangeLead {
	return domain.ChangeLead{
		LeadID:     c.LeadID,
		TargetDate: c.TargetDate,
		Status:     c.Status,
		UpdatedBy:  c.UpdatedBy,
	}
}
