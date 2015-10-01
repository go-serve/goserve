// This file only provides command to generate
// the templates subpackage from "data" directory

//go:generate go-bindata -o templates.go -ignore=gen.go -ignore=templates.go -pkg=assets -prefix=files/ files/...

package assets
