package dataloader

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	gqlmodel "github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewSkillByJobIdLoader(db *sqlx.DB) *generated.SkillByJobIdLoader {
	return generated.NewSkillByJobIdLoader(generated.SkillByJobIdLoaderConfig{
		Fetch: func(keys []string) ([][]*gqlmodel.Skill, []error) {

			jobIdSkillListMap := make(map[string]map[string]*gqlmodel.Skill)
			var result [][]*gqlmodel.Skill

			query, args, err := sqlx.In(`select job_id, g.id, g.created_by, g.value, g.time_created from milestones
			join milestoneskills m on milestones.id = m.milestone_id
			join globalskills g on m.skill_id = g.id
			where job_id in (?)`, keys)

			if err != nil {
				return nil, []error{err}
			}

			query = db.Rebind(query)
			rows, err := db.Queryx(query, args...)
			if err != nil {
				return nil, []error{err}
			}

			for rows.Next() {
				var jobId, skillId, createdBy, value, timeCreated string
				err := rows.Scan(&jobId, &skillId, &createdBy, &value, &timeCreated)
				if err != nil {
					return nil, []error{err}
				}
				_, ok := jobIdSkillListMap[jobId]
				if !ok {
					jobIdSkillListMap[jobId] = make(map[string]*gqlmodel.Skill)
				}

				jobIdSkillListMap[jobId][skillId] = &gqlmodel.Skill{
					ID:          skillId,
					CreatedBy:   createdBy,
					Value:       value,
					CreatedTime: timeCreated,
				}
			}

			for _, id := range keys {
				var temp []*gqlmodel.Skill
				for key := range jobIdSkillListMap[id] {
					temp = append(temp, jobIdSkillListMap[id][key])
				}
				result = append(result, temp)
			}

			return result, nil
		},
		Wait:     1 * time.Millisecond,
		MaxBatch: 100,
	})
}
