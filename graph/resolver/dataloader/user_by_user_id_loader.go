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
			resultChan := make(chan *FetchStruct)
			go func(result chan *FetchStruct) {
				rows, err := db.Queryx(query, args...)
				result <- &FetchStruct{
					rows: rows,
					err:  err,
				}
			}(resultChan)
			res := <-resultChan

			if res.err != nil {
				return nil, []error{err}
			}
			defer res.rows.Close()

			if err != nil {
				return nil, []error{err}
			}

			for res.rows.Next() {
				var tempUser dbmodel.User
				err := res.rows.StructScan(&tempUser)
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
		Wait:     5 * time.Millisecond,
		MaxBatch: 100,
	})

}
