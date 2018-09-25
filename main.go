package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"unicode/utf8"

	"github.com/tofuso/remipi/scancode"
)

var (
	textMessage = flag.String("s", "Hello World!", "キーボードに出力させる文字を指定してください。")
	dir         = flag.String("d", "/dev/hidg0", "デバイスファイル")
	talk        = flag.Bool("t", false, "指定すると対話モードで起動します。Ctr+Cで終了。")
	devf        *os.File
)

func main() {
	var err error
	flag.Parse()              //引数をパース
	devf, err = os.Open(*dir) //デバイスをオープン
	if err != nil {
		//デバイスを開く過程でエラーが発生
		fmt.Println(err)
		return
	}
	if !*talk {
		//一回だけ実行
		fmt.Println("入力された文字: ", *textMessage)
		err := run(*textMessage)
		if err != nil {
			fmt.Println(err)
			return
		}

	} else {
		//対話モード
		scanner := bufio.NewScanner(os.Stdin)
		quit := make(chan os.Signal)      //終了シグナルを作成
		signal.Notify(quit, os.Interrupt) //シグナルを設定
		//シグナルを受信するまで実行
		for scanner.Scan() {
			if s := scanner.Text(); utf8.RuneCountInString(s) > 0 {
				err := run(s)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
			select {
			case <-quit:
				return
			default:
			}
		}
	}
}

//キーボードに書き込む（開放も行われる）
func writekey(key scancode.Key) error {
	_, err := fmt.Fprintf(devf, "\\%X\\0\\%X\\0\\0\\0\\0\\0", key.Top, key.ID)
	if err != nil {
		return err
	}
	//開放
	_, err = fmt.Fprintf(devf, "\\%X\\0\\%X\\0\\0\\0\\0\\0", scancode.Open.Top, scancode.Open.ID)
	return err
}

//解析し、キーに打ち込む処理を行う
func run(s string) error {
	actf := false   //コマンド中か判定するフラグ
	var acts string //コマンドの内容を保存する
	for _, r := range s {
		if r == '|' && !actf {
			//コマンド開始
			actf = true //フラグを立てる
			acts = ""   //初期化
		} else if r == '|' && actf {
			//コマンド終了
			if actkey, ok := scancode.ActionMap[acts]; ok {
				//コマンドあり
				err := writekey(actkey)
				if err != nil {
					return err
				}
			} else {
				//該当コマンドなし
				fmt.Println("該当するコマンドがありませんでした。: ", acts)
			}
		} else if actf {
			//コマンド取得中
			acts += string(r)
		} else if key, ok := scancode.KeyMap[r]; ok {
			//通常のKeyであるなら
			err := writekey(key)
			if err != nil {
				return err
			}
		} else if skey, ok := scancode.JapaneaseKeyMap[r]; ok {
			//ひらがなのKeyであるなら日本語入力モードにする
			err := writekey(scancode.ChgIn)
			if err != nil {
				return err
			}
			//入力
			for _, k := range skey {
				err = writekey(k)
				if err != nil {
					return err
				}
			}
		} else {
			//該当する文字がない時
			fmt.Println("該当する文字がありません: ", r)
		}
	}
	return nil
}
