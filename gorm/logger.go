package gorm

import (
	"fmt"
	"github.com/make-money-fast/plugins/logger"
	"gorm.io/gorm"
	"strings"
	"time"
)

const keySqlState = "_state"

type loggerPlugin struct {
	logger *logger.Entry
}

func NewLoggerPlugin(logger *logger.Entry) gorm.Plugin {
	return &loggerPlugin{logger: logger}
}

func (p *loggerPlugin) Name() string {
	return "logger"
}

// Initialize registers all needed callbacks
func (p *loggerPlugin) Initialize(db *gorm.DB) (err error) {
	_ = db.Callback().Create().Before("gorm:create").Register("logger:before_create", p.before("insert"))
	_ = db.Callback().Create().After("gorm:create").Register("logger:after_create", p.after("insert"))
	_ = db.Callback().Update().Before("gorm:update").Register("logger:before_update", p.before("update"))
	_ = db.Callback().Update().After("gorm:update").Register("logger:after_update", p.after("update"))
	_ = db.Callback().Query().Before("gorm:query").Register("logger:before_query", p.before("select"))
	_ = db.Callback().Query().After("gorm:query").Register("logger:after_query", p.after("select"))
	_ = db.Callback().Delete().Before("gorm:delete").Register("logger:before_delete", p.before("delete"))
	_ = db.Callback().Delete().After("gorm:delete").Register("logger:after_delete", p.after("delete"))
	_ = db.Callback().Row().Before("gorm:row").Register("logger:before_row", p.before("row"))
	_ = db.Callback().Row().After("gorm:row").Register("logger:after_row", p.after("row"))
	_ = db.Callback().Raw().Before("gorm:raw").Register("logger:before_raw", p.before("raw"))
	_ = db.Callback().Raw().After("gorm:raw").Register("logger:after_raw", p.after("raw"))
	return
}

func (p *loggerPlugin) before(op string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if db == nil || db.Statement == nil || db.Statement.Context == nil {
			return
		}
		db.InstanceSet(keySqlState, &sqlState{
			beginTime: time.Now(),
			op:        op,
		})
	}
}

type sqlState struct {
	beginTime time.Time
	op        string
}

func (p *loggerPlugin) after(op string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		if db == nil || db.Statement == nil || db.Statement.Context == nil {
			return
		}
		v, ok := db.InstanceGet(keySqlState)
		if !ok || v == nil {
			return
		}
		before, ok := v.(*sqlState)
		if !ok || before == nil || before.op != op {
			return
		}

		latency := time.Since(before.beginTime)

		stmt := db.Statement.SQL.String()
		vars := db.Statement.Vars
		sql := formatSQL(stmt, vars)

		if db.Error != nil && !p.isErrorIgnorable(db.Error) {
			p.logger.Error(
				db.Statement.Context, "db exec failed",
				logger.Any("scene", "mysql_client"),
				logger.Err(db.Error),
				logger.Any("table", db.Statement.Table),
				logger.Any("op", op),
				logger.Any("latency", latency),
				logger.Any("sql", sql),
			)
		} else {

			p.logger.Debug(
				db.Statement.Context, "mysql execute succeeded",
				logger.Any("scene", "mysql_client"),
				logger.Any("table", db.Statement.Table),
				logger.Any("op", op),
				logger.Any("latency", latency),
				logger.Any("sql", sql),
			)
		}
	}
}

func (p *loggerPlugin) isErrorIgnorable(err error) bool {
	if err == gorm.ErrRecordNotFound {
		return true
	}
	return false
}

func formatSQL(sql string, vars []interface{}) string {
	for {
		if len(vars) == 0 {
			return sql
		}
		n := strings.Index(sql, "?")
		if n == -1 {
			return sql
		}
		var vv string
		v := vars[0]
		switch t := v.(type) {
		case int, int8, int16, int32, int64, uint32, uint8, uint16, uint64:
			vv = fmt.Sprintf("%d", v)
		case float64, float32:
			vv = fmt.Sprintf("%v", v)
		case time.Time:
			vv = fmt.Sprintf("'%s'", t.Format(`2006-01-02 15:04:05`))
		default:
			vv = fmt.Sprintf("'%s'", v)
		}
		sql = fmt.Sprintf("%s%v%s", sql[:n], vv, sql[n+1:])
		vars = vars[1:]
	}
}
