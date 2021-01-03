package holdem

import (
	"fmt"
	"qpoker/cards"
	"qpoker/models"
	"qpoker/qutils"
)

const (
	// MaxPlayerCount is the max amount of players for HoldEm
	MaxPlayerCount = 12

	// StatusInit init
	StatusInit = "init"
	// StatusReady ready
	StatusReady = "ready"
	// StatusActive active
	StatusActive = "active"
)

// GameManager holds game state and manages game flow
type GameManager struct {
	GameID   int64   `json:"game_id"`
	State    *HoldEm `json:"state"`
	Pot      *Pot    `json:"pot"`
	BigBlind int64   `json:"big_blind"`
	Status   string  `json:"status"`
	AllIn    bool    `json:"all_in"`

	gameHand        *models.GameHand
	gamePlayerHands map[int64]*models.GamePlayerHand
}

// GameOptions is a config object used to control game settings
type GameOptions struct {
	Capacity         int   `json:"capacity"`
	BigBlind         int64 `json:"big_blind"`
	DecisionTime     int   `json:"decision_time"`
	TimeBetweenHands int   `json:"time_between_hands"`
	BuyInMin         int64 `json:"buy_in_min"`
	BuyInMax         int64 `json:"buy_in_max"`
}

// NewGameOptions creates game options from options map
func NewGameOptions(options map[string]interface{}) GameOptions {
	return GameOptions{
		Capacity:         qutils.ToInt(options["capacity"]),
		BigBlind:         qutils.ToI64(options["big_blind"]),
		DecisionTime:     qutils.ToInt(options["decision_time"]),
		TimeBetweenHands: qutils.ToInt(options["time_between_hands"]),
		BuyInMin:         qutils.ToI64(options["buy_in_min"]),
		BuyInMax:         qutils.ToI64(options["buy_in_max"]),
	}
}

// NewGameManager returns a new GameManager
func NewGameManager(gameID int64, players []*Player, options GameOptions) (*GameManager, error) {
	if len(players) > MaxPlayerCount {
		return nil, fmt.Errorf("Invalid player count: %d", len(players))
	}

	if options.Capacity > 0 && len(players) > int(options.Capacity) {
		return nil, fmt.Errorf("Player count (%d) is greater than capacity (%d)", len(players), options.Capacity)
	}

	table := NewTable(options.Capacity, players)
	gm := &GameManager{
		GameID:   gameID,
		State:    NewHoldEm(table),
		BigBlind: options.BigBlind,
		Status:   StatusInit,
	}

	return gm, nil
}

// AddPlayer adds player to game
func (g *GameManager) AddPlayer(player *Player) error {
	err := g.State.Table.AddPlayer(player)

	g.UpdateStatus(StatusReady)

	return err
}

// RemovePlayer removes player from game
func (g *GameManager) RemovePlayer(playerID int64) error {
	err := g.State.Table.RemovePlayer(playerID)

	g.UpdateStatus(StatusInit)

	return err
}

// AddChips adds chips to Player
func (g *GameManager) AddChips(playerID, amount int64) {
	g.GetPlayer(playerID).Stack += amount
	g.UpdateStatus(StatusReady)
}

// NextHand moves game manager on to next hand
func (g *GameManager) NextHand() error {
	g.Status = StatusReady
	g.AllIn = false

	err := g.State.Deal()
	if err != nil {
		return err
	}

	// Save hand and player hands
	err = g.StartHand()
	if err != nil {
		return err
	}

	g.Pot = NewPot(g.State.Table.GetActivePlayers())

	// little
	g.State.Table.GetActivePlayer().LittleBlind = true
	g.playerBet(NewActionBet(qutils.MinInt64(g.BigBlind/2, g.State.Table.GetActivePlayer().Stack)))
	g.State.Table.ActivateNextPlayer(g.GetPlayerActions)

	// big blind
	g.State.Table.GetActivePlayer().BigBlind = true
	g.playerBet(NewActionBet(qutils.MinInt64(g.BigBlind, g.State.Table.GetActivePlayer().Stack)))

	g.UpdateStatus(StatusActive)
	g.ProcessAction()

	return nil
}

// NextRound advances gameplay to next round
func (g *GameManager) NextRound() {
	g.State.Advance()
	g.State.Table.NextRound(g.GetPlayerActions)
	g.Pot.ClearBets()
}

func (g *GameManager) cardsToStringArray(cardObjects []cards.Card) []string {
	stringCards := make([]string, len(cardObjects))

	for i, c := range cardObjects {
		stringCards[i] = c.ToString()
	}

	return stringCards
}

func (g *GameManager) isComplete() bool {
	playersRemaining := g.State.Table.GetActivePlayers()

	if len(playersRemaining) <= 1 {
		return true
	}

	if g.isRoundComplete() && g.State.State == StateRiver {
		return true
	}

	if g.IsAllIn() && g.State.State == StateRiver {
		return true
	}

	return false
}

