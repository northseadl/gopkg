package sp_api_proxy

import (
	"fmt"
	"strings"
)

const (
	RegionNa = "na"
	RegionEU = "eu"
	RegionFE = "fe"
)

type ProxyConfig struct {
	host   string
	region string
}

func NewProxyBaseUrl(host string) *ProxyConfig {
	return &ProxyConfig{
		host:   host,
		region: "na",
	}
}

func (p *ProxyConfig) UseRegion(region string) *ProxyConfig {
	p.region = region
	return p
}

func (p *ProxyConfig) BuildUri(uri string) string {
	uri = strings.TrimLeft(uri, "/")
	return fmt.Sprintf("http://%s/proxy/sp-api/%s/%s", p.host, p.region, uri)
}
