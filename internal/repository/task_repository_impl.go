package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/entity"
)

type taskRepositoryImpl struct {
	DB *sql.DB
}

func NewRepositoryTask(db *sql.DB) *taskRepositoryImpl {
	return &taskRepositoryImpl{DB: db}
}

func (repository *taskRepositoryImpl) Insert(ctx context.Context, task entity.Task) (entity.Task, error) {
	query := "INSERT INTO task (content, completed, priority, timestamp) VALUES ($, $, $, $)"
	res, err := repository.DB.ExecContext(ctx, query, task.Content, task.Completed, task.Timestamp, task.Priority)
	if err != nil {
		return task, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return task, err
	}

	task.Id = uint16(id)
	return task, nil
}

func (repository *taskRepositoryImpl) FindAll(ctx context.Context) ([]entity.Task, error) {
	query := "select id, content, completed, timestamp, priority from"
	rows, err := repository.DB.QueryContext(ctx, query)

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

func (repository *taskRepositoryImpl) FindById(ctx context.Context, id uint16) (entity.Task, error) {
	query := "select id, content, completed, timestamp, priority from task where id = $"
	rows, err := repository.DB.QueryContext(ctx, query, id)
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

// update and delete not implemented yet
func (repository *taskRepositoryImpl) Update(ctx context.Context, newTask entity.Task, id uint16) (entity.Task, error) {
	query := "select id, content, completed, timestamp, priority from task where id = $"
	rows, err := repository.DB.QueryContext(ctx, query, id)
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
func (repository *taskRepositoryImpl) Delete(ctx context.Context, id uint16) error {
	query := "select id, content, completed, timestamp, priority from task where id = $"
	rows, err := repository.DB.QueryContext(ctx, query, id)
	// task := entity.Task{}
	if err != nil {
		return err
	}
	//ada
	defer rows.Close()
	return nil
	// if rows.Next() {
	// 	err := rows.Scan(&task.Id, &task.Content, &task.Completed, &task.Timestamp, &task.Priority)
	// 	if err != nil {
	// 		return task, err
	// 	}
	// 	return task, nil
	// } else {
	// 	return task, errors.New("Id " + strconv.Itoa(int(id)) + " not found!")
	// }
}
