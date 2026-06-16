package libModel

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

const (
	EqualTo            = "="
	GreaterThan        = ">"
	GreaterThanEqualTo = ">="
	LessThan           = "<"
	EqualToEqualTo     = "<="
	Inequality         = "!="
	IN                 = "in"
	OR                 = "or"
	NOTIN              = "not in" // 新增 NOT IN
	Like               = "like"
	AllLike            = "AllLike"
	Between            = "between"
	GTE                = "gte"
	LTE                = "lte"
)

// 定义一个类型，用于接收闭包
type GormFunc func(*gorm.DB) *gorm.DB

type WhereBuildMap struct {
	Sql  string
	Vals []interface{}
}

// WhereBuild 组装查询条件
func (wb *WhereBuildMap) WhereBuild(condition, key string, value interface{}) {
	if wb.Sql != "" {
		wb.Sql += " AND "
	}
	switch condition {
	case IN:
		wb.Sql += fmt.Sprint(key, " ", IN, " (?) ")
		wb.Vals = append(wb.Vals, value)
	case NOTIN:
		wb.Sql += fmt.Sprintf("%s NOT IN (?)", key)
		wb.Vals = append(wb.Vals, value)
	case Like:
		// 单字段模糊匹配
		wb.Sql += fmt.Sprintf("%s LIKE ?", key)
		wb.Vals = append(wb.Vals, fmt.Sprintf("%v%%", value))
	case AllLike:
		// 单字段模糊匹配
		wb.Sql += fmt.Sprintf("%s LIKE ?", key)
		wb.Vals = append(wb.Vals, fmt.Sprintf("%%%v%%", value))
	case Between:
		switch v := value.(type) {
		case []int64:
			if len(v) == 2 {
				wb.Sql += fmt.Sprintf("%s BETWEEN FROM_UNIXTIME(?) AND FROM_UNIXTIME(?)", key)
				wb.Vals = append(wb.Vals, v[0], v[1])
			} else {
				panic("Between condition requires a slice with two int64 values")
			}
		default:
			panic("Between condition requires []int64")
		}
	case OR:
		// 把刚才可能拼的 " AND " 去掉，改成 " OR "
		wb.Sql = strings.TrimSuffix(wb.Sql, " AND ")
		wb.Sql += fmt.Sprintf(" OR %s IN (?)", key)
		wb.Vals = append(wb.Vals, value)
	case GTE:
		wb.Sql += fmt.Sprintf("%s >= FROM_UNIXTIME(?)", key)
		wb.Vals = append(wb.Vals, value)
	case LTE:
		wb.Sql += fmt.Sprintf("%s <= FROM_UNIXTIME(?)", key)
		wb.Vals = append(wb.Vals, value)
	default:
		wb.Sql += fmt.Sprint(key, condition, " ? ")
		wb.Vals = append(wb.Vals, value)
	}
}

// WhereBuild 组装查询条件
func (wb *WhereBuildMap) CustomWhereBuild(sql string, value interface{}) {
	wb.Sql += sql
	wb.Vals = append(wb.Vals, value)
}
