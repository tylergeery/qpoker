package holdem

const (
	// ActionBet is a bet action
	ActionBet = "Bet"
	// ActionCall is a call action
	ActionCall = "Call"
	// ActionCheck is a check action
	ActionCheck = "Check"
	// ActionFold is a fold action
	ActionFold = "Fold"
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
