package imagemanager

import (
	"io"
	"time"
)

type Image struct {
	Name      string
	File      io.Reader
	Tag       string
	Size      int64
	DateAdded time.Time
}

// id int(11) not null primary key auto_increment,
// uri varchar(255) not null,
// title varchar(100) not null,
// tag varchar(100),
// date_added datetime not null
