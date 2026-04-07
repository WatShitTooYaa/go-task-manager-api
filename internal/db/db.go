package db

import (
	"context"

	// repo "github.com/WatShitTooYaa/go-task-manager-api/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

// type PostgresDB struct {
// 	DB *pgxpool.Pool
// }

// func NewDatabase(ctx context.Context, urlDb string) (*PostgresDB, error) {
func NewDatabase(ctx context.Context, urlDb string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(urlDb)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 20
	config.MinConns = 5

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

// func (db *PostgresDB) BeginTx(ctx context.Context) pgx.Tx {
// 	tx, err := db.DB.Begin(ctx)
// 	if err != nil {
// 		return nil
// 	}
// 	return tx
// }

// func (db *PostgresDB) Insert(ctx context.Context, task entity.Task) (entity.Task, error) {
// 	query := `
// 	INSERT INTO task (content, completed, timestamp, priority)
// 	VALUES ($1, $2, $3, $4)
// 	RETURNING id
// 	`
// 	id := 0
// 	row := db.DB.QueryRow(ctx, query, task.Content, task.Completed, task.Timestamp, task.Priority)

// 	err := row.Scan(&id)
// 	if err != nil {
// 		return task, err
// 	}
// 	// id, err := res.
// 	// if err != nil {
// 	// 	return task, err
// 	// }

// 	task.Id = uint16(id)
// 	return task, nil
// }

// func (db *PostgresDB) FindAll(ctx context.Context) ([]entity.Task, error) {
// 	query := "SELECT id, content, completed, timestamp, priority FROM task"
// 	rows, err := db.DB.Query(ctx, query)

// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	// tasks := []entity.Task{}
// 	var tasks []entity.Task

// 	for rows.Next() {
// 		task := entity.Task{}
// 		err := rows.Scan(&task.Id, &task.Content, &task.Completed, &task.Timestamp, &task.Priority)
// 		if err != nil {
// 			return nil, err
// 		}
// 		tasks = append(tasks, task)
// 	}
// 	//ada
// 	return tasks, nil
// }

// func (db *PostgresDB) FindById(ctx context.Context, id int32) (entity.Task, error) {
// 	query := "select id, content, completed, timestamp, priority from task where id = $"
// 	rows, err := db.DB.Query(ctx, query, id)
// 	task := entity.Task{}
// 	if err != nil {
// 		return task, err
// 	}
// 	//ada
// 	defer rows.Close()
// 	if rows.Next() {
// 		err := rows.Scan(&task.Id, &task.Content, &task.Completed, &task.Timestamp, &task.Priority)
// 		if err != nil {
// 			return task, err
// 		}
// 		return task, nil
// 	} else {
// 		return task, errors.New("Id " + strconv.Itoa(int(id)) + " not found!")
// 	}
// }
