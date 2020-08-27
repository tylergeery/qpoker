package models

import (
	"fmt"
	"qpoker/qutils"
	"strconv"
)

// GameOptions handles getting/setting game specific options
type GameOptions struct {
	GameID     int64                  `json:"game_id"`
	GameTypeID int64                  `json:"game_type_id`
	Options    map[string]interface{} `json:"options"`
}

// GameOption contains information to create game option
type GameOption struct {
	GameTypeGameOptionID int64       `json:"game_type_game_option_id"`
	GameTypeID           int64       `json:"game_type_id"`
	GameOptionID         int64       `json:"game_option_id"`
	IsActive             bool        `json:"is_active"`
	DefaultValue         interface{} `json:"default_value"`
	Name                 string      `json:"name"`
	Label                string      `json:"label"`
	Type                 string      `json:"type"`
}

func getValueByType(valueType, value string) interface{} {
	switch valueType {
	case "boolean":
		v, err := strconv.ParseBool(value)
		if err != nil {
			fmt.Printf("Parse boolean err: %s\n", err)
		}
		return v
	case "number":
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Printf("Parse float err: %s\n", err)
		}
		return v
	case "integer":
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			fmt.Printf("Parse integer err: %s\n", err)
		}
		return v
	default:
		return value
	}
}

func getStringForType(valueType string, value interface{}) string {
	switch valueType {
	case "boolean":
		return fmt.Sprintf("%b", value)
	case "number":
		return fmt.Sprintf("%f", value)
	case "integer":
		return fmt.Sprintf("%d", qutils.ToI64(value))
	default:
		return value.(string)
	}
}

// GetGameOptionsForGame returns options for game
func GetGameOptionsForGame(gameID int64, gameTypeID int64) (GameOptions, error) {
	options := GameOptions{
		GameID:     gameID,
		GameTypeID: gameTypeID,
		Options:    map[string]interface{}{},
	}

	rows, err := ConnectToDB().Query(fmt.Sprintf(`
		SELECT go.name, go.type, COALESCE(gov.value, gtgo.default_value)
		FROM game_type_game_option gtgo
		LEFT JOIN game_option go ON go.id = gtgo.game_option_id
		LEFT JOIN game_game_option_value gov
			ON (gov.game_type_game_option_id = gtgo.id AND gov.game_id = $1)
		WHERE game_type_id = $2
		ORDER BY go.created_at ASC
	`), gameID, gameTypeID)

	if err != nil {
		return options, err
	}

	defer rows.Close()
	for rows.Next() {
		var name, valueType, value string
		rows.Scan(&name, &valueType, &value)
		options.Options[name] = getValueByType(valueType, value)
	}

	return options, err
}

// GetGameOptionRecordsForGameType returns all GameOption records for game type
func GetGameOptionRecordsForGameType(gameTypeID int64) ([]GameOption, error) {
	records := []GameOption{}

	rows, err := ConnectToDB().Query(fmt.Sprintf(`
		SELECT
			gtgo.id, gtgo.game_type_id, gtgo.game_option_id,
			gtgo.is_active, gtgo.default_value,
			go.name, go.label, go.type
		FROM game_type_game_option gtgo
		LEFT JOIN game_option go ON go.id = gtgo.game_option_id
		WHERE game_type_id = $1
		ORDER BY go.created_at ASC
	`), gameTypeID)
	if err != nil {
		return records, err
	}

	defer rows.Close()
	for rows.Next() {
		record, defaultValue := GameOption{}, ""
		rows.Scan(
			&record.GameTypeGameOptionID, &record.GameTypeID, &record.GameOptionID,
			&record.IsActive, &defaultValue, &record.Name, &record.Label, &record.Type)
		record.DefaultValue = getValueByType(record.Type, defaultValue)
		records = append(records, record)
	}

	return records, err
}

func (g GameOptions) validate() error {
	buyInMin, buyInMinOk := g.Options["buy_in_min"]
	buyInMax, buyInMaxOk := g.Options["buy_in_max"]

	if buyInMinOk && buyInMaxOk {
		bMax, bMin := qutils.ToI64(buyInMax), qutils.ToI64(buyInMin)
		if bMax < bMin {
			return fmt.Errorf("Game buy in max (%d) cannot be less than min (%d)", buyInMax, buyInMin)
		}
	}

	return nil
}

// Save writes the Game object to the database
func (g GameOptions) Save() error {
	// Validate all option types
	err := g.validate()
	if err != nil {
		return err
	}

	records, err := GetGameOptionRecordsForGameType(g.GameTypeID)
	if err != nil {
		return err
	}

	// TODO: remove any that no longer apply

	// save all options
	for i := range records {
		value, ok := g.Options[records[i].Name]
		if !ok {
			value = records[i].DefaultValue
		}

		recordValue := getStringForType(records[i].Type, value)
		g.save(records[i].GameTypeGameOptionID, recordValue)
	}

	return nil
}

func (g GameOptions) save(gameTypeGameOptionID int64, value string) error {
	_, err := ConnectToDB().Exec(`
		INSERT INTO game_game_option_value (game_id, game_type_game_option_id, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (game_id, game_type_game_option_id) DO UPDATE SET value = $4, updated_at = NOW();
	`, g.GameID, gameTypeGameOptionID, value, value)

	return err
}
