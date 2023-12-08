// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/example/helloworld": {
            "get": {
                "description": "do ping",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "example"
                ],
                "summary": "ping example",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/login": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "summary": "登录用户",
                "parameters": [
                    {
                        "description": "用户名和密码",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.LoginParam"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/user/signup": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "用户模块"
                ],
                "summary": "注册用户",
                "parameters": [
                    {
                        "description": "用户名和密码",
                        "name": "object",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.SignupParam"
                        }
                    }
                ],
                "responses": {}
            }
        }
    },
    "definitions": {
        "models.LoginParam": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "description": "binding:\"required\"表示必须要传",
                    "type": "string"
                }
            }
        },
        "models.SignupParam": {
            "type": "object",
            "required": [
                "email",
                "gender",
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "gender": {
                    "type": "integer"
                },
                "password": {
                    "type": "string"
                },
                "username": {
                    "description": "binding:\"required\"表示必须要传，gin的validator使用",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "latest",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "go_web_app项目接口文档",
	Description:      "go_web_app项目接口文档server端api文档",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}