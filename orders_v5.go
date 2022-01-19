package ingram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type LineType string

const (
	Position LineType = "P"
	Comment           = "C"
)

type OrderDetailRequest struct {
	OrderNumber    string `validate:"required"`
	CustomerNumber string `validate:"required"`
	ISOCountryCode string `validate:"required"`
}

type OrderDetailResponseServiceResponse struct {
	ServiceResponse OrderDetailServiceResponse `json:"serviceresponse"`
}

type OrderDetailServiceResponse struct {
	ResponsePreamble    ResponsePreamble    `json:"responsepreamble"`
	OrderDetailResponse OrderDetailResponse `json:"orderdetailresponse"`
}

type OrderDetailResponse struct {
	OrderNumber            string    `json:"ordernumber"`
	OrderType              string    `json:"ordertype"`
	CustomerOrderNumber    string    `json:"customerordernumber"`
	EndUserPoNumber        string    `json:"enduserponumber"`
	OrderStatus            string    `json:"orderstatus"`
	EntryTimestamp         time.Time `json:"entrytimestamp"`
	EntryMethodDescription string    `json:"entrymethoddescription"`
	OrderTotalValue        float64   `json:"ordertotalvalue"`
	OrderSubTotal          float64   `json:"ordersubtotal"`
	FreightAmount          string    `json:"freightamount"`
	CurrencyCode           string    `json:"currencycode"`
	TotalWeight            string    `json:"totalweight"`
	TotalTax               string    `json:"totaltax"`
	BillToAddress          OrderDetailAddress
	ShipToAddress          OrderDetailAddress
	Lines                  []OrderDetailLine        `json:"lines"`
	CommentLines           []OrderDetailCommentLine `json:"commentlines"`
	MiscFeeLines           []OrderDetailMiscFeeLine `json:"miscfeeline"`
	ExtendedSpecs          []ExtendedSpec           `json:"extendedspecs"`
}

type OrderDetailAddress struct {
	Suffix       string `json:"suffix"`
	Name         string `json:"name"`
	Attention    string `json:"attention"`
	AddressLine1 string `json:"addressline1"`
	AddressLine2 string `json:"addressline2"`
	AddressLine3 string `json:"addressline3"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postalcode"`
	CountryCode  string `json:"countrycode"`
}

type OrderDetailLine struct {
	LineNumber             string  `json:"linenumber"`
	GlobalLineNumber       string  `json:"globallinenumber"`
	OrderSuffix            string  `json:"ordersuffix"`
	ERPOrderNumber         string  `json:"erpordernumber"`
	LineStatus             string  `json:"linestatus"`
	PartNumber             string  `json:"partnumber"`
	ManufacturerPartNumber string  `json:"manufacturerpartnumber"`
	VendorName             string  `json:"vendorname"`
	VendorCode             string  `json:"vendorcode"`
	PartDescription1       string  `json:"partdescription1"`
	PartDescription2       string  `json:"partdescription2"`
	UnitWeight             string  `json:"unitweight"`
	UnitPrice              float64 `json:"unitprice"`
	ExtendedPrice          float64 `json:"extendedprice"`
	TaxAmount              float64 `json:"taxamount"`
	RequestedQuantity      string  `json:"requestedquantity"`
	ConfirmedQuantity      string  `json:"confirmedquantity"`
	BackorderQuantity      string  `json:"backorderquantity"`
}

type OrderDetailCommentLine struct {
	CommentText1 string `json:"commenttext1"`
	CommentText2 string `json:"commenttext2"`
}

type OrderDetailMiscFeeLine struct {
	Description  string `json:"description"`
	ChargeAmount string `json:"chargeamount"`
}