func (g *GameManager) isRoundComplete() bool {
	playersRemaining := g.State.Table.GetActivePlayers()
	nextPlayer := playersRemaining[0]

	// Check if everyone has called/checked the next player
	for _, player := range playersRemaining {
		// Everyone gets a chance to play
		if player.State == PlayerStateInit {
			return false
		}

		// If anybody besides the next player bet, keep going
		if player.ID != nextPlayer.ID && player.State == PlayerStateBet {
			return false
		}
	}

	if g.State.State != StateDeal {
		return true
	}

	if !nextPlayer.BigBlind {
		return true
	}

	return g.Pot.MaxBet() != g.BigBlind
}

func (g *GameManager) canBet() bool {
	player := g.State.Table.GetActivePlayer()
	playerBet := g.Pot.PlayerBets[player.ID]

	if player.Stack <= int64(0) {
		return false
	}

	// check if other players have chips left to bet
	if g.countWithEffectiveStackAbove(g.Pot.MaxBet()) <= 1 {
		return false
	}

	// If user has current max bet, they can raise (first round only)
	if g.Pot.MaxBet() > int64(0) && g.Pot.MaxBet() == playerBet {
		return true
	}

	if g.Pot.MaxBet() >= player.Stack {
		return false
	}

	return true
}

func (g *GameManager) canCall() bool {
	if g.State.Table.GetActivePlayer().Stack <= int64(0) {
		return false
	}

	return g.Pot.MaxBet() > g.Pot.PlayerBets[g.State.Table.GetActivePlayer().ID]
}

func (g *GameManager) canCheck() bool {
	player := g.State.Table.GetActivePlayer()
	if player.Stack <= int64(0) {
		return false
	}

	playersTotal, players := int64(0), g.State.Table.GetActivePlayers()
	for i := range players {
		if players[i].ID == g.State.Table.GetActivePlayer().ID {
			continue
		}

		playersTotal += players[i].Stack
	}

	// Cannot check if nobody else has a stack
	if playersTotal <= int64(0) {
		return false
	}

	return g.Pot.MaxBet() == g.Pot.PlayerBets[g.State.Table.GetActivePlayer().ID]
}

func (g *GameManager) canFold() bool {
	if g.State.Table.GetActivePlayer().Stack <= int64(0) {
		return false
	}

	return g.Pot.MaxBet() > g.Pot.PlayerBets[g.State.Table.GetActivePlayer().ID]
}

func (g *GameManager) playerBet(action Action) error {
	if g.Status == StatusActive {
		if !g.canBet() {
			return fmt.Errorf("Cannot bet")
		}

		minBet := qutils.MinInt64(g.BigBlind, g.State.Table.GetActivePlayer().Stack)
		if action.Amount < minBet {
			return fmt.Errorf("Insufficient bet amount")
		}
	}

	player := g.State.Table.GetActivePlayer()
	amount := qutils.MinInt64(action.Amount, player.Stack)

	g.Pot.AddBet(player.ID, amount)
	player.Stack -= amount
	player.State = PlayerStateBet

	return nil
}

func (g *GameManager) playerCall(action Action) error {
	if !g.canCall() {
		return fmt.Errorf("Cannot call")
	}

	player := g.State.Table.GetActivePlayer()
	amount := qutils.MinInt64(
		g.Pot.MaxBet()-g.Pot.PlayerBets[player.ID],
		g.State.Table.GetActivePlayer().Stack,
	)

	g.Pot.AddBet(player.ID, amount)
	player.Stack -= amount
	player.State = PlayerStateCall

	return nil
}

func (g *GameManager) playerCheck(action Action) error {
	if !g.canCheck() {
		return fmt.Errorf("Cannot check")
	}

	g.State.Table.GetActivePlayer().State = PlayerStateCheck

	return nil
}

func (g *GameManager) playerFold(action Action) error {
	if !g.canFold() {
		return fmt.Errorf("Cannot fold")
	}

	g.State.Table.GetActivePlayer().State = PlayerStateFold

	return nil
}

// PlayerAction performs an action for player
func (g *GameManager) PlayerAction(playerID int64, action Action) (bool, error) {
	if g.isComplete() {
		return false, fmt.Errorf("Game hand is already complete")
	}

	if playerID != g.State.Table.GetActivePlayer().ID {
		return false, fmt.Errorf("User (%d) must wait for player (%d) to act", playerID, g.State.Table.GetActivePlayer().ID)
	}

	actionMap := map[string]func(action Action) error{
		ActionBet:   g.playerBet,
		ActionCall:  g.playerCall,
		ActionCheck: g.playerCheck,
		ActionFold:  g.playerFold,
	}

	err := actionMap[action.Name](action)
	if err != nil {
		return false, err
	}

	return g.ProcessAction()
}

// ProcessAction handles the post-processing of an action
func (g *GameManager) ProcessAction() (bool, error) {
	if g.isComplete() {
		err := g.EndHand()
		g.State.Table.GetActivePlayer().SetPlayerActions(nil)
		return true, err
	}

	if g.isRoundComplete() {
		g.NextRound()
		return false, nil
	}

	g.State.Table.ActivateNextPlayer(g.GetPlayerActions)
	if g.IsAllIn() {
		g.NextRound()
		return false, nil
	}

	if !g.State.Table.GetActivePlayer().HasOptions() {
		return g.ProcessAction()
	}

	return false, nil
}

