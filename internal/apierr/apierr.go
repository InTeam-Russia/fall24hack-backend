package apierr

type ApiError struct {
	Error string `json:"error"`
}

func New(error string) ApiError {
	return ApiError{Error: error}
}

var InternalServer = ApiError{
	Error: "INTERNAL_SERVER",
}

var InvalidJSON = ApiError{
	Error: "INVALID_JSON",
}

var NotFound = ApiError{
	Error: "NOT_FOUND",
}

var WrongCredentials = ApiError{
	Error: "WRONG_CREDENTIALS",
}

var CookieNotExists = ApiError{
	Error: "COOKIE_NOT_EXISTS",
}

var SessionNotFound = ApiError{
	Error: "SESSION_NOT_FOUND",
}

var SessionExpired = ApiError{
	Error: "SESSION_EXPIRED",
}

var UserNotFound = ApiError{
	Error: "USER_NOT_FOUND",
}

var InvalidPageIndex = ApiError{
	Error: "INVALID_PAGE_INDEX",
}

var InvalidPageSize = ApiError{
	Error: "INVALID_PAGE_SIZE",
}

var InvalidID = ApiError{
	Error: "INVALID_ID",
}

var InvalidSearchType = ApiError{
	Error: "INVALID_SEARCH_TYPE",
}

var InvalidPollType = ApiError{
	Error: "INVALID_POLL_TYPE",
}

var NoRadioAnswers = ApiError{
	Error: "NO_RADIO_ANSWERS",
}
