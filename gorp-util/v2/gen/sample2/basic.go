package sample2

import (
	gorpUtil "github.com/zkcrescent/chaos/gorp-util"
)

type (
	// Abstract
	Base struct {
		ID          int64         `db:"id"`
		CreatedTime gorpUtil.Time `db:"created_time"`
	}

	// Abstract
	EnableBase struct {
		*Base
		RemovedTime gorpUtil.Time `db:"removed_time"` // 软删标记
	}
)
