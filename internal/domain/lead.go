package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LeadRepository interface {
	Create(ctx context.Context, lead Lead) (string, error)
	GetByID(ctx context.Context, leadID string) (*Lead, error)
	Search(ctx context.Context, filters LeadFilters) (PagingResult[Lead], error)
	Update(ctx context.Context, leadToUpdate Lead) error
	Delete(ctx context.Context, leadID string) error
	CreateBatch(ctx context.Context, leads []Lead) ([]string, error)
}

type Lead struct {
	LeadID          string
	FirstName       string
	LastName        string
	CompanyName     string
	LegalName       string
	LeadType        string
	Document        string
	DocumentType    DocumentType
	ShippingAddress Address
	BillingAddress  Address
	BusinessContact Contact
	PersonalContact Contact
	Tickets         []Ticket
	CreatedBy       string
	CreatedAt       time.Time
	UpdatedBy       string
	UpdatedAt       time.Time
	Active          bool
	Description     string
}

type EditLead struct {
	FirstName       *string
	LastName        *string
	CompanyName     *string
	LegalName       *string
	LeadType        *EntityType
	Document        *string
	DocumentType    *DocumentType
	ShippingAddress *Address
	BillingAddress  *Address
	BusinessContact *Contact
	PersonalContact *Contact
	Active          *bool
	UpdatedBy       string
	Description     *string
}

type LeadFilters struct {
	LeadID   []string
	State    []string
	Document []string
	LeadType []string
	Active   *bool
	PagingFilter
}

func NewLead(
	firstName,
	lastName,
	companyName,
	legalName,
	document,
	documentType,
	author string,
	personalContact,
	businessContact Contact,
	shippingAddress,
	billingAddress Address,
	description,
	leadType string,
) (Lead, error) {
	now := time.Now().UTC()

	leadID, err := uuid.NewUUID()
	if err != nil {
		return Lead{}, err
	}

	return Lead{
		LeadID:          leadID.String(),
		FirstName:       firstName,
		LastName:        lastName,
		CompanyName:     companyName,
		LegalName:       legalName,
		Document:        document,
		DocumentType:    DocumentType(documentType),
		LeadType:        leadType,
		ShippingAddress: shippingAddress,
		BillingAddress:  billingAddress,
		PersonalContact: personalContact,
		BusinessContact: businessContact,
		CreatedAt:       now,
		CreatedBy:       author,
		UpdatedAt:       now,
		UpdatedBy:       author,
		Active:          true,
		Description:     description,
	}, nil
}

func (p *Lead) MergeUpdate(updateLead EditLead) {
	p.UpdatedBy = updateLead.UpdatedBy
	p.UpdatedAt = time.Now().UTC()

	if updateLead.FirstName != nil {
		p.FirstName = *updateLead.FirstName
	}

	if updateLead.LastName != nil {
		p.LastName = *updateLead.LastName
	}

	if updateLead.CompanyName != nil {
		p.CompanyName = *updateLead.CompanyName
	}

	if updateLead.LegalName != nil {
		p.LegalName = *updateLead.LegalName
	}

	if updateLead.Document != nil {
		p.Document = *updateLead.Document
	}

	if updateLead.DocumentType != nil {
		p.DocumentType = *updateLead.DocumentType
	}

	if updateLead.ShippingAddress != nil {
		p.ShippingAddress = *updateLead.ShippingAddress
	}

	if updateLead.BillingAddress != nil {
		p.BillingAddress = *updateLead.BillingAddress
	}

	if updateLead.BusinessContact != nil {
		p.BusinessContact = *updateLead.BusinessContact
	}

	if updateLead.PersonalContact != nil {
		p.PersonalContact = *updateLead.PersonalContact
	}

	if updateLead.Active != nil {
		p.Active = *updateLead.Active
	}

	if updateLead.Description != nil {
		p.Description = *updateLead.Description
	}
}

func (p *Lead) GetRegion() int {
	return regions[p.ShippingAddress.State]
}
