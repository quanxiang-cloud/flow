package models

// BaseModel info
type BaseModel struct {
	ID string `json:"id"`

	CreatorID  string `json:"creatorId"`
	CreateTime string `json:"createTime"`
	ModifierID string `json:"modifierId"`
	ModifyTime string `json:"modifyTime"`

	CreatorName   string `gorm:"-" json:"creatorName"`
	CreatorAvatar string `gorm:"-" json:"creatorAvatar"`
	ModifierName  string `gorm:"-" json:"modifierName"`
}
