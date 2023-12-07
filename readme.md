### 集成一些常用的 golang 插件.

#### `gorm` 日志

```go
    var (
    db *gorm.DB
    )
    db.Use(NewLoggerPlugin(logger))
```

#### gin 日志

```go
    g.Use(middleware.Logger(logger, middleware.LogConfigure{}))
```

#### 系统日志. 
```go
    log := logger.NewLogger(
		logger.RotateFile("./log.log",100,28)
	)
	_ = log 
```