package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/repository/model"
	"github.com/jmoiron/sqlx"
	"log"
	"strings"
)

type SkillsRepoImpl struct {
	db *sqlx.DB
}

func NewSkillsRepo(db *sqlx.DB) *SkillsRepoImpl {
	return &SkillsRepoImpl{db: db}
}

var (
	ErrInvalidListLength = errors.New("input list must have atleast one item")
)

func (s *SkillsRepoImpl) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return s.db.BeginTxx(ctx, nil)
}

func (s *SkillsRepoImpl) CommitTx(ctx context.Context, tx *sqlx.Tx) (err error) {
	err = tx.Commit()
	if err != nil {
		err = tx.Rollback()
	}
	return err
}

func (s *SkillsRepoImpl) GetMatchingSkills(query string, limit *int) ([]*dbmodel.GlobalSkill, error) {
	rows, err := s.db.Queryx(skillsMatchingQuery, (query)+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return s.scanSkills(rows)
}

func (s *SkillsRepoImpl) GetByJobId(jobId string) ([]*dbmodel.GlobalSkill, error) {
	rows, err := s.db.Queryx(selectSkillsByJobIdQuery, jobId)
	if err != nil {
		return nil, err
	}
	return s.scanSkills(rows)
}

// findOrCreateSkills checks the input skills against those present in the database
// and maps existing skills to those in db and creates the one that do not exist
// then it returns a map of skills by the given skill value
// TODO: Refactor, duplicate method sig
func findOrCreateSkills(skills []string, userId string, tx *sqlx.Tx) (map[string]*dbmodel.GlobalSkill, error) {

	if len(skills) == 0 {
		return nil, ErrInvalidListLength
	}

	skillsMap := make(map[string]string)
	resultMap := make(map[string]*dbmodel.GlobalSkill)
	// put all skills in a map so that we only have unique ones
	for i, skill := range skills {
		skills[i] = strings.ToLower(skill)
		skillsMap[skill] = ""
	}

	// convert the map to string to use that as uniqueSkillsQueryArgs
	var skillsList []string
	for key := range skillsMap {
		skillsList = append(skillsList, key)
	}

	uniqueSkillsQuery, uniqueSkillsQueryArgs, err := sqlx.In(`select * from globalskills where value in (?)`, skillsList)
	if err != nil {
		log.Println("error while creating query for skills in db")
		return nil, err
	}
	uniqueSkillsQuery = tx.Rebind(uniqueSkillsQuery)

	// query the database to see if there are any skills that already exist
	uniqueSkillsRows, err := tx.Queryx(uniqueSkillsQuery, uniqueSkillsQueryArgs...)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Println("error while querying skill values")
			return nil, err
		}
	}

	// assign the database id to each skill in map
	// if a skill is not present in the database then id = 0
	for uniqueSkillsRows != nil && uniqueSkillsRows.Next() {
		tempSkill := &dbmodel.GlobalSkill{}
		err := uniqueSkillsRows.StructScan(tempSkill)
		if err != nil {
			return nil, err
		}
		skillsMap[strings.ToLower(tempSkill.Value)] = tempSkill.Id
		resultMap[strings.ToLower(tempSkill.Value)] = tempSkill
	}
	uniqueSkillsRows.Close()

	// check if there are any skills in skillMap that are not present in db
	// if there are skills that do not exist, set nonExistingSkills flag to true

	nonExistingSkills := false
	for k := range skillsMap {
		if skillsMap[k] == "" {
			nonExistingSkills = true
		}
	}

	// if there are skills that are not already present in the database
	// then create them
	if nonExistingSkills {
		var newSkillsQueryArgs []interface{}
		var newSkillsQueryValues []string

		// if a skill is not there in the database, ie. has id == 0 from our last fetch operation, then
		// add that to args
		for k := range skillsMap {
			if skillsMap[k] == "" {
				// this is done to prepare a query that inserts all the skill in a single statement
				newSkillsQueryValues = append(newSkillsQueryValues, "(?,?)")
				newSkillsQueryArgs = append(newSkillsQueryArgs, strings.ToLower(k), userId)
			}
		}

		newSkillsQueryStatement := fmt.Sprintf(`insert into globalskills(value, created_by) values %s returning *`, strings.Join(newSkillsQueryValues, ","))
		newSkillsQueryStatement = tx.Rebind(newSkillsQueryStatement)
		newSkillRows, err := tx.Queryx(newSkillsQueryStatement, newSkillsQueryArgs...)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// build a list of all skill ids
		for newSkillRows.Next() {
			insertedSkill := &dbmodel.GlobalSkill{}
			err := newSkillRows.StructScan(insertedSkill)
			if err != nil {
				return nil, err
			}
			resultMap[strings.ToLower(insertedSkill.Value)] = insertedSkill
		}
		newSkillRows.Close()
	}

	return resultMap, nil
}

func (s *SkillsRepoImpl) GetByUserId(userId string) ([]*dbmodel.GlobalSkill, error) {
	rows, err := s.db.Queryx(selectSkillsByUserIdQuery, userId)
	if err != nil {
		return nil, err
	}
	return s.scanSkills(rows)
}

