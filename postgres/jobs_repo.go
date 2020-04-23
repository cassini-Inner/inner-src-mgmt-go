package postgres

import (
	"context"
	"fmt"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

type JobsRepo struct {
	db *sqlx.DB
}

func NewJobsRepo(db *sqlx.DB) *JobsRepo {
	return &JobsRepo{db: db}
}

func (j *JobsRepo) CreateJob(ctx context.Context, input *gqlmodel.CreateJobInput, user *dbmodel.User) (*dbmodel.Job, error) {
	// create a new transaction for creating a job
	tx, err := j.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// insert the information into the job table
	var insertedJobId string
	err = tx.QueryRowContext(ctx, createJobQuery, input.Title, input.Desc, input.Difficulty, user.Id).Scan(&insertedJobId)

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	var insertedMilestoneIds []int
	stmt, valueArgs := getInsertMilestonesStatement(input, insertedJobId, user.Id)
	stmt = j.db.Rebind(stmt)
	// get the ids of newly inserted milestones
	milestonesInsertResult, err := tx.QueryContext(ctx, stmt, valueArgs...)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	for milestonesInsertResult.Next() {
		id := 0
		err := milestonesInsertResult.Scan(&id)
		if err != nil {
			_ = tx.Rollback()
			return nil, err
		}
		insertedMilestoneIds = append(insertedMilestoneIds, id)
	}

	milestonesInsertResult.Close()
	if err != nil {
		return nil, err
	}
	// build a map of unique skills from all the milestones
	skillsMap := make(map[string]int)
	var nonExistingSkills bool
	for _, milestone := range input.Milestones {
		for _, skill := range milestone.Skills {
			skillsMap[strings.ToLower(*skill)] = 0
		}
	}

	// convert the map to string to use that as uniqueSkillsQueryArgs
	var skillsList []string
	for key := range skillsMap {
		skillsList = append(skillsList, key)
	}
	uniqueSkillsQuery, uniqueSkillsQueryArgs, err := sqlx.In(`select value, id from globalskills where value in (?)`, skillsList)
	if err != nil {
		_ = tx.Rollback()
		log.Println("error while creating query for skills in db")
		return nil, err
	}
	uniqueSkillsQuery = j.db.Rebind(uniqueSkillsQuery)
	uniqueSkillsRows, err := tx.QueryContext(ctx, uniqueSkillsQuery, uniqueSkillsQueryArgs...)
	if err != nil {
		log.Println("error while querying skill values")
		_ = tx.Rollback()
		return nil, err
	}

	// assign the database id to each skill in map
	// if a skill is not present in the database then id = 0
	for uniqueSkillsRows.Next() {
		var skillValue string
		var skillId int
		err := uniqueSkillsRows.Scan(&skillValue, &skillId)
		if err != nil {
			return nil, err
		}
		skillsMap[strings.ToLower(skillValue)] = skillId
	}
	uniqueSkillsRows.Close()

	// check if there are any skills in skillMap that are not present in db
	for k := range skillsMap {
		if skillsMap[k] == 0 {
			nonExistingSkills = true
		}
	}

	// if there are skills that are not already present in the database
	// then create them
	if nonExistingSkills {
		var newSkillsQueryArgs []interface{}
		var newSkillsQueryValues []string
		for k := range skillsMap {
			if skillsMap[k] == 0 {
				newSkillsQueryValues = append(newSkillsQueryValues, "(?,?)")
				newSkillsQueryArgs = append(newSkillsQueryArgs, strings.ToLower(k), user.Id)
			}
		}

		newSkillsQueryStatement := fmt.Sprintf(`insert into globalskills(value, created_by) values %s returning id`, strings.Join(newSkillsQueryValues, ","))
		newSkillsQueryStatement = j.db.Rebind(newSkillsQueryStatement)
		uniqueSkillsRows, err := tx.QueryContext(ctx, newSkillsQueryStatement, newSkillsQueryArgs...)
		if err != nil {
			_ = tx.Rollback()
			log.Println(err)
			log.Println("could not create new skills")
			return nil, err
		}
		// build a list of all skill ids
		var insertedSkillsIds []int
		for uniqueSkillsRows.Next() {
			var id int
			err := uniqueSkillsRows.Scan(&id)
			if err != nil {
				_ = tx.Rollback()
				return nil, err
			}
			insertedSkillsIds = append(insertedSkillsIds, id)
		}
		uniqueSkillsRows.Close()

		// populate the skillMap with new skill values and ids
		skillQuery, skillsQueryArgs, err := sqlx.In(`select id, value from globalskills where id in (?)`, insertedSkillsIds)
		if err != nil {
			log.Println(err)
			_ = tx.Rollback()
			return nil, err
		}
		skillQuery = j.db.Rebind(skillQuery)
		newSkillRows, err := tx.QueryContext(ctx, skillQuery, skillsQueryArgs...)
		if err != nil {
			_ = tx.Rollback()
			log.Println(err)
			return nil, err
		}

		for newSkillRows.Next() {
			var id int
			var value string
			err = newSkillRows.Scan(&id, &value)
			if err != nil {
				return nil, err
			}
			skillsMap[value] = id
		}
		newSkillRows.Close()
	}

	// populate the milestoneskills table with new milestone_id's and skill_id's
	var milestoneSkillsArgs []interface{}
	var milestoneSkillsQuery []string
	for i, milestone := range input.Milestones {
		milestoneID := insertedMilestoneIds[i]
		for _, skill := range milestone.Skills {
			milestoneSkillsQuery = append(milestoneSkillsQuery, "(?, ?)")
			milestoneSkillsArgs = append(milestoneSkillsArgs, skillsMap[strings.ToLower(*skill)], milestoneID)
		}
	}
	newMilestoneSkillsQuery := j.db.Rebind(fmt.Sprintf(`insert into milestoneskills(skill_id, milestone_id) values %v returning id`, strings.Join(milestoneSkillsQuery, ",")))

	_, err = tx.ExecContext(ctx, newMilestoneSkillsQuery, milestoneSkillsArgs...)
	if err != nil {
		_ = tx.Rollback()
		log.Println("error while creating non existing milestoneskills")
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return j.GetById(insertedJobId)
}

func (j *JobsRepo) UpdateJob(input *gqlmodel.UpdateJobInput) (*dbmodel.Job, error) {
	panic("Not implemented")
}

func (j *JobsRepo) DeleteJob(jobId string) (*dbmodel.Job, error) {
	panic("Not implemented")
}

// Get the complete job details based on the job id
func (j *JobsRepo) GetById(jobId string) (*dbmodel.Job, error) {
	var job dbmodel.Job
	err := j.db.QueryRowx(selectJobByIdQuery, jobId).StructScan(&job)
	if err != nil {
		return nil, err
	}
	return &job, nil
}

// GetByUserId returns all jobs created by that user
func (j *JobsRepo) GetByUserId(userId string) ([]*dbmodel.Job, error) {

	rows, err := j.db.Queryx(selectJobsByUserIdQuery, userId)
	if err != nil {
		return nil, err
	}

	var jobs []*dbmodel.Job
	for rows.Next() {
		var job dbmodel.Job
		rows.StructScan(&job)
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

//TODO: Refactor this.
func (j *JobsRepo) GetStatsByUserId(userId string) (*gqlmodel.UserStats, error) {
	panic("not implemented")
}

//TODO: Add sorting order functionality
func (j *JobsRepo) GetAll(filters *gqlmodel.JobsFilterInput) ([]*dbmodel.Job, error) {
	var jobSkills []string
	for _, skill := range filters.Skills {
		jobSkills = append(jobSkills, strings.ToLower(*skill))
	}

	var jobStatuses []string
	for _, status := range filters.Status {
		jobStatuses = append(jobStatuses, strings.ToLower(status.String()))
	}

	query, args, err := sqlx.In(selectAllJobsWithFiltersQuery, jobSkills, jobStatuses)
	if err != nil {
		return nil, err
	}
	query = j.db.Rebind(query)

	rows, err := j.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	var result []*dbmodel.Job
	for rows != nil && rows.Next() {
		var tempJob dbmodel.Job
		rows.StructScan(&tempJob)
		result = append(result, &tempJob)
	}
	return result, nil
}

const (
	selectJobByIdQuery            = `SELECT * FROM jobs WHERE id = $1 and is_deleted=false`
	selectJobsByUserIdQuery       = `SELECT * FROM jobs WHERE created_by = $1 and is_deleted=false`
	selectAllJobsWithFiltersQuery = `select * from jobs where jobs.id in (
		select milestones.job_id from globalskills 
		join milestoneskills on globalskills.id = milestoneskills.skill_id
		join milestones on milestones.id = milestoneskills.milestone_id
		where globalskills.value in (?)
		)
		and jobs.status in (?)`
	createJobQuery = `insert into jobs(title, description, difficulty,created_by) values ($1, $2, $3, $4) returning jobs.id`
)

func getInsertMilestonesStatement(input *gqlmodel.CreateJobInput, insertedJobId, userId string) (string, []interface{}) {
	var valueStrings []string
	var valueArgs []interface{}
	for _, milestone := range input.Milestones {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, milestone.Title, milestone.Desc, milestone.Duration, milestone.Status.String(), milestone.Resolution, insertedJobId)
	}
	stmt := fmt.Sprintf("INSERT INTO milestones (title, description, duration, status, resolution, job_id) VALUES %s returning id",
		strings.Join(valueStrings, ", "))
	return stmt, valueArgs
}
