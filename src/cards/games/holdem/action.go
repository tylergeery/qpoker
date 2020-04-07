package holdem

const (
	ActionBet   = "Bet"
	ActionCall  = "Call"
	ActionCheck = "Check"
	ActionFold  = "Fold"
)

// Action holds the info regarding a holdem action
type Action struct {
	Name   string `json:"name"`
	Amount int64  `json:"amount"`
}

// ActionNewBet returns a new bet action
func ActionNewBet(amount int64) Action {
	return Action{ActionBet, amount}
}
