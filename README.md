# Golang 错误号处理

Conventional Golang Error

# Info

## 优势

1. 规范了错误号的使用方式
2. 方便快捷的错误定义和创建，提供了 `New`、`Newf`、`Neww` 三种创建错误方式
3. 约束了错误号的唯一性
4. 兼顾了 HTTP 返回和业务层逻辑错误号

## 劣势

无

# Usage

## 导入包
```go
import github.com/chenjinya/errors
```

## 定义错误号
可以在自己项目里定义 `ErrCode` ，建议使用 6 位错误号，以与默认错误号区分。
比如

```go
var CreateCardFailError = errors.NewErrCode(100101, http.StatusInternalServerError, "创建卡片失败")
```

## 使用错误号
也可以直接使用包里提供的默认错误号

```go
func CreateUser(ID int64) error {
	if ID == 0 {
        return errors.ParamError.New("缺少 ID", nil)
    }
    ...
}
```

使用 format 方法

```go
func CreateUser(ID int64) error {
	if ID == 0 {
        return errors.ParamError.Newf("缺少参数 %s", "ID") // 可以不传递 err，默认为 nil
    }
    ...
}
```

使用 format 方法带 `error`

```go
func CreateUser(ID int64) error {
	...
	err := dao.CreateUser(ID int64)
	if err != nil { 
		return errors.DbError.Newf("创建用户数据失败, ID: %v", ID, err) // 传递原始 error
    }
    ...
}
```

使用默认错误信息带 `error`

```go
func CreateUser(ID int64) error {
	...
	err := dao.CreateUser(ID int64)
	if err != nil { 
		return errors.DbError.Neww(err)
    }
    ...
}
```

## API 返回

结合 `gin` 使用：

```go
func NotOK(w http.ResponseWriter, err error) {
	statusCode := http.StatusInternalServerError
	code := statusCode
	message := "Unexpected Error"
	if baseError, ok := err.(errors.BaseErrorInterface); ok {
		statusCode = int(baseError.StatusCode())
		code = int(baseError.Code())
		message = baseError.Message()
	} else {
		if err != nil {
			message = err.Error()
		}
	}

	log.Errorf(nil, "[ERROR] %s", err)
}
```
