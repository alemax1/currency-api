package domain

const (
	ErrNothingUpdated       = "nothing updated"
	ErrNothingFound         = "nothing found"
	ErrInvalidCurrencyTypes = "cannot change currencies with equal types"
	ErrDuplicateValue       = "value already exists"
	ErrValueCannotBeZero    = "value cannot be zero"
)

type ErrType string

const (
	Client ErrType = "client"
)

type ServiceError struct {
	Message string
	Type    ErrType
}

func NewServiceError(msg string, t ErrType) error {
	return &ServiceError{
		Message: msg,
		Type:    t,
	}
}

func (s ServiceError) Error() string {
	return s.Message
}
