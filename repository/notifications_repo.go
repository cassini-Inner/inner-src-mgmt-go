package repository

import (
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
)

type NotificationsRepo interface {
	Repository
	Create(recipientId, senderId, notificationType, jobId string) (*dbmodel.Notification, error)
	CreateWithTx(tx *sqlx.Tx, recipientId, senderId, notificationType, jobId string) (*dbmodel.Notification, error)

	Get(notificationId string) (*dbmodel.Notification, error)
	GetAllByReceiverId(receiverId string, afterId *string, limit int) ([]*dbmodel.Notification, error)
	GetNotificationCountForReceiver(receiverId string) (count int, err error)
}
