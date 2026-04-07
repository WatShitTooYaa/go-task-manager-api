package repository

import (
	"context"
	"errors"
	"strconv"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
)

type taskRepositoryImplPool struct {
	DBPool *pgxpool.Pool
}

func NewRepositoryTaskPool(dbPool *pgxpool.Pool) TaskRepository {
	return &taskRepositoryImplPool{DBPool: dbPool}
}

func (pool *taskRepositoryImplPool) Insert(ctx context.Context, task entity.Task) (entity.Task, error) {
	query := `
	INSERT INTO task (content, completed, timestamp, priority)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	id := 0
	row := pool.DBPool.QueryRow(ctx, query, task.Content, task.Completed, task.Timestamp, task.Priority)

	err := row.Scan(&id)
	if err != nil {
		return task, err
	}
	// id, err := res.
	// if err != nil {
	// 	return task, err
	// }

	task.Id = uint16(id)
	return task, nil
}

func (pool *taskRepositoryImplPool) FindAll(ctx context.Context) ([]entity.Task, error) {
	query := "SELECT id, content, completed, timestamp, priority FROM task"
	rows, err := pool.DBPool.Query(ctx, query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// tasks := []entity.Task{}
	var tasks []entity.Task

	for rows.Next() {
		task := entity.Task{}
		err := rows.Scan(&task.Id, &task.Content, &task.Completed, &task.Timestamp, &task.Priority)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	//ada
	return tasks, nil
}

func (pool *taskRepositoryImplPool) FindById(ctx context.Context, id uint16) (entity.Task, error) {
	query := "select id, content, completed, timestamp, priority from task where id = $1"
	rows, err := pool.DBPool.Query(ctx, query, id)
	task := entity.Task{}
	if err != nil {
		return task, err
	}
	//ada
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&task.Id, &task.Content, &task.Completed, &task.Timestamp, &task.Priority)
		if err != nil {
			return task, err
		}
		return task, nil
	} else {
		return task, errors.New("Id " + strconv.Itoa(int(id)) + " not found!")
	}
}

func (pool *taskRepositoryImplPool) Update(ctx context.Context, newTask entity.Task, id uint16) (entity.Task, error) {
	query := `
	update task
	set content = $1,
		completed = $2,
		timestamp = $3,
		priority = $4
	where id = $5
	returning *
	`
	rows, err := pool.DBPool.Query(ctx, query, newTask.Content, newTask.Completed, newTask.Timestamp, newTask.Priority, id)
	task := entity.Task{}
	if err != nil {
		return task, err
	}
	//ada
	defer rows.Close()
	if rows.Next() {
		err := rows.Scan(&task.Id, &task.Content, &task.Completed, &task.Timestamp, &task.Priority)
		if err != nil {
			return task, err
		}
		return task, nil
	} else {
		return task, errors.New("Id " + strconv.Itoa(int(id)) + " not found!")
	}
}

func (pool *taskRepositoryImplPool) Delete(ctx context.Context, id uint16) error {
	query := `
	DELETE FROM task
	where id = $1
	`
	_, err := pool.DBPool.Exec(ctx, query, id)
	// task := entity.Task{}
	if err != nil {
		return err
	}
	//ada
	// defer rows.Close()
	return nil
	// if rows.Next() {
	// 	err := rows.Scan(&task.Id, &task.Content, &task.Completed, &task.Timestamp, &task.Priority)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return task, nil
	// } else {
	// 	return task, errors.New("Id " + strconv.Itoa(int(id)) + " not found!")
	// }
}

// func ()
