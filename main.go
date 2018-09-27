package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"unicode/utf8"

	"github.com/tofuso/remipi/scancode"
)

var (
	textMessage = flag.String("s", "Hello World!", "キーボードに出力させる文字を指定してください。")
	dir         = flag.String("d", "/dev/hidg0", "デバイスファイル")
	talk        = flag.Bool("t", false, "指定すると対話モードで起動します。Ctr+Cで終了。")
)

func main() {
	flag.Parse() //引数をパース
	if *talk {
		//対話モード
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			if s := scanner.Text(); utf8.RuneCountInString(s) > 0 {
				f, err := process(s)
				if err != nil {
					fmt.Println(err)
					return
				} else if f {
					//終了フラグが立った時
					return
				}
			}
		}
	} else {
		//一回だけ実行
		fmt.Println("入力された文字: ", *textMessage)
		f, err := process(*textMessage)
		if err != nil {
			fmt.Println(err)
			return
		} else if f {
			//終了フラグが立った時
			return
		}
	}
}

//キーボードに書き込む（開放も行われる）
func writekey(key scancode.Key) error {
	//ファイルにバイナリを書き込む
	var err error
	err = ioutil.WriteFile(*dir, []byte{key.Top, 0x0, key.ID, 0x0, 0x0, 0x0, 0x0, 0x0}, 0777)
	if err != nil {
		return err
	}
	//開放
	//ファイルにバイナリを書き込む
	err = ioutil.WriteFile(*dir, []byte{scancode.Open.Top, 0x0, scancode.Open.ID, 0x0, 0x0, 0x0, 0x0, 0x0}, 0777)
	return err
}

//解析し、キーに打ち込む処理を行う
func process(s string) (bool, error) {
	actf := false   //コマンド中か判定するフラグ
	var acts string //コマンドの内容を保存する
	for _, r := range s {
		if r == '_' && !actf {
			//コマンド開始
			actf = true //フラグを立てる
			acts = ""   //初期化
		} else if r == '_' && actf {
			//コマンド終了
			actf = false          //フラグを折る
			f, err := doaction(acts) //アクションを発生させる
			if err != nil{
				return false, err
			}else if f{
				//終了フラグが発生した時
				return true, nil
			}
		} else if actf {
			//コマンド取得中
			acts += string(r)
		} else if key, ok := scancode.KeyMap[r]; ok {
			//通常のKeyであるなら
			err := writekey(key)
			if err != nil {
				return false, err
			}
		} else if skey, ok := scancode.JapaneaseKeyMap[r]; ok {
			//ひらがなのKeyであるなら日本語入力モードにする
			err := writekey(scancode.ChgIn)
			if err != nil {
				return false, err
			}
			//入力
			for _, k := range skey {
				err = writekey(k)
				if err != nil {
					return false, err
				}
			}
			//戻す
			err = writekey(scancode.ChgIn)
			if err != nil {
				return false, err
			}

		} else {
			//該当する文字がない時
			fmt.Println("該当する文字がありません: ", r)
		}
	}
	return false, nil
}

// actsから対応するアクションを引き出す。
func doaction(acts string) (bool, error) {
	switch acts {
	case "quit":
		//終了する
		return true, nil
	case "sec":
		//一秒のウェイトを入れる
		time.Sleep(time.Second)
	/*使用不可 UsageIDが不明
	case "_":
		//__と入力したので_を返す
		err := writekey(scancode.KeyMap['_'])
		if err != nil{
			return false, err
		}
	*/
	default:
		//それら以外
		if actkey, ok := scancode.ActionMap[acts]; ok {
			//コマンドあり
			err := writekey(actkey)
			if err != nil {
				return false, err
			}
		} else {
			//該当コマンドなし
			fmt.Println("該当するコマンドがありませんでした。: ", acts)
		}
	}
	return false, nil
}
