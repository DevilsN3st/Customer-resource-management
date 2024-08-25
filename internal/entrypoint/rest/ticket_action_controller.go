package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type TicketActionController struct {
	ticketActionService application.TicketActionService
}

func NewTicketActionController(
	ticketActionService application.TicketActionService,
) TicketActionController {
	return TicketActionController{
		ticketActionService: ticketActionService,
	}
}

func (c *TicketActionController) ChangeOwner(ctx *gin.Context) {
	ticketID := ctx.Param("ticketID")
	if ticketID == "" {
		ctx.Error(domain.NewValidationError("ticket_id is required", nil))
		return
	}

	var changeOwnerDTO ChangeOwnerDTO
	if err := ctx.BindJSON(&changeOwnerDTO); err != nil {
		ctx.Error(err)
		return
	}

	changeOwner := mapChangeOwnerDTOToChangeOwner(changeOwnerDTO)

	err := c.ticketActionService.ChangeOwner(ctx.Request.Context(), ticketID, changeOwner)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *TicketActionController) ChangeStatus(ctx *gin.Context) {
	ticketID := ctx.Param("ticketID")
	if ticketID == "" {
		ctx.Error(domain.NewValidationError("ticket_id is required", nil))
		return
	}

	var changeStatusDTO ChangeStatusDTO
	if err := ctx.BindJSON(&changeStatusDTO); err != nil {
		ctx.Error(err)
		return
	}

	changeStatus, err := mapChangeStatusDTOToChangeStatus(changeStatusDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.ticketActionService.ChangeStatus(ctx.Request.Context(), ticketID, changeStatus)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *TicketActionController) ChangeLead(ctx *gin.Context) {
	ticketID := ctx.Param("ticketID")
	if ticketID == "" {
		ctx.Error(domain.NewValidationError("ticket_id is required", nil))
		return
	}

	var changeLeadDTO ChangeLeadDTO
	if err := ctx.BindJSON(&changeLeadDTO); err != nil {
		ctx.Error(err)
		return
	}

	changeLead := mapChangeLeadDTOToChangeLead(changeLeadDTO)

	err := c.ticketActionService.ChangeLead(ctx.Request.Context(), ticketID, changeLead)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *TicketActionController) DownloadReport(ctx *gin.Context) {
	ticketID := ctx.Param("ticketID")
	if ticketID == "" {
		ctx.Error(domain.NewValidationError("ticket_id is required", nil))
		return
	}

	report, filename, err := c.ticketActionService.GenerateReport(ctx.Request.Context(), ticketID)
	if err != nil {
		ctx.Error(err)
		return
	}

	contentType := fmt.Sprintf("application/vnd.openxmlformats-officedocument.wordprocessingml.document;%s.docx", filename)
	ctx.Data(http.StatusOK, contentType, report)
}
