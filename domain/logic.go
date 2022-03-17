package domain

type operatorService struct {
	operatorRepo OperatorRepository
}

func NewOperatorService(operatorRepo OperatorRepository) OperatorService {
	return &operatorService{
		operatorRepo,
	}
}

func (s *operatorService) InsertAction(operatorId string, action Action) error {
	return s.operatorRepo.InsertAction(operatorId, action)
}

func (s *operatorService) FindOperator(code string, opts string) (*Operator, error) {
	return s.operatorRepo.FindOperator(code, opts)
}
