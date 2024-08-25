package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateTenantDTO struct {
	CompanyName               string                    `json:"company_name"`
	LegalName                 string                    `json:"legal_name"`
	Document                  string                    `json:"document"`
	BusinessContact           ContactDTO                `json:"business_contact"`
	TenantPlatformTemplateDTO TenantPlatformTemplateDTO `json:"template"`
	CreatedBy                 string                    `json:"created_by"`
}

type TenantDTO struct {
	TenantID        string     `json:"tenant_id"`
	CompanyName     string     `json:"company_name"`
	LegalName       string     `json:"legal_name"`
	Document        string     `json:"document"`
	DocumentType    string     `json:"document_type"`
	BusinessContact ContactDTO `json:"business_contact"`
	CreatedBy       string     `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedBy       string     `json:"updated_by"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Active          bool       `json:"active"`
}

type UpdateTenantDTO struct {
	CompanyName     *string     `json:"company_name"`
	LegalName       *string     `json:"legal_name"`
	Document        *string     `json:"document"`
	DocumentType    *string     `json:"document_type"`
	BusinessContact *ContactDTO `json:"business_contact"`
	UpdatedBy       string      `json:"updated_by"`
}

func mapTenantToTenantDTO(tenant domain.Tenant) TenantDTO {
	return TenantDTO{
		TenantID:        tenant.TenantID,
		CompanyName:     tenant.CompanyName,
		LegalName:       tenant.LegalName,
		Document:        tenant.Document,
		DocumentType:    string(tenant.DocumentType),
		BusinessContact: mapContactToContactDTO(tenant.BusinessContact),
		CreatedBy:       tenant.CreatedBy,
		CreatedAt:       tenant.CreatedAt,
		UpdatedBy:       tenant.UpdatedBy,
		UpdatedAt:       tenant.UpdatedAt,
		Active:          tenant.Active,
	}
}

func mapCreateTenantDTOToTenant(tenantDTO CreateTenantDTO) (domain.Tenant, error) {
	tenantPlatformTemplate, err := mapTenantPlatformTemplateDTOToTenantPlatformTemplate(tenantDTO.TenantPlatformTemplateDTO)
	if err != nil {
		return domain.Tenant{}, err
	}

	return domain.NewTenant(
		tenantDTO.LegalName,
		tenantDTO.CompanyName,
		tenantDTO.Document,
		tenantDTO.CreatedBy,
		mapContactDTOToContact(tenantDTO.BusinessContact),
		tenantPlatformTemplate,
	)
}

func mapTenantsToTenantDTOs(tenants []domain.Tenant) []TenantDTO {
	tenantDTOs := make([]TenantDTO, 0, len(tenants))
	for _, tenant := range tenants {
		tenantDTO := mapTenantToTenantDTO(tenant)
		tenantDTOs = append(tenantDTOs, tenantDTO)
	}

	return tenantDTOs
}

func mapUpdateTenantDTOToUpdateTenant(updateTenantDTO UpdateTenantDTO) domain.UpdateTenant {
	var parsedDocumentType *domain.DocumentType
	if updateTenantDTO.DocumentType != nil {
		documentType := domain.DocumentType(*updateTenantDTO.DocumentType)
		parsedDocumentType = &documentType
	}

	var parsedBusinessContact *domain.Contact
	if updateTenantDTO.BusinessContact != nil {
		businessContact := mapContactDTOToContact(*updateTenantDTO.BusinessContact)
		parsedBusinessContact = &businessContact
	}

	return domain.UpdateTenant{
		CompanyName:     updateTenantDTO.CompanyName,
		LegalName:       updateTenantDTO.LegalName,
		Document:        updateTenantDTO.Document,
		DocumentType:    parsedDocumentType,
		BusinessContact: parsedBusinessContact,
		UpdatedBy:       updateTenantDTO.UpdatedBy,
	}
}
