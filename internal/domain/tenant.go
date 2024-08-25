package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TenantRepository interface {
	Create(ctx context.Context, tenant Tenant) (string, error)
	GetByID(ctx context.Context, tenantID string) (*Tenant, error)
	Search(ctx context.Context, filters TenantFilters) (PagingResult[Tenant], error)
	Update(ctx context.Context, tenant Tenant) error
	Delete(ctx context.Context, tenantID string) error
}

type Tenant struct {
	TenantID        string
	CompanyName     string
	LegalName       string
	Document        string
	DocumentType    DocumentType
	BusinessContact Contact
	Template        TenantPlatformTemplate
	Tickets         []Ticket
	CreatedBy       string
	CreatedAt       time.Time
	UpdatedBy       string
	UpdatedAt       time.Time
	Active          bool
}

type UpdateTenant struct {
	CompanyName     *string
	LegalName       *string
	Document        *string
	DocumentType    *DocumentType
	BusinessContact *Contact
	UpdatedBy       string
}

func (c *Tenant) MergeUpdate(newTenant UpdateTenant) {
	now := time.Now().UTC()
	c.UpdatedAt = now
	c.UpdatedBy = newTenant.UpdatedBy

	if newTenant.CompanyName != nil {
		c.CompanyName = *newTenant.CompanyName
	}

	if newTenant.LegalName != nil {
		c.LegalName = *newTenant.LegalName
	}

	if newTenant.Document != nil {
		c.Document = *newTenant.Document
	}

	if newTenant.DocumentType != nil {
		c.DocumentType = *newTenant.DocumentType
	}

	if newTenant.BusinessContact != nil {
		c.BusinessContact = *newTenant.BusinessContact
	}
}

type TenantFilters struct {
	TenantID    []string
	CompanyName []string
	Document    []string
	Active      *bool
	PagingFilter
}

func NewTenant(legalName, companyName, document, author string, businessContact Contact, platformTemplate TenantPlatformTemplate) (Tenant, error) {
	now := time.Now().UTC()

	tenantID, err := uuid.NewRandom()
	if err != nil {
		return Tenant{}, err
	}

	return Tenant{
		TenantID:        tenantID.String(),
		CompanyName:     companyName,
		LegalName:       legalName,
		Document:        document,
		DocumentType:    CNPJ,
		BusinessContact: businessContact,
		Template:        platformTemplate,
		CreatedBy:       author,
		CreatedAt:       now,
		UpdatedBy:       author,
		UpdatedAt:       now,
		Active:          true,
	}, nil
}
