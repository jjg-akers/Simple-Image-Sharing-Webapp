package imagemanager

import "io"

type ImageV1 struct {
	Meta *Meta `json:"Meta"`
	URI  string     `json:"url"`
	File io.Reader  `json:"File"`
}