package imagemanager

import (
	"context"
	"time"
)

type Meta struct {
	FileName    string    `json:"FileName"`
	Tag         string    `schema:"tag" json:"Tag"`
	Title       string    `schema:"title,required" json:"Title"`
	Description string    `schema:"description,required" json:"Description"`
	Size        int64     `json:"Size"`
	DateAdded   time.Time `json:"DateAdded"`
}

type MetaGetter interface {
	GetMeta(ctx context.Context, tags []string) ([]*Meta, error)
	GetRandom(ctx context.Context, n int) ([]*Meta, error)
}

type MetaSetter interface {
	SetMeta(ctx context.Context, meta *Meta) error
}

type MetaRepo interface {
	MetaGetter
	MetaSetter
}
