package middleware

import (
	"context"
	customErrors "github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCurrentUserFromContext(t *testing.T) {
	t.Run("throw error when no user is in context", func(t *testing.T) {
		ctx := context.Background()
		user, err := GetCurrentUserFromContext(ctx)
		assert.Nil(t, user)
		assert.NotNil(t, err)
		assert.Equal(t, customErrors.ErrNoUserInContext, err)
	})

	t.Run("should return user if present", func(t *testing.T) {
		ctx := context.Background()
		mockUser := &dbmodel.User{
			Id: "1",
		}
		ctxWithUser := context.WithValue(ctx, CurrentUserKey, mockUser)
		user, err := GetCurrentUserFromContext(ctxWithUser)
		assert.Nil(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, mockUser.Id, user.Id)
	})
}
