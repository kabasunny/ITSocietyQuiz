package repositories

import (
	"errors"
	"gin-practice/models"

	"gorm.io/gorm"
)

type IItemRepository interface {
	FindAll() (*[]models.Item, error)
	FindById(itemId uint, userId uint) (*models.Item, error)
	Create(newItem models.Item) (*models.Item, error)
	Update(updateItem models.Item) (*models.Item, error)
	Delete(itemId uint, userId uint) error
}

type itemRepository struct {
	db *gorm.DB
}

func (r *itemRepository) Create(newItem models.Item) (*models.Item, error) {
	result := r.db.Create(&newItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newItem, nil
}

func (r *itemRepository) Delete(itemId uint, userId uint) error {
	deleteItem, err := r.FindById(itemId, userId)
	if err != nil {
		return err
	}
	result := r.db.Delete(&deleteItem)
	// 論理削除が行われる。
	// 物理削除を行う場合は、r.db.Unscoped().Delete(&deleteItem)とする
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *itemRepository) FindAll() (*[]models.Item, error) {
	var items []models.Item
	result := r.db.Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return &items, nil
}

func (r *itemRepository) FindById(itemId uint, userId uint) (*models.Item, error) {
	var item models.Item
	result := r.db.First(&item, "id = ? AND user_id = ?", itemId, userId)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil, errors.New("Item not found")
		}
		return nil, result.Error
	}
	return &item, nil
}

func (r *itemRepository) Update(updateItem models.Item) (*models.Item, error) {
	result := r.db.Save(&updateItem)
	if result.Error != nil {
		return nil, result.Error
	}
	return &updateItem, nil
}

func NewItemRepository(db *gorm.DB) IItemRepository { //ファクトリ関数
	return &itemRepository{db: db}
}
