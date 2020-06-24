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

func (n NotificationsRepo) GetAllByReceiverId(receiverId string, afterId *string, limit int) ([]*dbmodel.Notification, error) {
	var rows *sqlx.Rows

	if afterId == nil {
		query := n.db.Rebind(getNotificationsByReceiverIdWithoutAfter)
		rows, err := n.db.Queryx(query, receiverId, limit)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		return scanNotificationRows(rows)
	}

	query := n.db.Rebind(getNotificationsByReceiverIdWithAfter)

	rows, err := n.db.Queryx(query, receiverId, *afterId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanNotificationRows(rows)
}

func scanNotificationRows(rows *sqlx.Rows) (result []*dbmodel.Notification, err error) {
	for rows.Next() {
		scannedNotification := &dbmodel.Notification{}
		err = rows.StructScan(scannedNotification)
		if err != nil {
			return nil, err
		}

		result = append(result, scannedNotification)
	}

	return result, nil
}

func (n NotificationsRepo) GetNotificationCountForReceiver(receiverId string) (count int, err error) {
	err = n.db.QueryRowx(getNotificationCountForRecipient, receiverId).Scan(&count)
	if err != nil {
		return 0, nil
	}

	return count, nil
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

	getNotificationsByReceiverIdWithAfter = "select * from notifications where recipient_id = ? and id < ? order by time_created desc fetch first ? rows only"

	getNotificationsByReceiverIdWithoutAfter = "select * from notifications where recipient_id = ? order by time_created desc fetch first ? rows only"

	getNotificationCountForRecipient = "select count(*) from notifications where recipient_id = $1"
)
