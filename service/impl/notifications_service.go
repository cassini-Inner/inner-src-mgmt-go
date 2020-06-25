package impl

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/repository"
)

type NotificationsService struct {
	notificationsRepo repository.NotificationsRepo
}

func NewNotificationsService(notificationsRepo repository.NotificationsRepo) *NotificationsService {
	return &NotificationsService{notificationsRepo: notificationsRepo}
}


func (n NotificationsService) GetUnreadNotificationCountForReceiver(recipientId string) (int, error) {
	return n.notificationsRepo.GetUnreadNotificationCountForReceiver(recipientId)
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

func (n NotificationsService) MarkAllNotificationsRead(ctx context.Context, recipientId string) ([]*gqlmodel.NotificationItem, error) {
	tx, err := n.notificationsRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	markedNotifications, err := n.notificationsRepo.MarkAllUserNotificationsReadWithTx(tx, recipientId)
	if err != nil {
		return nil, err
	}

	gqlNotifications := make([]*gqlmodel.NotificationItem, 0)
	for _, notification := range markedNotifications {
		tempGqlNotification := &gqlmodel.NotificationItem{}
		tempGqlNotification.MapDbToGql(notification)
		gqlNotifications = append(gqlNotifications, tempGqlNotification)
	}

	err = n.notificationsRepo.CommitTx(ctx, tx)
	if err != nil {
		return nil, err
	}

	return gqlNotifications, nil
}

func (n NotificationsService) MarkNotificationsReadById(ctx context.Context, recipientId string, notificationIds ...string) ([]*gqlmodel.NotificationItem, error) {
	tx, err := n.notificationsRepo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	markedNotifications, err := n.notificationsRepo.MarkUserNotificationsReadWithTx(tx, recipientId, notificationIds...)
	if err != nil {
		return nil, err
	}

	gqlNotifications := make([]*gqlmodel.NotificationItem, 0)
	for _, notification := range markedNotifications {
		tempGqlNotification := &gqlmodel.NotificationItem{}
		tempGqlNotification.MapDbToGql(notification)
		gqlNotifications = append(gqlNotifications, tempGqlNotification)
	}

	err = n.notificationsRepo.CommitTx(ctx, tx)
	if err != nil {
		return nil, err
	}

	return gqlNotifications, nil
}
