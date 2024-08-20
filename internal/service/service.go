package service

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/storeCustomer/internal/models"
	"github.com/tanush-128/openzo_backend/storeCustomer/internal/repository"
)

type StoreCustomerService interface {

	//CRUD
	CreateStoreCustomer(ctx *gin.Context, req models.StoreCustomer) (models.StoreCustomer, error)
	GetStoreCustomerByID(ctx *gin.Context, id string) (models.StoreCustomer, error)
	GetStoreCustomersByStoreID(ctx *gin.Context, storeID string) ([]models.StoreCustomer, error)
	GetStoreCustomersByUserID(ctx *gin.Context, user_data_id string) ([]models.StoreCustomer, error)
	// ChangeStoreCustomerStatus(ctx *gin.Context, id string, status string) (models.StoreCustomer, error)
	UpdateStoreCustomer(ctx *gin.Context, req models.StoreCustomer) (models.StoreCustomer, error)
	DeleteStoreCustomer(ctx *gin.Context, id string) error
}

type storeCustomerService struct {
	storeCustomerRepository repository.StoreCustomerRepository
	kafkaProducer           *kafka.Producer
}

func NewStoreCustomerService(storeCustomerRepository repository.StoreCustomerRepository,
	producer *kafka.Producer,

) StoreCustomerService {
	return &storeCustomerService{storeCustomerRepository: storeCustomerRepository, kafkaProducer: producer}
}

func (s *storeCustomerService) CreateStoreCustomer(ctx *gin.Context, req models.StoreCustomer) (models.StoreCustomer, error) {

	createdStoreCustomer, err := s.storeCustomerRepository.CreateStoreCustomer(req)
	if err != nil {
		return models.StoreCustomer{}, err // Propagate error
	}

	// Produce message to Kafka
	// topic := "storeCustomers"
	// storeCustomerMsg, _ := json.Marshal(createdStoreCustomer)
	// msg := &kafka.Message{
	// 	TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	// 	Value:          storeCustomerMsg,
	// }
	// s.kafkaProducer.Produce(msg, nil)

	return createdStoreCustomer, nil
}

func (s *storeCustomerService) UpdateStoreCustomer(ctx *gin.Context, req models.StoreCustomer) (models.StoreCustomer, error) {
	updatedStoreCustomer, err := s.storeCustomerRepository.UpdateStoreCustomer(req)
	if err != nil {
		return models.StoreCustomer{}, err
	}

	// topic := "storeCustomers"
	// storeCustomerMsg, _ := json.Marshal(updatedStoreCustomer)
	// msg := &kafka.Message{
	// 	TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	// 	Value:          storeCustomerMsg,
	// }
	// s.kafkaProducer.Produce(msg, nil)

	return updatedStoreCustomer, nil
}

func (s *storeCustomerService) GetStoreCustomersByUserID(ctx *gin.Context, user_id string) ([]models.StoreCustomer, error) {
	storeCustomers, err := s.storeCustomerRepository.GetStoreCustomersByUserID(user_id)
	if err != nil {
		return []models.StoreCustomer{}, err
	}

	return storeCustomers, nil
}

func (s *storeCustomerService) DeleteStoreCustomer(ctx *gin.Context, id string) error {
	err := s.storeCustomerRepository.DeleteStoreCustomer(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *storeCustomerService) GetStoreCustomerByID(ctx *gin.Context, id string) (models.StoreCustomer, error) {
	storeCustomer, err := s.storeCustomerRepository.GetStoreCustomerByID(id)
	if err != nil {
		return models.StoreCustomer{}, err
	}

	return storeCustomer, nil
}

func (s *storeCustomerService) GetStoreCustomersByStoreID(ctx *gin.Context, storeID string) ([]models.StoreCustomer, error) {
	storeCustomers, err := s.storeCustomerRepository.GetStoreCustomersByStoreID(storeID)
	if err != nil {
		return []models.StoreCustomer{}, err
	}

	return storeCustomers, nil
}
