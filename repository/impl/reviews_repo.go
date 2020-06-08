package impl

import (
	"context"
	"database/sql"
	"github.com/cassini-Inner/inner-src-mgmt-go/custom_errors"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
)

type ReviewRepoImpl struct {
	db *sqlx.DB
}

func NewReviewRepoImpl(db *sqlx.DB) *ReviewRepoImpl {
	return &ReviewRepoImpl{db: db}
}

func (r ReviewRepoImpl) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}

func (r ReviewRepoImpl) CommitTx(ctx context.Context, tx *sqlx.Tx) (err error) {
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
	}
	return nil
}

func (r ReviewRepoImpl) GetById(id string) (review *dbmodel.Review, err error) {
	review = &dbmodel.Review{}
	err = r.db.QueryRowx(selectReviewById, id).StructScan(review)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (r ReviewRepoImpl) GetByIdTx(tx *sqlx.Tx, id string) (review *dbmodel.Review, err error) {
	review = &dbmodel.Review{}
	err = tx.QueryRowx(selectReviewById, id).StructScan(review)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (r ReviewRepoImpl) GetForUserIdMilestoneId(milestoneId, userId string) (*dbmodel.Review, error) {
	review := &dbmodel.Review{}
	err := r.db.QueryRowx(selectReviewByMilestoneIdUserId, milestoneId, userId).StructScan(review)
	if err != nil {
		return nil, err
	}
	return review, nil
}

func (r ReviewRepoImpl) Add(tx *sqlx.Tx, review dbmodel.Review) (result *dbmodel.Review, err error) {
	result = &dbmodel.Review{}
	if err = tx.QueryRowx(insertReview, review.Rating, review.Remark, review.MilestoneId, review.UserId).StructScan(result); err != nil {
		pqErr := err.(*pq.Error)
		// 23505 is the error code for unique key constraint violations
		if pqErr.Code == "23505" {
			return nil, custom_errors.ErrAlreadyExists
		}
		return nil, err
	}
	return result, nil
}

func (r ReviewRepoImpl) Update(tx *sqlx.Tx, id string, review dbmodel.Review) (result *dbmodel.Review, err error) {
	result = &dbmodel.Review{}
	if err = tx.QueryRowx(updateReview, review.Rating, review.Remark, time.Now().UTC(), id).StructScan(result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r ReviewRepoImpl) Delete(tx *sqlx.Tx, id string) (review *dbmodel.Review, err error) {
	review = &dbmodel.Review{}
	if err = tx.QueryRowx(deleteReview, id).StructScan(review); err != nil {
		return nil, err
	}
	return review, nil
}

func (r ReviewRepoImpl) GetForUserId(userId string) ([]*dbmodel.Review, error) {
	var result []*dbmodel.Review
	rows, err := r.db.Queryx(selectReviewByUserId, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, custom_errors.ErrNoEntityMatchingId
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		tempReview := &dbmodel.Review{}
		err := rows.StructScan(tempReview)
		if err != nil {
			return nil, err
		}
		result = append(result, tempReview)
	}
	return result, nil
}

const (
	selectReviewById                = "select * from reviews where id = $1 and is_deleted=false"
	selectReviewByMilestoneIdUserId = "select * from reviews where milestone_id = $1 and user_id = $2"
	selectReviewByUserId            = "select * from reviews where user_id = $1"
	insertReview                    = "insert into reviews(rating, remark, milestone_id, user_id) values ($1, $2, $3, $4) returning *"
	updateReview                    = "update reviews set rating=$1, remark=$2, time_updated=$3 where id=$4 and is_deleted=false returning *"
	deleteReview                    = "update reviews set is_deleted=false where id=$1 returning *"
)
