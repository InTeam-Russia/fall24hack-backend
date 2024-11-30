package polls

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PGRepo struct {
	pg     *pgxpool.Pool
	logger *zap.Logger
}

func NewPGRepo(pg *pgxpool.Pool, logger *zap.Logger) Repo {
	return &PGRepo{pg, logger}
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

		polls = append(polls, poll)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return polls, nil
}

func (r *PGRepo) AddAnswer(userId int64, pollId int64, text string) error {
	query := `
		INSERT INTO user_answers (user_id, poll_id, text)
		VALUES ($1, $2, $3)
	`

	_, err := r.pg.Exec(context.Background(), query, userId, pollId, text)
	if err != nil {
		return err
	}

	return nil
}

func (r *PGRepo) CreatePoll(poll *CreateModel, authorId int64, cluster int) error {
	// TODO: Use transactions

	insertPollsQuery := `
		INSERT INTO polls (text, type, author_id, cluster)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	row := r.pg.QueryRow(
		context.Background(),
		insertPollsQuery,
		poll.Text,
		string(poll.Type),
		authorId,
		cluster,
	)

	var pollId int64
	err := row.Scan(&pollId)
	if err != nil {
		fmt.Println("ABOBA")
		return err
	}

	d := 1

	if poll.Type == RADIO {
		strBuilder := strings.Builder{}
		strBuilder.WriteString("INSERT INTO radio_answers (answer_id, poll_id, text) VALUES\n")

		args := make([]any, 0)

		for ansIdx, ans := range poll.Answers {
			strBuilder.WriteString(fmt.Sprintf("($%d, $%d, $%d)", d, d+1, d+2))
			if ansIdx != len(poll.Answers)-1 {
				strBuilder.WriteString(",\n")
			}
			d += 3

			args = append(args, ansIdx, pollId, ans)
		}

		insertRadioAnswersQuery := strBuilder.String()
		r.logger.Debug(insertRadioAnswersQuery)

		_, err = r.pg.Exec(context.Background(), insertRadioAnswersQuery, args...)
		if err != nil {
			return err
		}
	}

	return nil
}
