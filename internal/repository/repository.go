package repository

import (
	"log"

	"github.com/google/uuid"
	"github.com/tanush-128/openzo_backend/storeCustomer/internal/models"

	"gorm.io/gorm"
)

type StoreCustomerRepository interface {
	CreateStoreCustomer(StoreCustomer models.StoreCustomer) (models.StoreCustomer, error)
	GetStoreCustomerByID(id string) (models.StoreCustomer, error)
	GetStoreCustomersByStoreID(storeID string) ([]models.StoreCustomer, error)
	UpdateStoreCustomer(StoreCustomer models.StoreCustomer) (models.StoreCustomer, error)
	// GetStoreCustomersByUserID(ctx *gin.Context, user_data_id string) ([]models.StoreCustomer, error)
	GetStoreCustomersByUserID(user_data_id string) ([]models.StoreCustomer, error)

	DeleteStoreCustomer(id string) error
	GetUserIdByPhone(phone string) (string, error)
	// ChangeStoreCustomerStatus(ctx *gin.Context, id string, status string) (models.StoreCustomer, error)
	// Add more methods for other StoreCustomer operations (GetStoreCustomerByEmail, UpdateStoreCustomer, etc.)

}

type storeCustomerRepository struct {
	db *gorm.DB
}

func NewStoreCustomerRepository(db *gorm.DB) StoreCustomerRepository {

	return &storeCustomerRepository{db: db}
}

func (r *storeCustomerRepository) CreateStoreCustomer(storeCustomer models.StoreCustomer) (models.StoreCustomer, error) {
	storeCustomer.ID = uuid.New().String()
	if storeCustomer.PhoneNo != "" {
		user_id, _ := r.GetUserIdByPhone(storeCustomer.PhoneNo)

		storeCustomer.UserID = user_id
	}
	tx := r.db.Create(&storeCustomer)

	if tx.Error != nil {
		return models.StoreCustomer{}, tx.Error
	}

	return storeCustomer, nil
}

func (r *storeCustomerRepository) GetStoreCustomerByID(id string) (models.StoreCustomer, error) {
	var StoreCustomer models.StoreCustomer
	tx := r.db.Where("id = ?", id).First(&StoreCustomer)
	if tx.Error != nil {
		return models.StoreCustomer{}, tx.Error
	}

	return StoreCustomer, nil
}

func (r *storeCustomerRepository) GetStoreCustomersByStoreID(storeID string) ([]models.StoreCustomer, error) {
	var StoreCustomers []models.StoreCustomer

	tx := r.db.Where("store_id = ?", storeID).Find(&StoreCustomers)
	if tx.Error != nil {
		return []models.StoreCustomer{}, tx.Error
	}

	return StoreCustomers, nil
}

func (r *storeCustomerRepository) GetStoreCustomersByUserID(user_data_id string) ([]models.StoreCustomer, error) {
	var StoreCustomers []models.StoreCustomer
	tx := r.db.Where("user_id = ?", user_data_id).Find(&StoreCustomers)

	if tx.Error != nil {
		return []models.StoreCustomer{}, tx.Error
	}
	// log.Printf("StoreCustomers: %+v", StoreCustomers)
	for _, storeCustomer := range StoreCustomers {
		log.Printf("StoreCustomer: %+v", storeCustomer.ID)
	}
	return StoreCustomers, nil
}

func (r *storeCustomerRepository) UpdateStoreCustomer(storeCustomer models.StoreCustomer) (models.StoreCustomer, error) {

	if storeCustomer.PhoneNo != "" {
		user_id, _ := r.GetUserIdByPhone(storeCustomer.PhoneNo)

		storeCustomer.UserID = user_id
	}
	tx := r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&storeCustomer)
	if tx.Error != nil {
		return models.StoreCustomer{}, tx.Error
	}

	return storeCustomer, nil
}

func (r *storeCustomerRepository) DeleteStoreCustomer(id string) error {
	tx := r.db.Where("id = ?", id).Delete(&models.StoreCustomer{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (r *storeCustomerRepository) GetUserIdByPhone(phone string) (string, error) {
	var user_id string
	tx := r.db.Table("users").Select("id").Where("phone = ?", phone).First(&user_id)

	if tx.Error != nil {
		return "", tx.Error
	}

	return user_id, nil
}

// Implement other repository methods (GetStoreCustomerByID, GetStoreCustomerByEmail, UpdateStoreCustomer, etc.) with proper error handling
