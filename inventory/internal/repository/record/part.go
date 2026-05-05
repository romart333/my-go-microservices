package record

import "time"

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         int64
	PartType      string // HULL, ENGINE, SHIELD, WEAPON
	StockQuantity int64
	CreatedAt     time.Time
}
