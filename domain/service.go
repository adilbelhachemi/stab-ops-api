package domain

type OperatorService interface {
	InsertAction(operatorId string, action Action) error
	FindOperator(operatorId string, opts OperatorFilter) (*Operator, error)
	GetOperators(opts OperatorFilter) ([]*Operator, error)
	UpdateOperator(operatorId string, opts OperatorFilter) error
}
