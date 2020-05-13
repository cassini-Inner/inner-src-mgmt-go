package dataloader

import (
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/generated"
	"github.com/cassini-Inner/inner-src-mgmt-go/graph/model"
	"github.com/jmoiron/sqlx"
	"time"
)

func NewSkillByMilestoneIdLoader(db *sqlx.DB) *generated.SkillByMilestoneIdLoader {
	return generated.NewSkillByMilestoneIdLoader(generated.SkillByMilestoneIdLoaderConfig{
		Fetch: func(keys []string) ([][]*model.Skill, []error) {

			milestoneIdSkillListMap := make(map[string][]*model.Skill)
			var result [][]*model.Skill

			query, args, err := sqlx.In(`select milestone_id, g.id, g.created_by, g.value, g.time_created from milestones
join milestoneskills m on milestones.id = m.milestone_id
join globalskills g on m.skill_id = g.id
where milestone_id in (?)`, keys)
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

			if res.err != nil {
				return nil, []error{err}
			}

			for res.rows.Next() {
				var milestoneId, skillId, createdBy, value, timeCreated string
				err := res.rows.Scan(&milestoneId, &skillId, &createdBy, &value, &timeCreated)
				if err != nil {
					return nil, []error{err}
				}
				_, ok := milestoneIdSkillListMap[milestoneId]
				if !ok {
					milestoneIdSkillListMap[milestoneId] = make([]*model.Skill, 0)
				}

				milestoneIdSkillListMap[milestoneId] = append(milestoneIdSkillListMap[milestoneId], &model.Skill{
					ID:          skillId,
					CreatedBy:   createdBy,
					Value:       value,
					CreatedTime: timeCreated,
				})
			}

			for _, id := range keys {
				result = append(result, milestoneIdSkillListMap[id])
			}

			return result, nil
		},
		Wait:     5 * time.Millisecond,
		MaxBatch: 100,
	})
}
