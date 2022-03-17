package domain

type OperatorRepository interface {
	InsertAction(operatorId string, action Action) error
	FindOperator(operatorId string, opts string) (*Operator, error)
	InsertOperators(ops []interface{})
}
