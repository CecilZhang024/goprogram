package initializers

import "awesomeProject1/modules"

func SyncDB() {
	DB.AutoMigrate(&modules.Users{})

}
