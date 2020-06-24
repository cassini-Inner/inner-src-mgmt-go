package impl

import (
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/repository"
)

type NotificationsService struct {
	notificationsRepo repository.NotificationsRepo
}

func NewNotificationsService(notificationsRepo repository.NotificationsRepo) *NotificationsService {
	return &NotificationsService{notificationsRepo: notificationsRepo}
}

func (n NotificationsService) GetAllPaginated(receiverId string, afterId *string, limit int) (result []*gqlmodel.NotificationItem, err error) {

	notifications, err := n.notificationsRepo.GetAllByReceiverId(receiverId, afterId, limit)
	if err != nil {
		return nil, err
	}

	for _, notification := range notifications {
		gqlNotification := &gqlmodel.NotificationItem{}
		gqlNotification.MapDbToGql(notification)
		result = append(result, gqlNotification)
	}

	return result, nil
}

func (n NotificationsService) GetNotificationCountForReceiver(receiverId string) (int, error) {
	count, err := n.notificationsRepo.GetNotificationCountForReceiver(receiverId)
	if err != nil {
		return 0, err
	}
	return count, nil
}
