module github.com/loveuer/go-sqlite3/gormlite

go 1.21

toolchain go1.23.0

require (
	github.com/ncruces/go-sqlite3 v0.18.3
	gorm.io/gorm v1.25.11
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/ncruces/julianday v1.0.0 // indirect
	github.com/tetratelabs/wazero v1.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)

replace (
	github.com/ncruces/go-sqlite v0.18.3 => github.com/loveuer/go-sqlite v1.0.0
)
