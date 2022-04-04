package application

import (
	"strconv"

	"github.com/google/uuid"
)

type PriceDto float64

func (p PriceDto) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(p), 'f', 2, 64)), nil
}

type CartDto struct {
	Id         uuid.UUID `json:"id"`
	CustomerId uuid.UUID `json:"customer_id"`
	Items      []ItemDto `json:"items"`
}

type ItemDto struct {
	ProductId uuid.UUID
	UnitPrice PriceDto
	Quantity  int
}

type CustomerDto struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type ProductDto struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	UnitPrice PriceDto  `json:"unit_price"`
}
