package api

import (
	db "GoBankProject/db/sqlc"
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createAccountRequest struct {
	Owner string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD KES UGSH TSH EUR"`
}
var status int
func (s Server) createAccount(ctx *gin.Context)  {
	var req createAccountRequest
	if err := ctx.ShouldBindJSON(&req); err != nil{
		status = http.StatusBadRequest
		ctx.JSON(status, errorResponse(status, err))
		return
	}
	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Balance:  0,
		Currency: req.Currency,
	}
	account, err := s.store.CreateAccount(ctx, arg); if err != nil{
		status = http.StatusInternalServerError
		ctx.JSON(status, errorResponse(status, err))
		return
	}
	ctx.JSON(http.StatusOK, account)

}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}
func (s *Server) getAccount(ctx *gin.Context) {
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil{
		status = http.StatusBadRequest
		ctx.JSON(status, errorResponse(status, err))
		return
	}
	account, err := s.store.GetAccount(ctx, req.ID); if err != nil {
		if err == sql.ErrNoRows{
			status = http.StatusNotFound
			ctx.JSON(status, errorResponse(status, err))
			return
		}
		status = http.StatusInternalServerError
		ctx.JSON(status, errorResponse(status, err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}
type listAccountRequest struct {
	PageID int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}
func (s Server) listAccount(ctx *gin.Context)  {
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req); err != nil{
		status = http.StatusBadRequest
		ctx.JSON(status, errorResponse(status, err))
		return
	}
	arg := db.ListAccountParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	account, err := s.store.ListAccounts(ctx, arg); if err != nil {
		if err == sql.ErrNoRows{
			status = http.StatusNotFound
			ctx.JSON(status, errorResponse(status, err))
			return
		}
		status = http.StatusInternalServerError
		ctx.JSON(status, errorResponse(status, err))
		return
	}
	ctx.JSON(http.StatusOK, account)
}