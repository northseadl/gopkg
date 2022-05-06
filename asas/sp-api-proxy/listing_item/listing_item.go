package listing_item

import (
	"errors"
	"fmt"
	sp_api_proxy "go-pkg/asas/sp-api-proxy"
	"io"
	"net/http"
	"net/url"
)

const (
	API                = "Listings Items"
	Version            = "2021-08-01"
	getListingsItemUri = "/listings/2021-08-01/items/%s/%s" // sellerId, sku
	defaultLocal       = "en_US"
	includeData        = "summaries,attributes,issues,offers,fulfillmentAvailability,procurement"
)

type ListingItem struct {
	ProxyConfig *sp_api_proxy.ProxyConfig
}

func NewListingItemAPI(host string) ListingItem {
	return ListingItem{
		ProxyConfig: sp_api_proxy.NewProxyBaseUrl(host),
	}
}

func (l ListingItem) GetListingsItem(region string, marketplaceId string, sellerId string, sellerSku string) ([]byte, error) {
	iUrl := l.ProxyConfig.UseRegion(region).BuildUri(fmt.Sprintf(getListingsItemUri, sellerId, url.QueryEscape(sellerSku)))
	req, _ := http.NewRequest(http.MethodGet, iUrl, nil)
	query := req.URL.Query()
	query.Set("marketplaceIds", marketplaceId)
	query.Set("includedData", includeData)
	query.Set("issueLocale", defaultLocal)
	req.URL.RawQuery = query.Encode()
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	result, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("stats code %d, response %s", resp.StatusCode, string(result)))
	}
	return result, nil
}
