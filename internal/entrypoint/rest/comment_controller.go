package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type CommentController struct {
	commentService application.CommentService
}

func NewCommentController(commentService application.CommentService) CommentController {
	return CommentController{
		commentService: commentService,
	}
}

func (c *CommentController) CreateComment(ctx *gin.Context) {
	ticketID := ctx.Param("ticketID")
	if ticketID == "" {
		ctx.Error(domain.NewValidationError("ticketID is required", nil))
		return
	}

	var commentDTO CreateCommentDTO
	if err := ctx.ShouldBindJSON(&commentDTO); err != nil {
		ctx.Error(domain.NewValidationError("invalid request body", nil))
		return
	}

	comment, err := mapCreateCommentDTOToComment(commentDTO, ticketID)
	if err != nil {
		ctx.Error(err)
		return
	}

	commentID, err := c.commentService.Create(ctx, comment)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"comment_id": commentID})
}

func (c *CommentController) GetByID(ctx *gin.Context) {
	commentID := ctx.Param("commentID")
	if commentID == "" {
		ctx.Error(domain.NewValidationError("commentID is required", nil))
		return
	}

	comment, err := c.commentService.GetByID(ctx, commentID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mapCommentToCommentDTO(*comment))
}

func (c *CommentController) GetByTicketID(ctx *gin.Context) {
	ticketID := ctx.Param("ticketID")
	if ticketID == "" {
		ctx.Error(domain.NewValidationError("ticketID is required", nil))
		return
	}

	comments, err := c.commentService.GetByTicketID(ctx, ticketID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mapCommentsToCommentDTOs(comments))
}
