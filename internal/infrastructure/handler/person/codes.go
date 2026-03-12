package person

// Códigos padronizados da API.
const (
	CodeInvalidBody = "PERSON_INVALID_BODY"

	CodeNameRequired = "PERSON_NAME_REQUIRED"
	CodeNameTooShort = "PERSON_NAME_TOO_SHORT"

	CodeEmailRequired = "PERSON_EMAIL_REQUIRED"
	CodeEmailTooShort = "PERSON_EMAIL_TOO_SHORT"
	CodeEmailInvalid  = "PERSON_EMAIL_INVALID"

	CodePasswordRequired = "PERSON_PASSWORD_REQUIRED"
	CodePasswordTooShort = "PERSON_PASSWORD_TOO_SHORT"

	CodeInvalidCredentials = "PERSON_INVALID_CREDENTIALS"
	CodePersonDuplicated   = "PERSON_DUPLICATED"
	CodePersonNotFound     = "PERSON_NOT_FOUND"
	CodeInternalError      = "PERSON_INTERNAL_ERROR"

	CodeRegistred     = "PERSON_REGISTRED"
	CodeAuthenticated = "PERSON_AUTHENTICATED"
)
