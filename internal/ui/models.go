package ui

type FormField struct {
	Value string
	Err   error
}
type Form map[string]FormField
