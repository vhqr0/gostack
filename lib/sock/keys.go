package sock

func SockCreatorKey(family, typ uint32) uint64 {
	return (uint64(family) << 32) + uint64(typ)
}
