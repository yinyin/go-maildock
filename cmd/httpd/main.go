package main

import (
	"github.com/yinyin/go-maildock/cmd"
	"log"
	"net/http"
)

func main() {
	cfg, err := cmd.LoadConfigurationWithFlags()
	if nil != err {
		log.Printf("cannot load configuration: %v", err)
		return
	}
	handler := newMailDockDisplayHandlerFromConfiguration(cfg)
	if nil != err {
		log.Fatalf("failed on setting up mail dock display serving handler: %v", err)
		return
	}
	defer handler.Close()
	err = http.ListenAndServe(cfg.HTTPListenOn, handler)
	log.Fatalf("result of http.ListenAndServe(): %v", err)
}
