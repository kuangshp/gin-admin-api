package utils

import (
	"gorm.io/gen"
	"net/http"
	"strconv"
)

func Paginate(r *http.Request) func(db gen.Dao) gen.Dao {
	return func(db gen.Dao) gen.Dao {
		q := r.URL.Query()
		pageNumber, _ := strconv.Atoi(q.Get("pageNumber"))
		if pageNumber == 0 {
			pageNumber = 1
		}

		pageSize, _ := strconv.Atoi(q.Get("pageSize"))
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (pageNumber - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
