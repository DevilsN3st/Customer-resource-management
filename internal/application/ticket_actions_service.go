package application

import (
	"context"
	"github.com/icrxz/crm-api-core/internal/domain"
	"slices"
)

type ticketActionService struct {
	ticketRepository domain.TicketRepository
	commentService   CommentService
	reportService    ReportService
}

type TicketActionService interface {
	ChangeOwner(ctx context.Context, ticketID string, newOwner domain.ChangeOwner) error
	ChangeStatus(ctx context.Context, ticketID string, newStatus domain.ChangeStatus) error
	ChangeLead(ctx context.Context, ticketID string, newLead domain.ChangeLead) error
	GenerateReport(ctx context.Context, ticketID string) ([]byte, string, error)
}

func NewTicketActionService(
	ticketRepository domain.TicketRepository,
	commentService CommentService,
	reportService ReportService,
) TicketActionService {
	return &ticketActionService{
		ticketRepository: ticketRepository,
		commentService:   commentService,
		reportService:    reportService,
	}
}

func (c *ticketActionService) ChangeOwner(ctx context.Context, ticketID string, newOwner domain.ChangeOwner) error {
	crmTicket, err := c.ticketRepository.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	ticketUpdate := domain.TicketUpdate{
		OwnerID:   &newOwner.OwnerID,
		Status:    &newOwner.Status,
		UpdatedBy: newOwner.UpdatedBy,
	}
	crmTicket.MergeUpdate(ticketUpdate)

	return c.ticketRepository.Update(ctx, *crmTicket)
}

func (c *ticketActionService) ChangeStatus(ctx context.Context, ticketID string, newStatus domain.ChangeStatus) error {
	crmTicket, err := c.ticketRepository.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	ticketUpdate := domain.TicketUpdate{
		Status:    &newStatus.Status,
		UpdatedBy: newStatus.UpdatedBy,
	}
	crmTicket.MergeUpdate(ticketUpdate)

	if newStatus.Content != nil {
		err = c.createChangeStatusComment(ctx, ticketID, newStatus)
		if err != nil {
			return err
		}
	}

	return c.ticketRepository.Update(ctx, *crmTicket)
}

func (c *ticketActionService) ChangeLead(ctx context.Context, ticketID string, newLead domain.ChangeLead) error {
	crmTicket, err := c.ticketRepository.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	ticketUpdate := domain.TicketUpdate{
		LeadID:     &newLead.LeadID,
		Status:     &newLead.Status,
		TargetDate: &newLead.TargetDate,
		UpdatedBy:  newLead.UpdatedBy,
	}
	crmTicket.MergeUpdate(ticketUpdate)

	return c.ticketRepository.Update(ctx, *crmTicket)
}

func (c *ticketActionService) GenerateReport(ctx context.Context, ticketID string) ([]byte, string, error) {
	if ticketID == "" {
		return nil, "", domain.NewValidationError("ticket_id is required", nil)
	}

	crmTicket, err := c.ticketRepository.GetByID(ctx, ticketID)
	if err != nil {
		return nil, "", err
	}

	if !slices.Contains([]domain.TicketStatus{domain.REPORT, domain.PAYMENT, domain.RECEIPT, domain.CLOSED}, crmTicket.Status) {
		return nil, "", domain.NewValidationError("ticket is not in status REPORT", map[string]any{"status": crmTicket.Status})
	}

	return c.reportService.GenerateReport(ctx, *crmTicket)
}

func (c *ticketActionService) createChangeStatusComment(ctx context.Context, ticketID string, newStatus domain.ChangeStatus) error {
	var commentType domain.CommentType
	switch newStatus.Status {
	case domain.WAITING_LEAD:
		commentType = domain.CONTENT
	case domain.REPORT:
		commentType = domain.RESOLUTION
	case domain.REJECTED:
		commentType = domain.REJECTION
	}

	newComment, err := domain.NewComment(ticketID, *newStatus.Content, newStatus.UpdatedBy, commentType, newStatus.Attachments)
	if err != nil {
		return err
	}

	_, err = c.commentService.Create(ctx, newComment)
	return err
}
