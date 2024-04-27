package model

type MerchantResource struct {
	ID   int64  `db:"ID"`
	Name string `db:"Name"`
	Code string `db:"Code"`
}

type MerchantCleanResource struct {
	Name string `json:"name" validate:"required,alphanumunicode"`
	Code string `json:"code" validate:"required,alphanumunicode"`
}

func (m MerchantResource) Clean() MerchantCleanResource {
	return MerchantCleanResource{
		Name: m.Name,
		Code: m.Code,
	}
}
