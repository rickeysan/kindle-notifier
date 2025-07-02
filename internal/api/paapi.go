package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type PAAPIClient struct {
	AccessKey    string
	SecretKey    string
	PartnerTag   string
	Region       string
	ServiceName  string
	Host         string
	ContentType  string
	RequestURL   string
}

type SearchBooksResponse struct {
	ItemsResult struct {
		Items []struct {
			ASIN         string `json:"ASIN"`
			DetailPageURL string `json:"DetailPageURL"`
			Images struct {
				Primary struct {
					Large struct {
						URL string `json:"URL"`
					} `json:"Large"`
				} `json:"Primary"`
			} `json:"Images"`
			ItemInfo struct {
				Title struct {
					DisplayValue string `json:"DisplayValue"`
				} `json:"Title"`
				ByLineInfo struct {
					Contributors []struct {
						Name string `json:"Name"`
						Role string `json:"Role"`
					} `json:"Contributors"`
				} `json:"ByLineInfo"`
			} `json:"ItemInfo"`
		} `json:"Items"`
	} `json:"ItemsResult"`
}

func NewPAAPIClient() *PAAPIClient {
	return &PAAPIClient{
		AccessKey:    os.Getenv("PAAPI_ACCESS_KEY"),
		SecretKey:    os.Getenv("PAAPI_SECRET_KEY"),
		PartnerTag:   os.Getenv("PAAPI_PARTNER_TAG"),
		Region:       "us-west-2",
		ServiceName:  "ProductAdvertisingAPI",
		Host:         "webservices.amazon.co.jp",
		ContentType:  "application/json; charset=utf-8",
		RequestURL:   "https://webservices.amazon.co.jp/paapi5/searchitems",
	}
}

func (c *PAAPIClient) SearchKindleUnlimitedBooks() (*SearchBooksResponse, error) {
	payload := map[string]interface{}{
		"Keywords": "Kindle Unlimited",
		"SearchIndex": "Books",
		"Resources": []string{
			"Images.Primary.Large",
			"ItemInfo.Title",
			"ItemInfo.ByLineInfo",
		},
		"PartnerTag": c.PartnerTag,
		"PartnerType": "Associates",
		"Marketplace": "www.amazon.co.jp",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.RequestURL, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, err
	}

	// Set required headers
	amzDate := time.Now().UTC().Format("20060102T150405Z")
	req.Header.Set("Content-Type", c.ContentType)
	req.Header.Set("X-Amz-Date", amzDate)
	req.Header.Set("X-Amz-Target", "com.amazon.paapi5.v1.ProductAdvertisingAPIv1.SearchItems")

	// Sign the request
	signature := c.signRequest(req, payloadBytes, amzDate)
	req.Header.Set("Authorization", signature)

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SearchBooksResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *PAAPIClient) signRequest(req *http.Request, payload []byte, amzDate string) string {
	// Implementation of AWS Signature Version 4
	// This is a simplified version. You'll need to implement the full signing process
	dateStamp := amzDate[:8]
	
	// Create canonical request
	canonicalURI := "/"
	canonicalQueryString := ""
	canonicalHeaders := fmt.Sprintf("content-type:%s\nhost:%s\nx-amz-date:%s\n", 
		c.ContentType, c.Host, amzDate)
	signedHeaders := "content-type;host;x-amz-date"
	
	payloadHash := hash(payload)
	canonicalRequest := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		"POST", canonicalURI, canonicalQueryString, canonicalHeaders, signedHeaders, payloadHash)

	// Create string to sign
	algorithm := "AWS4-HMAC-SHA256"
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", dateStamp, c.Region, c.ServiceName)
	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s",
		algorithm, amzDate, credentialScope, hash([]byte(canonicalRequest)))

	// Calculate signature
	signingKey := getSignatureKey(c.SecretKey, dateStamp, c.Region, c.ServiceName)
	signature := hex.EncodeToString(hmacSHA256(signingKey, []byte(stringToSign)))

	// Create authorization header
	authorizationHeader := fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		algorithm, c.AccessKey, credentialScope, signedHeaders, signature)

	return authorizationHeader
}

func hash(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func hmacSHA256(key []byte, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func getSignatureKey(key, dateStamp, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+key), []byte(dateStamp))
	kRegion := hmacSHA256(kDate, []byte(region))
	kService := hmacSHA256(kRegion, []byte(service))
	kSigning := hmacSHA256(kService, []byte("aws4_request"))
	return kSigning
} 