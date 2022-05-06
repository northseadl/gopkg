package report

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/northseadl/gopkg/asas/sp_api_proxy"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	API                  = "Reports"
	Version              = "2021-06-30"
	createReportUri      = "/reports/2021-06-30/reports"
	getReportUri         = "/reports/2021-06-30/reports/%s"  // reportId
	getReportDocumentUri = "reports/2021-06-30/documents/%s" // reportDocumentId
)

type Report struct {
	ProxyConfig *sp_api_proxy.ProxyConfig
}

func NewReportAPI(host string) Report {
	return Report{
		ProxyConfig: sp_api_proxy.NewProxyBaseUrl(host),
	}
}

type Spec struct {
	ReportType     string
	StartDate      time.Time
	MarketPlaceIds []string
	Others         map[string]interface{}
}

func (r Report) FetchReport(region string, spec Spec) (io.Reader, error) {
	url := r.ProxyConfig.UseRegion(region).BuildUri(createReportUri)
	data := make(map[string]interface{})
	data["reportType"] = spec.ReportType
	data["dataStartTime"] = spec.StartDate.Format(time.RFC3339)
	data["marketplaceIds"] = spec.MarketPlaceIds
	if spec.Others != nil {
		for key, _ := range spec.Others {
			data[key] = spec.Others[key]
		}
	}
	jsonData, _ := json.Marshal(data)
	// create report
	resp, err := http.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusAccepted {
		return nil, errors.New(fmt.Sprintf("create report error, status: %s, reponse: %s", resp.Status, string(bodyBytes)))
	}
	result := make(map[string]interface{})
	_ = json.Unmarshal(bodyBytes, &result)
	if _, ok := result["reportId"]; !ok {
		return nil, errors.New("create report error, can't find reportId")
	}
	reportId := result["reportId"].(string)
	// get report
	url = r.ProxyConfig.UseRegion(region).BuildUri(fmt.Sprintf(getReportUri, reportId))
	// loop until processStatus = DONE
	for true {
		resp, err = http.Get(url)
		bodyBytes, _ = io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New(fmt.Sprintf("get report error, status: %s, reponse: %s", resp.Status, string(bodyBytes)))
		}
		_ = json.Unmarshal(bodyBytes, &result)
		if _, ok := result["processingStatus"]; !ok {
			return nil, errors.New("get report error, can't find processingStatus")
		}
		processingStatus := result["processingStatus"].(string)
		if processingStatus == "CANCELLED" || processingStatus == "FATAL" {
			return nil, errors.New(fmt.Sprintf("get report error, processingStatus: %s", processingStatus))
		}
		if processingStatus == "DONE" {
			break
		}
		if processingStatus != "IN_PROGRESS" && processingStatus != "IN_QUEUE" {
			return nil, errors.New(fmt.Sprintf("get report error, processingStatus: %s ???", processingStatus))
		}
		time.Sleep(time.Second * 2)
	}
	if _, ok := result["reportDocumentId"]; !ok {
		return nil, errors.New("get report error, can't find reportDocumentId")
	}
	reportDocumentId := result["reportDocumentId"].(string)
	// get document
	url = r.ProxyConfig.UseRegion(region).BuildUri(fmt.Sprintf(getReportDocumentUri, reportDocumentId))
	resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	bodyBytes, _ = io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("get report error, status: %s, reponse: %s", resp.Status, string(bodyBytes)))
	}
	_ = json.Unmarshal(bodyBytes, &result)
	if _, ok := result["url"]; !ok {
		return nil, errors.New("get report error, can't find url")
	}
	url = result["url"].(string)
	var compressionAlgorithm string
	if _, ok := result["compressionAlgorithm"]; ok {
		compressionAlgorithm = result["compressionAlgorithm"].(string)
	}
	resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ = io.ReadAll(resp.Body)
		return nil, errors.New(fmt.Sprintf("get document error, status: %s, reponse: %s", resp.Status, string(bodyBytes)))
	}
	ctp := strings.Trim(resp.Header.Get("Content-Type"), "\r\t ")

	// stream read
	switch compressionAlgorithm {
	case "":
		return charset.NewReader(resp.Body, ctp)
	case "GZIP":
		return gzip.NewReader(resp.Body)
	default:
		return nil, errors.New("get document error, compressionAlgorithm unsupported")
	}
}
