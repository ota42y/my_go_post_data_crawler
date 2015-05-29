package config

import (
	"gopkg.in/yaml.v2"
)

// MongodbDatabase is server setting
type MongodbDatabase struct {
	URL      string
	User     string
	Pass     string
	Database string
}

// GetDialURL return string for mgo.dial
func (m *MongodbDatabase) GetDialURL() string {
	url := m.URL

	if m.User != "" && m.Pass != "" {
		url = m.User + ":" + m.Pass + "@" + url
	}
	if m.Database != "" {
		url = url + "/" + m.Database
	}

	return url
}

// MongodbBackup is backup archive path and workspace
type MongodbBackup struct {
	ArchivePath string
	Workspace   string
}

// NewMongodbDatabaseFromData return MongodbDatabase from yaml data
func NewMongodbDatabaseFromData(data []byte) *MongodbDatabase {
	var d MongodbDatabase

	err := yaml.Unmarshal(data, &d)
	if err != nil {
		return nil
	}
	return &d
}

// NewMongodbBackupFromData return MongodbBackup from yaml data
func NewMongodbBackupFromData(data []byte) *MongodbBackup {
	var d MongodbBackup

	err := yaml.Unmarshal(data, &d)
	if err != nil {
		return nil
	}

	return &d
}
