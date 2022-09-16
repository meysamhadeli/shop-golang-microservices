package constants

import "time"

const (
	ConfigPath           = "CONFIG_PATH"
	Json                 = "json"
	MaxHeaderBytes       = 1 << 20
	StackSize            = 1 << 10 // 1 KB
	BodyLimit            = "2M"
	ReadTimeout          = 15 * time.Second
	WriteTimeout         = 15 * time.Second
	GzipLevel            = 5
	WaitShotDownDuration = 3 * time.Second
)

const (
	ErrBadRequest          = "Bad Request"
	ErrWrongCredentials    = "Wrong Credentials"
	ErrNotFound            = "Not Found"
	ErrUnauthorized        = "Unauthorized"
	ErrForbidden           = "Forbidden"
	ErrInternalServerError = "Internal Server Error"
	ErrDomain              = "Domain Model Error"
	ErrApplication         = "Application Service Error"
	ErrApi                 = "Api Error"
)
