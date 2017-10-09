package main

import (
	"context"
	maildock "github.com/yinyin/go-maildock"
	"github.com/yinyin/go-maildock/cmd"
	"github.com/yinyin/go-maildock/database"
	utilhttphandlers "github.com/yinyin/go-util-http-handlers"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type searchResult struct {
	QueryAddress string
	Records      []*maildock.MailRecord
}

type maildockDisplayHandler struct {
	dbCfg                database.Configuration
	staticContentHandler http.Handler
}

func (d *maildockDisplayHandler) handleSearch(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	c := strings.SplitN(path, "/", 3)
	if len(c) < 3 {
		http.NotFound(w, r)
		return
	}
	queryRecipientAddress := c[2]
	dbconn, err := d.dbCfg.OpenConnection()
	if nil != err {
		http.Error(w, "failed on connecting to database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dbconn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	records, err := dbconn.SearchForRecipient(ctx, queryRecipientAddress)
	var result = struct {
		QueryAddress string                 `json:"query_address"`
		Records      []*maildock.MailRecord `json:"records"`
	}{
		QueryAddress: queryRecipientAddress,
		Records:      records}
	utilhttphandlers.JSONResponse(w, result)
}

func (d *maildockDisplayHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.HasPrefix(path, "/search") {
		d.handleSearch(w, r)
	} else if nil != d.staticContentHandler {
		d.staticContentHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (d *maildockDisplayHandler) Close() {
	if nil != d.staticContentHandler {
		if closer, ok := d.staticContentHandler.(io.Closer); ok {
			closer.Close()
			log.Printf("closed static content handler")
		}
	}
}

func makeStaticContentHandler(cfg *cmd.Configuration) (handler http.Handler) {
	contentCfg := cfg.HTTPContent
	if "" != contentCfg.Path {
		if zipContentHandler, err := utilhttphandlers.NewZipArchiveContentServer(contentCfg.Path, contentCfg.Prefix, "index.html"); nil != err {
			log.Printf("failed on open zip archive as content source: %v", err)
		} else {
			log.Printf("use content from zip archive: (path=%v, prefix=%v)", contentCfg.Path, contentCfg.Prefix)
			return zipContentHandler
		}
	}
	if "" != contentCfg.ProxyTo {
		if targetURL, err := url.Parse(contentCfg.ProxyTo); nil != err {
			log.Printf("failed on parsing proxy target URL: %v", err)
		} else {
			log.Printf("use content from URL: %v", contentCfg.ProxyTo)
			return httputil.NewSingleHostReverseProxy(targetURL)
		}
	}
	return nil
}

func newMailDockDisplayHandlerFromConfiguration(cfg *cmd.Configuration) (h *maildockDisplayHandler) {
	return &maildockDisplayHandler{
		dbCfg:                cfg.Database.Config,
		staticContentHandler: makeStaticContentHandler(cfg),
	}
}
