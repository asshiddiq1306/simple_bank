package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/asshiddiq1306/simple_bank/db/sql"
	"github.com/asshiddiq1306/simple_bank/token"
	"github.com/gin-gonic/gin"
)

type transferTxReq struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) transferTxAPI(c *gin.Context) {
	var req transferTxReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.AuthPay)

	fromAccount, val := server.isValidCurrency(c, req.FromAccountID, req.Currency)
	if !val {
		return
	}

	if fromAccount.Owner != authPayload.Username {
		err := errors.New("this account doesn't belongs to auth user")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, val = server.isValidCurrency(c, req.ToAccountID, req.Currency)
	if !val {
		return
	}

	arg := db.TransferTxArg{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	transfer, err := server.store.TransferTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, transfer)

}

func (server *Server) isValidCurrency(c *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account1, err := server.store.GetAccountByID(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return account1, false
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return account1, false
	}

	if account1.Currency != currency {
		err := fmt.Errorf("currency mismatch %s X %s", account1.Currency, currency)
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return account1, false
	}
	return account1, true

}
