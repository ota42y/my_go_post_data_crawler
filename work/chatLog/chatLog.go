package chatLog

import (
  "encoding/json"
  "bufio"
  "fmt"
  "os"
  "regexp"
  "time"
  "log"
  "io/ioutil"
  "path"
)

type ChatLog struct{
  Nick string `json:nick`
  Date string `json:date`
  Message string `json:message`
  Channel string `json:channel`
  Type string `json type`
}

func (chat ChatLog) ToString() string {
  t, e := time.Parse("2006-01-02T15:04:05.999-0700", chat.Date)
  if e != nil {
    log.Fatal(e)
  }
  return t.Format("15:04") + " " + chat.Message
}

type Worker struct{
  logFolder string
  saveFolder string
}

func NewWorkerFromMap(config map[interface{}]interface{}) (worker *Worker) {
	chatLogConfig := config["chatLog"].(map[interface{}]interface{})
	logFolder := chatLogConfig["logFolder"].(string)
	saveFolder := chatLogConfig["saveFolder"].(string)
  rootDir := config["rootDir"].(string)

	return &Worker{
		logFolder: path.Join(rootDir, logFolder),
		saveFolder: path.Join(rootDir, saveFolder),
	}
}

// 今日のログを出力する
// 昨日のログファイルがあれば最新版にする(日付またぎ対策)

// 指定されたファイル名がディレクトリかどうか調べる
func IsDirectory(name string) (isDir bool,err error) {
	fInfo,err := os.Stat(name) // FileInfo型が返る。
	if err != nil {
		return false,err // もしエラーならエラー情報を返す
	}
	// ディレクトリかどうかチェック
	return fInfo.IsDir(),nil
}

func (worker *Worker) Work() {
  fmt.Printf("work chatLog\n")

  fileInfos,err := ioutil.ReadDir(worker.logFolder)

  if err != nil {
    fmt.Errorf("Directory cannot read %s\n",err)
    return
  }

  for _,fileInfo := range fileInfos {
    // *FileInfo型
    folderPath := path.Join(worker.logFolder, fileInfo.Name())
    flag, e := IsDirectory(folderPath)
    if e != nil {
      panic(e)
    }

    if flag {
      worker.saveRoomLog(worker.logFolder, fileInfo.Name())
    }
  }
}

func (worker *Worker) saveRoomLog(logFolder string, roomName string) {
  roomFolder := path.Join(logFolder, roomName)
  fileInfos,err := ioutil.ReadDir(roomFolder)

  if err != nil {
    fmt.Errorf("Directory cannot read %s\n",err)
    return
  }

  for _,fileInfo := range fileInfos {
    // *FileInfo型
    filePath := path.Join(roomFolder, fileInfo.Name())
    flag, e := IsDirectory(filePath)
    if e != nil {
      fmt.Printf("%v", e)
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
      }else {
        // .DS_Store
      }
    }
  }
}

func (worker *Worker) saveTodayLog(logDir string, roomName string, fileName string) {
  fmt.Printf("save chatLog %s\n", fileName)

  logs := getFilteredLog(path.Join(logDir, roomName, fileName))

  saveFolder := path.Join(worker.saveFolder, roomName)
  if _, err := os.Stat(saveFolder); os.IsNotExist(err) {
    os.MkdirAll(saveFolder, 0777)
  }

  worker.saveLogToFile(path.Join(saveFolder, fileName), logs)
}

func (worker *Worker) saveYesterdayLog(logDir string, roomName string, fileName string) {
  fmt.Printf("save chatLog %s\n", fileName)
  logs := getFilteredLog(path.Join(logDir, roomName, fileName))

  saveFolder := path.Join(worker.saveFolder, roomName)
  if _, err := os.Stat(saveFolder); os.IsNotExist(err) {
    os.MkdirAll(saveFolder, 0777)
  }

  // ファイルが存在する場合のみ、最新版にアップデートする
  saveFilePath := path.Join(saveFolder, fileName)
  if _, err := os.Stat(saveFilePath); !os.IsNotExist(err) {
    worker.saveLogToFile(saveFilePath, logs)
  }
}

func (worker *Worker) saveLogToFile(filePath string, logs []ChatLog) {
  f, err := os.Create(filePath)
  if err != nil {
    fmt.Printf("error %v\n", err)
    return
  }

  defer f.Close()

  for _, log := range logs {
    if _, err = f.WriteString(log.ToString() + "\n"); err != nil {
      fmt.Printf("%v\n", err)
    }
  }
}

func getFilteredLog(filepath string) []ChatLog {
  log_data := loadFile(filepath)

  logs := []ChatLog{}
  for _, log := range log_data {
    if filter(log) {
      logs = append(logs, log)
    }
  }
  return logs
}

func filter(log ChatLog) bool {
  if regexpFilter("!(ota42y)", log.Nick){
    return false
  }
  if regexpFilter("^@.*", log.Message){
    return false
  }
  if regexpFilter("^$", log.Message){
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

func loadFile(filepath string) []ChatLog {
  f, err := os.Open(filepath)
  if err != nil {
    fmt.Fprintf(os.Stderr, "File %s could not read: %v\n", filepath, err)
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
    fmt.Fprintf(os.Stderr, "File %s scan error: %v\n", filepath, err)
    return nil
  }

  return logs
}
