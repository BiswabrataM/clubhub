package configs

type DBConfig struct {
	Host   string
	Port   int
	User   string
	Pass   string
	Dbname string
}

var PgConfig = DBConfig{
	Host:   "localhost",
	Port:   5432,
	User:   "postgres",
	Pass:   "postgres",
	Dbname: "go4",
}
