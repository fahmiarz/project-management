package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/fahmiarz/project-management/config"
	"github.com/fahmiarz/project-management/models"
	"github.com/fahmiarz/project-management/models/types"
	"github.com/fahmiarz/project-management/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardService interface {
	Create(card *models.Card, listPublicID string) error
	Update(card *models.Card, listPublicID string) error
	Delete(id uint) error

	GetByID(id uint) (*models.Card,error)
	GetByPublicID(publicID string) (*models.Card, error)
	GetByListID(listPublicID string) ([]models.Card, error)
}

type cardService struct {
	cardRepo repositories.CardRepository
	listRepo repositories.ListRepository
	userRepo repositories.UserRepository
}

func NewCardService(
	cardRepo repositories.CardRepository,
	listRepo repositories.ListRepository,
	userRepo repositories.UserRepository,
) CardService {
	return &cardService{cardRepo, listRepo, userRepo}
}

func (s* cardService) 	Create(card *models.Card, listPublicID string) error {
	// 1.Ambil list dari listPublicID
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err !=nil {
		return fmt.Errorf("list not found: %w", err)
	}

	// 2. Set list internal id
	card.ListID = list.InternalID

	//3. Generate public_id jika belum ada
	if card.PublicID == uuid.Nil {
		card.PublicID = uuid.New()
	}
	card.CreatedAt = time.Now()

	// 4. Mulai Transaksi
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 5. Simpan card
	if err := tx.Create(card).Error; err !=nil {
		tx.Rollback()
		return fmt.Errorf("failed to create card: %w", err)
	}

	// 6. Update atau buat card_position
	var position models.CardPosition
	if err := tx.Model(&models.CardPosition{}).
		Where("list_internal_id = ?", list.InternalID).
		First(&position).Error; err !=nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Buat baru jika belum ada
				position = models.CardPosition{
					PublicID: uuid.New(),
					ListID : list.InternalID,
					CardOrder: types.UUIDArray{card.PublicID},
				}

				if err := tx.Create(&position).Error; err !=nil {
					tx.Rollback()
					return fmt.Errorf("failed to create card position: %w", err)
				}
			} else {
				tx.Rollback()
				return fmt.Errorf("failed to get card position: %w", err)
			}
		} else {
			// Tambah card baru ke urutan
			position.CardOrder = append(position.CardOrder, card.PublicID)
			if err := tx.Model(&models.CardPosition{}).
				Where("internal_id = ?", position.InternalID).
				Update("card_order", position.CardOrder).Error; err !=nil {
					tx.Rollback()
					return fmt.Errorf("failed to update card position: %w", err)
				}
		}
		// 7. Commit transaksi 
		if err := tx.Commit().Error; err !=nil {
			return fmt.Errorf("transaction commit failed: %w", err)
		}

		return nil
	}

	func (s *cardService) Update(card *models.Card, listPublicID string) error {}