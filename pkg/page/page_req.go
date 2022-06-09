package page

const (
	// Asc asc order
	Asc = "asc"
	// Desc desc order
	Desc = "desc"
)

// ReqPage page request base params
type ReqPage struct {
	Page   int         `json:"page"`
	Size   int         `json:"size"`
	Orders []OrderItem `json:"orders"`
}

// OrderItem order item
type OrderItem struct {
	Column    string `json:"column"`
	Direction string `json:"direction"` // asc|desc
}
