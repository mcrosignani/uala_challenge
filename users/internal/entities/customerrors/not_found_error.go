package customerrors

var NotFoundErr *NotFoundError

type NotFoundError struct{}

func (e *NotFoundError) Error() string {
	return "id not found"
}
