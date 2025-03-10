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
        "/song": {
            "post": {
                "description": "Создает новую песню с указанными данными",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Создание новой песни",
                "parameters": [
                    {
                        "description": "Данные песни",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.CreateSongParams"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Данные созданной песни",
                        "schema": {
                            "$ref": "#/definitions/entities.SongData"
                        }
                    },
                    "400": {
                        "description": "Неверный формат запроса",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "409": {
                        "description": "Песня уже существует",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "502": {
                        "description": "Ошибка внешнего сервиса",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/song/{id}": {
            "delete": {
                "description": "Удаляет песню по указанному ID",
                "tags": [
                    "songs"
                ],
                "summary": "Удаление песни",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Песня успешно удалена"
                    },
                    "400": {
                        "description": "Неверный формат ID",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            },
            "patch": {
                "description": "Обновляет информацию о песне по указанному ID",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Обновление данных песни",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Обновляемые данные песни",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.PatchSongParams"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "Песня успешно обновлена"
                    },
                    "400": {
                        "description": "Неверный формат запроса или ID",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Песня не найдена",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/song/{id}/lyrics": {
            "get": {
                "description": "Получить куплеты песни по ID песни с возможностью пагинации",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "lyrics"
                ],
                "summary": "Получить текст песни",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "С какого куплета начать",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Сколько куплетов вывести для пагинации",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Текст песни успешно получен",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.LyricsVerseData"
                            }
                        }
                    },
                    "400": {
                        "description": "Неверный запрос",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Текст песни не найден",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/songs": {
            "get": {
                "description": "Возвращает список песен с возможностью фильтрации",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Получение списка песен",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Название группы",
                        "name": "group",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Название песни",
                        "name": "song",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Дата релиза от (формат: 2006-01-02)",
                        "name": "release_date_from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Дата релиза до (формат: 2006-01-02)",
                        "name": "release_date_to",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "С какой песни выводить",
                        "name": "offset",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Сколько песен выводить",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Список песен",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entities.SongData"
                            }
                        }
                    },
                    "400": {
                        "description": "Неверный формат запроса",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "404": {
                        "description": "Песни не найдены",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entities.LyricsVerseData": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "index": {
                    "type": "integer"
                }
            }
        },
        "entities.SongData": {
            "type": "object",
            "properties": {
                "group": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "link": {
                    "type": "string"
                },
                "release_date": {
                    "type": "string"
                },
                "song": {
                    "type": "string"
                }
            }
        },
        "handlers.CreateSongParams": {
            "type": "object",
            "required": [
                "group",
                "song"
            ],
            "properties": {
                "group": {
                    "type": "string",
                    "minLength": 1
                },
                "song": {
                    "type": "string",
                    "minLength": 1
                }
            }
        },
        "handlers.ErrorResponse": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "string"
                }
            }
        },
        "handlers.PatchSongParams": {
            "type": "object",
            "properties": {
                "group": {
                    "type": "string",
                    "minLength": 1
                },
                "link": {
                    "type": "string",
                    "minLength": 1
                },
                "lyrics": {
                    "type": "string",
                    "minLength": 1
                },
                "release_date": {
                    "type": "string"
                },
                "song": {
                    "type": "string",
                    "minLength": 1
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
