package mongodb

import (
	"../../../config"
	"../../../lib/logger"
	"fmt"
	"github.com/mattn/go-shellwords"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
	"os/exec"
	"path"
	"time"
)

var tagName = "backup.mongodb"

// Mongodb is mongodb daily backup work
type Mongodb struct {
	backup *config.MongodbBackup
	mongo  *config.MongodbDatabase
	logger logger.Logger
}

// NewMongodb return Mongodb struct
func NewMongodb(bk *config.MongodbBackup, mongo *config.MongodbDatabase, logger logger.Logger) *Mongodb {
	return &Mongodb{
		backup: bk,
		mongo:  mongo,
		logger: logger,
	}
}

// Execute execute backup and convert zip file
func (m *Mongodb) Execute() {
	m.createExpireIndex()
	m.backupData()
}

func (m *Mongodb) repairDatabase() {
	session, err := mgo.Dial(m.mongo.GetDialURL())
	if err != nil {
		m.logger.Error(tagName, "database dial error %s", err)
		return
	}
	defer session.Close()

	result := bson.M{}
	d := session.DB(m.mongo.Database)
	err = d.Run(bson.D{{"repairDatabase", 1}}, &result)
	if err != nil {
		m.logger.Error(tagName, "repairDatabase error %s", err)
	}
	m.logger.Debug(tagName, "%s", result)
}

func (m *Mongodb) backupData() {
	m.logger.Debug(tagName, "bakup start")

	// create date
	start := time.Now().AddDate(0, 0, -3)
	start = start.Truncate(time.Hour * 24).Add(time.Hour * -9)
	end := time.Now()

	// reset workspace
	err := os.RemoveAll(m.backup.Workspace)
	if err != nil {
		m.logger.Error(tagName, "os.RemoveAll(%s) error", m.backup.Workspace)
		m.logger.Error(tagName, err.Error())
		return
	}
	saveFolder := path.Join(m.backup.Workspace, end.Format("2006-01-02"))
	if _, err := os.Stat(saveFolder); os.IsNotExist(err) {
		os.MkdirAll(saveFolder, 0777)
	}

	m.dump(start, end, saveFolder)
	m.archive(end, saveFolder)
}

func (m *Mongodb) getAllCollectionNames() []string {
	var names []string

	session, err := mgo.Dial(m.mongo.GetDialURL())
	if err != nil {
		m.logger.Error(tagName, "database dial error %s", err)
		return names
	}
	defer session.Close()

	d := session.DB(m.mongo.Database)
	rawNames, err := d.CollectionNames()
	if err != nil {
		m.logger.Error(tagName, "get collection names error %s", err)
		return names
	}

	// delete system.indexes
	for _, n := range rawNames {
		if n != "system.indexes" {
			names = append(names, n)
		}
	}

	return names
}

func (m *Mongodb) dump(start time.Time, end time.Time, saveFolder string) {
	startMsec := start.Unix() * 1000
	endMsec := end.Unix() * 1000

	collectionNames := m.getAllCollectionNames()

	for _, collectionName := range collectionNames {
		// mongodump --host localhost --db ${DB_NAME} -c ${COLLECTION_NAME} -q "{time : { \$gte : 20150427, \$lt : ISODate(\"2015-04-26T00:00:00+09:00\") } }"
		query := fmt.Sprintf("\"{time : { \\$gte :  new Date(%d), \\$lt :  new Date(%d) } }\"", startMsec, endMsec)
		cmd := fmt.Sprintf("mongodump --host %s --db %s -c %s -q %s -o %s", m.mongo.URL, m.mongo.Database, collectionName, query, saveFolder)
		m.logger.Debug(tagName, "parse %s", cmd)

		if m.mongo.User != "" && m.mongo.Pass != "" {
			cmd = fmt.Sprintf("%s -u %s -p %s", cmd, m.mongo.User, m.mongo.Pass)
		}
		
		args, err := shellwords.Parse(cmd)
		if err != nil {
			m.logger.Error(tagName, err.Error())
			return
		}

		m.logger.Debug(tagName, "execute %s", cmd)
		out, err := exec.Command(args[0], args[1:]...).Output()
		m.logger.Debug(tagName, "execute end %s", string(out))

		if err != nil {
			m.logger.Error(tagName, err.Error())
		}
	}
}

func (m *Mongodb) archive(today time.Time, saveFolder string) {
	archiveFolder := path.Join(m.backup.ArchivePath, today.Format("2006"), today.Format("01"))
	if _, err := os.Stat(archiveFolder); os.IsNotExist(err) {
		os.MkdirAll(archiveFolder, 0777)
	}

	filePath := path.Join(archiveFolder, today.Format("2006-01-02")+".zip")

	cmd := fmt.Sprintf("zip -r %s %s", filePath, saveFolder)
	m.logger.Debug(tagName, "parse %s", cmd)

	args, err := shellwords.Parse(cmd)
	if err != nil {
		m.logger.Error(tagName, err.Error())
	}

	m.logger.Debug(tagName, "execute %s", cmd)
	out, err := exec.Command(args[0], args[1:]...).Output()
	m.logger.Debug(tagName, "execute end %s", string(out))
	if err != nil {
		m.logger.Error(tagName, err.Error())
	}
}

func (m *Mongodb) createExpireIndex() {
	m.logger.Debug(tagName, "create ExpireInex")

	collections := m.getAllCollectionNames()
	for _, collection := range collections {
		m.addIndex(collection)
	}
}

func (m *Mongodb) addIndex(collection string) {
	m.logger.Debug(tagName, "add index "+collection)

	session, err := mgo.Dial(m.mongo.GetDialURL())
	if err != nil {
		m.logger.Error(tagName, err.Error())
		return
	}
	defer session.Close()

	c := session.DB(m.mongo.Database).C(collection)
	index := mgo.Index{
		Key:         []string{"time"},
		Background:  true,
		ExpireAfter: time.Hour * 24 * 7,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		m.logger.Error(tagName, err.Error())
	}
}
