package internal

type DontUniqueError struct{}

func (d *DontUniqueError) Error() string {
	return "Dont unique content"
}
