package types

type ISvc interface {
	Start()
	KeepAlive()
	Stop()
}
