package entity

type ErrUnauthorized struct{}

func (e *ErrUnauthorized) Error() string {
	return "unauthorized"
}

type ErrPermissionDenied struct{}

func (e *ErrPermissionDenied) Error() string {
	return "permission denied"
}
