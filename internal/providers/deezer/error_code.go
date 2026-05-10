package deezer

import (
	"fmt"
)

type ErrorCode int

const (
	QuoteExceededCode      ErrorCode = 4
	ItemsLimitExceededCode ErrorCode = 100
	PermissionCode         ErrorCode = 200
	TokenInvalidCode       ErrorCode = 300
	ParameterInvalidCode   ErrorCode = 500
	ParameterMissingCode   ErrorCode = 501
	QueryInvalidCode       ErrorCode = 600
	ServiceBusyCode        ErrorCode = 700
	DataNotFoundCode       ErrorCode = 800
	AccountNotAllowedCode  ErrorCode = 901
)

var errorCodeNames = map[ErrorCode]string{
	QuoteExceededCode:      "QuoteExceeded",
	ItemsLimitExceededCode: "ItemsLimitExceeded",
	PermissionCode:         "Permission",
	TokenInvalidCode:       "TokenInvalid",
	ParameterInvalidCode:   "ParameterInvalid",
	ParameterMissingCode:   "ParameterMissing",
	QueryInvalidCode:       "QueryInvalid",
	ServiceBusyCode:        "ServiceBusy",
	DataNotFoundCode:       "DataNotFound",
	AccountNotAllowedCode:  "AccountNotAllowed",
}

func (e ErrorCode) String() string {
	if name, ok := errorCodeNames[e]; ok {
		return name
	}

	return fmt.Sprintf("Unknown(%d)", e)
}

func (e ErrorCode) IsKnown() bool {
	_, ok := errorCodeNames[e]
	return ok
}
