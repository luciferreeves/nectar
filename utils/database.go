package utils

import "nectar/types"

// GetDefaultPort returns the default port for a given connection type
func GetDefaultPort(connType types.ConnectionType) string {
	if port, exists := DefaultPorts[connType]; exists {
		return port
	}
	return ""
}

// GetFieldCount returns the total number of fields for a connection type
func GetFieldCount(connType types.ConnectionType) int {
	if count, exists := FieldCounts[connType]; exists {
		return count
	}
	return 0
}

// GetInputIndex returns the input index for a given field based on connection type
func GetInputIndex(connType types.ConnectionType, fieldIndex int) int {
	if mapping, exists := FieldMappings[connType]; exists {
		if inputIndex, exists := mapping[fieldIndex]; exists {
			return inputIndex
		}
	}
	return -1
}

// NextConnectionType cycles to the next connection type
func NextConnectionType(current types.ConnectionType) types.ConnectionType {
	for i, connType := range ConnectionTypes {
		if connType == current {
			return ConnectionTypes[(i+1)%len(ConnectionTypes)]
		}
	}
	return current
}

// PrevConnectionType cycles to the previous connection type
func PrevConnectionType(current types.ConnectionType) types.ConnectionType {
	for i, connType := range ConnectionTypes {
		if connType == current {
			prevIndex := (i - 1 + len(ConnectionTypes)) % len(ConnectionTypes)
			return ConnectionTypes[prevIndex]
		}
	}
	return current
}
