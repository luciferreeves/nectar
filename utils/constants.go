package utils

import "nectar/types"

// Form field indices using iota for better readability
const (
	FieldConnectionType = iota
	FieldHost
	FieldPort
	FieldSSL
	FieldUser
	FieldPassword
	FieldConnectionName
	FieldColor
)

// SQLite-specific field indices (redefine to match the layout)
const (
	SQLiteFieldConnectionType = iota // 0: Connection Type
	SQLiteFieldDatabaseFile          // 1: Database File
	SQLiteFieldConnectionName        // 2: Connection Name
	SQLiteFieldColor                 // 3: Color
)

// Input field indices using iota
const (
	InputHost = iota
	InputPort
	InputUser
	InputPassword
	InputConnectionName
)

// Database connection defaults
var (
	DefaultPorts = map[types.ConnectionType]string{
		types.PostgreSQL: "5432",
		types.MySQL:      "3306",
		types.SQLite:     "",
	}

	ConnectionTypes = []types.ConnectionType{
		types.PostgreSQL,
		types.MySQL,
		types.SQLite,
	}

	// Total field counts for each database type
	FieldCounts = map[types.ConnectionType]int{
		types.SQLite:     4, // Connection Type, Database File, Connection Name, Color
		types.PostgreSQL: 8, // Connection Type, Host, Port, SSL, User, Password, Connection Name, Color
		types.MySQL:      8, // Same as PostgreSQL
	}
)

// Field mapping for SQLite (simplified structure)
var SQLiteFieldMapping = map[int]int{
	SQLiteFieldConnectionName: InputConnectionName,
}

// Field mapping for PostgreSQL and MySQL (full structure)
var NonSQLiteFieldMapping = map[int]int{
	FieldHost:           InputHost,
	FieldPort:           InputPort,
	FieldUser:           InputUser,
	FieldPassword:       InputPassword,
	FieldConnectionName: InputConnectionName,
}

// Complete field mappings for all database types
var FieldMappings = map[types.ConnectionType]map[int]int{
	types.SQLite:     SQLiteFieldMapping,
	types.PostgreSQL: NonSQLiteFieldMapping,
	types.MySQL:      NonSQLiteFieldMapping,
}
