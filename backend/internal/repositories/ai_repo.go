package repositories

import (
	"errors"

	"worldcup-mate/internal/database"
	"worldcup-mate/internal/models"

	"gorm.io/gorm"
)

func SaveGeneratedContent(content *models.AIGeneratedContent) error {
	var existing models.AIGeneratedContent
	err := database.DB.Where("cache_key = ?", content.CacheKey).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return database.DB.Create(content).Error
	}
	if err != nil {
		return err
	}
	content.ID = existing.ID
	return database.DB.Save(content).Error
}

func GetGeneratedContentByCacheKey(cacheKey string) (*models.AIGeneratedContent, error) {
	var content models.AIGeneratedContent
	err := database.DB.Where("cache_key = ?", cacheKey).First(&content).Error
	return &content, err
}

func SaveUsageLog(log *models.AIUsageLog) error {
	return database.DB.Create(log).Error
}

func CreateConversation(conv *models.AIConversation) error {
	return database.DB.Create(conv).Error
}

func UpdateConversation(conv *models.AIConversation) error {
	return database.DB.Save(conv).Error
}

func SaveMessage(msg *models.AIMessage) error {
	return database.DB.Create(msg).Error
}

func ListConversations(userID uint) ([]models.AIConversation, error) {
	var conversations []models.AIConversation
	err := database.DB.Where("user_id = ?", userID).
		Order("updated_at DESC").
		Find(&conversations).Error
	return conversations, err
}

func GetConversation(userID uint, conversationID uint) (*models.AIConversation, error) {
	var conv models.AIConversation
	err := database.DB.Where("id = ? AND user_id = ?", conversationID, userID).First(&conv).Error
	return &conv, err
}

func GetConversationWithMessages(userID uint, conversationID uint) (*models.AIConversation, []models.AIMessage, error) {
	conv, err := GetConversation(userID, conversationID)
	if err != nil {
		return nil, nil, err
	}
	var messages []models.AIMessage
	err = database.DB.Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Order("created_at ASC").
		Find(&messages).Error
	return conv, messages, err
}

func ListRecentMessages(userID uint, conversationID uint, limit int) ([]models.AIMessage, error) {
	var messages []models.AIMessage
	if limit <= 0 {
		limit = 8
	}
	err := database.DB.Where("conversation_id = ? AND user_id = ?", conversationID, userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, nil
}

func DeleteConversation(userID uint, conversationID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("conversation_id = ? AND user_id = ?", conversationID, userID).Delete(&models.AIMessage{}).Error; err != nil {
			return err
		}
		result := tx.Where("id = ? AND user_id = ?", conversationID, userID).Delete(&models.AIConversation{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
