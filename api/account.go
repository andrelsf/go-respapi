package api

import (
	"database/sql"
	"net/http"

	db "github.com/andrelsf/go-restapi/db/sqlc"
	"github.com/gin-gonic/gin"
)

func (server *Server) postAccount(ctx *gin.Context) {
	var req postAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(http.StatusUnprocessableEntity, err))
		return
	}

	params := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(http.StatusUnprocessableEntity, err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

func (server *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(http.StatusUnprocessableEntity, err))
		return
	}

	account, err := server.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(http.StatusNotFound, err))
			return
		}

		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(http.StatusUnprocessableEntity, err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (server *Server) getAllAccounts(ctx *gin.Context) {
	var req getAllAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(http.StatusUnprocessableEntity, err))
		return
	}

	params := db.ListAccountsParams{
		Limit:  req.Limit,
		Offset: (req.Offset - 1) * req.Limit,
	}

	accounts, err := server.store.ListAccounts(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, errorResponse(http.StatusUnprocessableEntity, err))
		return
	}

	ctx.JSON(http.StatusPartialContent, accounts)
}
