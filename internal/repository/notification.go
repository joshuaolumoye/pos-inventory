package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/joshuaolumoye/pos-backend/internal/domain"
)

type NotificationRepo struct {
	DB *sqlx.DB
}

func (r *NotificationRepo) CreateNotification(n *domain.Notification) error {
	query := `INSERT INTO notifications (id, business_id, product_id, notification_type, message, is_read, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, n.ID, n.BusinessID, n.ProductID, n.NotificationType, n.Message, n.IsRead, n.CreatedAt)
	return err
}

func (r *NotificationRepo) GetNotifications(businessID string, unreadOnly bool, limit, offset int) ([]*domain.Notification, error) {
	base := `SELECT id, business_id, product_id, notification_type, message, is_read, created_at FROM notifications WHERE business_id = ?`
	args := []interface{}{businessID}
	if unreadOnly {
		base += " AND is_read = 0"
	}
	base += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)
	rows, err := r.DB.Queryx(base, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var notifs []*domain.Notification
	for rows.Next() {
		n := &domain.Notification{}
		err := rows.Scan(&n.ID, &n.BusinessID, &n.ProductID, &n.NotificationType, &n.Message, &n.IsRead, &n.CreatedAt)
		if err != nil {
			return nil, err
		}
		notifs = append(notifs, n)
	}
	return notifs, nil
}

func (r *NotificationRepo) MarkNotificationRead(notificationID, businessID string) error {
	query := `UPDATE notifications SET is_read = 1 WHERE id = ? AND business_id = ?`
	_, err := r.DB.Exec(query, notificationID, businessID)
	return err
}

func (r *NotificationRepo) ExistsUnreadLowStockNotification(businessID, productID string) (bool, error) {
	query := `SELECT COUNT(1) FROM notifications WHERE business_id = ? AND product_id = ? AND notification_type = 'low_stock' AND is_read = 0`
	var count int
	err := r.DB.Get(&count, query, businessID, productID)
	return count > 0, err
}
