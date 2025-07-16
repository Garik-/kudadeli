package model

import "encoding/json"

type PaymentType int

const (
	PaymentTypeCash PaymentType = iota + 1 // 1
	PaymentTypeCard                        // 2
)

func (p PaymentType) String() string {
	switch p {
	case PaymentTypeCash:
		return "наличные"
	case PaymentTypeCard:
		return "карта"
	default:
		return "неизвестно"
	}
}

func (p PaymentType) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}
