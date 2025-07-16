package model

type PaymentType int

const (
	PaymentTypeCash PaymentType = iota + 1 // 1
	PaymentTypeCard                        // 2
)

func (p PaymentType) String() string {
	switch p {
	case PaymentTypeCash:
		return "Наличные"
	case PaymentTypeCard:
		return "Карта"
	default:
		return "Неизвестно"
	}
}
