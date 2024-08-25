package application

import (
	"context"
	"errors"
	"strconv"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type ticketService struct {
	customerService  CustomerService
	userService      UserService
	ticketRepository domain.TicketRepository
	productService   ProductService
}

type TicketService interface {
	CreateTicket(ctx context.Context, newTicket domain.CreateTicket) (string, error)
	GetTicketByID(ctx context.Context, ticketID string) (*domain.Ticket, error)
	SearchTickets(ctx context.Context, filters domain.TicketFilters) (domain.PagingResult[domain.Ticket], error)
	UpdateTicket(ctx context.Context, ticketID string, newTicket domain.TicketUpdate) error
}

func NewTicketService(
	customerService CustomerService,
	ticketRepository domain.TicketRepository,
	productService ProductService,
	userService UserService,
) TicketService {
	return &ticketService{
		customerService:  customerService,
		ticketRepository: ticketRepository,
		productService:   productService,
		userService:      userService,
	}
}

func (c *ticketService) CreateTicket(ctx context.Context, newTicket domain.CreateTicket) (string, error) {
	crmTicket := newTicket.Ticket
	customer, err := c.customerService.GetByID(ctx, crmTicket.CustomerID)
	if err != nil {
		return "", err
	}
	crmTicket.Region = customer.GetRegion()

	err = c.assignOwnerToNewTicket(ctx, &crmTicket)
	if err != nil {
		return "", err
	}

	productID, err := c.productService.CreateProduct(ctx, newTicket.Product)
	if err != nil {
		return "", err
	}
	crmTicket.ProductID = productID

	ticketID, err := c.ticketRepository.Create(ctx, crmTicket)
	if err != nil {
		return "", err
	}

	return ticketID, nil
}

func (c *ticketService) GetTicketByID(ctx context.Context, ticketID string) (*domain.Ticket, error) {
	if ticketID == "" {
		return nil, domain.NewValidationError("ticket id cannot be empty", nil)
	}

	return c.ticketRepository.GetByID(ctx, ticketID)
}

func (c *ticketService) SearchTickets(ctx context.Context, filters domain.TicketFilters) (domain.PagingResult[domain.Ticket], error) {
	return c.ticketRepository.Search(ctx, filters)
}

func (c *ticketService) assignOwnerToNewTicket(ctx context.Context, crmTicket *domain.Ticket) error {
	regionStringified := strconv.Itoa(crmTicket.Region)

	user, err := c.userService.Search(ctx, domain.UserFilters{
		Region: []string{regionStringified},
		Role:   []string{string(domain.OPERATOR)},
	})
	if err != nil {
		var customErr *domain.CustomError
		if !errors.As(err, &customErr) || !customErr.IsNotFound() {
			return domain.NewValidationError("user not found", nil)
		}
	}

	if len(user) > 0 {
		crmTicket.OwnerID = user[0].UserID
		crmTicket.Status = domain.CUSTOMER_INFO
	}

	return nil
}

func (c *ticketService) UpdateTicket(ctx context.Context, ticketID string, newTicket domain.TicketUpdate) error {
	if ticketID == "" {
		return domain.NewValidationError("ticket id cannot be empty", nil)
	}

	crmTicket, err := c.ticketRepository.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	crmTicket.MergeUpdate(newTicket)

	return c.ticketRepository.Update(ctx, *crmTicket)
}
