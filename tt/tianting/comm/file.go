package comm

import (
	"fmt"
	"os"
	"path"
	"time"
)

type FileLogger struct {
	Level       LogLevel
	filePath    string
	fileName    string
	fileObj     *os.File
	errFileObj  *os.File
	maxFileSize int64
}

func NewFileLogger(levelStr, fp, fn string, maxSize int64) (logger *FileLogger) {
	LogLevel, err := parseLogLevel(levelStr)
	if err != nil {
		fmt.Println(err)
	}
	fl := &FileLogger{
		Level:       LogLevel,
		filePath:    fp,
		fileName:    fn,
		maxFileSize: maxSize,
	}
	err = fl.initFile()
	if err != nil {
		fmt.Println(err)
	}

	return
}

func (f *FileLogger) initFile() error {
	fullFileName := path.Join(f.filePath, f.fileName)
	fileObj, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file failed,err:%v\n", err)
		return err
	}
	errFileObj, err := os.OpenFile(fullFileName+".err", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open err log file failed,err:%v\n", err)
		return err
	}
	f.fileObj = fileObj
	f.errFileObj = errFileObj
	return nil
}
func (f *FileLogger) enable(logLevel LogLevel) bool {
	return f.Level <= logLevel
}

func (f *FileLogger) checkSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed,err:%v\n", err)
		return false
	}
	return fileInfo.Size() >= f.maxFileSize
}
func (f *FileLogger) splitFile(file *os.File) (*os.File, error) {
	//需要切割文件,拿到日志文件状态
	nowStr := time.Now().Format("20060102150405000")
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed,err:%v\n", err)
		return nil, err
	}
	logName := path.Join(f.filePath, fileInfo.Name())         //拿到当前的日志文件完整路径
	NewLogName := fmt.Sprintf("%s/%s.bak%s", logName, nowStr) //拼接一个备份的名字
	//1. 关闭当前日志文件
	file.Close()
	//2. rename文件，或备份

	os.Rename(logName, NewLogName)
	//3. 打开一个新的日志文件
	fileObj, err := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open new file failed:%v\n", err)
		return nil, err
	}
	//4. 将打开新的日志对象赋值给新的f.fileObj
	return fileObj, nil

}
func (f *FileLogger) log(lv LogLevel, format string, a ...interface{}) {
	if f.enable(lv) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now()
		funcName, fileName, lineNo := getInfo(3)
		if f.checkSize(f.fileObj) {
			newFile, err := f.splitFile(f.fileObj)
			if err != nil {
				return
			}
			f.fileObj = newFile
		}
		fmt.Fprintf(f.fileObj, "[%s] [%s] [%s:%s:%d]\n", now.Format("2006-01-02 15:04:05"), getLogString(lv), fileName, funcName, lineNo, msg)
		if lv >= ERROR {
			if f.checkSize(f.errFileObj) {
				newFile, err := f.splitFile(f.errFileObj)
				if err != nil {
					return
				}
				f.errFileObj = newFile
			}
			fmt.Fprintf(f.errFileObj, "[%s] [%s] [%s:%s:%d]\n", now.Format("2006-01-02 15:04:05"), getLogString(lv), fileName, funcName, lineNo, msg)

		}
	}
}

func (f *FileLogger) Debug(format string, a ...interface{}) {
	f.log(DEBUG, format, a...)
}

func (f *FileLogger) Info(format string, a ...interface{}) {
	f.log(INFO, format, a...)
}

func (f *FileLogger) Warning(format string, a ...interface{}) {
	f.log(WARNING, format, a...)
}

func (f *FileLogger) Error(format string, a ...interface{}) {
	f.log(ERROR, format, a...)
}

func (f *FileLogger) Fatal(format string, a ...interface{}) {
	f.log(FATAL, format, a...)
}

func (f *FileLogger) Close() {
	f.fileObj.Close()
	f.errFileObj.Close()
}
