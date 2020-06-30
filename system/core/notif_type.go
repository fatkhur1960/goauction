package core

// NotifType notification type
type NotifType int

const (
	// GotBidder type for user bid
	GotBidder NotifType = iota
	// GotMessage type when user got chat message
	GotMessage NotifType = iota
	// BidClosed type when product bid is closed
	BidClosed NotifType = iota
	// GotWinner type when product has winner
	GotWinner NotifType = iota
	// WinBid type when user had win the bid
	WinBid NotifType = iota
)
