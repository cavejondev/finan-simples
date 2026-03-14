package transaction

// Códigos padronizados da API.
const (

	// REQUEST
	CodeInvalidBody          = "TRANSACTION_INVALID_BODY"
	CodeInvalidTransactionID = "TRANSACTION_INVALID_ID"

	// VALIDATION
	CodeInvalidAmount       = "TRANSACTION_INVALID_AMOUNT"
	CodeAccountRequired     = "TRANSACTION_ACCOUNT_REQUIRED"
	CodeSubcategoryRequired = "TRANSACTION_SUBCATEGORY_REQUIRED"

	// NOT FOUND
	CodeTransactionNotFound = "TRANSACTION_NOT_FOUND"
	CodeAccountNotFound     = "TRANSACTION_ACCOUNT_NOT_FOUND"
	CodeSubcategoryNotFound = "TRANSACTION_SUBCATEGORY_NOT_FOUND"
	CodeCategoryNotFound    = "TRANSACTION_CATEGORY_NOT_FOUND"

	// BUSINESS RULES
	CodeTransferAccountSame = "TRANSACTION_TRANSFER_ACCOUNT_SAME"

	// AUTH
	CodeUnauthorized = "TRANSACTION_UNAUTHORIZED"

	// INTERNAL
	CodeInternalError = "TRANSACTION_INTERNAL_ERROR"

	// SUCCESS
	CodeTransactionCreated = "TRANSACTION_CREATED"
	CodeTransactionsListed = "TRANSACTIONS_LISTED"
	CodeTransactionFound   = "TRANSACTION_FOUND"
)
