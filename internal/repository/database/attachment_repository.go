package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/icrxz/crm-api-core/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type attachmentRepository struct {
	db *mongo.Client
}

func NewAttachmentRepository(db *mongo.Client) domain.AttachmentRepository {
	return &attachmentRepository{
		db: db,
	}
}

func (r *attachmentRepository) attachmentCollection(ctx context.Context) *mongo.Collection {
	userCollection := GetCollection(r.db, "user")
	return userCollection
}

func (r *attachmentRepository) Save(ctx context.Context, attachment domain.Attachment) error {
	return nil
}

func (r *attachmentRepository) SaveBatch(ctx context.Context, attachments []domain.Attachment) error {
	if len(attachments) == 0 {
		return nil
	}

	attachmentsDTO := mapAttachmentsToAttachmentsDTO(attachments)

	_, err := r.attachmentCollection(ctx).InsertMany(ctx, attachmentsDTO)
	if err != nil {
		return err
	}

	return nil
}

func (r *attachmentRepository) GetByID(ctx context.Context, attachmentID string) (domain.Attachment, error) {
	if attachmentID == "" {
		return domain.Attachment{}, domain.NewValidationError("attachment_id is required", nil)
	}

	var attachmentDTO AttachmentDTO
	err := r.attachmentCollection(ctx).FindOne(context.TODO(), bson.D{{"_id", attachmentID}}).Decode(&attachmentDTO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Attachment{}, domain.NewNotFoundError("attachment not found", map[string]any{"attachment_id": attachmentID})
		}

		return domain.Attachment{}, err
	}

	return mapAttachmentDTOToAttachment(attachmentDTO), nil
}

func (r *attachmentRepository) GetByCommentID(ctx context.Context, commentID string) ([]domain.Attachment, error) {
	if commentID == "" {
		return nil, domain.NewValidationError("comment_id is required", nil)
	}

	var attachmentsDTO []AttachmentDTO
	err := r.attachmentCollection(ctx).FindOne(context.TODO(), bson.D{{"comment_id", commentID}}).Decode(&attachmentsDTO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	attachments := mapAttachmentsDTOToAttachments(attachmentsDTO)

	return attachments, nil
}
