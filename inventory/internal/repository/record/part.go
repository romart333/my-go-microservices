package record

import (
	"time"

	"github.com/google/uuid"
)

type Part struct {
	UUID          uuid.UUID
	Name          string
	Description   string
	Price         int64
	PartType      string // HULL, ENGINE, SHIELD, WEAPON
	StockQuantity int64
	CreatedAt     time.Time
}
