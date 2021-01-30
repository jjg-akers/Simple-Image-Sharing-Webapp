package meta

import (
	"context"
	"time"
)

type Meta struct {
	FileName    string
	Tag         string `schema:"tag"`
	Title       string `schema:"title,required`
	Description string `schema:"description,required"`
	Size        int64
	DateAdded   time.Time
}

type Getter interface {
	Get(ctx context.Context, tags []string) ([]*Meta, error)
}

type Setter interface {
	Set(ctx context.Context, meta *Meta) error
}

type GetterSetter interface {
	Getter
	Setter
}
