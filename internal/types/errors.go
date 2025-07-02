package types

var (
	SuccessCode = 0
)

type ErrorCode int

const (
	ErrSystemBusy     ErrorCode = 10
	ErrServiceTimeout ErrorCode = 11
	ErrDatabaseError  ErrorCode = 12
	ErrCacheError     ErrorCode = 13
	ErrInvalidParam   ErrorCode = 14
	ErrInternalError  ErrorCode = 15

	// 认证授权 (1000-1099)
	ErrUnauthorized     ErrorCode = 1000
	ErrInvalidSignature ErrorCode = 1001
	ErrExpiredToken     ErrorCode = 1002
	ErrPermissionDenied ErrorCode = 1003
	ErrInvalidAuth      ErrorCode = 1004

	// 账户资产 (1100-1199)
	ErrAccountNotFound     ErrorCode = 1100
	ErrInsufficientBalance ErrorCode = 1101
	ErrFrozenAccount       ErrorCode = 1102
	ErrBalanceNotEnough    ErrorCode = 1103
	ErrInvalidAddress      ErrorCode = 1104
	ErrWithdrawLocked      ErrorCode = 1105

	// 交易订单 (1200-1299)
	ErrInvalidOrder   ErrorCode = 1200
	ErrOrderNotFound  ErrorCode = 1201
	ErrOrderExists    ErrorCode = 1202
	ErrMarketClosed   ErrorCode = 1203
	ErrPriceInvalid   ErrorCode = 1204
	ErrAmountTooSmall ErrorCode = 1205

	// 风控系统 (1300-1399)
	ErrWithdrawalLimit  ErrorCode = 1300
	ErrAntiPhishingCode ErrorCode = 1301
	ErrTradeBan         ErrorCode = 1302
)

type ErrorMeta struct {
	Code           ErrorCode
	DefaultMessage string
}

var errRegistry = map[ErrorCode]*ErrorMeta{
	// 系统类
	ErrSystemBusy:     {Code: ErrSystemBusy, DefaultMessage: "System busy"},
	ErrServiceTimeout: {Code: ErrServiceTimeout, DefaultMessage: "Service timeout"},
	ErrDatabaseError:  {Code: ErrDatabaseError, DefaultMessage: "Database error"},
	ErrCacheError:     {Code: ErrCacheError, DefaultMessage: "Cache error"},
	ErrInvalidParam:   {Code: ErrInvalidParam, DefaultMessage: "Invalid param"},
	ErrInternalError:  {Code: ErrInternalError, DefaultMessage: "Internal error"},

	// 认证授权类
	ErrUnauthorized:     {Code: ErrUnauthorized, DefaultMessage: "Unauthorized"},
	ErrInvalidSignature: {Code: ErrInvalidSignature, DefaultMessage: "Invalid signature"},
	ErrExpiredToken:     {Code: ErrExpiredToken, DefaultMessage: "Expired token"},
	ErrPermissionDenied: {Code: ErrPermissionDenied, DefaultMessage: "Permission denied"},
	ErrInvalidAuth:      {Code: ErrInvalidAuth, DefaultMessage: "Invalid authentication"},

	// 账户资产类
	ErrAccountNotFound:  {Code: ErrAccountNotFound, DefaultMessage: "Account not found"},
	ErrBalanceNotEnough: {Code: ErrBalanceNotEnough, DefaultMessage: "Insufficient balance"},
	ErrInvalidAddress:   {Code: ErrInvalidAddress, DefaultMessage: "Invalid address"},
	ErrWithdrawLocked:   {Code: ErrWithdrawLocked, DefaultMessage: "Withdrawal locked"},

	// 交易订单类
	ErrInvalidOrder:   {Code: ErrInvalidOrder, DefaultMessage: "Invalid order"},
	ErrOrderNotFound:  {Code: ErrOrderNotFound, DefaultMessage: "Order not found"},
	ErrOrderExists:    {Code: ErrOrderExists, DefaultMessage: "Order already exists"},
	ErrMarketClosed:   {Code: ErrMarketClosed, DefaultMessage: "Market closed"},
	ErrPriceInvalid:   {Code: ErrPriceInvalid, DefaultMessage: "Invalid price"},
	ErrAmountTooSmall: {Code: ErrAmountTooSmall, DefaultMessage: "Amount too small"},
}

func RegisterError(code ErrorCode, defaultMessage string) {
	errRegistry[code] = &ErrorMeta{Code: code, DefaultMessage: defaultMessage}
}

func GetErrorMsg(code ErrorCode) string {
	if _, ok := errRegistry[code]; ok {
		return errRegistry[code].DefaultMessage
	}
	return "Unknown error"
}
