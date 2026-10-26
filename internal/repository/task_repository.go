package repository

import (
	"context"

	"gorm.io/gorm"
	"taskflow/internal/model"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *TaskRepository) GetAll(ctx context.Context, userID uint) ([]model.Task, error) {
	var tasks []model.Task

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tasks).Error
	return tasks, err

}

func (r *TaskRepository) GetPage(ctx context.Context, userID uint, limit, offset int) ([]model.Task, error) {
	var tasks []model.Task

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&tasks).Error

	return tasks, err
}

func (r *TaskRepository) GetByID(ctx context.Context, userID, taskID uint) (*model.Task, error) {
	var task model.Task

	err := r.db.WithContext(ctx).
		Where("user_id = ? AND id = ?", userID, taskID).
		First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil

}

func (r *TaskRepository) Update(ctx context.Context, task *model.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

func (r *TaskRepository) Delete(ctx context.Context, userID, taskID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&model.Task{}, taskID).Error
}
