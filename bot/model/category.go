package model

import "encoding/json"

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
		return "материалы"
	case CategoryLabor:
		return "работа/оплата мастерам"
	case CategoryTools:
		return "инструменты"
	case CategoryFurniture:
		return "мебель и техника"
	case CategoryUnexpected:
		return "прочее/непредвиденное"
	default:
		return "неизвестно"
	}
}

func (c Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}
