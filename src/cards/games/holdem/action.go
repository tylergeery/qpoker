package holdem

const (
	// ActionBet is a bet action
	ActionBet = "bet"
	// ActionCall is a call action
	ActionCall = "call"
	// ActionCheck is a check action
	ActionCheck = "check"
	// ActionFold is a fold action
	ActionFold = "fold"
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
