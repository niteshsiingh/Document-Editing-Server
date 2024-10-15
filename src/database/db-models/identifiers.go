package dbmodels

type Base struct {
	B int `json:"_b"`
}

type IdentifierID struct {
	ID   uint                   `json:"id"`
	Base map[string]interface{} `gorm:"serializer:json" json:"_base"`
	C    []int                  `gorm:"serializer:json" json:"_c"`
	D    []int                  `gorm:"serializer:json" json:"_d"`
	S    []int                  `gorm:"serializer:json" json:"_s"`
}

type Identifier struct {
	Element string       `gorm:"type:string;not null" json:"elem"`
	ID      IdentifierID `gorm:"embedded" json:"id"`
}
