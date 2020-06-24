package resolver

import (
	"context"
	"encoding/base64"
	"github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/resolver/dataloader"
	"github.com/cassini-Inner/inner-src-mgmt-go/middleware"
)

func (r *queryResolver) ViewerNotifications(ctx context.Context, limit int, after *string) (*gqlmodel.NotificationConnection, error) {
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, custom_errors.ErrUserNotAuthenticated
	}
	var afterString *string
	if after != nil {
		if *after == "" {
			return nil, custom_errors.ErrInvalidCursor
		}
		decodedAfter, err := base64.StdEncoding.DecodeString(*after)
		if err != nil {
			return nil, custom_errors.ErrInvalidCursor
		}
		cursorStr := string(decodedAfter)
		afterString = &cursorStr
	}

	notifications, err := r.NotificationsService.GetAllPaginated(user.Id, afterString, limit+1)
	if err != nil {
		return nil, err
	}

	totalNotificationsForUser, err := r.NotificationsService.GetNotificationCountForReceiver(user.Id)

	if err != nil {
		return nil, err
	}

	notificationEdges := make([]*gqlmodel.NotificationEdge, 0)
	for i, notification := range notifications {
		if i < limit {
			notificationEdges = append(notificationEdges, &gqlmodel.NotificationEdge{
				Node:   notification,
				Cursor: base64.StdEncoding.EncodeToString([]byte(notification.ID)),
			})
		}
	}

	var endCursor *string
	if len(notificationEdges) > 0 {
		endCursor = &notificationEdges[len(notificationEdges)-1].Cursor
	}

	return &gqlmodel.NotificationConnection{
		TotalCount: totalNotificationsForUser,
		Edges:      notificationEdges,
		PageInfo: &gqlmodel.PageInfo{
			HasNextPage: len(notifications) > limit,
			EndCursor:   endCursor,
		},
	}, nil
}

func (r *notificationItemResolver) Recipient(ctx context.Context, obj *gqlmodel.NotificationItem) (*gqlmodel.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.Recipient.ID)
}

func (r *notificationItemResolver) Sender(ctx context.Context, obj *gqlmodel.NotificationItem) (*gqlmodel.User, error) {
	return dataloader.GetUserByUserIdLoader(ctx).Load(obj.Sender.ID)
}

func (r *notificationItemResolver) Job(ctx context.Context, obj *gqlmodel.NotificationItem) (*gqlmodel.Job, error) {
	return dataloader.GetJobByJobIdLoader(ctx).Load(obj.Job.ID)
}

func (r *mutationResolver) MarkAllViewerNotificationsRead(ctx context.Context) ([]*gqlmodel.NotificationItem, error) {
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, custom_errors.ErrUserNotAuthenticated
	}

	return r.NotificationsService.MarkAllNotificationsRead(ctx, user.Id)
}

func (r *mutationResolver) MarkViewerNotificationsRead(ctx context.Context, ids []string) ([]*gqlmodel.NotificationItem, error) {
	user, err := middleware.GetCurrentUserFromContext(ctx)
	if err != nil {
		return nil, custom_errors.ErrUserNotAuthenticated
	}

	return r.NotificationsService.MarkNotificationsReadById(ctx, user.Id, ids...)
}
