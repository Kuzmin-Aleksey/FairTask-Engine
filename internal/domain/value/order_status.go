package value

type OrderStatus string

const (
	OrderStatusProcessed OrderStatus = "processed"
	OrderStatusAwait     OrderStatus = "await"
	OrderStatusAccept    OrderStatus = "accept"
	OrderStatusReject    OrderStatus = "reject"
)
