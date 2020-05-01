package dataloader

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewUserByUserIdLoader(db *sqlx.DB) *generated.UserLoader {
	return generated.NewUserLoader(generated.UserLoaderConfig{
		Fetch: func(keys []string) ([]*gqlmodel.User, []error) {
			users := make(map[string]dbmodel.User)
			var result []*gqlmodel.User
			query, args, err := sqlx.In(`select * from users where users.id in (?) and users.is_deleted = false`, keys)
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

}


