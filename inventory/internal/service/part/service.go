package part

type PartService struct {
	partRepository PartRepository
}

func NewPartService(partRepository PartRepository) *PartService {
	return &PartService{partRepository: partRepository}
}