func (i *Ingram) OrderDetail(ctx context.Context, orderDetail *OrderDetailRequest) (*OrderDetailResponseServiceResponse, error) {
	err := i.validate.Struct(orderDetail)
	if err != nil {
		return nil, err
	}

	err = i.checkAndUpdateToken(ctx)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(fmt.Sprintf("%s/resellers/v5/orders/%s", i.endpoint, orderDetail.OrderNumber))
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Add("customernumber", orderDetail.CustomerNumber)
	q.Add("isocountrycode", orderDetail.ISOCountryCode)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", i.token.AccessToken))
	req.Header.Set("Accept", "*/*")

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

	var response OrderDetailResponseServiceResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

type createOrderV5 struct {
	OrderCreateRequest OrderCreateRequest `json:"ordercreaterequest"`
}

type OrderCreateRequest struct {
	RequestPreamble    RequestPreamble    `json:"requestpreamble"`
	OrderCreateDetails OrderCreateDetails `json:"ordercreatedetails"`
}

type RequestPreamble struct {
	ISOCountryCode string `json:"isocountrycode"`
	CustomerNumber string `json:"customernumber"`
}

type OrderCreateDetails struct {
	CustomerPurchaseOrderNumber    string         `json:"customerponumber" validate:"min=1,max=18"`
	OrderType                      string         `json:"ordertype,omitempty"`
	EndUserOrderNumber             string         `json:"enduserordernumber,omitempty" validate:"max=18"`
	BillToSuffix                   string         `json:"billtosuffix,omitempty" validate:"max=3"`
	ShipToSuffix                   string         `json:"shiptosuffix,omitempty" validate:"max=3"`
	ShipToAddress                  *ShipToAddress `json:"shiptoaddress,omitempty"`
	CarrierCode                    string         `json:"carriercode,omitempty" validate:"max=2"`
	ThirdPartyFreightAccountNumber string         `json:"thirdpartyfrieghtaccountnumber,omitempty"`
	SpecialBidNumber               string         `json:"specialbidnumber,omitempty"`
	Lines                          []Line         `json:"lines,omitempty"`
	ExtendedSpecs                  []ExtendedSpec `json:"extendedspecs,omitempty"`
}

type ShipToAddress struct {
	Attention    string `json:"attention,omitempty" validate:"max=35"`
	AddressLine1 string `json:"addressline1" validate:"required,max=35"`
	AddressLine2 string `json:"addressline2" validate:"required,max=35"`
	AddressLine3 string `json:"addressline3,omitempty" validate:"max=35"`
	City         string `json:"city" validate:"required,max=21"`
	State        string `json:"state" validate:"max=2"`
	PostalCode   string `json:"postalcode" validate:"required,max=9"`
	CountryCode  string `json:"countrycode" validate:"required,max=2"`
}

type Line struct {
	LineType             LineType       `json:"linetype,omitempty"`
	LineNumber           string         `json:"linenumber,omitempty"`
	IngramPartNumber     string         `json:"ingrampartnumber,omitempty"`
	Quantity             int            `json:"quantity"`
	VendorPartNumber     string         `json:"vendorpartnumber,omitempty"`
	CustomerPartNumber   string         `json:"customerpartnumber,omitempty"`
	UPCCode              string         `json:"UPCCode,omitempty"`
	WareHouseID          string         `json:"warehouseid,omitempty"`
	EndUserPrice         float64        `json:"enduserprice,omitempty"`
	UnitPrice            float64        `json:"unitprice,omitempty"`
	EndUser              *EndUser       `json:"enduser,omitempty"`
	ProductExtendedSpecs []ExtendedSpec `json:"productextendedspecs,omitempty"`
}

type EndUser struct {
	ID              string `json:"id,omitempty"`
	AddressLine1    string `json:"addressline1,omitempty"`
	AddressLine2    string `json:"addressline2,omitempty"`
	AddressLine3    string `json:"addressline3,omitempty"`
	City            string `json:"city,omitempty"`
	State           string `json:"state,omitempty"`
	PostalCode      string `json:"postalcode,omitempty"`
	CountryCode     string `json:"countrycode,omitempty"`
	PhoneNumber     string `json:"phonenumber,omitempty"`
	ExtensionNumber string `json:"extensionnumber,omitempty"`
	FaxNumber       string `json:"faxnumber,omitempty"`
	Email           string `json:"email,omitempty"`
}

type ExtendedSpec struct {
	AttributeName  string `json:"attributename"`
	AttributeValue string `json:"attributevalue"`
}

type OrderCreateResponseServiceResponse struct {
	ServiceResponse OrderServiceResponse `json:"serviceresponse"`
}

type OrderServiceResponse struct {
	ResponsePreamble ResponsePreamble `json:"responsepreamble"`
	OrderSummary     OrderSummary     `json:"ordersummary"`
}

type ResponsePreamble struct {
	ResponseStatus  string `json:"responsestatus"`
	StatusCode      string `json:"statuscode"`
	ResponseMessage string `json:"responsemessage"`
}

type OrderSummary struct {
	OrderCreateResponses []OrderCreateResponse `json:"ordercreateresponse"`
}

type OrderCreateResponse struct {
	NumberOfLinesWithSuccess string                    `json:"numberoflineswithsuccess"`
	NumberOfLinesWithError   string                    `json:"numberoflineswitherror"`
	NumberOfLinesWithWarning string                    `json:"numberoflineswithwarning"`
	GlobalOrderID            string                    `json:"globalorderid"`
	OrderType                string                    `json:"ordertype"`
	OrderTimestamp           string                    `json:"ordertimestamp"`
	InvoicingSystemOrderID   string                    `json:"invoicingsystemorderid"`
	TaxAmount                float64                   `json:"taxamount"`
	FreightAmount            float64                   `json:"freightamount"`
	OrderAmount              float64                   `json:"orderamount"`
	Lines                    []OrderCreateResponseLine `json:"lines"`
}

type OrderCreateResponseLine struct {
	LineType         string `json:"linetype"`
	GlobalLineNumber string `json:"globallinenumber"`
	PartNumber       string `json:"partnumber"`
	GlobalSKUID      string `json:"globalskuid"`
	LineNumber       string `json:"linenumber"`
}

func (i *Ingram) CreateOrderV5(ctx context.Context, order *OrderCreateRequest) (*OrderCreateResponseServiceResponse, error) {
	err := i.validate.Struct(order)
	if err != nil {
		return nil, err
	}

	err = i.checkAndUpdateToken(ctx)
	if err != nil {
		return nil, err
	}

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(createOrderV5{
		OrderCreateRequest: *order,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, i.endpoint+"/resellers/v5/orders", b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", i.token.AccessToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")

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
		return nil, errors.New(res.Status)
	}

	var response OrderCreateResponseServiceResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
