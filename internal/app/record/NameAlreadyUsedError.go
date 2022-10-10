package record

type NameAlreadyUsedError struct{}

func (e NameAlreadyUsedError) Error() string {
	return "Name already used"
}
