# remipi
---
remipiはremote input for raspberry piの略…だと思います。  
RaspberryPiを遠隔入力キーボードとして使う際に入力を支援するためのプログラムです。  

## 使い方
---

### インストール方法

このページを参照してください。

### 引数

`[Text]`に打ち込みたい文字列を打ち込む。

```
remipi -s [Text]
```

USBキーボードのバイナリを送信するデバイスファイルを指定する。  
（デフォルトは`/dev/hidg0`）

```
remipi -d /dev/usb
```

対話モードで打ち込む。

```
remipi -t
```

### 機能

打ち込む際に`!`で特定のキーワードを囲むことで、特殊な操作ができます。  
例：

```
!win-r!			//windowsキー+Rキー  
!win-r!notepad!enter!	//メモ帳の起動
```

#### 特殊操作一覧

```
open		//キーをすべて開放します
right		//十字キー右
left		//十字キー左
up		//十字キー上
down		//十字キー下
enter		//enterキー
esc		//escapeキー
back		//バックスペースキー
capslock	//capslockキー
insert		//insertキー
delete		//deleteキー
home		//homeキー
end		//endキー
pageup		//pageupキー
pagedown	//pagedownキー
numlock		//numlockキー
printscreen	//printscreenキー
zen-or-han	//全角半角キー
win-r		//windowsキーとrキーを同時押し

quit		//対話モードを終了させる
sec		//一秒間待機
```

## 注意事項
---

現在利用可能な入力  
アルファベット（小文字）  
アルファベット（大文字）  
数字
キーボードから予測変換無しで入力可能な記号や空白  
（ただし、`\`と`|`と`_`を除く）  
ひらがな  
前述の記号の全角入力のもの