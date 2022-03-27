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

func (s *operatorService) FindOperator(code string, opts OperatorFilter) (*Operator, error) {
	return s.operatorRepo.FindOperator(code, opts)
}

func (s *operatorService) GetOperators(opts OperatorFilter) ([]*Operator, error) {
	return s.operatorRepo.GetOperators(opts)
}

func (s *operatorService) UpdateOperator(operatorId string, opts OperatorFilter) (*Operator, error) {
	return s.operatorRepo.UpdateOperator(operatorId, opts)
}
