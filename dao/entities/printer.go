package entities

type Printer struct {
	SpaceAsset
	Name  string `gorm:"type:varchar(64)"`
	Sn    string `gorm:"type:varchar(64)"`
	Type  string `gorm:"type:varchar(64)"`
	Model string `gorm:"type:varchar(64)"`
}
