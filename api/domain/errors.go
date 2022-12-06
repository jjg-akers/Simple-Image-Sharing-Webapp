package domain

import "errors"

var ErrNotFound = errors.New("No images found in db for given tag")
