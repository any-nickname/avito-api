package main

import "avito-rest-api/internal/app"

// @Title 			Сервис по работе с сегментами
// @Version 		1.0
// @Description 	Сервис предоставляет функционал для работы с сегментами и пользователями для целей аналитики

// @Contact.name 	Грищенко Владимир
// @Contact.email 	vladimirxsky@gmail.com

// @Host 			localhost:8080
// @BasePath 		/

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
