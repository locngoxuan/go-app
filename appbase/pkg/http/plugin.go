package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"go-app/appbase/pkg/interfaces"
	"go-app/appbase/plugin"
	glog "log"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func NewHttp(ctx *cli.Context) interfaces.HTTP {
	opt := HttpRouterOption{
		Bind: ctx.String(interfaces.FlagHttpBind),
		Port: ctx.Int(interfaces.FlagHttpPort),
	}
	return &HttpPlugin{
		HttpRouterOption: opt,
	}
}

type FwdToZeroWriter struct {
}

func (fw *FwdToZeroWriter) Write(p []byte) (n int, err error) {
	log.Error().Msg(string(p))
	return len(p), nil
}

type HttpRouterOption struct {
	Port       int
	Bind       string
	TLSEnabled bool
	TLSKey     string
	TLSCert    string
}

type HttpPlugin struct {
	plugin.EmptyPlugin
	HttpRouterOption
	httpServer *http.Server
	router     http.Handler
}

// SetRouter implements interfaces.HTTP.
func (h *HttpPlugin) SetRouter(router http.Handler) {
	h2s := &http2.Server{}
	h.router = h2c.NewHandler(router, h2s)
}

func (h *HttpPlugin) startHttpServerBackground() error {
	var err error
	if h.TLSEnabled {
		//Achieving a Perfect SSL with go: https://blog.bracebin.com/achieving-perfect-ssl-labs-score-with-go
		h.httpServer.TLSConfig = &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
		h.httpServer.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler))
		log.Info().Str("tls_cert", h.TLSCert).Str("tls_key", h.TLSKey).Msg("tls files configuration")
		err = h.httpServer.ListenAndServeTLS(h.TLSCert, h.TLSKey)
	} else {
		err = h.httpServer.ListenAndServe()
	}
	if err != nil {
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}
	return nil
}

func (h *HttpPlugin) Start() error {
	bindingAddr := fmt.Sprintf("%s:%d", h.Bind, h.Port)
	log.Info().Str("address", bindingAddr).Msg("starting http server...")
	h.httpServer = &http.Server{
		Addr:     bindingAddr,
		Handler:  h.router,
		ErrorLog: glog.New(&FwdToZeroWriter{}, "", 0),
	}
	go func() {
		err := h.startHttpServerBackground()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	}()
	return nil
}

func (h *HttpPlugin) Stop() error {
	err := h.httpServer.Shutdown(context.Background())
	if err != nil {
		log.Error().Err(err).Msg("failed to stop http server")
	}
	return h.httpServer.Close()
}
