package postgres

import (
	dbmodel "github.com/cassini-Inner/inner-src-mgmt-go/postgres/model"
	"github.com/jmoiron/sqlx"
)

type SkillsRepo struct {
	db *sqlx.DB
}

func NewSkillsRepo(db *sqlx.DB) *SkillsRepo {
	return &SkillsRepo{db: db}
}

func (s *SkillsRepo) GetByJobId(jobId string) ([]*dbmodel.GlobalSkill, error) {
	rows, err := s.db.Queryx(getSkillsByJobIdQuery, jobId)
	if err != nil {
		return nil, err
	}
	return s.scanSkills(rows)
}

func (s *SkillsRepo) GetByUserId(userId string) ([]*dbmodel.GlobalSkill, error) {
	rows, err := s.db.Queryx(getSkillsByUserIdQuery, userId)
	if err != nil {
		return nil, err
	}
	return s.scanSkills(rows)
}

func (s *SkillsRepo) GetByMilestoneId(milestoneId string) ([]*dbmodel.GlobalSkill, error) {
	rows, err := s.db.Queryx(getSkillsByMilestoneIdQuery, milestoneId)
	if err != nil {
		return nil, err
	}
	return s.scanSkills(rows)
}

func (s *SkillsRepo) scanSkills(rows *sqlx.Rows) ([]*dbmodel.GlobalSkill, error) {
	var result []*dbmodel.GlobalSkill
	var skill dbmodel.GlobalSkill

	for rows != nil && rows.Next() {
		rows.StructScan(&skill)
		result = append(result, &skill)
	}
	return result, nil
}

const (
	getSkillsByJobIdQuery = `select distinct (globalskills.id), globalskills.created_by,
		globalskills.value,
		globalskills.time_created
		from jobs
		join milestones on jobs.id = milestones.job_id	 
		join milestoneskills on milestoneskills.milestone_id = milestones.id
		join globalskills on milestoneskills.skill_id = globalskills.id
		where jobs.id = $1
		order by globalskills.value`

	getSkillsByMilestoneIdQuery = `select
		distinct (globalskills.id),
			globalskills.created_by,
			globalskills.value,
			globalskills.time_created
		from milestones
		join milestoneskills on milestoneskills.milestone_id = milestones.id
		join globalskills on milestoneskills.skill_id = globalskills.id
		where milestones.id = $1
		order by globalskills.value`

	getSkillsByUserIdQuery = `select
		distinct (globalskills.id),
			globalskills.created_by,
			globalskills.value,
			globalskills.time_created
		from users
		join userskills on userskills.user_id = users.id
		join globalskills on globalskills.id = userskills.id
		where users.id = $1`
)
