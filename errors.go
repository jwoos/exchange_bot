package main


type GenericError struct {
	message string
}


func (ge *GenericError) Error() string {
	return ge.message
}
