package model

import (
	"encoding/json"
)

type Category byte

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

func (c Category) IsValid() bool {
	switch c {
	case CategoryMaterials,
		CategoryLabor,
		CategoryTools,
		CategoryFurniture,
		CategoryUnexpected:
		return true

	default:
		return false
	}
}

func Categories() []Category {
	return []Category{
		CategoryMaterials,
		CategoryLabor,
		CategoryTools,
		CategoryFurniture,
		CategoryUnexpected,
	}
}

func (c Category) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}
