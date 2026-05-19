package pagination

import (
	"math"

	"gorm.io/gorm"
)

type Pagination struct {
	Page    int
	PerPage int
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetPerPage()
}

func (p *Pagination) GetPerPage() int {
	if p.PerPage <= 0 {
		return 10
	}
	return p.PerPage
}

func (p *Pagination) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

func Paginate(page, perPage int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		p := &Pagination{Page: page, PerPage: perPage}
		return db.Offset(p.GetOffset()).Limit(p.GetPerPage())
	}
}

func GetTotalPage(totalData int64, perPage int) int {
	p := &Pagination{PerPage: perPage}
	return int(math.Ceil(float64(totalData) / float64(p.GetPerPage())))
}
