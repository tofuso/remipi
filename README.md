# remipi
---
remipiはremote input for raspberry piの略…だと思います。  
RaspberryPiを遠隔入力キーボードとして使う際に入力を支援するためのプログラムです。  

## Usage
---

### 引数

`Word`という文字列を打ち込む

```
remipi -s Word
```

[hardpass-sendHID](https://github.com/girst/hardpass-sendHID)のバイナリが存在するディレクトリを指定  
（デフォルトは~/hardpress-sendHID）

```
remipi -s Word -d /usr/bin
```

対話モードで打ち込む（`Ctr+C`で終了）

```
remipi -t
```
