package domain

type OperatorSerializer interface {
	Decode(input []byte) (*Operator, error)
	Encode(input *Operator) ([]byte, error)
}
