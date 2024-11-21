package models

type Dl struct {
	ID      int     `gorm:"primaryKey"`
	Code    string  `gorm:"size:64;not null;unique"`
	Title   string  `gorm:"size:64;not null;unique"`
	Version float64 `gorm:"default:1.00"`
}

type SL struct {
	ID      int    `gorm:"primaryKey"`
	Code    string `gorm:"size:64;not null;unique"`
	Title   string `gorm:"size:64;not null;unique"`
	HasDl   bool
	Version float64 `gorm:"default:1.00"`
}

type Voucher struct {
	ID      int     `gorm:"primaryKey"`
	Number  string  `gorm:"size:64;not null;unique"`
	Version float64 `gorm:"default:1.00"`
}

type VoucherItem struct {
	ID        int     `gorm:"primaryKey"`
	VoucherID int     `gorm:"not null"`
	DlID      *int    `gorm:"null"`
	SlID      int     `gorm:"not null"`
	Dl        Dl      
	Sl        SL      
	Voucher   Voucher 
	Debit     float64 `gorm:"check:debit >= 0"`
	Credit    float64 `gorm:"check:credit >= 0"`
}
