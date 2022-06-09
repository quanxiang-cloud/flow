package options

import (
	"gorm.io/gorm"
)

// Opt options interface
type Opt interface {
	SetDB(db *gorm.DB)
}

// Options type options functions
type Options func(Opt)

// WithDB set db client to OPT
func WithDB(db *gorm.DB) Options {
	return func(o Opt) {
		o.SetDB(db)
	}
}
