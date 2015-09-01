package dropbox

import (
	"../../../config/backup"
	"../../../lib/logger"
	"fmt"
	"github.com/mattn/go-shellwords"

	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"
)

var tagName = "backup.dropbox"

// Dropbox is dropbox daily backup work
type Dropbox struct {
	dropbox *backup.Dropbox
	logger  logger.Logger
}

// NewDropbox return Dropbox struct
func NewDropbox(bk *backup.Dropbox, logger logger.Logger) *Dropbox {
	if bk == nil {
		return nil
	}

	return &Dropbox{
		dropbox: bk,
		logger:  logger,
	}
}

// Execute execute backup
func (d *Dropbox) Execute() {
	now := time.Now()
	d.backupData(now)
	d.removeOldBackup()
}

func (d *Dropbox) backupData(now time.Time) {
	d.logger.Debug(tagName, "bakup start")

	saveFolder := path.Join(d.dropbox.Dst, now.Format("2006-01-02"))
	if _, err := os.Stat(saveFolder); os.IsNotExist(err) {
		os.MkdirAll(saveFolder, 0777)
	}

	// rsync -av --delete c++/ c++_backup
	cmd := fmt.Sprintf("rsync -av %s %s", d.dropbox.Src, saveFolder)

	d.logger.Debug(tagName, "parse %s", cmd)
	args, err := shellwords.Parse(cmd)
	if err != nil {
		d.logger.Error(tagName, err.Error())
		return
	}

	d.logger.Debug(tagName, "execute %s", cmd)
	out, err := exec.Command(args[0], args[1:]...).Output()
	d.logger.Debug(tagName, "execute end %s", string(out))

	if err != nil {
		d.logger.Error(tagName, err.Error())
	}
}

func (d *Dropbox) removeOldBackup() {
	fileInfos, err := ioutil.ReadDir(d.dropbox.Dst)

	if err != nil {
		d.logger.Error(tagName, fmt.Sprintf("Directory cannot read %s", err))
		return
	}

	backups := make(map[time.Time]string)

	for _, fileInfo := range fileInfos {
		folderPath := path.Join(d.dropbox.Dst, fileInfo.Name())
		flag, e := IsDirectory(folderPath)
		if e != nil {
			d.logger.Error(tagName, fmt.Sprintf("fileInfo error %s", err))
			return
		}

		if flag {
			// directory
			t, e := time.Parse("2006-01-02", fileInfo.Name())
			if e == nil {
				backups[t] = folderPath
			}
		}
	}

	for d.dropbox.BackupNum < len(backups) {
		// delete old data

		t := time.Now().AddDate(1, 0, 0)
		for key := range backups {
			if t.After(key) {
				t = key
			}
		}
		folderPath, ok := backups[t]
		if !ok {
			// no folder
			break
		}

		os.RemoveAll(folderPath)
		delete(backups, t)
	}
}

// IsDirectory return filepath is directory
func IsDirectory(filepath string) (isDir bool, err error) {
	fInfo, err := os.Stat(filepath) // FileInfo型が返る。
	if err != nil {
		return false, err // もしエラーならエラー情報を返す
	}
	// ディレクトリかどうかチェック
	return fInfo.IsDir(), nil
}
