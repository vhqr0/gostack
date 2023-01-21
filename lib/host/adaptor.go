package host

import (
	"io"
)

var (
	ifaceAdaptorCreatorMap map[string]IfaceAdaptorCreator = make(map[string]IfaceAdaptorCreator)
)

type (
	IfaceAdaptor        io.ReadWriteCloser
	IfaceAdaptorCreator func(map[string]string) (IfaceAdaptor, error)
)

func RegisterIfaceAdaptor(typ string, creator IfaceAdaptorCreator) {
	ifaceAdaptorCreatorMap[typ] = creator
}

func NewIfaceAdaptor(typ string, args map[string]string) (IfaceAdaptor, error) {
	if creator, ok := ifaceAdaptorCreatorMap[typ]; ok {
		return creator(args)
	}
	return nil, &InvalidIfaceAdaptorTyp{Typ: typ}
}

func IfaceAdaptorRequireArg(arg string, args map[string]string) (string, error) {
	if v, ok := args[arg]; ok {
		return v, nil
	}
	return "", &RequiredIfaceAdaptorArg{Arg: arg}
}
