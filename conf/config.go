package conf

var GlobalConfig = Config{
	Server: Server{
		Addr: ":8088",
	},
	Redis: Redis{
		Addr: "14.103.235.1:6380",
		Pass: "151536123456",
	},
	SQL: SQL{
		DriverName: "postgres",
		Src:        "user=postgres dbname=bs_2025_2_12 password=151536123456 host=118.190.152.69 port=15431 sslmode=disable",
		Var:        "Dollar",
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

var Grok = LLMApi{
	Model:  "grok-3",
	Url:    "https://api-proxy.me/xai/v1/chat/completions",
	Apikey: "xai-4Kag7Eqy8UNK1zCoXEUuELtwLIgTJ4DmXqWDryuVzSsAf30YgsZ05wRPTtCqmoVkJXqwMsC75A4mIgyR",
}

var Claude = LLMApi{
	Model:  "discount:claude-sonnet-4-20250514",
	Url:    "https://yourapi.cn/v1/chat/completions",
	Apikey: "sk-vgKgdYdwigIjkcRG4IzoRWpKuD67AHiC9ol0BBnYIRYRNEeh",
}
