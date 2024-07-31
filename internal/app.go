package app

import (
	L "kblswitch/internal/log"
	S "kblswitch/internal/svc"
	T "kblswitch/internal/types"
)

type App struct {
	log T.ILog
	svc T.ISvc
}

func NewApp() *App {
	log := L.NewLogFprintf("kblswitch", "TRACE")
	svc := S.NewKBLSwitch(log)
	return &App{
		log: log,
		svc: svc,
	}
}

func (a *App) KeepAlive() {
	a.svc.KeepAlive()
}

func (a *App) Start() func(err error) {
	a.svc.Start()
	a.log.LogInfo("started")
	return func(err error) {
		a.svc.Stop()
		if err != nil {
			a.log.LogError(err, "stoped with error")
		} else {
			a.log.LogInfo("stoped")
		}
	}
}
