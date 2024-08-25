package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type LeadDTO struct {
	LeadID          string    `db:"lead_id"`
	FirstName       string    `db:"first_name"`
	LastName        string    `db:"last_name"`
	CompanyName     string    `db:"company_name"`
	LegalName       string    `db:"legal_name"`
	LeadType        string    `db:"lead_type"`
	Document        string    `db:"document"`
	DocumentType    string    `db:"document_type"`
	ShippingAddress string    `db:"shipping_address"`
	ShippingCity    string    `db:"shipping_city"`
	ShippingState   string    `db:"shipping_state"`
	ShippingZipCode string    `db:"shipping_zip_code"`
	ShippingCountry string    `db:"shipping_country"`
	BillingAddress  string    `db:"billing_address"`
	BillingCity     string    `db:"billing_city"`
	BillingState    string    `db:"billing_state"`
	BillingZipCode  string    `db:"billing_zip_code"`
	BillingCountry  string    `db:"billing_country"`
	PersonalPhone   string    `db:"personal_phone"`
	BusinessPhone   string    `db:"business_phone"`
	PersonalEmail   string    `db:"personal_email"`
	BusinessEmail   string    `db:"business_email"`
	Region          *int      `db:"region"`
	CreatedBy       string    `db:"created_by"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedBy       string    `db:"updated_by"`
	UpdatedAt       time.Time `db:"updated_at"`
	Active          bool      `db:"active"`
	Description     *string   `db:"description"`
}

func mapLeadToLeadDTO(lead domain.Lead) LeadDTO {
	return LeadDTO{
		LeadID:          lead.LeadID,
		FirstName:       lead.FirstName,
		LastName:        lead.LastName,
		CompanyName:     lead.CompanyName,
		LegalName:       lead.LegalName,
		LeadType:        lead.LeadType,
		Document:        lead.Document,
		DocumentType:    string(lead.DocumentType),
		ShippingAddress: lead.ShippingAddress.Address,
		ShippingCity:    lead.ShippingAddress.City,
		ShippingState:   lead.ShippingAddress.State,
		ShippingZipCode: lead.ShippingAddress.ZipCode,
		ShippingCountry: lead.ShippingAddress.Country,
		BillingAddress:  lead.BillingAddress.Address,
		BillingCity:     lead.BillingAddress.City,
		BillingState:    lead.BillingAddress.State,
		BillingZipCode:  lead.BillingAddress.ZipCode,
		BillingCountry:  lead.BillingAddress.Country,
		PersonalPhone:   lead.PersonalContact.PhoneNumber,
		BusinessPhone:   lead.BusinessContact.PhoneNumber,
		PersonalEmail:   lead.PersonalContact.Email,
		BusinessEmail:   lead.BusinessContact.Email,
		CreatedBy:       lead.CreatedBy,
		CreatedAt:       lead.CreatedAt,
		UpdatedBy:       lead.UpdatedBy,
		UpdatedAt:       lead.UpdatedAt,
		Active:          lead.Active,
		Description:     &lead.Description,
	}
}

func mapLeadDTOToLead(leadDTO LeadDTO) domain.Lead {
	var descriptionString string
	if leadDTO.Description != nil {
		descriptionString = *leadDTO.Description
	}

	return domain.Lead{
		LeadID:       leadDTO.LeadID,
		FirstName:    leadDTO.FirstName,
		LastName:     leadDTO.LastName,
		CompanyName:  leadDTO.CompanyName,
		LegalName:    leadDTO.LegalName,
		LeadType:     leadDTO.LeadType,
		Document:     leadDTO.Document,
		DocumentType: domain.DocumentType(leadDTO.DocumentType),
		ShippingAddress: domain.Address{
			Address: leadDTO.ShippingAddress,
			City:    leadDTO.ShippingCity,
			State:   leadDTO.ShippingState,
			ZipCode: leadDTO.ShippingZipCode,
			Country: leadDTO.ShippingCountry,
		},
		BillingAddress: domain.Address{
			Address: leadDTO.BillingAddress,
			City:    leadDTO.BillingCity,
			State:   leadDTO.BillingState,
			ZipCode: leadDTO.BillingZipCode,
			Country: leadDTO.BillingCountry,
		},
		PersonalContact: domain.Contact{
			PhoneNumber: leadDTO.PersonalPhone,
			Email:       leadDTO.PersonalEmail,
		},
		BusinessContact: domain.Contact{
			PhoneNumber: leadDTO.BusinessPhone,
			Email:       leadDTO.BusinessEmail,
		},
		CreatedBy:   leadDTO.CreatedBy,
		CreatedAt:   leadDTO.CreatedAt,
		UpdatedBy:   leadDTO.UpdatedBy,
		UpdatedAt:   leadDTO.UpdatedAt,
		Active:      leadDTO.Active,
		Description: descriptionString,
	}
}

func mapLeadDTOsToLeads(leadDTOs []LeadDTO) []domain.Lead {
	leads := make([]domain.Lead, 0, len(leadDTOs))
	for _, leadDTO := range leadDTOs {
		lead := mapLeadDTOToLead(leadDTO)
		leads = append(leads, lead)
	}

	return leads
}

func mapLeadsToLeadDTOs(leads []domain.Lead) []interface{} {
	leadDTOs := make([]interface{}, 0, len(leads))
	for _, lead := range leads {
		leadDTO := mapLeadToLeadDTO(lead)
		leadDTOs = append(leadDTOs, leadDTO)
	}

	return leadDTOs
}
