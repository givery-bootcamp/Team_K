package controllers

import (
	"errors"
	"myapp/internal/controllers/request"
	"myapp/internal/controllers/response"
	"myapp/internal/repositories"

	"github.com/gin-gonic/gin"
)

func GetRepliesByPostId(ctx *gin.Context) {
	req := request.GetRepliesRequest{}
	if err := ctx.ShouldBindQuery(&req); err != nil {
		handleError(ctx, 400, errors.New("get replies query param is invalid"))
	}
	repository := repositories.NewReplyRepository(DB(ctx))
	replies, err := repository.GetRepliesByPostId(req.PostId)
	if err != nil {
		handleError(ctx, 500, err)
	} else if replies != nil {
		res := response.GetRepliesByPostIdResponse{}
		for _, reply := range replies {
			res = append(res, response.NewReplyResponse(reply))
		}
		ctx.JSON(200, res)
	} else {
		handleError(ctx, 404, errors.New("not found"))
	}
}

func PostReply(ctx *gin.Context) {
	req := request.PostReplyRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleError(ctx, 400, errors.New("post reply request body is invalid"))
	}
	repository := repositories.NewReplyRepository(DB(ctx))
	reply, err := repository.CreateReply(req.PostId, req.UserId, req.Content)
	if err != nil {
		handleError(ctx, 500, errors.New("create reply record failed"))
	}
	ctx.JSON(201, response.NewReplyResponse(reply))
}
