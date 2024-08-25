package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateLeadDTO struct {
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	CompanyName     string     `json:"company_name"`
	LegalName       string     `json:"legal_name"`
	LeadType        string     `json:"lead_type"`
	Document        string     `json:"document"`
	DocumentType    string     `json:"document_type"`
	ShippingAddress AddressDTO `json:"shipping"`
	BillingAddress  AddressDTO `json:"billing"`
	PersonalContact ContactDTO `json:"personal_contact"`
	BusinessContact ContactDTO `json:"business_contact"`
	CreatedBy       string     `json:"created_by"`
	Description     string     `json:"description"`
}

type LeadDTO struct {
	LeadID          string     `json:"lead_id"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	CompanyName     string     `json:"company_name"`
	LegalName       string     `json:"legal_name"`
	LeadType        string     `json:"lead_type"`
	Document        string     `json:"document"`
	DocumentType    string     `json:"document_type"`
	ShippingAddress AddressDTO `json:"shipping"`
	BillingAddress  AddressDTO `json:"billing"`
	PersonalContact ContactDTO `json:"personal_contact"`
	BusinessContact ContactDTO `json:"business_contact"`
	Region          int        `json:"region"`
	Tickets         []any      `json:"tickets"`
	CreatedBy       string     `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedBy       string     `json:"updated_by"`
	UpdatedAt       time.Time  `json:"updated_at"`
	Active          bool       `json:"active"`
	Description     string     `json:"description"`
}

type EditLeadDTO struct {
	FirstName       *string     `json:"first_name"`
	LastName        *string     `json:"last_name"`
	CompanyName     *string     `json:"company_name"`
	LegalName       *string     `json:"legal_name"`
	LeadType        *string     `json:"lead_type"`
	Document        *string     `json:"document"`
	DocumentType    *string     `json:"document_type"`
	ShippingAddress *AddressDTO `json:"shipping"`
	BillingAddress  *AddressDTO `json:"billing"`
	PersonalContact *ContactDTO `json:"personal_contact"`
	BusinessContact *ContactDTO `json:"business_contact"`
	Active          *bool       `json:"active"`
	UpdatedBy       string      `json:"updated_by"`
	Description     *string     `json:"description"`
}

func mapLeadToLeadDTO(lead domain.Lead) LeadDTO {
	return LeadDTO{
		LeadID:          lead.LeadID,
		FirstName:       lead.FirstName,
		LastName:        lead.LastName,
		CompanyName:     lead.CompanyName,
		LegalName:       lead.LegalName,
		LeadType:        string(lead.LeadType),
		Document:        lead.Document,
		DocumentType:    string(lead.DocumentType),
		ShippingAddress: mapAddressToAddressDTO(lead.ShippingAddress),
		BillingAddress:  mapAddressToAddressDTO(lead.BillingAddress),
		Region:          lead.GetRegion(),
		PersonalContact: mapContactToContactDTO(lead.PersonalContact),
		BusinessContact: mapContactToContactDTO(lead.BusinessContact),
		CreatedBy:       lead.CreatedBy,
		CreatedAt:       lead.CreatedAt,
		UpdatedBy:       lead.UpdatedBy,
		UpdatedAt:       lead.UpdatedAt,
		Active:          lead.Active,
		Description:     lead.Description,
	}
}

func mapCreateLeadDTOToLead(leadDTO CreateLeadDTO) (domain.Lead, error) {
	return domain.NewLead(
		leadDTO.FirstName,
		leadDTO.LastName,
		leadDTO.CompanyName,
		leadDTO.LegalName,
		leadDTO.Document,
		leadDTO.DocumentType,
		leadDTO.CreatedBy,
		mapContactDTOToContact(leadDTO.PersonalContact),
		mapContactDTOToContact(leadDTO.BusinessContact),
		mapAddressDTOToAddress(leadDTO.ShippingAddress),
		mapAddressDTOToAddress(leadDTO.BillingAddress),
		leadDTO.Description,
		leadDTO.LeadType,
	)
}

func mapLeadsToLeadDTOs(leads []domain.Lead) []LeadDTO {
	leadDTOs := make([]LeadDTO, 0, len(leads))
	for _, lead := range leads {
		leadDTO := mapLeadToLeadDTO(lead)
		leadDTOs = append(leadDTOs, leadDTO)
	}

	return leadDTOs
}

func mapEditLeadDTOToEditLead(editLeadDTO EditLeadDTO) domain.EditLead {
	var parsedLeadType *domain.EntityType
	if editLeadDTO.LeadType != nil {
		leadType := domain.EntityType(*editLeadDTO.LeadType)
		parsedLeadType = &leadType
	}

	var parsedDocumentType *domain.DocumentType
	if editLeadDTO.DocumentType != nil {
		documentType := domain.DocumentType(*editLeadDTO.DocumentType)
		parsedDocumentType = &documentType
	}

	var parsedShippingAddress *domain.Address
	if editLeadDTO.ShippingAddress != nil {
		shippingAddress := mapAddressDTOToAddress(*editLeadDTO.ShippingAddress)
		parsedShippingAddress = &shippingAddress
	}

	var parsedBillingAddress *domain.Address
	if editLeadDTO.BillingAddress != nil {
		billingAddress := mapAddressDTOToAddress(*editLeadDTO.BillingAddress)
		parsedBillingAddress = &billingAddress
	}

	var parsedPersonalContact *domain.Contact
	if editLeadDTO.PersonalContact != nil {
		personalContact := mapContactDTOToContact(*editLeadDTO.PersonalContact)
		parsedPersonalContact = &personalContact
	}

	var parsedBusinessContact *domain.Contact
	if editLeadDTO.BusinessContact != nil {
		businessContact := mapContactDTOToContact(*editLeadDTO.BusinessContact)
		parsedBusinessContact = &businessContact
	}

	return domain.EditLead{
		FirstName:       editLeadDTO.FirstName,
		LastName:        editLeadDTO.LastName,
		CompanyName:     editLeadDTO.CompanyName,
		LegalName:       editLeadDTO.LegalName,
		LeadType:        parsedLeadType,
		Document:        editLeadDTO.Document,
		DocumentType:    parsedDocumentType,
		ShippingAddress: parsedShippingAddress,
		BillingAddress:  parsedBillingAddress,
		PersonalContact: parsedPersonalContact,
		BusinessContact: parsedBusinessContact,
		Active:          editLeadDTO.Active,
		UpdatedBy:       editLeadDTO.UpdatedBy,
		Description:     editLeadDTO.Description,
	}
}
