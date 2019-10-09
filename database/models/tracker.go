package models

import (
	"github.com/jinzhu/gorm"
)

// Tracker - Model representation of an autodl tracker file
type Tracker struct {
	gorm.Model
	Tracker string `sql:"type:varchar(256);not null"`
}

// // general
// Id   int    `gorm:"primary_key;not null"`
// Name string `sql:"type:text"`

// // subscriber
// IsSubscriber      bool `sql:"default:false"`
// SubscribedPackage Package
// SubscribedExpiry  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`

// // logic bools
// IsVip   bool `sql:"default:false"`
// IsAdmin bool `sql:"default:false"`

// // timestamps
// LastSeen  time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
// CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
