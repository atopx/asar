# asar

`asar` 是一个用于压缩和解压缩 ASAR（Atom Shell Archive）格式文件的 Go 库。


## 安装

```bash
go get -u github.com/atopx/asar@latest
```

## 使用

### 导入包

```go
import "github.com/atopx/asar"
```

### 压缩目录

`Pack` 函数用于将指定目录压缩成 ASAR 文件。

```go
err := asar.Pack("path/to/directory", "path/to/destination.asar")
if err != nil {
    log.Fatal(err)
}
```

### 解压缩 ASAR 文件

`Unpack` 函数用于将 ASAR 文件解压缩到指定目录。

```go
err := asar.Unpack("path/to/source.asar", "path/to/destination")
if err != nil {
    log.Fatal(err)
}
```

## 示例

以下是一个完整的示例，用于将目录压缩成 ASAR 文件，然后解压缩该 ASAR 文件。

```go
package main

import (
    "log"
    "github.com/atopx/asar"
)

func main() {
    // 压缩目录
    err := asar.Pack("path/to/directory", "path/to/destination.asar")
    if err != nil {
        log.Fatal(err)
    }

    // 解压缩 ASAR 文件
    err = asar.Unpack("path/to/source.asar", "path/to/destination")
    if err != nil {
        log.Fatal(err)
    }
}
```

## 贡献

欢迎提交 Issue 和 Pull Request。

## 许可证

[MIT](./LICENSE)
