package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/asshiddiq1306/simple_bank/db/sql"
	"github.com/asshiddiq1306/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createNewUserReq struct {
	Username       string `json:"username" binding:"required,alphanum"`
	HashedPassword string `json:"hashed_password" binding:"required,min=6"`
	FullName       string `json:"full_name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
}

type newUserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func userResp(user db.User) newUserResponse {
	return newUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createNewUserAPI(c *gin.Context) {
	var req createNewUserReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashPassword, err := util.HashedPassword(req.HashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateNewUserArgs{
		Username:       req.Username,
		HashedPassword: hashPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateNewUser(c, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := userResp(user)

	c.JSON(http.StatusOK, rsp)

}

type userLoginReq struct {
	Username       string `json:"username" binding:"required,alphanum"`
	HashedPassword string `json:"hashed_password" binding:"required,min=6"`
}

type userLoginResp struct {
	AccessToken string          `json:"access_token"`
	User        newUserResponse `json:"user"`
}

func (server *Server) userLoginAPI(c *gin.Context) {
	var req userLoginReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByUsername(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.HashedPassword, user.HashedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.TokenAccessDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := userLoginResp{
		AccessToken: accessToken,
		User:        userResp(user),
	}

	c.JSON(http.StatusOK, rsp)
}
