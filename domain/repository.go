package domain

type OperatorRepository interface {
	InsertAction(operatorId string, action Action) error
	FindOperator(operatorId string, opts OperatorFilter) (*Operator, error)
	UpdateOperator(operatorId string, opts OperatorFilter) error
	GetOperators(opts OperatorFilter) ([]*Operator, error)
	InsertOperators(ops []interface{})
}
