package models

type ClientService struct {
	ID              int    `gorm:"primaryKey;column:id"`
	Keyword         string `gorm:"column:keyword"`
	Shortcode       string `gorm:"column:shortcode"`
	TelcoID         string `gorm:"column:telcoid"`
	AdsID           string `gorm:"column:ads_id"`
	ClientPartnerID string `gorm:"column:client_partner_id;not null"`
	WapAocRefID     string `gorm:"column:wap_aoc_refid"`
	WapAocID        string `gorm:"column:wap_aoc_id"`
	WapAocMedia     string `gorm:"column:wap_aoc_media"`
	PostbackURL     string `gorm:"column:postback_url"`
	DNURL           string `gorm:"column:dn_url"`
	PostbackCounter int    `gorm:"column:postback_counter"` // Nullable integer
}
