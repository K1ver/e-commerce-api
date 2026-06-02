package docs

import (
	_ "embed"

	"github.com/swaggo/swag"
)

//go:embed swagger.json
var swaggerSpec []byte

// SwaggerInfo exported metadata
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Title:            "E-Commerce API",
	Description:      "E-Commerce API with JWT, roles (admin/seller/buyer), cart, orders and payments.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  string(swaggerSpec),
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
