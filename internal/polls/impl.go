package polls

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGRepo struct {
	pg *pgxpool.Pool
}

func NewPGRepo(pg *pgxpool.Pool) Repo {
	return &PGRepo{pg}
}

func (r *PGRepo) GetUncompletedPolls(pageIndex int, pageSize int, userId int64) ([]Model, error) {
	query := `
		SELECT p.id, p.text, p.type, p.author_id, p.cluster, ARRAY_REMOVE(ARRAY_AGG(ra.text), NULL) AS answer_text
		FROM polls p
		LEFT JOIN user_answers ua ON p.id = ua.poll_id AND ua.user_id = $1
		LEFT JOIN radio_answers ra ON p.id = ra.poll_id
		WHERE ua.poll_id IS NULL
		GROUP BY p.id
		LIMIT $2 OFFSET $3
	`

	polls := make([]Model, 0)

	rows, err := r.pg.Query(context.Background(), query, userId, pageSize, pageIndex*pageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var poll Model

		if err := rows.Scan(&poll.Id, &poll.Text, &poll.Type, &poll.AuthorID, &poll.Cluster, &poll.Answers); err != nil {
			return nil, err
		}

		// if poll.Type == RADIO && answerText != nil {
		// 	poll.Answers = append(poll.Answers, *answerText)
		// }

		polls = append(polls, poll)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return polls, nil
}
