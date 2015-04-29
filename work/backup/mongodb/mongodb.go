package mongodb

import (
	"../../../config"
	"fmt"
	"github.com/mattn/go-shellwords"
	"gopkg.in/mgo.v2"
	"os"
	"os/exec"
	"path"
	"time"
)

// Mongodb is mongodb daily backup work
type Mongodb struct {
	backup *config.MongodbBackup
	mongo  *config.MongodbDatabase
}

// NewMongodb return Mongodb struct
func NewMongodb(c *config.MongodbBackup, mongo *config.MongodbDatabase) *Mongodb {
	return &Mongodb{
		backup: c,
		mongo:  mongo,
	}
}

// Execute execute backup and convert zip file
func (m *Mongodb) Execute() {
	m.backupData()
	m.createExpireIndex()
}

func (m *Mongodb) backupData() {
	// create date
	start := time.Now().AddDate(0, 0, -3)
	start = start.Truncate(time.Hour * 24).Add(time.Hour * -9)
	end := time.Now()

	// reset workspace
	err := os.RemoveAll(m.backup.Workspace)
	if err != nil {
		panic(err)
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
		return names
	}
	defer session.Close()

	d := session.DB(m.mongo.Database)
	rawNames, err := d.CollectionNames()
	if err != nil {
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
	// mongodump --host localhost --db ${DB_NAME} --collection ${COLLECTION_NAME} -q "{time : { \$gte : 20150427, \$lt : ISODate(\"2015-04-26T00:00:00+09:00\") } }"
	query := fmt.Sprintf("\"{time : { \\$gte :  new Date(%d), \\$lt :  new Date(%d) } }\"", start.Unix(), end.Unix())
	cmd := fmt.Sprintf("mongodump --host %s --db %s -q %s -o %s", m.mongo.GetDialURL(), m.mongo.Database, query, saveFolder)

	fmt.Println(cmd)
	args, err := shellwords.Parse(cmd)
	if err != nil {
		panic(err)
	}

	out, err := exec.Command(args[0], args[1:]...).Output()
	fmt.Println(string(out))
	if err != nil {
		fmt.Println(err)
		fmt.Println("error")
		panic(err)
	}
}

func (m *Mongodb) archive(today time.Time, saveFolder string) {
	archiveFolder := path.Join(m.backup.ArchivePath, today.Format("2006"), today.Format("01"))
	if _, err := os.Stat(archiveFolder); os.IsNotExist(err) {
		os.MkdirAll(archiveFolder, 0777)
	}

	filePath := path.Join(archiveFolder, today.Format("2006-01-02")+".zip")

	cmd := fmt.Sprintf("zip -r %s %s", filePath, saveFolder)
	fmt.Println(cmd)
	args, err := shellwords.Parse(cmd)
	if err != nil {
		panic(err)
	}

	out, err := exec.Command(args[0], args[1:]...).Output()
	fmt.Println(string(out))
	if err != nil {
		fmt.Println(err)
		fmt.Println("error")
		panic(err)
	}
}

func (m *Mongodb) createExpireIndex() {
	collections := m.getAllCollectionNames()
	for _, collection := range collections {
		m.addIndex(collection)
	}
}

func (m *Mongodb) addIndex(collection string) {
	session, err := mgo.Dial(m.mongo.GetDialURL())
	if err != nil {
		panic(err)
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
		panic(err)
	}
}
