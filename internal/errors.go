package utils

import "errors"

var PgDuplicateKeyMsg = "duplicate key value violates unique constraint"
var PgNoRowsMsg = "no rows in result set"

const (
	// Popular errors
	InternalServerError ErrorType = 500
	BadRequest                    = 400
	Unauthorized                  = 401
	NotFound                      = 404
	Conflict                      = 409
	Forbidden                     = 403
	NotImplemented                = 501
	NotAcceptable                 = 406
	MethodNotAllowed              = 405

	// Custom errors
	MailRequired                   = 1
	PasswordRequired               = 2
	PasswordNotMatch               = 3
	CodeRequired                   = 4
	RefreshTokenRequired           = 5
	AccessTokenRequired            = 6
	OldPasswordRequired            = 7
	NewPasswordRequired            = 8
	ConfirmPasswordRequired        = 9
	InvalidCode                    = 10
	ExpiredCode                    = 11
	TooManyRequests                = 12
	PasswordIncorrect              = 13
	ChoiceOtherPassword            = 14
	EmailNotRegistered             = 15
	QuicklyRequest                 = 16
	ValidationJSONFailure          = 17
	ValidationTokenFailure         = 18
	ExistUser                      = 19
	CodeInvalid                    = 20
	UsernameRequired               = 21
	ExistUserName                  = 22
	MatchCodeIsExpired             = 23
	MatchCodeIsIncorrect           = 24
	UserAlreadyMatched             = 25
	MatchCodeIsNotExactly          = 26
	UserNotMatch                   = 27
	TitleRequired                  = 28
	BodyRequired                   = 29
	ValueTooLong                   = 30
	UserInRelationship             = 31
	PlayerIdRequired               = 32
	UserCannotMatchWithHimself     = 33
	KeyStringRequired              = 34
	UserIDRequired                 = 35
	AccountIsNotNeedToCancelDelete = 36
	PassCodeRequired               = 37
	UserDeleted                    = 38
)

func (e ErrorResponse) Error() string {
	switch e.ErrorType {
	case InternalServerError:
		return "Internal Server Error"
	case BadRequest:
		return "Bad Request"
	case Unauthorized:
		return "Unauthorized"
	case NotFound:
		return "Not Found"
	case Conflict:
		return "Conflict"
	case Forbidden:
		return "Forbidden"
	case NotImplemented:
		return "Not Implemented"
	case NotAcceptable:
		return "Not Acceptable"
	case MethodNotAllowed:
		return "Method Not Allowed"
	case MailRequired:
		return "Mail Required"
	case PasswordRequired:
		return "Password Required"
	case PasswordNotMatch:
		return "Password Not Match"
	case CodeRequired:
		return "Code Required"
	case RefreshTokenRequired:
		return "Refresh Token Required"
	case AccessTokenRequired:
		return "Access Token Required"
	case OldPasswordRequired:
		return "Old Password Required"
	case NewPasswordRequired:
		return "New Password Required"
	case ConfirmPasswordRequired:
		return "Confirm Password Required"
	case InvalidCode:
		return "Invalid Code"
	case ExpiredCode:
		return "Expired Code"
	case TooManyRequests:
		return "You have been locked out for security reasons. Please try again later."
	case PasswordIncorrect:
		return "Password Incorrect"
	case ChoiceOtherPassword:
		return "Password has been used. Please choose another password."
	case EmailNotRegistered:
		return "Email not registered"
	case QuicklyRequest:
		return "Too many requests. Please try again later."
	case ValidationJSONFailure:
		return "validation of verification data json failed"
	case ValidationTokenFailure:
		return "validation of verification token failed"
	case ExistUser:
		return "User already exists"
	case CodeInvalid:
		return "Verify code invalid"
	case UsernameRequired:
		return "Username required"
	case ExistUserName:
		return "Username already exists"
	case MatchCodeIsExpired:
		return "match code is expired"
	case MatchCodeIsIncorrect:
		return "match code is incorrect"
	case UserAlreadyMatched:
		return "user already matched"
	case MatchCodeIsNotExactly:
		return "match code is not exactly"
	case UserNotMatch:
		return "user not match"
	case TitleRequired:
		return "title required"
	case BodyRequired:
		return "body required"
	case ValueTooLong:
		return "value too long"
	case UserInRelationship:
		return "user in a relationship"
	case PlayerIdRequired:
		return "player id required"
	case UserCannotMatchWithHimself:
		return "user cannot match with himself"
	case KeyStringRequired:
		return "key string required"
	case UserIDRequired:
		return "user id required"
	case AccountIsNotNeedToCancelDelete:
		return "account is not need to cancel delete"
	case PassCodeRequired:
		return "pass code required"
	case UserDeleted:
		return "user is on the deletion schedule."
	default:
		return "Unknown Error"
	}
}

type CustomErrorWrapper struct {
	Message string `json:"message"` // Human-readable message for clients
	Code    int    `json:"-"`       // HTTP Status code. We use `-` to skip json marshaling.
	Err     error  `json:"-"`       // The original error. Same reason as above.
}

func NewErrorWrapper(code int, err error, message string) CustomErrorWrapper {
	return CustomErrorWrapper{
		Message: message,
		Code:    code,
		Err:     err,
	}
}

// Returns Message if Err is nil. You can handle custom implementation of your own.
func (err CustomErrorWrapper) Error() string {
	// guard against panics
	if err.Err != nil {
		return err.Err.Error()
	}
	return err.Message
}

func (err CustomErrorWrapper) Unwrap() error {
	return err.Err // Returns inner error
}

// Dig Returns the innermost CustomErrorWrapper
func (err CustomErrorWrapper) Dig() CustomErrorWrapper {
	var ew CustomErrorWrapper
	if errors.As(err.Err, &ew) {
		// Recursively digs until wrapper error is not in which case it will stop
		return ew.Dig()
	}
	return err
}

type ErrorType int

type ErrorResponse struct {
	ErrorType ErrorType
}

func NewErrorResponse(errorType ErrorType) ErrorResponse {
	return ErrorResponse{
		ErrorType: errorType,
	}
}
