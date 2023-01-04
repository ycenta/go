package payment

import (
	"errors"
	Product "go/src/product"

	"gorm.io/gorm"
)

type Repository interface {
	Create(payment Payment) (Payment, error)
	GetAll() ([]Payment, error)
	GetById(id int) (Payment, error)
	Update(id int, input InputPayment) (Payment, error)
	Delete(id int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(payment Payment) (Payment, error) {
	//verify that product exists
	err := r.db.Where("id = ?", payment.ProductID).First(&payment.Product).Error
	if err != nil {
		return payment, err
	}

	err = r.db.Create(&payment).Error
	if err != nil {
		return payment, err
	}

	return payment, nil
}

func (r *repository) GetAll() ([]Payment, error) {
	var payments []Payment
	//preload => load products linked
	err := r.db.Preload("Product").Find(&payments).Error
	if err != nil {
		return payments, err
	}

	return payments, nil
}

func (r *repository) GetById(id int) (Payment, error) {
	var payment Payment

	//preload => load products linked
	err := r.db.Preload("Product").Where(&Payment{ID: id}).First(&payment).Error
	if err != nil {
		return payment, err
	}

	return payment, nil
}

func (r *repository) Update(id int, input InputPayment) (Payment, error) {
	payment, err := r.GetById(id)
	if err != nil {
		return payment, err
	}

	//get the linked product
	var product Product.Product
	err = r.db.Where(&Product.Product{ID: input.ProductID}).First(&product).Error
	if err != nil {
		return payment, err
	}

	payment.ProductID = input.ProductID
	payment.PricePaid = input.PricePaid
	err = r.db.Save(&payment).Error
	// err = r.db.Model(&payment).Updates(Payment{ProductID: input.ProductID, PricePaid: input.PricePaid}).Error
	if err != nil {
		return payment, err
	}

	return payment, nil
}

func (r *repository) Delete(id int) error {
	payment := &Payment{ID: id}
	tx := r.db.Delete(payment)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("Payment not found")
	}

	return nil
}
