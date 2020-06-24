package model

const (
	ApplicationStatusWithdrawn = "withdrawn"
	ApplicationStatusPending   = "pending"
	ApplicationStatusRejected  = "rejected"
	ApplicationStatusAccepted  = "accepted"

	NotificationTypeApplicationCreated   = "NOTIFICATION_APPLICATION_CREATED"
	NotificationTypeApplicationAccepted  = "NOTIFICATION_APPLICATION_ACCEPTED"
	NotificationTypeApplicationRejected  = "NOTIFICATION_APPLICATION_REJECTED"
	NotificationTypeApplicationWithdrawn = "NOTIFICATION_APPLICATION_WITHDRAWN"
	NotificationTypeApplicationRemoved   = "NOTIFICATION_APPLICATION_REMOVED"
	NotificationTypeCommentAdded         = "NOTIFICATION_COMMENT_ADDED"
	NotificationTypeMilestoneCompleted   = "NOTIFICATION_MILESTONE_COMPLETED"

	MilestoneStatusCompleted = "completed"
	MilestoneStatusOngoing   = "ongoing"
	MilestoneStatusOpen      = "open"
)
