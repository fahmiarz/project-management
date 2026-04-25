package services

import (
	"errors"
	"fmt"
	"sort"
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

	GetByID(id uint) (*models.Card, error)
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

func (s *cardService) Create(card *models.Card, listPublicID string) error {
	// 1.Ambil list dari listPublicID
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found: %w", err)
	}

	// 2. Set list internal id
	card.ListID = list.InternalID

	// 3. Generate public_id jika belum ada
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
	if err := tx.Create(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create card: %w", err)
	}

	// 6. Update atau buat card_position
	var position models.CardPosition
	if err := tx.Model(&models.CardPosition{}).
		Where("list_internal_id = ?", list.InternalID).
		First(&position).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Buat baru jika belum ada
			position = models.CardPosition{
				PublicID:  uuid.New(),
				ListID:    list.InternalID,
				CardOrder: types.UUIDArray{card.PublicID},
			}

			if err := tx.Create(&position).Error; err != nil {
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
			Update("card_order", position.CardOrder).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update card position: %w", err)
		}
	}

	// 7. Commit transaksi
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

func (s *cardService) Update(card *models.Card, listPublicID string) error {
	// ambil card lama
	existingCard, err := s.cardRepo.FindByPublicID(card.PublicID.String())
	if err != nil {
		return fmt.Errorf("card not found: %w", err)
	}

	// ambil list baru
	newList, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return fmt.Errorf("list not found : %w", err)
	}

	// mulai trx
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// jika pindah list -> hapus dari posisi list lama & tambah ke list baru
	if existingCard.ListID != newList.InternalID {
		// hapus dari list lama
		var oldPos models.CardPosition
		if err := tx.Where("list_internal_id = ?", existingCard.ListID).First(&oldPos).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return fmt.Errorf("failed to get old card position: %w", err)
			}
		} else {
			filtered := make(types.UUIDArray, 0, len(oldPos.CardOrder))
			for _, id := range oldPos.CardOrder {
				if id != existingCard.PublicID {
					filtered = append(filtered, id)
				}
			}
			// update
			if err := tx.Model(&models.CardPosition{}).
				Where("internal_id = ?", oldPos.InternalID).
				Update("card_order", types.UUIDArray(filtered)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update old card position: %w", err)
			}
		}

		// tambah ke list baru
		var newPos models.CardPosition
		res := tx.Where("list_internal_id = ?", newList.InternalID).First(&newPos)
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			newPos = models.CardPosition{
				PublicID:  uuid.New(),
				ListID:    newList.InternalID,
				CardOrder: types.UUIDArray{existingCard.PublicID},
			}
			if err := tx.Create(&newPos).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to create card position for new list: %w", err)
			}
		} else if res.Error == nil {
			updateOrder := append(newPos.CardOrder, existingCard.PublicID)
			if err := tx.Model(&models.CardPosition{}).
				Where("internal_id = ?", newPos.InternalID).
				Update("card_order", types.UUIDArray(updateOrder)).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to update card position: %w", err)
			}
		} else {
			tx.Rollback()
			return fmt.Errorf("failed to get new card position: %w", res.Error)
		}
	}

	// update data card
	card.InternalID = existingCard.InternalID
	card.PublicID = existingCard.PublicID
	card.ListID = existingCard.ListID

	if err := tx.Save(card).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update card: %w", err)
	}

	// commit trx
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction commit failed: %w", err)
	}

	return nil
}

func (s *cardService) Delete(id uint) error {
	return s.cardRepo.Delete(id)
}

func (s *cardService) GetByListID(listPublicID string) ([]models.Card, error) {
	//verifikasi listnya ada
	list, err := s.listRepo.FindByPublicID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("list not found: %w", err)
	}
	//ambil card position
	position, err := s.cardRepo.FindCardPositionByListID(list.InternalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card position : %w", err)
	}
	//ambil semua card di list tersebut
	cards, err := s.cardRepo.FindByListID(listPublicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cards : %w", err)
	}
	//sorting
	if position !=nil && len (position.CardOrder) > 0 {
		cards = sortCardByPosition(cards, position.CardOrder)
	}
	return cards, nil
}
func sortCardByPosition( cards []models.Card, order [] uuid.UUID) []models.Card {
	//buat map untuk pencarian cepat
	orderMap := make(map[uuid.UUID]int)
	for i, id := range order {
		orderMap[id] = i
	}
	defaultIndex := len(order)
	//sorting slice
	sort.SliceStable(cards, func(i, j int) bool {
		idxI , okI := orderMap[cards[i].PublicID]
		if !okI {
			idxI = defaultIndex
		}
		idxJ, okJ := orderMap[cards[j].PublicID]
		if !okJ {
			idxJ = defaultIndex
		}
		if idxI == idxJ {
			return cards[i].CreatedAt.Before(cards[j].CreatedAt)
		}
		return idxI < idxJ
	})
	return cards
}

func (s *cardService) GetByID(id uint) (*models.Card, error) {
	return s.cardRepo.FindByID(id)
}

func (s *cardService) GetByPublicID(publicID string) (*models.Card, error) {
	return s.cardRepo.FindByPublicID(publicID)
}
