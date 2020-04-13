package model

type Applications struct {
	PendingCount  *int           `json:"pendingCount"`
	AcceptedCount *int           `json:"acceptedCount"`
	RejectedCount *int           `json:"rejectedCount"`
	Applications  []*Application `json:"applications"`
}
