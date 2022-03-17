package json

import (
	"encoding/json"
	"github.com/pkg/errors"
	"stablex/domain"
)

type Serializer struct{}

func (r *Serializer) Decode(input []byte) (*domain.Operator, error) {
	redirect := &domain.Operator{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}
	return redirect, nil
}

func (r *Serializer) Encode(input *domain.Operator) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}
	return rawMsg, nil
}
