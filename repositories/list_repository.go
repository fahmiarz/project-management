package repositories

import (
	"github.com/fahmiarz/project-management/config"
	"github.com/fahmiarz/project-management/models"
	"github.com/google/uuid"
)

//Buat interface
type ListRepository interface {
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePosition(boardPublicID string, position []string) error
	GetCardPosition(listPublicID string)([]uuid.UUID, error)
	FindByBoardID(boardID string) ([]models.List, error)
	FindByPublicID(publicID string) (*models.List, error)
	FindByID(id uint) (*models.List, error)
}

//Buat struct
type listRepository struct{

}

//Buat constructor
func NewListRepository() ListRepository {
	return &listRepository{}
}

//Buat Fungsi create
func (r *listRepository) Create(list *models.List) error {
	return config.DB.Create(list).Error
}

//Buat fungsi update
func (r *listRepository) Update(list *models.List) error {
	return config.DB.Model(&models.List{}).
	Where("public_id = ?", list.PublicID).Updates(map[string]interface{}{
		"title" : list.Title,
	}).Error
}

//Buat fungsi delete
func (r *listRepository) Delete(id uint) error {
	return config.DB.Delete(&models.List{}, id).Error
}

//Buat fungsi update position
func (r *listRepository) UpdatePosition(boardPublicID string, position []string) error {
	return config.DB.Model(&models.ListPosition{}). 
	Where("board_internal_id = (Select internal_id FROM boards Where public_id = ? )", boardPublicID).
	Update("list_order", position).Error
}

//Buat fungsi get card position
func (r* listRepository) GetCardPosition(listPublicID string)([]uuid.UUID, error) {
	var position models.CardPosition
	err := config.DB.Joins("JOIN lists ON lists.internal_id = card_position.list_internal_id").
	Where("list.public_id =?", listPublicID).Error
	return position.CardOrder, err
}

//Buat fungsi findbyboard id
func (r* listRepository) FindByBoardID(boardID string) ([]models.List, error) {
	var list []models.List
	err := config.DB.Where("board_public_id = ?", boardID).Order("internal_id ASC").Find(&list).Error
	return list, err
}

//buat fungsi find by public id
func (r* listRepository) FindByPublicID(publicID string) (*models.List, error) {
	var list models.List
	err := config.DB.Where("public_id = ?", publicID).First(&list).Error
	return &list, err
}

func (r* listRepository) FindByID(id uint) (*models.List, error) {
	var list models.List
	err := config.DB.First(&list, id).Error
	return &list, err
}


