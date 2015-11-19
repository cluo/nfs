package nfsmain

import (
	"log"
	"net/http"

	"github.com/Unknwon/goconfig"
	"github.com/mangalaman93/nfs/voip"
)

var (
	apps map[string]AppLine
)

func Start(config *goconfig.ConfigFile) error {
	port, err := config.GetValue("CONTROLLER", "port")
	if err != nil {
		return err
	}

	vl, err := voip.NewVoipLine(config)
	if err != nil {
		return err
	}
	apps = make(map[string]AppLine)
	apps[vl.GetDB()] = vl
	log.Println("[INFO] registered db:", vl.GetDB(), "with VoipLine instance")

	h, err := NewHandler(config, apps)
	if err != nil {
		return err
	}
	go http.ListenAndServe(":"+port, h)
	go vl.Start()
	log.Println("[INFO] listening for data over line protocol on port", port)
	return nil
}

func Stop() {
	for _, app := range apps {
		app.Stop()
	}

	// TODO: for now, we don't know how to stop http
	// listener and we do not want to get into
	// complications of creating listener ourselves
	log.Println("[INFO] exiting control loop")
}
