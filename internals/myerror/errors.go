package myerror

type MultiLoggerInterface interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Sync() error
}

type MyErrors struct {
	Text string
}

func (e *MyErrors) Error() string {
	return e.Text
}

type CheckError struct {
	RequestId int
	Logger    MultiLoggerInterface
}

type ResultError struct {
	Status  int    `json:"status"`
	Explain string `json:"explain,omitempty"`
}

// Error of server
const (
	ErrDB              = "database is not responding"
	ErrAtoi            = "func Atoi convert string in int"
	IntNil             = 0
	ErrNotStringAndInt = "expected type string or int"
	ErrUnmarshal       = "unmarshal json"
	ErrMarshal         = "marshaling in json"
	ErrCheck           = "err check"
	ErrEncode          = "Encode"
	ErrInternal        = "err internal"
)

// Error of main
const (
	MCreateDBNotConnect           = "db not connect"
	MCreateDBTransactionNotCreate = "transaction setup not create"
	MCreateDBFileNotFound         = "createtables.sql not found"
	MCreateDBFileNotCreate        = "table not create"
	MCreateDBNotCommit            = "transaction setup not commit"
)

// Error of LinkShort
const (
	LSHCreateLinkShortTransactionNotCreate = "transaction Create Link Short not create CreateLinkShortPostgres"
	LSHCreateLinkShortNotInsert            = "Link short not insert CreateLinkShortPostgres"
	LSHCreateLinkShortNotCommit            = "Link short not commit CreateLinkShortPostgres"
	LSHCreateLinkShortNotInsertUniqueDB    = "ERROR: duplicate key value violates unique constraint \"link_link_key\" (SQLSTATE 23505)"
	LSHCreateLinkShortNotInsertUnique      = "link is not unique CreateLinkShortPostgres"
	LSHCreateLinkShortAppNotGenerate       = "link is not generate CreateLinkShortPostgres"

	LSHTakeLinkShortTransactionNotCreate = "transaction Take Link Short not create"
	LSHTakeLinkShortNotFound             = "link full not found"
	LSHTakeLinkShortNotScan              = "link full not scan"
	LSHTakeLinkShortNotCommit            = "Link full not commit"

	LSHTakeLinkShortNotFoundRedis          = "link full not found"
	LSHCreateLinkShortExistsRedis          = "link full exists"
	LSHCreateLinkShortNotSetFullLinkRedis  = "link is not set fullLink"
	LSHCreateLinkShortNotSetShortLinkRedis = "link is not set shortLink"
	LSHCreateLinkShortNilConn              = "connect is nil"
	LSHTakeLinkShortNilConn                = "connect is nil"
)
