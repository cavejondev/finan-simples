package account

// Códigos padronizados da API.
const (
	CodeInvalidBody      = "ACCOUNT_INVALID_BODY"
	CodeInvalidAccountID = "ACCOUNT_INVALID_ID"

	CodeNameRequired = "ACCOUNT_NAME_REQUIRED"
	CodeNameTooShort = "ACCOUNT_NAME_TOO_SHORT"

	CodeAccountNotFound = "ACCOUNT_NOT_FOUND"
	CodeAccountClosed   = "ACCOUNT_CLOSED"

	CodeAccountNameDuplicated = "ACCOUNT_NAME_DUPLICATED"

	CodeUnauthorized = "ACCOUNT_UNAUTHORIZED"

	CodeInternalError = "ACCOUNT_INTERNAL_ERROR"

	CodeAccountCreated = "ACCOUNT_CREATED"
	CodeAccountUpdated = "ACCOUNT_UPDATED"
	CodeAccountsListed = "ACCOUNT_LISTED"
	CodeAccountFound   = "ACCOUNT_FOUND"
)
