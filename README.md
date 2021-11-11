# SudyCookieParser
SudyTechCookieParser


```shell
go get github.com/xingxp/SudyCookieParser
```
it is very easy to use this lib
```go
import (
	"github.com/sirupsen/logrus"
	"github.com/xingxp/SudyCookieParser"
)

func main() {

	var cookie = "BMCoPL4BAAAAIgkAANO4Aet8AQAAgO42..."
	sudyCookieInfo := sudy_cookie_parser.NewSudyCookieInfo(cookie)
	logrus.Info(sudyCookieInfo)

}
```
