package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type CommentRepository interface {
	Create(ctx context.Context, comment Comment) (string, error)
	GetByID(ctx context.Context, commentID string) (*Comment, error)
	GetByTicketID(ctx context.Context, ticketID string) ([]Comment, error)
}

type Comment struct {
	CommentID   string
	TicketID    string
	Content     string
	CommentType CommentType
	Attachments []Attachment
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedBy   string
	UpdatedAt   time.Time
}

type CommentType string

const (
	CONTENT    CommentType = "Content"
	COMMENT    CommentType = "Comment"
	RESOLUTION CommentType = "Resolution"
	REJECTION  CommentType = "Rejection"
)

func NewComment(
	ticketID string,
	content string,
	createdBy string,
	commentType CommentType,
	attachments []Attachment,
) (Comment, error) {
	now := time.Now().UTC()
	commentID, err := uuid.NewUUID()
	if err != nil {
		return Comment{}, err
	}

	return Comment{
		CommentID:   commentID.String(),
		CommentType: commentType,
		TicketID:    ticketID,
		Content:     content,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
		CreatedAt:   now,
		UpdatedAt:   now,
		Attachments: attachments,
	}, nil
}
