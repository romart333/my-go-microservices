package input

type CreateOrderInput struct {
	HullUUID   string
	EngineUUID string
	ShieldUUID *string
	WeaponUUID *string
}

func (r *CreateOrderInput) PartUUIDs() []string {
	uuids := []string{r.HullUUID, r.EngineUUID}
	if r.ShieldUUID != nil {
		uuids = append(uuids, *r.ShieldUUID)
	}
	if r.WeaponUUID != nil {
		uuids = append(uuids, *r.WeaponUUID)
	}
	return uuids
}
