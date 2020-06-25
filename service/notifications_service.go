package service

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

type NotificationsService interface {
	GetAllPaginated(receiverId string, afterId *string, limit int) ([]*gqlmodel.NotificationItem, error)
	GetNotificationCountForReceiver(receiverId string) (count int, err error)
	MarkAllNotificationsRead(ctx context.Context, id string) ([]*gqlmodel.NotificationItem, error)
	MarkNotificationsReadById(ctx context.Context, recipientId string, notificationIds ...string) ([]*gqlmodel.NotificationItem, error)
	GetUnreadNotificationCountForReceiver(recipientId string) (int, error)
}
