package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
	"net/http"
	"strconv"
	"strings"
)

type LeadController struct {
	leadService application.LeadService
}

func NewLeadController(leadService application.LeadService) LeadController {
	return LeadController{
		leadService: leadService,
	}
}

func (c *LeadController) CreateLead(ctx *gin.Context) {
	var leadDTO *CreateLeadDTO
	err := ctx.BindJSON(&leadDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	lead, err := mapCreateLeadDTOToLead(*leadDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	leadID, err := c.leadService.Create(ctx.Request.Context(), lead)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"lead_id": leadID})
}

func (c *LeadController) GetLead(ctx *gin.Context) {
	leadID := ctx.Param("leadID")
	if leadID == "" {
		ctx.Error(domain.NewValidationError("param leadID cannot be empty", nil))
		return
	}

	lead, err := c.leadService.GetByID(ctx.Request.Context(), leadID)
	if err != nil {
		ctx.Error(err)
		return
	}

	leadDTO := mapLeadToLeadDTO(*lead)

	ctx.JSON(http.StatusOK, leadDTO)
}

func (c *LeadController) SearchLeads(ctx *gin.Context) {
	filters, err := c.parseQueryToFilters(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	leads, err := c.leadService.Search(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	searchResult := mapSearchResultToSearchResultDTO(leads, mapLeadsToLeadDTOs)

	ctx.JSON(http.StatusOK, searchResult)
}

func (c *LeadController) UpdateLead(ctx *gin.Context) {
	leadID := ctx.Param("leadID")
	if leadID == "" {
		ctx.Error(domain.NewValidationError("param leadID cannot be empty", nil))
		return
	}

	var editLeadDTO *EditLeadDTO
	err := ctx.BindJSON(&editLeadDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	editLead := mapEditLeadDTOToEditLead(*editLeadDTO)

	err = c.leadService.Update(ctx.Request.Context(), leadID, editLead)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *LeadController) DeleteLead(ctx *gin.Context) {
	leadID := ctx.Param("leadID")
	if leadID == "" {
		ctx.Error(domain.NewValidationError("param leadID cannot be empty", nil))
		return
	}

	err := c.leadService.Delete(ctx.Request.Context(), leadID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *LeadController) CreateBatch(ctx *gin.Context) {
	author := ctx.GetHeader("X-Author")
	if author == "" {
		ctx.Error(domain.NewValidationError("header X-Author cannot be empty", nil))
		return
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(err)
		return
	}

	file, err := fileHeader.Open()
	defer file.Close()

	if !strings.Contains(fileHeader.Filename, ".csv") {
		ctx.Error(domain.NewValidationError("file must be a csv", nil))
		return
	}

	result, err := c.leadService.CreateBatch(ctx, file, author)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"lead_ids": result})
}

func (c *LeadController) parseQueryToFilters(ctx *gin.Context) (domain.LeadFilters, error) {
	filters := domain.LeadFilters{
		PagingFilter: domain.PagingFilter{
			Limit:  10,
			Offset: 0,
		},
	}

	if documents := ctx.QueryArray("document"); len(documents) > 0 {
		filters.Document = documents
	}

	if leadTypes := ctx.QueryArray("lead_type"); len(leadTypes) > 0 {
		filters.LeadType = leadTypes
	}

	if leadIDs := ctx.QueryArray("lead_id"); len(leadIDs) > 0 {
		filters.LeadID = leadIDs
	}

	if states := ctx.QueryArray("state"); len(states) > 0 {
		filters.State = states
	}

	if active := ctx.Query("active"); active != "" {
		isActive := active == "true"
		filters.Active = &isActive
	}

	validationErr := make([]error, 0)
	if limitParam := ctx.Query("limit"); limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err != nil {
			validationErr = append(validationErr, domain.NewValidationError("limit must be a number", nil))
		} else {
			filters.PagingFilter.Limit = parsedLimit
		}
	}

	if offsetParam := ctx.Query("offset"); offsetParam != "" {
		parsedOffset, err := strconv.Atoi(offsetParam)
		if err != nil {
			validationErr = append(validationErr, domain.NewValidationError("offset must be a number", nil))
		} else {
			filters.PagingFilter.Offset = parsedOffset
		}
	}

	if len(validationErr) > 0 {
		return domain.LeadFilters{}, errors.Join(validationErr...)
	}

	return filters, nil
}
