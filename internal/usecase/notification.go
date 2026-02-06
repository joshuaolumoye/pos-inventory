package usecase

import (
	"time"

	"github.com/joshuaolumoye/pos-backend/internal/domain"
	"github.com/joshuaolumoye/pos-backend/pkg/utils"
)

type NotificationUsecase struct {
	NotificationRepo domain.NotificationRepository
}

func (u *NotificationUsecase) CreateLowStockNotification(businessID, productID, productName string, quantity int) error {
	exists, err := u.NotificationRepo.ExistsUnreadLowStockNotification(businessID, productID)
	if err != nil {
		return err
	}
	if exists {
		return nil // Don't create duplicate unread notification
	}
	msg := productName + " is running low, only " + utils.Itoa(quantity) + " items left."
	n := &domain.Notification{
		ID:               utils.GenerateUUID(),
		BusinessID:       businessID,
		ProductID:        productID,
		NotificationType: "low_stock",
		Message:          msg,
		IsRead:           false,
		CreatedAt:        time.Now().Unix(),
	}
	return u.NotificationRepo.CreateNotification(n)
}

func (u *NotificationUsecase) GetNotifications(businessID string, unreadOnly bool, limit, offset int) ([]*domain.Notification, error) {
	return u.NotificationRepo.GetNotifications(businessID, unreadOnly, limit, offset)
}

func (u *NotificationUsecase) MarkNotificationRead(notificationID, businessID string) error {
	return u.NotificationRepo.MarkNotificationRead(notificationID, businessID)
}
