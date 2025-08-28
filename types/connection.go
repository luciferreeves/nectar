package types

type ConnectionType int

const (
	PostgreSQL ConnectionType = iota
	MySQL
	SQLite
)

func (ct ConnectionType) String() string {
	switch ct {
	case PostgreSQL:
		return "PostgreSQL"
	case MySQL:
		return "MySQL"
	case SQLite:
		return "SQLite"
	default:
		return "Unknown"
	}
}

type Connection struct {
	Name         string
	Type         ConnectionType
	Host         string
	Port         string
	User         string
	Password     string
	Database     string
	DatabaseFile string
	EnableSSL    bool
	Color        string
}
