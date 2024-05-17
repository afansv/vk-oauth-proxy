package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/caarlos0/env/v11"

	"github.com/afansv/vk-oauth-proxy/modifiers"
	"github.com/afansv/vk-oauth-proxy/proxy"
	"github.com/afansv/vk-oauth-proxy/store"
)

// ProxyRequestHandler handles the http request using proxy
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

type Config struct {
	UserEmailStoreTTL time.Duration `env:"USER_EMAIL_STORE_TTL" envDefault:"1m"`
	OauthUpstreamHost string        `env:"OAUTH_UPSTREAM_HOST" envDefault:"https://oauth.vk.com"`
	APIUpstreamHost   string        `env:"API_UPSTREAM_HOST" envDefault:"https://api.vk.com"`
	OAuthProxyAddr    string        `env:"OAUTH_PROXY_ADDR" envDefault:":9090"`
	APIProxyAddr      string        `env:"API_PROXY_ADDR" envDefault:":9091"`
}

func main() {
	cfg := &Config{}
	opts := env.Options{
		Prefix: "VOP_",
	}
	// Load env vars.
	if err := env.ParseWithOptions(cfg, opts); err != nil {
		log.Fatal(err)
	}
	// dependencies
	userEmailStore := store.NewUserEmail(cfg.UserEmailStoreTTL)
	oauthResponseModifier := modifiers.NewOAuthResponseModifier(userEmailStore)
	apiResponseModifier := modifiers.NewAPIResponseModifier(userEmailStore)

	// initialize a reverse proxies and pass the actual backend server urls here
	oauthProxy, err := proxy.NewProxy(cfg.OauthUpstreamHost, oauthResponseModifier.Modify)
	if err != nil {
		log.Fatal(err)
	}
	apiProxy, err := proxy.NewProxy(cfg.APIUpstreamHost, apiResponseModifier.Modify)
	if err != nil {
		log.Fatal(err)
	}

	oauthProxyServer := http.NewServeMux()
	oauthProxyServer.HandleFunc("/", ProxyRequestHandler(oauthProxy))

	apiProxyServer := http.NewServeMux()
	apiProxyServer.HandleFunc("/", ProxyRequestHandler(apiProxy))

	go userEmailStore.Start()
	go func() {
		log.Println("starting oauth proxy at", cfg.OAuthProxyAddr)
		log.Fatal(http.ListenAndServe(cfg.OAuthProxyAddr, oauthProxyServer))
	}()
	go func() {
		log.Println("starting api proxy at", cfg.APIProxyAddr)
		log.Fatal(http.ListenAndServe(cfg.APIProxyAddr, apiProxyServer))
	}()
	select {}
}
