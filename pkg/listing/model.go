package listing

import (
	"time"
)

type Ad struct {
	Id            int
	Url           string
	Created       time.Time
	Price         string
	SqFootage     string
	Street        string
	Neighbourhood string
	Location      string
}
