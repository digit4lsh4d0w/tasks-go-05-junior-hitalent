package database

type DatabaseError struct {
	DSN       string
	Operation string
	Err       error
}

func (e *DatabaseError) Error() string {
	return "database error during " + e.Operation + ": " + e.Err.Error()
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}
