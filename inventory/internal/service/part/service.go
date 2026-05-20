package part

type service struct {
	partRepository PartRepository
}

func NewPartService(partRepository PartRepository) *service {
	return &service{partRepository: partRepository}
}
