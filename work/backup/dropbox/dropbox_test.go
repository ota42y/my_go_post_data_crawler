package dropbox

import (
	"fmt"
	"os"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"../../../config/backup"
	"../../../test"
	"../../../util"
)

func removeTestData() {
	os.RemoveAll("./workDir/")
}

func createTestData() {
	removeTestData()
	os.MkdirAll("./workDir/test/", 0777)
	os.MkdirAll("./workDir/backup/", 0777)

	createTestFile("./workDir/test/test.txt")
}

func createBackupData() {
	removeTestData()
	os.MkdirAll("./workDir/backup/", 0777)
	os.MkdirAll("./workDir/backup/2001-01-01/", 0777)
	os.MkdirAll("./workDir/backup/2001-01-02/", 0777)
	os.MkdirAll("./workDir/backup/2001-01-03/", 0777)

	createTestFile("./workDir/backup/2001-01-01/test.txt")
	createTestFile("./workDir/backup/2001-01-02/test.txt")
	createTestFile("./workDir/backup/2001-01-03/test.txt")
}

func countBacukupDataNum() int {
	count := 0
	if checkTestData("2001-01-01") {
		count++
	}
	if checkTestData("2001-01-02") {
		count++
	}
	if checkTestData("2001-01-03") {
		count++
	}
	return count
}

func createTestFile(filepath string) {
	f, err := os.Create(filepath)
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()
	f.WriteString("test file")
}

func isTestFileExist(filepath string) bool {
	isDir, err := IsDirectory(filepath)
	if isDir || err != nil {
		return false
	}

	text := string(util.LoadFile(filepath))
	return text == "test file"
}

func checkTestData(folderName string) bool {
	return isTestFileExist(fmt.Sprintf("./workDir/backup/%s/test.txt", folderName))
}

func TestBackup(t *testing.T) {
	configStr := `Src = "./workDir/test/"
Dst = "./workDir/backup"
BackupNum = 3
`
	c := backup.NewDropbox(configStr)

	testLogger := &test.Logger{}
	d := NewDropbox(c, testLogger)

	Convey("backup data", t, func() {
		So(d, ShouldNotBeNil)

		Convey("collect", func() {
			createTestData()
			now := time.Now()
			d.backupData(now)
			So(checkTestData(now.Format("2006-01-02")), ShouldBeTrue)
			removeTestData()
		})

		Convey("auto backup delete", func() {
			Convey("no delete", func() {
				Convey("same num", func() {
					createBackupData()
					d.removeOldBackup()
					So(countBacukupDataNum(), ShouldEqual, 3)
					removeTestData()
				})

				Convey("big num", func() {
					d.dropbox.BackupNum = 4
					createBackupData()
					d.removeOldBackup()
					So(countBacukupDataNum(), ShouldEqual, 3)
					removeTestData()
				})
			})
			Convey("delete", func() {
				Convey("delete collect", func() {
					d.dropbox.BackupNum = 2
					createBackupData()
					d.removeOldBackup()
					So(countBacukupDataNum(), ShouldEqual, 2)
					So(checkTestData("2001-01-01"), ShouldBeFalse)
					removeTestData()
				})

				Convey("can't delete", func() {
					d.dropbox.BackupNum = -1
					createBackupData()
					d.removeOldBackup()
					So(countBacukupDataNum(), ShouldEqual, 0)
					So(checkTestData("2001-01-01"), ShouldBeFalse)
					removeTestData()
				})
			})
		})
	})
}
