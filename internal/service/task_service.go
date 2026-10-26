package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"taskflow/internal/cache"
	"taskflow/internal/model"
	"taskflow/internal/repository"
)

type TaskService struct {
	repo  *repository.TaskRepository
	cache *cache.RedisCache
}

func NewTaskService(repo *repository.TaskRepository, redisCache *cache.RedisCache) *TaskService {
	return &TaskService{repo: repo, cache: redisCache}
}

func (s *TaskService) Create(ctx context.Context, userID uint, title, description, priority string) (*model.Task, error) {
	if title == "" {
		return nil, errors.New("название задачи обязательно")
	}

	if priority == "" {
		priority = "medium"
	}
	if priority != "low" && priority != "medium" && priority != "high" {
		return nil, errors.New("недопустимый приоритет задачи")
	}

	task := &model.Task{
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      "pending",
		Priority:    priority,
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	s.invalidateUserCache(ctx, userID)
	return task, nil
}

func (s *TaskService) GetAll(ctx context.Context, userID uint) ([]model.Task, error) {
	return s.GetPage(ctx, userID, 1, 20)
}

func (s *TaskService) GetPage(ctx context.Context, userID uint, page, pageSize int) ([]model.Task, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	key := fmt.Sprintf("tasks:user:%d:page:%d:size:%d", userID, page, pageSize)

	if s.cache != nil {
		value, err := s.cache.Get(ctx, key)
		if err == nil {
			var tasks []model.Task
			if json.Unmarshal([]byte(value), &tasks) == nil {
				return tasks, nil
			}
		}
	}

	tasks, err := s.repo.GetPage(ctx, userID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		data, err := json.Marshal(tasks)
		if err == nil {
			_ = s.cache.Set(ctx, key, string(data), 30*time.Second)
		}
	}

	return tasks, nil
}

func (s *TaskService) GetByID(ctx context.Context, userID, taskID uint) (*model.Task, error) {
	return s.repo.GetByID(ctx, userID, taskID)
}

func (s *TaskService) Update(ctx context.Context, userID, taskID uint, title, description, status, priority string) (*model.Task, error) {
	task, err := s.repo.GetByID(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}

	if title != "" {
		task.Title = title
	}
	task.Description = description

	if status != "" {
		if status != "pending" && status != "in_progress" && status != "completed" {
			return nil, errors.New("недопустимый статус задачи")
		}
		task.Status = status
	}

	if priority != "" {
		if priority != "low" && priority != "medium" && priority != "high" {
			return nil, errors.New("недопустимый приоритет задачи")
		}
		task.Priority = priority
	}

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	s.invalidateUserCache(ctx, userID)
	return task, nil
}

func (s *TaskService) Delete(ctx context.Context, userID, taskID uint) error {
	if err := s.repo.Delete(ctx, userID, taskID); err != nil {
		return err
	}

	s.invalidateUserCache(ctx, userID)
	return nil
}

func (s *TaskService) invalidateUserCache(ctx context.Context, userID uint) {
	if s.cache == nil {
		return
	}

	keys := make([]string, 0, 20)
	for page := 1; page <= 10; page++ {
		for _, size := range []int{10, 20} {
			keys = append(keys, fmt.Sprintf("tasks:user:%d:page:%d:size:%d", userID, page, size))
		}
	}

	_ = s.cache.Delete(ctx, keys...)
}
