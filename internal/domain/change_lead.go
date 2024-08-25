package domain

import "time"

type ChangeLead struct {
	LeadID     string
	TargetDate time.Time
	Status     TicketStatus
	UpdatedBy  string
}
