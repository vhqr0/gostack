package host

type (
	InvalidIfaceAdaptorTyp struct{ Typ string }
	RequiredIfaceAdaptorArg struct{ Arg string }
)

func (err *InvalidIfaceAdaptorTyp) Error() string {
	return "invalid iface adaptor typ: " + err.Typ
}

func (err *RequiredIfaceAdaptorArg) Error() string {
	return "required iface adaptor arg: " + err.Arg
}
