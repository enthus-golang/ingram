package ingram

import "time"

type WebhookEvent string
type WebhookEventType string
type WebhookLineStatus string

const (
	UpdateEvent WebhookEvent = "im::updated"

	OrderShipped  WebhookEventType = "IM::order_shipped"
	OrderInvoiced WebhookEventType = "IM::order_invoiced"
	OrderHold     WebhookEventType = "IM::order_hold"
	OrderVoided   WebhookEventType = "IM::order_voided"

	LineShipped    WebhookLineStatus = "IM::SHIPPED"
	LineSalesHold  WebhookLineStatus = "IM::SALES_HOLD"
	LineOnlineHold WebhookLineStatus = "IM::IM_ONLINE_HOLD"
)

type Webhook struct {
	Topic          string          `json:"topic"`
	Event          string          `json:"event"`
	EventTimeStamp time.Time       `json:"eventTimeStamp"`
	EventID        string          `json:"eventId"`
	Resource       WebhookResource `json:"resource"`
}

type WebhookResource struct {
	EventType           WebhookEventType      `json:"eventType"`
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
