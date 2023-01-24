package host

import (
	"github.com/vhqr0/gostack/lib/util"
)

func init() {
	RegisterAdaptor("tap", NewTAPAdaptor)
}

func NewTAPAdaptor(args map[string]string) (Adaptor, error) {
	if name, err := AdaptorRequireArg("Name", args); err != nil {
		return nil, err
	} else {
		return util.OpenTAP(name)
	}
}
