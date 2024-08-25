package database

import (
	"context"
	"database/sql"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type commentRepository struct {
	db *mongo.Client
}

func NewCommentRepository(db *mongo.Client) domain.CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (r *commentRepository) commentCollection(ctx context.Context) *mongo.Collection {
	userCollection := GetCollection(r.db, "user")
	return userCollection
}

func (r *commentRepository) Create(ctx context.Context, comment domain.Comment) (string, error) {
	commentDTO := mapCommentToCommentDTO(comment)

	_, err := r.commentCollection(ctx).InsertOne(ctx, commentDTO)

	if err != nil {
		return "", err
	}

	return comment.CommentID, nil
}

func (r *commentRepository) GetByID(ctx context.Context, commentID string) (*domain.Comment, error) {
	if commentID == "" {
		return nil, domain.NewValidationError("commentID is required", nil)
	}

	var commentDTO CommentDTO
	err := r.commentCollection(ctx).FindOne(context.TODO(), bson.D{{"_id", commentID}}).Decode(&commentDTO)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no comment found with this id", map[string]any{"comment_id": commentID})
		}
		return nil, err
	}

	comment := mapCommentDTOToComment(commentDTO)

	return &comment, nil
}

func (r *commentRepository) GetByTicketID(ctx context.Context, ticketID string) ([]domain.Comment, error) {
	if ticketID == "" {
		return nil, domain.NewValidationError("ticketID is required", nil)
	}

	var commentDTOs []CommentDTO
	err := r.commentCollection(ctx).FindOne(context.TODO(), bson.D{{"ticket_id", ticketID}}).Decode(&commentDTOs)
	if err != nil {
		return nil, err
	}

	comments := mapCommentDTOsToComments(commentDTOs)

	return comments, nil
}
