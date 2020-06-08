alter table reviews
    add constraint review_unique_milestoneId_userId unique (milestone_id, user_id);