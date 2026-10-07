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

// ApiaryRegistrationKey 返回一次蜂场登记使用的唯一性键。
func ApiaryRegistrationKey(id string) string {
	return strings.ToLower(id)
}

// SameApiaryIdentity 判断两个入口编号是否指向同一座蜂场。
func SameApiaryIdentity(left, right string) bool {
	return ApiaryRegistrationKey(left) == ApiaryRegistrationKey(right)
}

// RegistrationMatches 判断已有登记是否与本次提交描述同一份业务事实。
func (a Apiary) RegistrationMatches(name, district, keeper string) bool {
	return a.Name == TrimText(name) &&
		a.District == TrimText(district) &&
		a.Keeper == TrimText(keeper)
}
