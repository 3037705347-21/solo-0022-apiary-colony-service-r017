package domain

import (
	"strings"
	"time"
)

// Apiary 表示一座蜂场。
type Apiary struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	District  string    `json:"district"`
	Keeper    string    `json:"keeper"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewApiary(id, name, district, keeper string, createdAt time.Time) (Apiary, error) {
	if err := ValidateIdentifier("id", id, 3, 40); err != nil {
		return Apiary{}, err
	}
	if err := ValidateText("name", name, 1, 60); err != nil {
		return Apiary{}, err
	}
	if err := ValidateText("district", district, 1, 60); err != nil {
		return Apiary{}, err
	}
	if err := ValidateText("keeper", keeper, 1, 40); err != nil {
		return Apiary{}, err
	}
	return Apiary{
		ID:        TrimText(id),
		Name:      TrimText(name),
		District:  TrimText(district),
		Keeper:    TrimText(keeper),
		CreatedAt: createdAt.UTC(),
	}, nil
}

// ApiaryRegistrationKey 返回一次蜂场登记使用的唯一性键，
// 编号规范化为小写后再判重，避免大小写差异绕过唯一性。
func ApiaryRegistrationKey(id string) string {
	return strings.ToLower(strings.TrimSpace(id))
}
