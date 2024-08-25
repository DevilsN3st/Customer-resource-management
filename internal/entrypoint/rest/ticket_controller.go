package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type TicketController struct {
	ticketService application.TicketService
}

func NewTicketController(
	ticketService application.TicketService,
) TicketController {
	return TicketController{
		ticketService: ticketService,
	}
}

func (c *TicketController) CreateTicket(ctx *gin.Context) {
	var createTicketDTO *CreateTicketDTO
	if err := ctx.BindJSON(&createTicketDTO); err != nil {
		ctx.Error(err)
		return
	}

	newTicket, err := mapCreateTicketDTOToCreateTicket(*createTicketDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	ticketID, err := c.ticketService.CreateTicket(ctx.Request.Context(), newTicket)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"ticket_id": ticketID})
}

func (c *TicketController) GetTicket(ctx *gin.Context) {
	ticketID := ctx.Param("ticketID")
	if ticketID == "" {
		ctx.Error(domain.NewValidationError("ticket_id is required", nil))
		return
	}

	crmTicket, err := c.ticketService.GetTicketByID(ctx.Request.Context(), ticketID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ticketDTO := mapTicketToTicketDTO(*crmTicket)

	ctx.JSON(http.StatusOK, ticketDTO)
}

func (c *TicketController) SearchTickets(ctx *gin.Context) {
	filters := c.parseQueryToFilters(ctx)

	tickets, err := c.ticketService.SearchTickets(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	searchResult := mapSearchResultToSearchResultDTO(tickets, mapTicketsToTicketDTOs)

	ctx.JSON(http.StatusOK, searchResult)
}

func (c *TicketController) parseQueryToFilters(ctx *gin.Context) domain.TicketFilters {
	filters := domain.TicketFilters{
		PagingFilter: domain.PagingFilter{
			Limit:  10,
			Offset: 0,
		},
	}

	if ownerIDs := ctx.QueryArray("owner_id"); len(ownerIDs) > 0 {
		filters.OwnerID = ownerIDs
	}

	if tenantIDs := ctx.QueryArray("tenant_id"); len(tenantIDs) > 0 {
		filters.TenantID = tenantIDs
	}

	if leadIDs := ctx.QueryArray("lead_id"); len(leadIDs) > 0 {
		filters.LeadID = leadIDs
	}

	if customerIDs := ctx.QueryArray("customer_id"); len(customerIDs) > 0 {
		filters.CustomerID = customerIDs
	}

	if status := ctx.QueryArray("status"); len(status) > 0 {
		filters.Status = status
	}

	if region := ctx.QueryArray("region"); len(region) > 0 {
		filters.Region = region
	}

	if limit := ctx.Query("limit"); limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err == nil {
			filters.Limit = parsedLimit
		}
	}

	if offset := ctx.Query("offset"); offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err == nil {
			filters.Offset = parsedOffset
		}
	}

	return filters
}

func (c *TicketController) UpdateTicket(ctx *gin.Context) {
	ticketID := ctx.Param("ticketID")
	if ticketID == "" {
		ctx.Error(domain.NewValidationError("ticket_id is required", nil))
		return
	}

	var updateTicketDTO *UpdateTicketDTO
	if err := ctx.BindJSON(&updateTicketDTO); err != nil {
		ctx.Error(err)
		return
	}

	ticketUpdate := mapUpdateTicketDTOToUpdateTicket(*updateTicketDTO)

	err := c.ticketService.UpdateTicket(ctx.Request.Context(), ticketID, ticketUpdate)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
