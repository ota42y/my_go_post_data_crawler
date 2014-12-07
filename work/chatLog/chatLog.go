package chatLog

import (
	"./../../lib/logger"
	"./../../lib/evernote"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"time"
	"gopkg.in/yaml.v2"
)

type ChatLog struct {
	Nick    string `json:nick`
	Date    string `json:date`
	Message string `json:message`
	Channel string `json:channel`
	Type    string `json type`
}

func (chat ChatLog) ToString() string {
	t, e := time.Parse("2006-01-02T15:04:05.999-0700", chat.Date)
	if e != nil {
		return "error " + chat.Message
	}
	return t.Format("15:04") + " " + chat.Message
}

type Worker struct {
	logFolder  string
	saveFolder string
	logger     *logger.MyLogger
	evernote *evernote.Sender
}

type Setting struct{
	LogFolder string
	SaveFolder string
}

func New(buf []byte, logger *logger.MyLogger, evernote *evernote.Sender) (worker *Worker) {
	s := Setting{}
	err := yaml.Unmarshal(buf, &s)
	if err != nil {
		return nil
	}

	return &Worker{
		logFolder:  s.LogFolder,
		saveFolder: s.SaveFolder,
		logger:     logger,
		evernote: evernote,
	}
}

// 今日のログを出力する
// 昨日のログファイルがあれば最新版にする(日付またぎ対策)

// 指定されたファイル名がディレクトリかどうか調べる
func IsDirectory(name string) (isDir bool, err error) {
	fInfo, err := os.Stat(name) // FileInfo型が返る。
	if err != nil {
		return false, err // もしエラーならエラー情報を返す
	}
	// ディレクトリかどうかチェック
	return fInfo.IsDir(), nil
}

func (worker *Worker) Work() {
	worker.logger.LogPrint("chat_log", "work")

	fileInfos, err := ioutil.ReadDir(worker.logFolder)
	if err != nil {
		worker.logger.ErrorPrint("chat_log", fmt.Sprintf("Directory cannot read %s", err))
		return
	}

	for _, fileInfo := range fileInfos {
		// *FileInfo型
		folderPath := path.Join(worker.logFolder, fileInfo.Name())
		flag, e := IsDirectory(folderPath)
		if e != nil {
			worker.logger.ErrorPrint("chat_log", fmt.Sprintf("%s", err))
			return
		}

		if flag {
			worker.saveRoomLog(worker.logFolder, fileInfo.Name())
		}
	}
}

func (worker *Worker) saveRoomLog(logFolder string, roomName string) {
	roomFolder := path.Join(logFolder, roomName)
	fileInfos, err := ioutil.ReadDir(roomFolder)

	if err != nil {
		worker.logger.ErrorPrint("chat_log", fmt.Sprintf("Directory cannot read %s", err))
		return
	}

	for _, fileInfo := range fileInfos {
		// *FileInfo型
		filePath := path.Join(roomFolder, fileInfo.Name())
		flag, e := IsDirectory(filePath)
		if e != nil {
			worker.logger.ErrorPrint("chat_log", fmt.Sprintf("%s", err))
			return
		}

		if !flag {
			// ログファイルなので保存する
			t, e := time.Parse("2006-01-02.txt", fileInfo.Name())
			if e == nil {
				today := time.Now()
				if today.Year() == t.Year() && today.Month() == t.Month() && today.Day() == t.Day() {
					worker.saveTodayLog(logFolder, roomName, fileInfo.Name())
				}

				yesterday := today.AddDate(0, 0, -1)
				if yesterday.Year() == t.Year() && yesterday.Month() == t.Month() && yesterday.Day() == t.Day() {
					worker.saveYesterdayLog(logFolder, roomName, fileInfo.Name())
				}
			} else {
				// .DS_Store
			}
		}
	}
}

func (worker *Worker) saveTodayLog(logDir string, roomName string, fileName string) {
	logs := worker.getFilteredLog(path.Join(logDir, roomName, fileName))

	saveFolder := path.Join(worker.saveFolder, roomName)
	if _, err := os.Stat(saveFolder); os.IsNotExist(err) {
		os.MkdirAll(saveFolder, 0777)
	}

	worker.logger.LogPrint("chat_log", fmt.Sprintf("save chat log %s", fileName))
	worker.saveLogToFile(path.Join(saveFolder, fileName), logs)
}

func (worker *Worker) saveYesterdayLog(logDir string, roomName string, fileName string) {
	logs := worker.getFilteredLog(path.Join(logDir, roomName, fileName))

	saveFolder := path.Join(worker.saveFolder, roomName)
	if _, err := os.Stat(saveFolder); os.IsNotExist(err) {
		os.MkdirAll(saveFolder, 0777)
	}

	// ファイルが存在する場合、削除してEvernoteに送信する(何度も送らないためファイルがあるときのみ)
	saveFilePath := path.Join(saveFolder, fileName)
	if _, err := os.Stat(saveFilePath); !os.IsNotExist(err) {
		worker.logger.LogPrint("chat_log", fmt.Sprintf("send chat log %s to evernote", fileName))

		// 10kb
		byteStr := make([]byte, 0, 1024* 10)
		for _, v := range logs {
			byteStr = append(byteStr, v.ToString()...)
			byteStr = append(byteStr, '\n')
		}

		worker.evernote.SendNote(fileName, string(byteStr))

		// 送信後にファイルは消す(次回以降送らないため)
		os.Remove(saveFilePath)
	}
}

func (worker *Worker) saveLogToFile(filePath string, logs []ChatLog) {
	f, err := os.Create(filePath)
	if err != nil {
		worker.logger.ErrorPrint("chat_log", fmt.Sprintf("error %s", err))
		return
	}

	defer f.Close()

	for _, log := range logs {
		if _, err = f.WriteString(log.ToString() + "\n"); err != nil {
			worker.logger.ErrorPrint("chat_log",fmt.Sprintf("error %s", err))
		}
	}
}

func (worker *Worker) getFilteredLog(filepath string) []ChatLog {
	log_data := worker.loadFile(filepath)

	logs := []ChatLog{}
	for _, log := range log_data {
		if filter(log) {
			logs = append(logs, log)
		}
	}
	return logs
}

func filter(log ChatLog) bool {
	if regexpFilter("!(ota42y)", log.Nick) {
		return false
	}
	if regexpFilter("^@.*", log.Message) {
		return false
	}
	if regexpFilter("^$", log.Message) {
		return false
	}

	return true
}

func regexpFilter(reg string, text string) bool {
	if m, _ := regexp.MatchString(reg, text); !m {
		return false
	}
	return true
}

func (worker *Worker) loadFile(filepath string) []ChatLog {
	f, err := os.Open(filepath)
	if err != nil {
		worker.logger.ErrorPrint("chat_log", fmt.Sprintf("File %s could not read: %v\n", filepath, err))
		return nil
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)

	logs := []ChatLog{}
	for scanner.Scan() {
		// 1行読み込む
		line := scanner.Text()

		// ログ構造体にデータを入れる
		var chat_log ChatLog
		if err := json.Unmarshal([]byte(line), &chat_log); err != nil {
			panic(err)
		}

		logs = append(logs, chat_log)
	}

	if serr := scanner.Err(); serr != nil {
		worker.logger.ErrorPrint("chat_log", fmt.Sprintf("File %s scan error: %v\n", filepath, err))
		return nil
	}

	return logs
}
