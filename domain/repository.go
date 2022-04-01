package domain

import "stablex/auth"

type OperatorRepository interface {
	InsertAction(operatorId string, action Action) error
	FindOperator(operatorId string, opts OperatorFilter) (*Operator, error)
	UpdateOperator(operatorId string, opts OperatorFilter) (*Operator, error)
	GetOperators(opts OperatorFilter) ([]*Operator, error)
	InsertOperators(ops []interface{})
}

type AuthRepository interface {
	CreateAuth(userid string, td *auth.TokenDetails) error
}
