package domain

type OperatorService interface {
	InsertAction(operatorId string, action Action) error
	FindOperator(operatorId string, opts string) (*Operator, error)
}