package sock

type (
	SockCreator interface {
		NewSock(family, typ uint32) (Sock, error)
	}
	SockFactory struct{ creatorMap map[uint64]SockCreator }
)

func (f *SockFactory) Register(family, typ uint32, creator SockCreator) {
	f.creatorMap[SockCreatorKey(family, typ)] = creator
}

func (f *SockFactory) NewSock(family, typ uint32) (Sock, error) {
	if creator, ok := f.creatorMap[SockCreatorKey(family, typ)]; ok {
		return creator.NewSock(family, typ)
	}
	return nil, &InvalidSockArgs{Family: family, Typ: typ}
}

func NewSockFactory() *SockFactory {
	f := &SockFactory{
		creatorMap: make(map[uint64]SockCreator),
	}
	return f
}
