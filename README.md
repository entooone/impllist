# impllist

https://github.com/nu50218/impls に開発を移行したため Archive する。

ユーザ定義型から実装している interface の一覧を取得する。また、ユーザ定義の interface からその interface を実装している型の一覧を取得する

## TODO

- [ ] ユーザ定義型からその型が実装している interface の一覧を取得する
- [ ] ユーザ定義の interface からその interface を実装している型の一覧を取得する
- [ ] 探索するパッケージや対象をフィルタできるようにする
    - パッケージの指定
    - 公開された interface のみを対象とする
    - 調べるユーザ定義型が定義されたパッケージを含める
    - 調べるユーザ定義型が定義されたパッケージが import しているものを含める
- [ ] カーソルの位置を取得してその下にある型について一覧を取得する

## Install

```
$ go get -u github.com/entooone/impllist
```

## How to use

### `os.File` が実装している interface の一覧を取得する

```
$ impllist os.File
```
