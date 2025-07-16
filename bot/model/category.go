package model

type Category int

const (
	CategoryMaterials  Category = iota + 1 // 1
	CategoryLabor                          // 2
	CategoryTools                          // 3
	CategoryFurniture                      // 4
	CategoryUnexpected                     // 5
)

func (c Category) String() string {
	switch c {
	case CategoryMaterials:
		return "Материалы"
	case CategoryLabor:
		return "Работа/оплата мастерам"
	case CategoryTools:
		return "Инструменты"
	case CategoryFurniture:
		return "Мебель и техника"
	case CategoryUnexpected:
		return "Прочее/непредвиденное"
	default:
		return "Неизвестно"
	}
}
