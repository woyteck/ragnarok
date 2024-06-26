package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"woyteck.pl/ragnarok/models"
	"woyteck.pl/ragnarok/types"
)

type BoilerMemoriesStore struct {
	db *sql.DB
}

func NewBoilerMemoriesStore(db *sql.DB) *BoilerMemoriesStore {
	return &BoilerMemoriesStore{
		db: db,
	}
}

func (s *BoilerMemoriesStore) Truncate(ctx context.Context) error {
	_, err := models.Memories().DeleteAll(ctx, s.db)

	return err
}

func (s *BoilerMemoriesStore) GetMemoryByUUID(ctx context.Context, id uuid.UUID) (*types.Memory, error) {
	memory, err := models.FindMemory(ctx, s.db, id.String())
	if err != nil {
		return nil, err
	}

	return s.memoryModelToType(memory), nil
}

func (s *BoilerMemoriesStore) GetMemoryBySource(ctx context.Context, source string) (bool, *types.Memory, error) {
	count, err := models.Memories().Count(ctx, s.db)
	if err != nil {
		return false, nil, err
	}
	if count == 0 {
		return false, nil, nil
	}

	memory, err := models.Memories(qm.Where("source=?", source)).One(ctx, s.db)
	if err != nil {
		return false, nil, err
	}

	return true, s.memoryModelToType(memory), nil
}

func (s *BoilerMemoriesStore) GetMemories(ctx context.Context, allFields bool, limit int, offset int) ([]*types.Memory, error) {
	memories, err := models.Memories().All(ctx, s.db)
	if err != nil {
		return nil, err
	}

	results := []*types.Memory{}
	for _, memory := range memories {
		results = append(results, s.memoryModelToType(memory))
	}

	return results, nil
}

func (s *BoilerMemoriesStore) memoryModelToType(memory *models.Memory) *types.Memory {
	result := &types.Memory{
		ID: uuid.MustParse(memory.UUID),
	}
	if memory.CreatedAt.Valid {
		result.CreatedAt = &memory.CreatedAt.Time
	}
	if memory.UpdatedAt.Valid {
		result.UpdatedAt = &memory.UpdatedAt.Time
	}
	if memory.DeletedAt.Valid {
		result.DeletedAt = &memory.DeletedAt.Time
	}
	if memory.MemoryType != "" {
		result.MemoryType = memory.MemoryType
	}
	if memory.Source.Valid {
		result.Source = memory.Source.String
	}
	if memory.Content.Valid {
		result.Content = memory.Content.String
	}

	return result
}

func (s *BoilerMemoriesStore) InsertMemory(ctx context.Context, memory *types.Memory) error {
	var m models.Memory

	return m.Insert(ctx, s.db, boil.Infer())
}

func (s *BoilerMemoriesStore) UpdateMemory(ctx context.Context, memory *types.Memory) error {
	m, err := models.FindMemory(ctx, s.db, memory.ID.String())
	if err != nil {
		return err
	}

	m.UpdatedAt = null.TimeFrom(time.Now())
	if memory.DeletedAt != nil {
		m.DeletedAt = null.TimeFrom(*memory.DeletedAt)
	}
	m.MemoryType = memory.MemoryType
	m.Source = null.StringFrom(memory.Source)
	m.Content = null.StringFrom(memory.Content)

	return nil
}
