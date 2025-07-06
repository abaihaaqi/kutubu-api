package repository

import (
	"github.com/abaihaaqi/kutubu-api/model"

	"gorm.io/gorm"
)

type BookRepository interface {
	GetAll(userID uint) ([]model.Book, error)
	GetByID(id uint, userID uint) (*model.Book, error)
	Create(book *model.Book) error
	Update(book *model.Book) error
	Delete(book *model.Book) error
	Search(keyword string, userID uint) ([]model.Book, error)
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db}
}

func (r *bookRepository) GetAll(userID uint) ([]model.Book, error) {
	var books []model.Book
	err := r.db.Where("user_id = ?", userID).Find(&books).Error
	return books, err
}

func (r *bookRepository) GetByID(id uint, userID uint) (*model.Book, error) {
	var book model.Book
	err := r.db.Where("user_id = ?", userID).First(&book, id).Error
	return &book, err
}

func (r *bookRepository) Create(book *model.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) Update(book *model.Book) error {
	return r.db.Save(book).Error
}

func (r *bookRepository) Delete(book *model.Book) error {
	return r.db.Delete(book).Error
}

func (r *bookRepository) Search(keyword string, userID uint) ([]model.Book, error) {
	var books []model.Book
	query := "%" + keyword + "%"
	err := r.db.Where("(user_id = ?) AND (title ILIKE ? OR author ILIKE ?)", userID, query, query).Find(&books).Error
	return books, err
}
