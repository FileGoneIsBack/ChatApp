package server

import (
	"login/core/models/antiflood"
	"login/core/models/functions"
	"login/core/models/validation"
	"login/core/models"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"errors"
	"github.com/gorilla/mux"
)



type Server struct {
	mux           http.Server
	router        *mux.Router
	routes        map[string]*Route
	config        *Config
	TemplateCache *functions.TemplateCache 
}

func NewServer(config *Config) *Server {
	return &Server{
		mux: http.Server{
			Addr:         config.Addr,
			Handler:      nil,
			TLSConfig: &tls.Config{
				MinVersion: tls.VersionTLS13,
				CipherSuites: []uint16{
					tls.TLS_AES_128_GCM_SHA256,
					tls.TLS_AES_256_GCM_SHA384,
					tls.TLS_CHACHA20_POLY1305_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					// Add other supported cipher suites
				},
			},
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		},
		router:        mux.NewRouter(),
		routes:        make(map[string]*Route),
		config:        config,
		TemplateCache: functions.NewTemplateCache(), 
	}
}

func (s *Server) ListenAndServe() error {
    s.router.Use(antiflood.Limit(
        5,
        5*time.Second,
        antiflood.WithKeyFuncs(antiflood.KeyByRealIP, antiflood.KeyByEndpoint),
        antiflood.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusTooManyRequests)
            w.Write([]byte(`{"status": "error", "message":"you have been ratelimited!"}`))
        }),
    ))

    s.router.Use(antiflood.Limit(
        750,
        1*time.Minute,
        antiflood.WithKeyFuncs(antiflood.KeyByRealIP),
        antiflood.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusTooManyRequests)
            w.Write([]byte(`{"status": "error", "message":"you have been ratelimited!"}`))
        }),
    ))

    s.mux.Handler = s.router

    var err error
    if models.Config.Secure {
        cert := models.Config.Cert
        key := models.Config.Key
        if cert == "" || key == "" {
            return errors.New("certificate or key is empty")
        }

        s.mux.Addr = strings.Split(s.config.Addr, ":")[0] + ":443"
        Sanitize.LogMessage("success", "Server is running on HTTPS on "+s.mux.Addr, nil)
        Sanitize.LogMessage("info", "Listening with "+fmt.Sprint(s.Subrouters())+" subrouters and "+fmt.Sprint(s.Routes())+" routes.", nil)

        err = s.mux.ListenAndServeTLS(cert, key)
    } else {
        Sanitize.LogMessage("warning", "Server is running on HTTP on "+s.mux.Addr, nil)
        Sanitize.LogMessage("info", "Listening with "+fmt.Sprint(s.Subrouters())+" subrouters and "+fmt.Sprint(s.Routes())+" routes.", nil)

        err = s.mux.ListenAndServe()
    }

    // Handle errors and filter specific ones
    if err != nil {
        if !strings.Contains(err.Error(), "TLS handshake error") {
            log.Printf("Server error: %v", err)
        }
    }

    return err
}

func (s *Server) Subrouters() int {
	subs := 0
	for _, sub := range s.routes {
		if sub.Subrouter {
			subs++
		}
	}
	return subs
}
func (s *Server) Routes() int {
	subs := 0
	for _, sub := range s.routes {
		if !sub.Subrouter && sub.Handler != nil {
			subs++
		}
	}
	return subs
}
