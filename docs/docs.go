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
        "/": {
            "get": {
                "description": "Renders an HTML page with all stored metrics",
                "produces": [
                    "text/html"
                ],
                "tags": [
                    "Metrics"
                ],
                "summary": "Display metrics in HTML",
                "responses": {
                    "200": {
                        "description": "HTML page with metrics",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/updates": {
            "post": {
                "description": "Accepts a JSON array of metrics and updates them in the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Metrics"
                ],
                "summary": "Update multiple metrics",
                "parameters": [
                    {
                        "description": "Array of metrics to update",
                        "name": "metrics",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dto.Metric"
                            }
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Metric"
                            }
                        }
                    },
                    "400": {
                        "description": "Invalid input",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "401": {
                        "description": "Missing or invalid Hash header",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/value/{type}/{name}": {
            "get": {
                "description": "Returns the value of a specified metric by its type and name",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Metrics"
                ],
                "summary": "Retrieve a metric value",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Metric type",
                        "name": "type",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Metric name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "number"
                        }
                    },
                    "404": {
                        "description": "Metric not found",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/{type}/{name}/{value}": {
            "post": {
                "description": "Accepts a metric value from URL parameters and updates or creates the metric",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Metrics"
                ],
                "summary": "Update or create a single metric",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Metric type (e.g., gauge, counter)",
                        "name": "type",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Metric name",
                        "name": "name",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "number",
                        "description": "Metric value",
                        "name": "value",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Metric"
                        }
                    },
                    "400": {
                        "description": "Invalid Value",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dto.Metric": {
            "type": "object",
            "properties": {
                "delta": {
                    "description": "Изменение для метрик типа counter",
                    "type": "integer"
                },
                "name": {
                    "description": "Название метрики",
                    "type": "string"
                },
                "type": {
                    "description": "Тип метрики (counter или gauge)",
                    "type": "string"
                },
                "value": {
                    "description": "Значение для метрик типа gauge",
                    "type": "number"
                }
            }
        },
        "models.Metric": {
            "type": "object",
            "properties": {
                "created_at": {
                    "description": "Время создания метрики",
                    "type": "string"
                },
                "delta": {
                    "description": "Изменение для метрик типа counter",
                    "type": "integer"
                },
                "id": {
                    "description": "Уникальный идентификатор метрики",
                    "type": "integer"
                },
                "name": {
                    "description": "Название метрики",
                    "type": "string"
                },
                "type": {
                    "description": "Тип метрики (counter или gauge)",
                    "type": "string"
                },
                "value": {
                    "description": "Значение для метрик типа gauge",
                    "type": "number"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}