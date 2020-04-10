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

// NewActionBet returns a new bet action
func NewActionBet(amount int64) Action {
	return Action{ActionBet, amount}
}

// NewActionFold returns a new fold action
func NewActionFold() Action {
	return Action{Name: ActionFold}
}

// NewActionCall returns a new fold action
func NewActionCall() Action {
	return Action{Name: ActionCall}
}

// NewActionCheck returns a new fold action
func NewActionCheck() Action {
	return Action{Name: ActionCheck}
}
