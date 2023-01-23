package host

import (
	"github.com/vhqr0/gostack/lib/util"
)

func init() {
	RegisterIfaceAdaptor("tap", NewTAPIfaceAdaptor)
}

func NewTAPIfaceAdaptor(args map[string]string) (IfaceAdaptor, error) {
	if name, err := IfaceAdaptorRequireArg("Name", args); err != nil {
		return nil, err
	} else {
		return util.OpenTAP(name)
	}
}
