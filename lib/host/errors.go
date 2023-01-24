package host

type InvalidIfaceName struct{ Name string }

func (err *InvalidIfaceName) Error() string { return "invalid iface name: " + err.Name }

type (
	InvalidAdaptorTyp  struct{ Typ string }
	RequiredAdaptorArg struct{ Arg string }
)

func (err *InvalidAdaptorTyp) Error() string  { return "invalid adaptor typ: " + err.Typ }
func (err *RequiredAdaptorArg) Error() string { return "required adaptor arg: " + err.Arg }
