package secrets

type Postgres struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"db_name"`
}

type Hasher struct {
	Time    uint32 `mapstructure:"time"`
	Memory  uint32 `mapstructure:"memory"`
	Threads uint8  `mapstructure:"threads"`
	KeyLen  uint32 `mapstructure:"key_len"`
	SaltLen uint32 `mapstructure:"salt_len"`
}

type Secrets struct {
	Postgres Postgres
	Hasher   Hasher
}
