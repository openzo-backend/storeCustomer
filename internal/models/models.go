package models

import "time"

type StoreCustomer struct {
	ID        string    `json:"id" gorm:"primaryKey;size:36;unique"`
	StoreID   string    `json:"store_id" gorm:"size:36"`
	UserID    string    `json:"user_id" gorm:"size:36"`
	PhoneNo   string    `json:"phone_no"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
