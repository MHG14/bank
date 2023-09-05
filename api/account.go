package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/mhg14/simplebank/db/sqlc"
	"github.com/mhg14/simplebank/token"
)

// Create account handler
type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (s *Server) createAccountHandler(ctx *gin.Context) {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}

	account, err := s.store.CreateAccount(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(pqErr))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

///////////////////////////////////////////////////////////

// Get account handler
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccountHandler(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := s.store.GetAccount(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if account.Owner != authPayload.Username {
		err := errors.New("account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}

///////////////////////////////////////////////////////////

// List accounts handler
type listAccountsRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) listAccountsHandler(ctx *gin.Context) {
	var req listAccountsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
		Owner: authPayload.Username,
	}

	accounts, err := s.store.ListAccounts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

///////////////////////////////////////////////////////////

// Update accounts handler
type UpdateAccountParam struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type UpdateAccountBody struct {
	Balance int64 `json:"balance" binding:"required,min=0"`
}

func (s *Server) updateAccountHandler(ctx *gin.Context) {
	var param UpdateAccountParam
	var body UpdateAccountBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateAccountParams{
		ID:      param.ID,
		Balance: body.Balance,
	}

	account, err := s.store.UpdateAccount(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

///////////////////////////////////////////////////////////

// Delete account handler
type deleteAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) deleteAccountHandler(ctx *gin.Context) {
	var req deleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := s.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	rows, err := result.RowsAffected()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if rows == 0 {
		ctx.JSON(http.StatusNotFound, errorResponse(errors.New(sql.ErrNoRows.Error())))
		return
	}
	ctx.JSON(http.StatusOK, "deleted successfully")
}
