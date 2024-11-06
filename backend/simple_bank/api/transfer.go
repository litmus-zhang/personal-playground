package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/litmus-zhang/simple_bank/bank"
)

type transferRequest struct {
	FromAccount int64  `json:"from" binding:"required,min=1"`
	ToAccount   int64  `json:"to" binding:"required,min=1"`
	Amount      int64  `json:"amount" binding:"required,gt=0"`
	Currency    string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if server.validAccount(ctx, req.FromAccount, req.Currency) == false {
		return
	}
	if server.validAccount(ctx, req.ToAccount, req.Currency) == false {
		return
	}

	arg := bank.TransferTxParams{
		FromAccountID: req.FromAccount,
		ToAccountID:   req.ToAccount,
		Amount:        req.Amount,
	}
	res, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, res)

}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account %d currency mismatch: %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return account.Currency == currency

}
