package dataloader

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewUserAverageRatingLoader(db *sqlx.DB) *generated.UserAverageRatingLoader {
	return generated.NewUserAverageRatingLoader(generated.UserAverageRatingLoaderConfig{
		Fetch: func(keys []string) ([]*int, []error) {
			userRatingMap := make(map[string]int)
			result := make([]*int, 0)
			query, args, err := sqlx.In(
				`select round(avg(reviews.rating)), user_id from reviews where reviews.user_id in (?) group by reviews.user_id`,
				keys)

			if err != nil {
				return nil, []error{err}
			}

			query = db.Rebind(query)

			rows, err := db.Queryx(query, args...)
			if err != nil {
				return nil, []error{err}
			}

			for rows.Next() {
				var rating int
				var userId string
				err := rows.Scan(&rating, &userId)
				if err != nil {
					return nil, []error{err}
				}
				userRatingMap[userId] = rating
			}

			for _, id := range keys {
				val, ok := userRatingMap[id]
				if !ok {
					result = append(result, nil)
				} else {
					result = append(result, &val)

				}
			}

			return result, nil
		},
		Wait:     5 * time.Millisecond,
		MaxBatch: 100,
	})
}
