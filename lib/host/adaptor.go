package host

import (
	"io"
)

var (
	adaptorCreatorMap map[string]AdaptorCreator = make(map[string]AdaptorCreator)
)

type (
	Adaptor        io.ReadWriteCloser
	AdaptorCreator func(map[string]string) (Adaptor, error)
)

func RegisterAdaptor(typ string, creator AdaptorCreator) {
	adaptorCreatorMap[typ] = creator
}

func NewAdaptor(typ string, args map[string]string) (Adaptor, error) {
	if creator, ok := adaptorCreatorMap[typ]; ok {
		return creator(args)
	}
	return nil, &InvalidAdaptorTyp{Typ: typ}
}

func AdaptorRequireArg(arg string, args map[string]string) (string, error) {
	if v, ok := args[arg]; ok {
		return v, nil
	}
	return "", &RequiredAdaptorArg{Arg: arg}
}
