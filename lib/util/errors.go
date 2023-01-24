package util

type FeatUnsupported struct{ Feat string }

func (err *FeatUnsupported) Error() string {
	return "feat unsupported: " + err.Feat
}

type IfaceNameOverflow struct{ Name string }

func (err *IfaceNameOverflow) Error() string {
	return "iface name overflow: " + err.Name
}

type InvalidIPStr struct{ Str string }

func (err *InvalidIPStr) Error() string {
	return "invalid ip str: " + err.Str
}
