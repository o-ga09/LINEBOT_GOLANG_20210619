package util

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
	"os"
)

func Tohash(originalstr string) string {
	salt := os.Getenv("SALT")
	hashstr := []byte(originalstr + salt)

	hash_sha256 := sha256.Sum256(hashstr)

	return hex.EncodeToString(hash_sha256[:])
}

func isMatch(hash,userid string) bool{
	salt := os.Getenv("SALT")
	hashstr := []byte(userid + salt)
	hash_sha256 := sha256.Sum256(hashstr)

	if hash == hex.EncodeToString(hash_sha256[:]) {
		return true
	}
	return false
}

func LoggingSettings(logFile string) {
    //_=error
    //os.O_RDWR　READ　WRITE 読み書き両方する時
    //os.O_CREATE 存在しなかった場合新規ファイルを作成する場合
    //os.O_APPEND  ファイルに追記したいとき
    //0666 
    //引数: ファイルのパス, フラグ, パーミッション(わからなければ0666でおっけーです)
    //上記モード指定。読み込む、作成、権限(０６６６＝読み書き）を設定。
    logfile, _ := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    //stdout 画面上に出る出力 をlogfileに書き込む
    multiLogFile := io.MultiWriter(os.Stdout, logfile)
    //フォーマット指定
    //日付、時間、短いエラーの名前
    log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
    //ログファイルの出力先を変更   
    log.SetOutput(multiLogFile)
}