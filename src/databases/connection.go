package databases

import (
	"clubhub/configs"
	"clubhub/src/databases/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var e error

func InitializePgDb() {
	log.Println("DB: initializing postgres database")

	var pgDbConfig = configs.PgConfig

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", pgDbConfig.Host, pgDbConfig.User, pgDbConfig.Pass, pgDbConfig.Dbname, pgDbConfig.Port)
	DB, e = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if e != nil {
		panic(e)
	}

	dbgorm, err := DB.DB()
	if err != nil {
		panic(err)
	}

	dbgorm.Ping()

	models.Sync(DB)

	log.Println("DB: initialized postgres database")

}
