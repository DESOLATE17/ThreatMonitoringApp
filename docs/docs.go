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
        "/api/check-auth": {
            "get": {
                "description": "Retrieves user information based on the provided user context",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Check user authentication",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.User"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/logout": {
            "post": {
                "description": "Logs out the user by blacklisting the access token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Logout",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    }
                }
            }
        },
        "/api/monitoring-request-threats/threats/{threatId}": {
            "delete": {
                "description": "Deletes a threat from a request based on the user ID and threat ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MonitoringRequests"
                ],
                "summary": "Delete threat from request",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Threat ID",
                        "name": "threatId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/api/monitoring-requests": {
            "get": {
                "description": "Retrieves a list of monitoring requests based on the provided parameters",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MonitoringRequests"
                ],
                "summary": "Get list of monitoring requests",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Monitoring request status",
                        "name": "status",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Start date in the format '2006-01-02T15:04:05Z'",
                        "name": "start_date",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "End date in the format '2006-01-02T15:04:05Z'",
                        "name": "end_date",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.MonitoringRequest"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "delete": {
                "description": "Deletes a monitoring request for the given user ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MonitoringRequests"
                ],
                "summary": "Delete monitoring request by user ID",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/api/monitoring-requests/client": {
            "put": {
                "description": "Updates the status of a monitoring request by client on formated",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MonitoringRequests"
                ],
                "summary": "Update monitoring request status by client",
                "parameters": [
                    {
                        "description": "New status of the monitoring request",
                        "name": "newStatus",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.NewStatus"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/api/monitoring-requests/{id}": {
            "get": {
                "description": "Retrieves a monitoring request with the given ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MonitoringRequests"
                ],
                "summary": "Get monitoring request by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Monitoring Request ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/api/signIn": {
            "post": {
                "description": "Authenticates a user and generates an access token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "User sign-in",
                "parameters": [
                    {
                        "description": "User information",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserLogin"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/api/signUp": {
            "post": {
                "description": "Creates a new user account",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Authentication"
                ],
                "summary": "Sign up a new user",
                "parameters": [
                    {
                        "description": "User information",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.UserSignUp"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "409": {
                        "description": "Conflict",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/api/threats": {
            "get": {
                "description": "Retrieves a list of threats based on the provided query.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Threats"
                ],
                "summary": "Get threats list",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Query string to filter threats",
                        "name": "query",
                        "in": "query"
                    },
                    {
                        "type": "number",
                        "description": "LowPrice to filter threats",
                        "name": "lowPrice",
                        "in": "query"
                    },
                    {
                        "type": "number",
                        "description": "HighPrice string to filter threats",
                        "name": "highPrice",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            },
            "post": {
                "description": "Add a new threat with image, name, description, summary, count, and price",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Threats"
                ],
                "summary": "Add new threat",
                "parameters": [
                    {
                        "type": "file",
                        "description": "Threat image",
                        "name": "image",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Threat name",
                        "name": "name",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Threat description",
                        "name": "description",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Threat summary",
                        "name": "summary",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "Threat count",
                        "name": "count",
                        "in": "formData",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Threat price",
                        "name": "price",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    }
                }
            }
        },
        "/api/threats/request/{threatId}": {
            "post": {
                "description": "Adds a threat to a monitoring request",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Threats"
                ],
                "summary": "Add threat to request",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Threat ID",
                        "name": "threatId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/api/threats/{id}": {
            "get": {
                "description": "Retrieves a threat by its ID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Threats"
                ],
                "summary": "Get threat by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Threat ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Threat"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            },
            "put": {
                "description": "Updates a threat with the given ID",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Threats"
                ],
                "summary": "Update threat by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "name",
                        "name": "name",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "description",
                        "name": "description",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "count",
                        "name": "count",
                        "in": "formData"
                    },
                    {
                        "type": "string",
                        "description": "price",
                        "name": "price",
                        "in": "formData"
                    },
                    {
                        "type": "file",
                        "description": "image",
                        "name": "image",
                        "in": "formData"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            },
            "delete": {
                "description": "Deletes a threat with the given ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Threats"
                ],
                "summary": "Delete threat by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Threat ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        },
        "/monitoring-requests/admin/{requestId}": {
            "put": {
                "description": "Updates the status of a monitoring request with the given ID on \"accepted\"/\"closed\"/\"canceled\"",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MonitoringRequests"
                ],
                "summary": "Update monitoring request status by ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Request ID",
                        "name": "requestId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "New request status",
                        "name": "newRequestStatus",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.NewStatus"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    }
                }
            }
        }
    },
    "definitions": {
        "models.MonitoringRequest": {
            "type": "object",
            "properties": {
                "admin": {
                    "type": "string"
                },
                "adminId": {
                    "type": "integer"
                },
                "creationDate": {
                    "type": "string"
                },
                "creator": {
                    "type": "string"
                },
                "endingDate": {
                    "type": "string"
                },
                "formationDate": {
                    "type": "string"
                },
                "requestId": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "models.NewStatus": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                }
            }
        },
        "models.Threat": {
            "type": "object",
            "properties": {
                "count": {
                    "type": "integer"
                },
                "description": {
                    "type": "string"
                },
                "image": {
                    "type": "string"
                },
                "isDeleted": {
                    "type": "boolean"
                },
                "name": {
                    "type": "string"
                },
                "price": {
                    "type": "integer"
                },
                "summary": {
                    "type": "string"
                },
                "threatId": {
                    "type": "integer"
                }
            }
        },
        "models.User": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "isAdmin": {
                    "type": "boolean"
                },
                "login": {
                    "type": "string",
                    "maxLength": 64
                },
                "name": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 8
                },
                "registrationDate": {
                    "type": "string"
                },
                "userId": {
                    "type": "integer"
                }
            }
        },
        "models.UserLogin": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string",
                    "maxLength": 64
                },
                "password": {
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 8
                }
            }
        },
        "models.UserSignUp": {
            "type": "object",
            "required": [
                "login",
                "password"
            ],
            "properties": {
                "login": {
                    "type": "string",
                    "maxLength": 64
                },
                "password": {
                    "type": "string",
                    "maxLength": 64,
                    "minLength": 8
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3001",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "ThreatMonitoringApp",
	Description:      "App for serving threats monitoring requests",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
