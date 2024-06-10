package config

const (
	DEFAULT_DB_USER     = "user"
	DEFAULT_DB_PASSWORD = "123456"
	DEFAULT_DB_NAME     = "testdb"
	DEFAULT_DB_HOST     = "localhost"
	DEFAULT_DB_PORT     = "5543"
	DEFAULT_DB_SSLMODE  = "disable"

	SERVER_PORT = "8085"
	DEFAULT_URL = "http://localhost:8085"

	USER_FIRST_NAME = "Yakir"
	USER_LAST_NAME  = "Nisim"
	USER_PHONE      = "123-456-7890"
	USER_EMAIL      = "Yakir@gmail.com"
	USER_PASSWORD   = "123"

	DEFAULT_COMPANY_NAME    = "Tech Innovators Inc."
	DEFAULT_COMPANY_ADDRESS = "1234 Innovation Way, Silicon Valley, CA 94043"
	DEFAULT_COMPANY_LOGO    = "logo.png"
)

var TestConfig TConfig

type TConfig struct {
	TestDatabase TestDatabase
	TestServer   TestServer
	User         User
	Company      Company
}

type TestDatabase struct {
	UserName string
	Password string
	DbName   string
	Host     string
	Port     string
	SSLMode  string
}

type TestServer struct {
	Port string
	URL  string
}

type User struct {
	FirstName string
	LastName  string
	Phone     string
	Email     string
	Password  string
}

type Company struct {
	Name    string
	Address string
	Logo    string
}

func init() {
	TestConfig.TestDatabase.loadConfig()
	TestConfig.TestServer.loadConfig()
	TestConfig.User.loadConfig()
	TestConfig.Company.loadConfig()
}

func (td *TestDatabase) loadConfig() {
	td.UserName = DEFAULT_DB_USER
	td.Password = DEFAULT_DB_PASSWORD
	td.DbName = DEFAULT_DB_NAME
	td.Host = DEFAULT_DB_HOST
	td.Port = DEFAULT_DB_PORT
	td.SSLMode = DEFAULT_DB_SSLMODE
}

func (ts *TestServer) loadConfig() {
	ts.Port = SERVER_PORT
	ts.URL = DEFAULT_URL
}

func (u *User) loadConfig() {
	u.FirstName = USER_FIRST_NAME
	u.LastName = USER_LAST_NAME
	u.Phone = USER_PHONE
	u.Email = USER_EMAIL
	u.Password = USER_PASSWORD
}

func (c *Company) loadConfig() {
	c.Name = DEFAULT_COMPANY_NAME
	c.Address = DEFAULT_COMPANY_ADDRESS
	c.Logo = DEFAULT_COMPANY_LOGO
}
