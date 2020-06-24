package service

import (
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
)

type NotificationsService interface {
	GetAllPaginated(receiverId string, afterId *string, limit int) ([]*gqlmodel.NotificationItem, error)
	GetNotificationCountForReceiver(receiverId string) (count int, err error)
}
