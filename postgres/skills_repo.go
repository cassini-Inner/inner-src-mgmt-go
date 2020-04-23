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
	rows, err := s.db.Queryx(selectSkillsByJobIdQuery, jobId)
	if err != nil {
		return nil, err
	}
	return s.scanSkills(rows)
}

func (s *SkillsRepo) GetByUserId(userId string) ([]*dbmodel.GlobalSkill, error) {
	rows, err := s.db.Queryx(selectSkillsByUserIdQuery, userId)
	if err != nil {
		return nil, err
	}
	return s.scanSkills(rows)
}

func (s *SkillsRepo) GetByMilestoneId(milestoneId string) ([]*dbmodel.GlobalSkill, error) {
	rows, err := s.db.Queryx(selectSkillsByMilestoneIdQuery, milestoneId)
	if err != nil {
		return nil, err
	}
	return s.scanSkills(rows)
}

func (s *SkillsRepo) scanSkills(rows *sqlx.Rows) ([]*dbmodel.GlobalSkill, error) {
	var result []*dbmodel.GlobalSkill

	for rows != nil && rows.Next() {
		var skill dbmodel.GlobalSkill
		err := rows.StructScan(&skill)
		if err != nil {
			return nil, err
		}
		result = append(result, &skill)
	}
	return result, nil
}

const (
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

	selectSkillsByUserIdQuery = `select
		distinct (globalskills.id),
			globalskills.created_by,
			globalskills.value,
			globalskills.time_created
		from users 
		join userskills on userskills.user_id = users.id and userskills.is_deleted = false
		join globalskills on globalskills.id = userskills.id
		where users.id = $1 and users.is_deleted = false`
)
