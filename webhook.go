package ingram

import "time"

type WebhookEvent string

const (
	UpdateEvent WebhookEvent = "im::updated"
)

type WebhookEventType string

const (
	OrderShipped  WebhookEventType = "im::order_shipped"
	OrderInvoiced                  = "im::order_invoiced"
	OrderHold                      = "im::order_hold"
	OrderVoided                    = "im::order_voided"
)

type WebhookLineStatus string

const (
	LineShipped WebhookLineStatus = "im::shipped"
)

type Webhook struct {
	Topic          string          `json:"topic"`
	Event          string          `json:"event"`
	EventTimeStamp time.Time       `json:"evenTtimeStamp"`
	EventID        string          `json:"eventId"`
	Resource       WebhookResource `json:"resource"`
}

type WebhookResource struct {
	EventType           string                `json:"eventType"`
	OrderNumber         string                `json:"orderNumber"`
	CustomerPoNumber    string                `json:"customerPoNumber"`
	OrderEntryTimeStamp time.Time             `json:"orderEntryTimeStamp"`
	Lines               []WebhookResourceLine `json:"lines"`
}

type WebhookResourceLine struct {
	LineNumber          string                      `json:"lineNumber"`
	SubOrderNumber      string                      `json:"subOrderNumber"`
	LineStatus          string                      `json:"lineStatus"`
	IngramPartNumber    string                      `json:"ingramPartNumber"`
	VendorPartNumber    string                      `json:"vendorPartNumber"`
	RequestedQuantity   string                      `json:"requestedQuantity"`
	ShippedQuantity     string                      `json:"shippedQuantity"`
	BackOrderedQuantity string                      `json:"backOrderedQuantity"`
	ShipmentDetails     WebhookShipmentDetail       `json:"shipmentDetails"`
	SerialNumberDetails []WebhookSerialNumberDetail `json:"serialNumberDetails"`
}

type WebhookShipmentDetail struct {
	ShipmentDate        *string                        `json:"shipmentDate"`
	ShipFromWarehouseID string                         `json:"shipFromWarehouseId"`
	WarehouseName       string                         `json:"warehouseName"`
	CarrierCode         string                         `json:"carrierCode"`
	CarrierName         string                         `json:"carrierName"`
	PackageDetails      []WebhookShipmentPackageDetail `json:"packageDetails"`
}

type WebhookShipmentPackageDetail struct {
	CartonNumber   string `json:"cartonNumber"`
	QuantityInbox  string `json:"quantityInbox"`
	TrackingNumber string `json:"trackingNumber"`
}

type WebhookSerialNumberDetail struct {
	SerialNumber string `json:"serialNumber"`
}
