package meta

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
