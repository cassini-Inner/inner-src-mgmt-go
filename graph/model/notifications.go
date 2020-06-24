package model

import dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"

type NotificationItem struct {
	ID          string           `json:"id"`
	Recipient   *User            `json:"recipient"`
	Sender      *User            `json:"sender"`
	Type        NotificationType `json:"type"`
	Read        bool             `json:"read"`
	Job         *Job             `json:"job"`
	TimeCreated string           `json:"timeCreated"`
}

func (n *NotificationItem) MapDbToGql(notification *dbmodel.Notification) {
	n.ID = notification.Id
	n.Recipient = &User{ID: notification.RecipientId}
	n.Sender = &User{ID: notification.SenderId}
	n.Read = notification.Read
	n.Job = &Job{ID: notification.JobId}
	n.TimeCreated = notification.TimeCreated

	switch notification.Type {
	case dbmodel.NotificationTypeApplicationCreated:
		n.Type = NotificationTypeApplicationCreated
		break
	case dbmodel.NotificationTypeApplicationAccepted:
		n.Type = NotificationTypeApplicationAccepted
		break
	case dbmodel.NotificationTypeApplicationRejected:
		n.Type = NotificationTypeApplicationRejected
		break
	case dbmodel.NotificationTypeApplicationWithdrawn:
		n.Type = NotificationTypeApplicationWithdrawn
		break
	case dbmodel.NotificationTypeApplicationRemoved:
		n.Type = NotificationTypeApplicationRemoved
		break
	case dbmodel.NotificationTypeCommentAdded:
		n.Type = NotificationTypeCommentAdded
		break
	case dbmodel.NotificationTypeMilestoneCompleted:
		n.Type = NotificationTypeMilestoneCompleted
		break
	}
}
