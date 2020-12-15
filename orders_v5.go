package ingram

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type LineType string

const (
	Position LineType = "P"
	Comment           = "C"
)

type createOrderV5 struct {
	OrderCreateRequest OrderCreateRequest `json:"ordercreaterequest"`
}

type OrderCreateRequest struct {
	RequestPreamble    RequestPreamble    `json:"requestpreamble"`
	OrderCreateDetails OrderCreateDetails `json:"ordercreatedetails"`
}

type RequestPreamble struct {
	ISOCountryCode string
	CustomerNumber string
}

type OrderCreateDetails struct {
	CustomerPurchaseOrderNumber    string         `json:"customerponumber" validate:"min=1,max=18"`
	OrderType                      string         `json:"ordertype" validate:"required"`
	EndUserOrderNumber             string         `json:"enduserordernumber,omitempty" validate:"max=18"`
	BillToSuffix                   string         `json:"billtosuffix,omitempty" validate:"max=3"`
	ShipToSuffix                   string         `json:"shiptosuffix,omitempty" validate:"max=3"`
	ShipToAddress                  ShipToAddress  `json:"shiptoaddress"`
	CarrierCode                    string         `json:"carriercode,omitempty" validate:"max=2"`
	ThirdPartyFreightAccountNumber string         `json:"thirdpartyfrieghtaccountnumber,omitempty"`
	SpecialBidNumber               string         `json:"specialbidnumber,omitempty"`
	Lines                          []Line         `json:"lines"`
	ExtendedSpecs                  []ExtendedSpec `json:"extendedspecs"`
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
	EndUser              *EndUser       `json:"enduser"`
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

func (i *Ingram) CreateOrderV5(ctx context.Context, order *OrderCreateRequest) error {
	err := i.validate.Struct(order)
	if err != nil {
		return err
	}

	err = i.checkAndUpdateToken(ctx)
	if err != nil {
		return err
	}

	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(createOrderV5{
		OrderCreateRequest: *order,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, i.endpoint+"/resellers/v5/orders", b)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	return nil
}
