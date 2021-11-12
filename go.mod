module DevIntergTest

go 1.15

//go.mod部屬方式採local檔案管理

replace github.com/gin-gonic/gin v1.7.2 => ./pkg/src/github.com/gin-gonic/gin@v1.7.2

replace gorm.io/driver/mysql v1.1.1 => ./pkg/src/gorm.io/driver/mysql@v1.1.1

replace gorm.io/gorm v1.21.11 => ./pkg/src/gorm.io/gorm@v1.21.11

replace github.com/golang/gddo v0.0.0 => ./pkg/src/github.com/golang/gddo@v0.0.0-20210115222349-20d68f94ee1f

replace github.com/thinkoner/openssl v0.0.0 => ./pkg/src/github.com/thinkoner/openssl

replace github.com/mattn/go-isatty v0.0.12 => ./pkg/src/github.com/mattn/go-isatty@v0.0.12

replace github.com/go-sql-driver/mysql v1.6.0 => ./pkg/src/github.com/go-sql-driver/mysql@v1.6.0

replace github.com/ugorji/go v1.1.7 => ./pkg/src/github.com/ugorji/go@v1.1.7

replace github.com/jinzhu/inflection v1.0.0 => ./pkg/src/github.com/jinzhu/inflection@v1.0.0

replace github.com/gin-contrib/sse v0.1.0 => ./pkg/src/github.com/gin-contrib/sse@v0.1.0

replace github.com/golang/protobuf v1.3.3 => ./pkg/src/github.com/golang/protobuf@v1.3.3

replace github.com/go-playground/validator/v10 v10.4.1 => ./pkg/src/github.com/go-playground/validator/v10@v10.4.1

replace github.com/jinzhu/now v1.1.2 => ./pkg/src/github.com/jinzhu/now@v1.1.2

replace gopkg.in/yaml.v2 v2.2.8 => ./pkg/src/gopkg.in/yaml.v2@v2.2.8

replace github.com/ugorji/go/codec v1.1.7 => ./pkg/src/github.com/ugorji/go/codec@v1.1.7

replace golang.org/x/sys v0.0.0-20200116001909-b77594299b42 => ./pkg/src/golang.org/x/sys@v0.0.0-20200116001909-b77594299b42

replace github.com/go-playground/universal-translator v0.17.0 => ./pkg/src/github.com/go-playground/universal-translator@v0.17.0

replace github.com/leodido/go-urn v1.2.0 => ./pkg/src/github.com/leodido/go-urn@v1.2.0

replace golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 => ./pkg/src/golang.org/x/crypto@v0.0.0-20200622213623-75b288015ac9

replace github.com/go-playground/locales v0.13.0 => ./pkg/src/github.com/go-playground/locales@v0.13.0

// 取消不使用
// replace github.com/fdiminsions/finance.lib => ./pkg/src/finance.lib/
// replace pkg/shuttle/database_shuttle => ./pkg/src/shuttle/database_shuttle
// replace github.com/go-sql-driver/mysql v1.6.0 => ./pkg/src/github.com/go-sql-driver/mysql@v1.6.0
// replace github.com/gorilla/mux v1.8.0 => ./pkg/src/github.com/gorilla/mux@v1.8.0

require (
	github.com/gin-gonic/gin v1.7.2
	github.com/golang/gddo v0.0.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/thinkoner/openssl v0.0.0
	gorm.io/driver/mysql v1.1.1
	gorm.io/gorm v1.21.11
)
