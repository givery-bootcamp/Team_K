package controllers

import (
	"errors"
	"myapp/internal/controllers/request"
	"myapp/internal/controllers/response"
	"myapp/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Login(ctx *gin.Context) {
	req := request.LoginRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handleError(ctx, 400, errors.New("login request body is invalid"))
		return
	}

	log.Println("req:", req)

	loginRepository := repositories.NewLoginRepository(DB(ctx))
	result, err := loginRepository.Login(req.Email)
	if err != nil {
		handleError(ctx, 404, err)
	}
	if result.Password != req.Password {
		handleError(ctx, 401, errors.New("password is invalid"))
		return
	}

	sessionId, err := uuid.NewRandom()
	if err != nil {
		handleError(ctx, 500, err)
		return
	}
	userId := result.Id
	sessionRepository := repositories.NewSessionRepository(Redis(ctx))
	if err := sessionRepository.Set(ctx, sessionId.String(), userId); err != nil {
		handleError(ctx, 500, err)
		return
	}
	userRepository := repositories.NewUserRepository(DB(ctx))
	user, err := userRepository.FindByIdUser(userId)
	if err != nil {
		handleError(ctx, 500, err)
		return
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionId.String(),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		MaxAge:   3600,
	}
	res := response.NewUserResponse(user)
	ctx.SetSameSite(http.SameSiteNoneMode)
	http.SetCookie(ctx.Writer, cookie)
	ctx.JSON(200, res)
}
