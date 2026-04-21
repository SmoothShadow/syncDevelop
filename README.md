## 实现文本、截图和文件实时同步

### 用于内网多台电脑之间随时同步数据,抛弃通讯软件在线文档这种需要将数据上传到公网的服务

#### 使用方法

```
# windows
GOOS=windows GOARCH=amd64 go build -buildvcs=false -o bin/syncdevelop-win.exe ./cmd/

# mac
GOOS=darwin GOARCH=amd64 go build -buildvcs=false -o bin/syncdevelop-mac ./cmd/
```
