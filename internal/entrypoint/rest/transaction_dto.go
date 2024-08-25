package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateTransactionDTO struct {
	Type        domain.TransactionType `json:"type"`
	Value       float64                `json:"value"`
	CreatedBy   string                 `json:"created_by"`
	Description string                 `json:"description"`
}

type TransactionDTO struct {
	TransactionID string                   `json:"transaction_id"`
	Type          domain.TransactionType   `json:"type"`
	Value         float64                  `json:"value"`
	TicketID      string                   `json:"ticket_id"`
	Status        domain.TransactionStatus `json:"status"`
	AttachmentID  string                   `json:"attachment_id"`
	CreatedBy     string                   `json:"created_by"`
	CreatedAt     time.Time                `json:"created_at"`
	UpdatedBy     string                   `json:"updated_by"`
	UpdatedAt     time.Time                `json:"updated_at"`
	Description   string                   `json:"description"`
}

type TransactionUpdateDTO struct {
	Status       *domain.TransactionStatus `json:"status"`
	AttachmentID *string                   `json:"attachment_id"`
	Value        *float64                  `json:"value"`
	UpdatedBy    string                    `json:"updated_by"`
}

func mapCreateTransactionDTOToTransaction(transactionDTO CreateTransactionDTO, ticketID string) (domain.Transaction, error) {
	return domain.NewTransaction(
		transactionDTO.Type,
		transactionDTO.Value,
		ticketID,
		transactionDTO.CreatedBy,
		transactionDTO.Description,
	)
}

func mapTransactionToTransactionDTO(transaction domain.Transaction) TransactionDTO {
	return TransactionDTO{
		TransactionID: transaction.TransactionID,
		Type:          transaction.Type,
		Value:         transaction.Value,
		TicketID:      transaction.TicketID,
		Status:        transaction.Status,
		AttachmentID:  transaction.AttachmentID,
		CreatedBy:     transaction.CreatedBy,
		CreatedAt:     transaction.CreatedAt,
		UpdatedBy:     transaction.UpdatedBy,
		UpdatedAt:     transaction.UpdatedAt,
		Description:   transaction.Description,
	}
}

func mapTransactionsToTransactionsDTO(transactions []domain.Transaction) []TransactionDTO {
	transactionsDTO := make([]TransactionDTO, 0, len(transactions))
	for _, transaction := range transactions {
		transactionsDTO = append(transactionsDTO, mapTransactionToTransactionDTO(transaction))
	}
	return transactionsDTO
}

func mapTransactionUpdateDTOToTransactionUpdate(transactionUpdateDTO TransactionUpdateDTO) domain.TransactionUpdate {
	return domain.TransactionUpdate{
		Status:       transactionUpdateDTO.Status,
		AttachmentID: transactionUpdateDTO.AttachmentID,
		Value:        transactionUpdateDTO.Value,
		UpdatedBy:    transactionUpdateDTO.UpdatedBy,
	}
}
