package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type TenantDTO struct {
	TenantID      string    `db:"tenant_id"`
	CompanyName   string    `db:"company_name"`
	LegalName     string    `db:"legal_name"`
	Document      string    `db:"document"`
	DocumentType  string    `db:"document_type"`
	BusinessPhone string    `db:"business_phone"`
	BusinessEmail string    `db:"business_email"`
	CreatedBy     string    `db:"created_by"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedBy     string    `db:"updated_by"`
	UpdatedAt     time.Time `db:"updated_at"`
	Active        bool      `db:"active"`
}

func mapTenantToTenantDTO(tenant domain.Tenant) TenantDTO {
	return TenantDTO{
		TenantID:      tenant.TenantID,
		CompanyName:   tenant.CompanyName,
		LegalName:     tenant.LegalName,
		Document:      tenant.Document,
		DocumentType:  string(tenant.DocumentType),
		BusinessPhone: tenant.BusinessContact.PhoneNumber,
		BusinessEmail: tenant.BusinessContact.Email,
		CreatedBy:     tenant.CreatedBy,
		CreatedAt:     tenant.CreatedAt,
		UpdatedBy:     tenant.UpdatedBy,
		UpdatedAt:     tenant.UpdatedAt,
		Active:        tenant.Active,
	}
}

func mapTenantDTOToTenant(tenantDTO TenantDTO) domain.Tenant {
	return domain.Tenant{
		TenantID:     tenantDTO.TenantID,
		CompanyName:  tenantDTO.CompanyName,
		LegalName:    tenantDTO.LegalName,
		Document:     tenantDTO.Document,
		DocumentType: domain.DocumentType(tenantDTO.DocumentType),
		BusinessContact: domain.Contact{
			PhoneNumber: tenantDTO.BusinessPhone,
			Email:       tenantDTO.BusinessEmail,
		},
		CreatedBy: tenantDTO.CreatedBy,
		CreatedAt: tenantDTO.CreatedAt,
		UpdatedBy: tenantDTO.UpdatedBy,
		UpdatedAt: tenantDTO.UpdatedAt,
		Active:    tenantDTO.Active,
	}
}

func mapTenantDTOsToTenants(tenantDTOs []TenantDTO) []domain.Tenant {
	tenants := make([]domain.Tenant, 0, len(tenantDTOs))
	for _, tenantDTO := range tenantDTOs {
		tenant := mapTenantDTOToTenant(tenantDTO)
		tenants = append(tenants, tenant)
	}

	return tenants
}
