package conf

var GlobalConfig = Config{
	Server: Server{
		Addr: ":8088",
	},
	Redis: Redis{
		Addr: "",
		Pass: "",
	},
	SQL: SQL{
		DriverName: "",
		Src:        "",
		Var:        "",
	},
}

type Config struct {
	Server Server
	Redis  Redis
	SQL    SQL
}

type Server struct {
	Addr string
}

type Redis struct {
	Addr string
	Pass string
}

type SQL struct {
	DriverName string
	Src        string
	Var        string
}

type LLMApi struct {
	Model  string
	Url    string
	Apikey string
}
