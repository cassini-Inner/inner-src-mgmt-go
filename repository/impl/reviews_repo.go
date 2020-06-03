package impl

import (
	"context"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
	"time"
)

type ReviewRepoImpl struct {
	db *sqlx.DB
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

func (r ReviewRepoImpl) Add(tx *sqlx.Tx, review dbmodel.Review) (result *dbmodel.Review, err error) {
	result = &dbmodel.Review{}
	if err = tx.QueryRowx(insertReview, review.Rating, review.Remark, review.MilestoneId, review.UserId).StructScan(result); err != nil {
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

const (
	selectReviewById = "select * from reviews where id = $1 and is_deleted=false"
	insertReview     = "insert into reviews(rating, remark, milestone_id, user_id) values ($1, $2, $3, $4) returning *"
	updateReview     = "update reviews set rating=$1, remark=$2, time_updated=$3 where id=$4 and is_deleted=false returning *"
	deleteReview     = "delete from reviews where id=$1 and is_deleted=false returning *"
)
