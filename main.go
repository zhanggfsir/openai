package main

import (
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	_ "openai-backend/routers"
	_ "openai-backend/services"
	"openai-backend/utils/commands"
	_ "openai-backend/utils/redis"
)

func main() {

	commands.Bootstrap()

	beego.Run()

}

