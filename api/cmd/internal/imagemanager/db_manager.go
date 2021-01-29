package imagemanager

import "context"

type Inserter interface {
	Insert(ctx context.Context, imageName, tag string) error
}

// type Searcher interface {
// 	Search(ctx context.Context, tag string) ([]string, error)
// }

type DBImageManager interface {
	Inserter
	Searcher
}
