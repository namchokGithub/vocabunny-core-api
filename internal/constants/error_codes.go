package constants

const (
	// Success
	CodeSuccess      = "2000"
	CodeCreated      = "2001"
	CodeUpdated      = "2002"
	CodeDeleted      = "2003"
	CodePublished    = "2004"
	CodeLoginSuccess = "2005"

	// Validation
	CodeValidationFailed   = "4001"
	CodeInvalidRequestBody = "4002"
	CodeInvalidPagination  = "4003"
	CodeInvalidQueryParam  = "4004"

	// Auth/Identity
	CodeInvalidToken       = "4101"
	CodeExpiredToken       = "4102"
	CodeInvalidCredentials = "4103"
	CodeInvalidScope       = "4104"
	CodePermissionDenied   = "4301"

	// Content
	CodeSectionNotFound     = "4201"
	CodeLessonNotFound      = "4202"
	CodeUnitNotFound        = "4203"
	CodeQuestionSetNotFound = "4204"
	CodeQuestionNotFound    = "4205"
	CodeInvalidInclude      = "4206"
	CodeDuplicateSlug       = "4207"
	CodeInvalidPublishState = "4208"

	// Conflict
	CodeDuplicateResource  = "4901"
	CodeConcurrentConflict = "4902"

	// Internal/System
	CodeInternalError       = "5001"
	CodeDBTransactionFailed = "5002"
	CodeUnexpectedRepo      = "5003"
	CodeExternalUnavailable = "5004"
)
