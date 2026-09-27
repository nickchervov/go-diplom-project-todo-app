package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/nickchervov/go-diplom-project/internal/domain"
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	db *sqlx.DB
}

func migrationUp(db *sqlx.DB) error {
	driver, err := sqlite3.WithInstance(db.DB, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("creating db driver for migrate: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://pkg/migrations", "sqlite", driver)
	if err != nil {
		return fmt.Errorf("create migrate: %w", err)
	}

	if err := m.Up(); err != nil {
		return fmt.Errorf("up migrate: %w", err)
	}
	return nil
}

func New(ctx context.Context, dbFile string) (*Sqlite, error) {
	_, err := os.Stat(dbFile)
	var isNotInstalled bool
	if err != nil {
		isNotInstalled = true
	}

	if isNotInstalled {
		file, err := os.Create(dbFile)
		if err != nil {
			return nil, fmt.Errorf("creating db file: %w", err)
		}
		defer file.Close()

		db, err := sqlx.ConnectContext(ctx, "sqlite", dbFile)
		if err != nil {
			return nil, fmt.Errorf("connection to db: %w", err)
		}

		if err := migrationUp(db); err != nil {
			return nil, fmt.Errorf("migration up: %w", err)
		}
		return &Sqlite{db: db}, nil
	}

	db, err := sqlx.ConnectContext(ctx, "sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("connection to db: %w", err)
	}
	return &Sqlite{db: db}, nil
}

func (s *Sqlite) AddTask(ctx context.Context, task domain.Task) (int64, error) {
	query := "INSERT INTO scheduler (title, comment, repeat, date) VALUES ($1, $2, $3, $4)"
	res, err := s.db.ExecContext(ctx, query, task.Title, task.Comment, task.Repeat, task.Date)
	if err != nil {
		return 0, fmt.Errorf("adding task in db: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting id of last inserted row: %w", err)
	}
	return id, nil
}

func (s *Sqlite) GetTasks(ctx context.Context, limit int) ([]domain.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT $1"
	var tasks []domain.Task
	if err := s.db.SelectContext(ctx, &tasks, query, limit); err != nil {
		return nil, fmt.Errorf("getting tasks: %w", err)
	}
	if tasks == nil {
		return []domain.Task{}, nil
	}
	return tasks, nil
}

func (s *Sqlite) GetTasksByDate(ctx context.Context, date string, limit int) ([]domain.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE date = $1 LIMIT $2"
	var tasks []domain.Task
	if err := s.db.SelectContext(ctx, &tasks, query, date, limit); err != nil {
		return nil, fmt.Errorf("getting tasks by date: %w", err)
	}
	if tasks == nil {
		return []domain.Task{}, nil
	}
	return tasks, nil
}

func (s *Sqlite) GetTasksByTitleOrComment(ctx context.Context, search string, limit int) ([]domain.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE $1 OR comment LIKE $1 ORDER BY date LIMIT $2"
	var tasks []domain.Task
	if err := s.db.SelectContext(ctx, &tasks, query, fmt.Sprint("%"+search+"%"), limit); err != nil {
		return nil, fmt.Errorf("getting tasks by searching: %w", err)
	}
	if tasks == nil {
		return []domain.Task{}, nil
	}
	return tasks, nil
}

func (s *Sqlite) GetTask(ctx context.Context, id string) (domain.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = $1"
	var task domain.Task
	if err := s.db.GetContext(ctx, &task, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Task{}, domain.ErrTaskNotFound
		}
		return domain.Task{}, fmt.Errorf("getting task: %w", err)
	}
	return task, nil
}

func (s *Sqlite) UpdateTask(ctx context.Context, id string, task domain.Task) error {
	query := "UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 WHERE id = $5"
	res, err := s.db.ExecContext(ctx, query, task.Date, task.Title, task.Comment, task.Repeat, id)
	if err != nil {
		return fmt.Errorf("updating task: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting count of affected rows: %w", err)
	}
	if count == 0 {
		return domain.ErrTaskNotFound
	}
	return nil
}

func (s *Sqlite) UpdateDate(ctx context.Context, id, date string) error {
	query := "UPDATE scheduler SET date = $1 WHERE id = $2"
	res, err := s.db.ExecContext(ctx, query, date, id)
	if err != nil {
		return fmt.Errorf("updating date: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting count of affected rows: %w", err)
	}
	if count == 0 {
		return domain.ErrTaskNotFound
	}
	return nil
}

func (s *Sqlite) DeleteTask(ctx context.Context, id string) error {
	query := "DELETE FROM scheduler WHERE id = $1"
	res, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting task: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting count of affected rows: %w", err)
	}
	if count == 0 {
		return domain.ErrTaskNotFound
	}
	return nil
}

func (s *Sqlite) Close() {
	s.db.Close()
}
