package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"unicode/utf8"
)

var (
	textMessage = flag.String("s", "Hello World!", "キーボードに出力させる文字を指定してください。")
	binDir      = flag.String("d", "~/hardpass-sendHID", "scanバイナリが存在するディレクトリ。")
	talk        = flag.Bool("t", false, "指定すると対話モードで起動します。Ctr+Cで終了。")
)

func main() {
	flag.Parse()
	if !*talk {
		//一回だけ実行
		fmt.Println("入力された文字: ", *textMessage)
		err := write(*textMessage)
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
				write(s)
			}
			select {
			case <-quit:
				return
			default:
			}
		}
	}
}

//キーボードに書き込む
func write(str string) error {
	l := "echo -n \"" + str + "\" | sudo " + *binDir + "/scan /dev/hidg0 1 2"
	err := exec.Command(l).Run()
	//_, err := fmt.Println(l)
	return err
}
