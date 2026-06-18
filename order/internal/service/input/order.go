package input

import "github.com/google/uuid"

type CreateOrderInput struct {
	HullUUID   uuid.UUID
	EngineUUID uuid.UUID
	ShieldUUID *uuid.UUID
	WeaponUUID *uuid.UUID
}

func (i *CreateOrderInput) PartUUIDs() uuid.UUIDs {
	uuids := uuid.UUIDs{i.HullUUID, i.EngineUUID}
	if i.ShieldUUID != nil {
		uuids = append(uuids, *i.ShieldUUID)
	}
	if i.WeaponUUID != nil {
		uuids = append(uuids, *i.WeaponUUID)
	}
	return uuids
}
