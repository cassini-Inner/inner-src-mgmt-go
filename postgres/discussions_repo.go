package postgres

import (
	"fmt"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
)

type DiscussionsRepo struct {
	db *sqlx.DB
}

func NewDiscussionsRepo(db *sqlx.DB) *DiscussionsRepo {
	return &DiscussionsRepo{db: db}
}

//TODO: Implement
func (d *DiscussionsRepo) CreateComment(jobId, comment, userId string) (*dbmodel.Discussion, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	insertedCommentId := 0
	err = tx.QueryRow(`insert into discussions(job_id, created_by, content) values ($1,$2, $3) returning id`, jobId, userId, comment).Scan(&insertedCommentId)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	var id, job, createdBy, content, timeCreated, timeUpdated string
	err = tx.QueryRow(`select id, job_id, created_by,content, time_created, time_updated from discussions where id = $1`, insertedCommentId).Scan(&id, &job, &createdBy, &content, &timeCreated, &timeUpdated)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &dbmodel.Discussion{
		Id:          id,
		JobId:       job,
		CreatedBy:   createdBy,
		Content:     content,
		TimeCreated: timeCreated,
		TimeUpdated: timeUpdated,
	}, nil
}

func (d *DiscussionsRepo) UpdateComment(discussionId, content string) (*dbmodel.Discussion, error) {
	var id, jobId, createdBy, discContent, timeCreated, timeUpdated string
	var isDeleted bool
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(updateDiscussionById, content, discussionId).Scan(&id, &jobId, &createdBy, &discContent, &timeCreated, &timeUpdated, &isDeleted)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &dbmodel.Discussion{
		Id:          id,
		JobId:       jobId,
		CreatedBy:   createdBy,
		Content:     discContent,
		TimeCreated: timeCreated,
		TimeUpdated: timeUpdated,
		IsDeleted:   false,
	}, nil
}

func (d *DiscussionsRepo) DeleteComment(discussionId string) error {
	tx, err:= d.db.Begin()
	if err != nil {
		return err
	}
	result, err := tx.Exec(deleteDiscussionById, discussionId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	fmt.Printf("%+v", result)
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (d *DiscussionsRepo) GetByJobId(jobId string) ([]*dbmodel.Discussion, error) {
	rows, err := d.db.Queryx(getDiscussionsByJobId, jobId)
	if err != nil {
		return nil, err
	}

	var result []*dbmodel.Discussion
	for rows != nil && rows.Next() {
		var discussion dbmodel.Discussion
		rows.StructScan(&discussion)
		result = append(result, &discussion)
	}

	return result, nil
}

func (d *DiscussionsRepo) GetById(discussionId string) (*dbmodel.Discussion, error) {
	var discussion dbmodel.Discussion
	err := d.db.QueryRowx(getDiscussionById, discussionId).StructScan(&discussion)
	if err != nil {
		return nil, err
	}

	return &discussion, nil
}

const (
	getDiscussionsByJobId = `select * from discussions where job_id = $1 and is_deleted=false order by time_created`
	getDiscussionById     = `select * from discussions where id = $1 and is_deleted = false`
	updateDiscussionById  = `update discussions set content = $1 where id = $2 and is_deleted = false returning *`
	deleteDiscussionById = `update discussions set is_deleted = true where id = $1`
)
