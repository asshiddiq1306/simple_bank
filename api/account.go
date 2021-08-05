package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/asshiddiq1306/simple_bank/db/sql"
	"github.com/asshiddiq1306/simple_bank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createNewAccountReq struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createNewAccountAPI(c *gin.Context) {
	var req createNewAccountReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.AuthPay)
	arg := db.CreateNewAccountArgs{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := server.store.CreateNewAccount(c, arg)
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

	c.JSON(http.StatusOK, account)
}

type getAccountByIDReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccountByIDAPI(c *gin.Context) {
	var req getAccountByIDReq
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccountByID(c, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.AuthPay)

	if account.Owner != authPayload.Username {
		err := errors.New("this account doesn't belongs to auth user")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, account)
}

type getListAccountsReq struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) getListAccountsAPI(c *gin.Context) {
	var req getListAccountsReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.AuthPay)

	arg := db.GetListAccountsArgs{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := server.store.GetListAccounts(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, accounts)

}
