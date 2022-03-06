package api

type postAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=BRL USD EUR"`
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type getAllAccountsRequest struct {
	Offset int32 `form:"offset" binding:"required,min=1"`
	Limit  int32 `form:"limit" binding:"required,min=5,max=10"`
}
