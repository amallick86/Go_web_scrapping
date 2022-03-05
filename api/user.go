package api

import (
	db "Go_web_scrapping/db/sqlc"
	"Go_web_scrapping/util"
	"database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// create user account  request
type createAccountReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

// create user account  response
type createAccountRes struct {
	Id        int32     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at" `
}

//CreateUser handles request for user creation
// @Summary add a new user
// @Tags User
// @ID CreateUser
// @Accept json
// @Produce json
// @Param data body createAccountReq true "create user"
// @Success 201 {object} createAccountRes
// @Failure 400 {object} Err
// @Failure 500 {object} Err
// @Router /users [post]
func (server *Server) createUser(ctx *gin.Context) {
	var req createAccountReq
	var res createAccountRes
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse((err)))
		return
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Username: req.Username,
		Password: hashedPassword,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res = createAccountRes{
		Id:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
	ctx.JSON(http.StatusCreated, res)
}

// login  request
type loginReq struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

func newUserResponse(user db.User) createAccountRes {
	return createAccountRes{
		Id:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	}
}

// login response
type loginRes struct {
	User  createAccountRes `json:"user"`
	Token string           `json:"token"`
}

//Login handles request for user creation
// @Summary login
// @Tags User
// @ID Login
// @Accept json
// @Produce json
// @Param data body loginReq true "Login request"
// @Success 200 {object} loginRes
// @Failure 400 {object} Err
// @Failure 500 {object} Err
// @Router /login [post]
func (server *Server) login(ctx *gin.Context) {
	var req loginReq
	var res loginRes
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse((err)))
		return
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken, err := server.tokenMaker.CreateToken(
		user.Username,
		user.ID,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res = loginRes{
		User:  newUserResponse(user),
		Token: accessToken,
	}
	ctx.JSON(http.StatusOK, res)
}
