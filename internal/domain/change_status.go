package domain

type ChangeStatus struct {
	Status      TicketStatus
	UpdatedBy   string
	Content     *string
	Attachments []Attachment
}
