package entity

import (
	"code.rocketnine.space/tslocum/gohan"
	"code.rocketnine.space/tslocum/pretzel-tycoon/component"
)

func NewOnceEntity() gohan.Entity {
	e := gohan.NewEntity()
	e.AddComponent(&component.Once{})
	return e
}
