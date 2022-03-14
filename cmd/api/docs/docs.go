package docs

import (
	_ "embed"
	"github.com/swaggo/swag"
)

//go:embed swagger.yaml
var doc string

func init() {
	swaggerInfo := &swag.Spec{
		Version:          "",
		Host:             "",
		BasePath:         "",
		Schemes:          []string{},
		Title:            "",
		Description:      "",
		InfoInstanceName: "swagger",
		SwaggerTemplate:  doc,
	}

	swag.Register(swaggerInfo.InstanceName(), swaggerInfo)
}
