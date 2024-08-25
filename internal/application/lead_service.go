package application

import (
	"context"
	"encoding/csv"
	"io"
	"strings"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type leadService struct {
	leadRepository domain.LeadRepository
}

type LeadService interface {
	Create(ctx context.Context, lead domain.Lead) (string, error)
	GetByID(ctx context.Context, leadID string) (*domain.Lead, error)
	Update(ctx context.Context, leadID string, editLead domain.EditLead) error
	Delete(ctx context.Context, leadID string) error
	Search(ctx context.Context, filters domain.LeadFilters) (domain.PagingResult[domain.Lead], error)
	CreateBatch(ctx context.Context, file io.Reader, createdBy string) ([]string, error)
}

func NewLeadService(leadRepository domain.LeadRepository) LeadService {
	return &leadService{
		leadRepository: leadRepository,
	}
}

func (s *leadService) Create(ctx context.Context, lead domain.Lead) (string, error) {
	return s.leadRepository.Create(ctx, lead)
}

func (s *leadService) GetByID(ctx context.Context, leadID string) (*domain.Lead, error) {
	if leadID == "" {
		return nil, domain.NewValidationError("leadID cannot be empty", nil)
	}

	return s.leadRepository.GetByID(ctx, leadID)
}

func (s *leadService) Update(ctx context.Context, leadID string, editLead domain.EditLead) error {
	if leadID == "" {
		return domain.NewValidationError("leadID cannot be empty", nil)
	}

	lead, err := s.leadRepository.GetByID(ctx, leadID)
	if err != nil {
		return err
	}

	lead.MergeUpdate(editLead)

	return s.leadRepository.Update(ctx, *lead)
}

func (s *leadService) Delete(ctx context.Context, leadID string) error {
	if leadID == "" {
		return domain.NewValidationError("leadID cannot be empty", nil)
	}

	return s.leadRepository.Delete(ctx, leadID)
}

func (s *leadService) Search(ctx context.Context, filters domain.LeadFilters) (domain.PagingResult[domain.Lead], error) {
	return s.leadRepository.Search(ctx, filters)
}

func (s *leadService) CreateBatch(ctx context.Context, file io.Reader, createdBy string) ([]string, error) {
	fileCSV := csv.NewReader(file)

	leadsRows, err := s.readCSV(fileCSV)
	if err != nil {
		return nil, err
	}

	columnsIndex := s.getColumnHeadersIndex(leadsRows[0])
	leads, err := s.buildLead(leadsRows[1:], columnsIndex, createdBy)
	if err != nil {
		return nil, err
	}

	return s.leadRepository.CreateBatch(ctx, leads)
}

func (s *leadService) readCSV(fileCSV *csv.Reader) ([][]string, error) {
	csvRows := make([][]string, 0)

	for {
		row, err := fileCSV.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		csvRows = append(csvRows, row)
	}

	return csvRows, nil
}

func (s *leadService) buildLead(csvRows [][]string, columnsIndex map[string]int, author string) ([]domain.Lead, error) {
	leads := make([]domain.Lead, 0, len(csvRows))

	for _, row := range csvRows {
		phone := row[columnsIndex["Telefone"]]
		if strings.TrimSpace(phone) != "" {
			phone = strings.ReplaceAll(phone, "(", "")
			phone = strings.ReplaceAll(phone, ")", "")
			phone = "(55) +" + phone
		}

		personalContact := domain.Contact{
			PhoneNumber: phone,
		}

		shippingAddress := domain.Address{
			City:    row[columnsIndex["Cidade"]],
			State:   row[columnsIndex["Estado"]],
			Country: "brazil",
		}

		document := row[columnsIndex["Documento"]]
		document = strings.ReplaceAll(document, ".", "")
		document = strings.ReplaceAll(document, "-", "")
		document = strings.ReplaceAll(document, "/", "")

		documentType := ""
		if len(document) == 11 {
			documentType = "CPF"
		} else if len(document) == 14 {
			documentType = "CNPJ"
		}

		description := row[columnsIndex["Observacoes"]]

		newLead, err := domain.NewLead(
			row[columnsIndex["Nome"]],
			row[columnsIndex["Sobrenome"]],
			"",
			"",
			document,
			documentType,
			author,
			personalContact,
			domain.Contact{},
			shippingAddress,
			domain.Address{},
			description,
			row[columnsIndex["Tipo"]],
		)
		if err != nil {
			return nil, err
		}

		leads = append(leads, newLead)
	}
	return leads, nil
}

func (s *leadService) getColumnHeadersIndex(header []string) map[string]int {
	columnsIndex := make(map[string]int)
	for i, column := range header {
		columnsIndex[column] = i
	}
	return columnsIndex
}
