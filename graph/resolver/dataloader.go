package resolver

import (
	"context"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

const userLoaderKey = "userloader"

func DataloaderMiddleware(db *sqlx.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userloader := NewUserLoader(UserLoaderConfig{
			Fetch: func(keys []string) ([]*gqlmodel.User, []error) {
				 users :=  make(map[string]dbmodel.User)
				var result []*gqlmodel.User
				query, args, err := sqlx.In(`select * from users where users.id in (?)`, keys)
				if err != nil {
					return nil, []error{err}
				}

				query = db.Rebind(query)
				rows, err := db.Queryx(query, args...)
				if err != nil {
					return nil, []error{err}
				}

				for rows.Next() {
					var tempUser dbmodel.User
					err := rows.StructScan(&tempUser)
					if err != nil {
						return nil, []error{err}
					}
					users[tempUser.Id] = tempUser
				}

				for _, id := range keys {
					var tempGqlUser gqlmodel.User
					tempGqlUser.MapDbToGql(users[id])
					result = append(result, &tempGqlUser)
				}

				return result, nil
			},
			Wait:     1 * time.Millisecond,
			MaxBatch: 100,
		})

		ctx := context.WithValue(r.Context(), userLoaderKey, userloader)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserLoader(ctx context.Context) *UserLoader {
	userLoader := ctx.Value(userLoaderKey).(*UserLoader)
	return userLoader
}
