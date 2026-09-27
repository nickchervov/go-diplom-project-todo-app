package service

import (
	"context"

	"github.com/nickchervov/go-diplom-project/internal/domain"
)

type AdapterInterface interface {
	AddTask(ctx context.Context, task domain.Task) (int64, error)
	GetTasks(ctx context.Context, limit int) ([]domain.Task, error)
	GetTasksByDate(ctx context.Context, date string, limit int) ([]domain.Task, error)
	GetTasksByTitleOrComment(ctx context.Context, search string, limit int) ([]domain.Task, error)
	GetTask(ctx context.Context, id string) (domain.Task, error)
	UpdateTask(ctx context.Context, id string, task domain.Task) error
	UpdateDate(ctx context.Context, id, date string) error
	DeleteTask(ctx context.Context, id string) error
	Close()
}

type SchedulerService struct {
	repo     AdapterInterface
	password string
}

func New(repo AdapterInterface, pass string) *SchedulerService {
	return &SchedulerService{repo: repo, password: pass}
}