func (s *SkillsRepoImpl) GetByMilestoneId(milestoneId string) ([]*dbmodel.GlobalSkill, error) {
	rows, err := s.db.Queryx(selectSkillsByMilestoneIdQuery, milestoneId)
	if err != nil {
		return nil, err
	}
	return s.scanSkills(rows)
}

func (s *SkillsRepoImpl) scanSkills(rows *sqlx.Rows) ([]*dbmodel.GlobalSkill, error) {
	var result []*dbmodel.GlobalSkill

	for rows != nil && rows.Next() {
		skill := &dbmodel.GlobalSkill{}
		err := rows.StructScan(skill)
		if err != nil {
			return nil, err
		}
		result = append(result, skill)
	}
	return result, nil
}

func (s *SkillsRepoImpl) GetAll() ([]*dbmodel.GlobalSkill, error) {
	var result []*dbmodel.GlobalSkill

	rows, err := s.db.Queryx(selectAllSkills)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		 skill := &dbmodel.GlobalSkill{}
		err := rows.StructScan(skill)
		if err != nil {
			return nil, err
		}
		result = append(result, skill)
	}

	return result, nil
}

func (s *SkillsRepoImpl) FindOrCreateSkills(ctx context.Context, tx *sqlx.Tx, skillsList []string, userId string) (skillsMap map[string]*dbmodel.GlobalSkill, err error) {
	// build a map of unique skills from all the milestones
	skillsMap, err = findOrCreateSkills(skillsList, userId, tx)
	if err != nil {
		return nil, err
	}

	return skillsMap, nil
}

func (s *SkillsRepoImpl) MapSkillsToMilestones(ctx context.Context, tx *sqlx.Tx, skillsMap map[string]*dbmodel.GlobalSkill, input *gqlmodel.CreateJobInput, insertedMilestones []*dbmodel.Milestone) (err error) {
	var milestoneSkillsArgs []interface{}
	var milestoneSkillsQuery []string
	for i, milestone := range input.Milestones {
		milestoneID := insertedMilestones[i].Id
		for _, skill := range milestone.Skills {
			milestoneSkillsQuery = append(milestoneSkillsQuery, "(?, ?)")
			milestoneSkillsArgs = append(milestoneSkillsArgs, skillsMap[strings.ToLower(*skill)].Id, milestoneID)
		}
	}
	newMilestoneSkillsQuery := tx.Rebind(fmt.Sprintf(`insert into milestoneskills(skill_id, milestone_id) values %v returning id`, strings.Join(milestoneSkillsQuery, ",")))

	_, err = tx.ExecContext(ctx, newMilestoneSkillsQuery, milestoneSkillsArgs...)
	if err != nil {
		log.Println("error while creating non existing milestoneskills")
		return err
	}
	return nil
}

func (s *SkillsRepoImpl) AddSkillsToUserSkills(tx *sqlx.Tx, skills map[string]*dbmodel.GlobalSkill, userId string) error {
	if len(skills) == 0 {
		return nil
	}
	var newUserskillsValue []string
	var newUserskillsArgs []interface{}

	for key := range skills {
		newUserskillsValue = append(newUserskillsValue, "(?, ?)")
		newUserskillsArgs = append(newUserskillsArgs, userId, skills[key].Id)
	}

	stmt := tx.Rebind(fmt.Sprintf(insertIntoUserskillsquery, strings.Join(newUserskillsValue, ",")))

	rows, err := tx.Queryx(stmt, newUserskillsArgs...)
	defer rows.Close()

	if err != nil {
		return err
	}

	return nil
}

const (
	selectAllSkills          = `select * from globalskills`
	selectSkillsByJobIdQuery = `select distinct (globalskills.id), globalskills.created_by,
		globalskills.value,
		globalskills.time_created
		from jobs
		join milestones on jobs.id = milestones.job_id and milestones.is_deleted = false
		join milestoneskills on milestoneskills.milestone_id = milestones.id and milestoneskills.is_deleted = false
		join globalskills on milestoneskills.skill_id = globalskills.id
		where jobs.id = $1 and jobs.is_deleted = false
		order by globalskills.value`

	selectSkillsByMilestoneIdQuery = `select
		distinct (globalskills.id),
			globalskills.created_by,
			globalskills.value,
			globalskills.time_created
		from milestoneskills
		join globalskills on milestoneskills.skill_id = globalskills.id and milestoneskills.is_deleted = false
		where milestoneskills.milestone_id = $1
		order by globalskills.value`

	selectSkillsByUserIdQuery = `select distinct(g.id), g.created_by, g.value, g.time_created from users join userskills u on users.id = u.user_id and users.id = $1 join globalskills g on u.skill_id = g.id and u.is_deleted = false`

	deleteSkillsFromUserskillsByUserIdQuery = `update userskills set is_deleted = true where user_id = $1`

	insertIntoUserskillsquery = `insert into userskills(user_id, skill_id) values %v returning id`

	skillsMatchingQuery = `select * from globalskills where value like $1 limit $2`
)
