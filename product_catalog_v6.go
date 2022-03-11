package ingram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type PriceAndAvailabilityRequest struct {
	IncludeAvailability bool   `validate:"required"`
	IncludePricing      bool   `validate:"required"`
	CustomerNumber      string `validate:"required"`
	ISOCountryCode      string `validate:"required"`
	IngramPartNumber    string `validate:"required"`
}

type PriceAndAvailabilityResponse struct {
	ProductStatusCode         string       `json:"productStatusCode"`
	ProductStatusMessage      string       `json:"productStatusMessage"`
	IngramPartNumber          string       `json:"ingramPartNumber"`
	VendorPartNumber          string       `json:"vendorPartNumber"`
	CustomerPartNumber        string       `json:"customerPartNumber"`
	UPC                       string       `json:"upc"`
	PartNumberType            string       `json:"partNumberType"`
	VendorName                string       `json:"vendorName"`
	VendorNumber              string       `json:"vendorNumber"`
	Description               string       `json:"description"`
	ProductClass              string       `json:"productClass"`
	UOM                       string       `json:"UOM"`
	ProductStatus             string       `json:"productStatus"`
	AcceptBackOrder           bool         `json:"acceptBackOrder"`
	ProductAuthorized         bool         `json:"productAuthorized"`
	ReturnableProduct         bool         `json:"returnableProduct"`
	EndUserInfoRequired       bool         `json:"endUserInfoRequired"`
	GovtSpecialPriceAvailable bool         `json:"govtSpecialPriceAvailable"`
	GovtProgramType           string       `json:"govtProgramType"`
	GovtEndUserType           string       `json:"govtEndUserType"`
	Availability              Availability `json:"availability"`
	Pricing                   Pricing      `json:"pricing"`
}

type Availability struct {
	Available               bool                      `json:"available"`
	TotalAvailability       int64                     `json:"totalAvailability"`
	AvailabilityByWarehouse []AvailabilityByWarehouse `json:"availabilityByWarehouse"`
}

type AvailabilityByWarehouse struct {
	Location               string `json:"location"`
	WarehouseId            string `json:"warehouseId"`
	QuantityAvailable      int64  `json:"quantityAvailable"`
	QuantityBackordered    int64  `json:"quantityBackordered"`
	QuantityBackorderedEta string `json:"quantityBackorderedEta"`
}

type Pricing struct {
	CurrencyCode               string  `json:"currencyCode"`
	RetailPrice                float64 `json:"retailPrice"`
	MapPrice                   float64 `json:"mapPrice"`
	CustomerPrice              float64 `json:"customerPrice"`
	SpecialBidPricingAvailable bool    `json:"specialBidPricingAvailable"`
	WebDiscountsAvailable      bool    `json:"webDiscountsAvailable"`
}

type Product struct {
	IngramPartNumber     string                `json:"ingramPartNumber"`
	VendorPartNumber     string                `json:"vendorPartNumber"`
	CustomerPartNumber   string                `json:"customerPartNumber"`
	UPC                  string                `json:"upc"`
	QuantityRequested    string                `json:"quantityRequested"`
	AdditionalAttributes []AdditionalAttribute `json:"additionalAttributes"`
}

type AdditionalAttribute struct {
	AttributeName  string `json:"attributeName"`
	AttributeValue string `json:"attributeValue"`
}

func (i *Ingram) PriceAndAvailability(ctx context.Context, priceAndAvailabilityRequest *PriceAndAvailabilityRequest) ([]PriceAndAvailabilityResponse, error) {
	err := i.validate.Struct(priceAndAvailabilityRequest)
	if err != nil {
		return nil, err
	}

	err = i.checkAndUpdateToken(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(fmt.Sprintf("%s/resellers/v6/catalog/priceandavailability", i.endpoint))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("includeAvailability", strconv.FormatBool(priceAndAvailabilityRequest.IncludeAvailability))
	q.Add("includePricing", strconv.FormatBool(priceAndAvailabilityRequest.IncludePricing))
	u.RawQuery = q.Encode()

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(struct {
		Products []Product `json:"products"`
	}{
		Products: []Product{{
			IngramPartNumber: priceAndAvailabilityRequest.IngramPartNumber,
		}},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", i.token.AccessToken))
	req.Header.Set("Accept", "*/*")
	req.Header.Add("IM-CustomerNumber", priceAndAvailabilityRequest.CustomerNumber)
	req.Header.Add("IM-CountryCode", priceAndAvailabilityRequest.ISOCountryCode)
	req.Header.Add("IM-CorrelationID", uuid.NewV4().String()) // some random uuid to trick ingram

	if i.logger != nil {
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		i.logger.Printf(string(b))
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if i.logger != nil {
		b, err := httputil.DumpResponse(res, true)
		if err != nil {
			return nil, err
		}
		i.logger.Printf(string(b))
	}

	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("%s: %s", res.Status, string(body))
	}

	var response []PriceAndAvailabilityResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
