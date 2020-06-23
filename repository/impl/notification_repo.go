package impl

import (
	"context"
	customErrors "github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
)

type NotificationsRepo struct {
	db *sqlx.DB
}

func NewNotificationsRepo(db *sqlx.DB) *NotificationsRepo {
	return &NotificationsRepo{db: db}
}

func (n NotificationsRepo) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := n.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (n NotificationsRepo) CommitTx(ctx context.Context, tx *sqlx.Tx) error {
	err := tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (n NotificationsRepo) Create(recipientId, senderId, notificationType, jobId string) (*dbmodel.Notification, error) {
	if !validNotificationType(notificationType) {
		return nil, customErrors.ErrInvalidNotificationType
	}
	newNotification := &dbmodel.Notification{}
	err := n.db.QueryRowx(createNotificationQuery, recipientId, senderId, notificationType, jobId).StructScan(newNotification)
	if err != nil {
		return nil, err
	}

	return newNotification, nil
}

func (n NotificationsRepo) CreateWithTx(tx *sqlx.Tx, recipientId, senderId, notificationType, jobId string) (*dbmodel.Notification, error) {
	if !validNotificationType(notificationType) {
		return nil, customErrors.ErrInvalidNotificationType
	}

	newNotification := &dbmodel.Notification{}
	err := tx.QueryRowx(createNotificationQuery, recipientId, senderId, notificationType, jobId).StructScan(newNotification)

	if err != nil {
		return nil, err
	}

	return newNotification, nil
}

func (n NotificationsRepo) Get(notificationId string) (*dbmodel.Notification, error) {
	fetchedNotification := &dbmodel.Notification{}
	err := n.db.QueryRowx(getNotificationByIdQuery, notificationId).StructScan(fetchedNotification)

	if err != nil {
		return nil, err
	}

	return fetchedNotification, err
}

func (n NotificationsRepo) GetAllByReceiverId(receiverId string) ([]*dbmodel.Notification, error) {
	var result []*dbmodel.Notification
	rows, err := n.db.Queryx(getNotificationsByReceiverId, receiverId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notification *dbmodel.Notification
		err = rows.StructScan(notification)
		if err != nil {
			return nil, err
		}
		result = append(result, notification)
	}

	return result, nil
}

func validNotificationType(notificationType string) bool {
	if notificationType == dbmodel.NotificationTypeApplicationCreated ||
		notificationType == dbmodel.NotificationTypeApplicationAccepted ||
		notificationType == dbmodel.NotificationTypeApplicationRejected ||
		notificationType == dbmodel.NotificationTypeApplicationWithdrawn ||
		notificationType == dbmodel.NotificationTypeApplicationRemoved ||
		notificationType == dbmodel.NotificationTypeCommentAdded ||
		notificationType == dbmodel.NotificationTypeMilestoneCompleted {
		return true
	}
	return false
}

const (
	createNotificationQuery = "insert into notifications(recipient_id, sender_id, type, job_id) values ($1,$2, $3, $4) returning *"

	getNotificationByIdQuery = "select * from notification where id = $1"

	getNotificationsByReceiverId = "select * from notifications where receiver_id = $1"
)
