package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
	"net/http"
	"strconv"
)

type TenantController struct {
	tenantService application.TenantService
}

func NewTenantController(tenantService application.TenantService) TenantController {
	return TenantController{
		tenantService: tenantService,
	}
}

func (c *TenantController) CreateTenant(ctx *gin.Context) {
	var tenantDTO *CreateTenantDTO
	err := ctx.BindJSON(&tenantDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	tenant, err := mapCreateTenantDTOToTenant(*tenantDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	tenantID, err := c.tenantService.Create(ctx.Request.Context(), tenant)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(201, gin.H{"tenant_id": tenantID})
}

func (c *TenantController) UpdateTenant(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")
	if tenantID == "" {
		ctx.Error(domain.NewValidationError("param tenantID cannot be empty", nil))
		return
	}

	var updateTenantDTO *UpdateTenantDTO
	err := ctx.BindJSON(&updateTenantDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	updateTenant := mapUpdateTenantDTOToUpdateTenant(*updateTenantDTO)

	if err = c.tenantService.Update(ctx.Request.Context(), tenantID, updateTenant); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *TenantController) GetTenant(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")
	if tenantID == "" {
		ctx.Error(domain.NewValidationError("param tenantID cannot be empty", nil))
		return
	}

	tenant, err := c.tenantService.GetByID(ctx.Request.Context(), tenantID)
	if err != nil {
		ctx.Error(err)
		return
	}

	tenantDTO := mapTenantToTenantDTO(*tenant)

	ctx.JSON(http.StatusOK, tenantDTO)
}

func (c *TenantController) SearchTenants(ctx *gin.Context) {
	filters := c.parseQueryToFilters(ctx)

	tenants, err := c.tenantService.Search(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	tenantResponse := mapSearchResultToSearchResultDTO(tenants, mapTenantsToTenantDTOs)

	ctx.JSON(http.StatusOK, tenantResponse)
}

func (c *TenantController) DeleteTenant(ctx *gin.Context) {
	tenantID := ctx.Param("tenantID")
	if tenantID == "" {
		ctx.Error(domain.NewValidationError("param tenantID cannot be empty", nil))
		return
	}

	err := c.tenantService.Delete(ctx.Request.Context(), tenantID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *TenantController) parseQueryToFilters(ctx *gin.Context) domain.TenantFilters {
	filters := domain.TenantFilters{
		PagingFilter: domain.PagingFilter{
			Limit:  10,
			Offset: 0,
		},
	}

	if documents := ctx.QueryArray("document"); len(documents) > 0 {
		filters.Document = documents
	}

	if tenantIDs := ctx.QueryArray("tenant_id"); len(tenantIDs) > 0 {
		filters.TenantID = tenantIDs
	}

	if companyNames := ctx.QueryArray("company_name"); len(companyNames) > 0 {
		filters.CompanyName = companyNames
	}

	if active := ctx.Query("active"); active != "" {
		activeBool, err := strconv.ParseBool(active)
		if err != nil {
			return filters
		}
		filters.Active = &activeBool
	}

	if limitParam := ctx.Query("limit"); limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err == nil {
			filters.PagingFilter.Limit = parsedLimit
		}
	}

	if offsetParam := ctx.Query("offset"); offsetParam != "" {
		parsedOffset, err := strconv.Atoi(offsetParam)
		if err == nil {
			filters.PagingFilter.Offset = parsedOffset
		}
	}

	return filters
}