// GetPlayerActions returns allowed active player actions
func (g *GameManager) GetPlayerActions() map[string]bool {
	return map[string]bool{
		"can_bet":   g.canBet(),
		"can_call":  g.canCall(),
		"can_check": g.canCheck(),
		"can_fold":  g.canFold(),
	}
}

// ShowVisibleCards marks remaining players as visible cards
func (g *GameManager) ShowVisibleCards() {
	activePlayers := g.State.Table.GetActivePlayers()

	if len(activePlayers) <= 1 {
		return
	}

	for i := range activePlayers {
		activePlayers[i].CardsVisible = true
	}
}

// StartHand saves the initial game state to the DB
func (g *GameManager) StartHand() error {
	g.gameHand = &models.GameHand{GameID: g.GameID}
	err := g.gameHand.Save()
	if err != nil {
		return err
	}

	g.gamePlayerHands = map[int64]*models.GamePlayerHand{}
	for _, player := range g.State.Table.GetActivePlayers() {
		g.gamePlayerHands[player.ID] = &models.GamePlayerHand{
			GameHandID: g.gameHand.ID,
			Cards:      g.cardsToStringArray(player.Cards),
			PlayerID:   player.ID,
			Starting:   player.Stack,
		}
		err = g.gamePlayerHands[player.ID].Save()
		if err != nil {
			return err
		}
	}

	return nil
}

// EndHand saves the ending game state to the DB
func (g *GameManager) EndHand() error {
	g.ShowVisibleCards()
	payouts := g.Pot.GetPayouts(g.State.GetWinningIDs())

	g.gameHand.Board = g.cardsToStringArray(g.State.Board)
	g.gameHand.Payouts = payouts
	g.gameHand.Bets = g.Pot.PlayerTotals

	err := g.gameHand.Save()
	if err != nil {
		return err
	}

	playerIDs := []int64{}
	for playerID := range g.gamePlayerHands {
		playerIDs = append(playerIDs, playerID)
	}

	players := g.State.Table.GetPlayersFromIDs(playerIDs)
	for _, player := range players {
		hand := g.gamePlayerHands[player.ID]

		if amount, ok := payouts[player.ID]; ok {
			player.Stack += amount
		}
		hand.Ending = player.Stack

		err = hand.Save()
		if err != nil {
			fmt.Printf("Error saving user hand: %s\n", err)
		}
	}

	return nil
}

// UpdateStatus updates game status to most appropriate
func (g *GameManager) UpdateStatus(status string) {
	switch {
	case g.Status == StatusInit && status == StatusReady:
		if len(g.State.Table.GetReadyPlayers()) > 1 {
			g.Status = status
		}
		break
	case g.Status == StatusReady && status == StatusInit:
		if len(g.State.Table.GetActivePlayers()) <= 1 {
			g.Status = status
		}
	case g.Status == StatusReady && status == StatusActive:
		if len(g.State.Table.GetActivePlayers()) > 1 {
			g.Status = status
		}
		break
	default:
		break
	}
}

// GetVisibleCards returns client representation of cards for those visible
func (g *GameManager) GetVisibleCards(playerID int64) map[int64][]cards.Card {
	visibleCards := map[int64][]cards.Card{}
	allIn := g.IsAllIn()

	for _, player := range g.State.Table.GetAllPlayers() {
		if player.Cards == nil {
			continue
		}

		if (allIn && player.IsActive()) || player.ID == playerID || player.CardsVisible {
			visibleCards[player.ID] = player.Cards
		}
	}

	return visibleCards
}

// GetPlayer returns a player from a table
func (g *GameManager) GetPlayer(playerID int64) *Player {
	return g.State.Table.GetPlayerByID(playerID)
}

// IsAllIn determines whether a round is all in and should proceed automagically
func (g *GameManager) IsAllIn() bool {
	g.AllIn = g.isAllIn()

	return g.AllIn
}

func (g *GameManager) isAllIn() bool {
	if g.Status != StatusActive {
		return false
	}

	playersWithStack := g.countWithStackAbove(int64(0))
	if playersWithStack == 0 {
		return true
	}

	if playersWithStack > 1 {
		return false
	}

	// Check if last person has accepted bet
	players := g.State.Table.GetActivePlayers()
	for i := range players {
		if players[i].Stack > int64(0) {
			return g.Pot.PlayerBets[players[i].ID] == g.Pot.MaxBet()
		}
	}

	// should never get here
	return true
}

func (g *GameManager) countWithStackAbove(amount int64) int {
	countWithStack := 0
	players := g.State.Table.GetActivePlayers()

	for i := range players {
		if players[i].Stack > amount {
			countWithStack++
		}
	}

	return countWithStack
}

func (g *GameManager) countWithEffectiveStackAbove(amount int64) int {
	countWithStack := 0
	players := g.State.Table.GetActivePlayers()

	for i := range players {
		if (players[i].Stack + g.Pot.PlayerBets[players[i].ID]) > amount {
			countWithStack++
		}
	}

	return countWithStack
}
