# understanding-go-context

Zenn の書籍 [(2021, Zenn) よくわかる context の使い方](https://zenn.dev/hsaki/books/golang-context) をやるリポジトリ

## goroutine

> A "go" statement starts the execution of a function call as an independent concurrent thread of control, or goroutine, within the same address space.

訳

> "go" 文は、同じアドレス空間内の独立した並行スレッド（goroutine と呼ぶ）として、関数呼び出しの実行を開始する。

ref: https://go.dev/ref/spec#Go_statements

> **warning**
> 
> goroutine で並行に実行しても、並列に実行されるとは限らないことに注意。

## context の役割

`Context` 型の主な役割は 3 つ：

- **処理の締切を伝達**
- **キャンセル信号の伝播**
- **リクエストスコープ値の伝達**

## context の意義

context が役立つのは、1 つの処理が複数の goroutine を跨いで行われる場合。

- 例
  - HTTP リクエストを受け付けて、DB に保存する処理では、まず main goroutine がリクエストを受け付け、リクエストを処理するために goroutine を起動する。リクエストハンドラなかで DB に接続してデータを取得する処理のために別の goroutine を起動する。
  - > DBからのデータ取得のために複数個のゴールーチンを立てるというのは、例えば「複数個あるDBレプリカ全てにリクエストを送り、一番早くに結果が返ってきたものを採用する」といったときなどが考えられます。
    > Go公式ブログの ["Go Concurrency Patterns: Timing out, moving on"](https://zenn.dev/hsaki/books/golang-context/viewer/definition#context%E3%81%AE%E6%84%8F%E7%BE%A9:~:text=%22Go%20Concurrency%20Patterns%3A%20Timing%20out%2C%20moving%20on%22) にも、そのようなパターンについて言及されてます。

このように、

- Go のプログラマがそのことについて意識していなくても、ライブラリの仕様上複数の goroutine 上に処理が跨る
- 一つの処理を行うために、いくつもの goroutine が木構造的に積み上がっていく（下図参照）

というのが珍しくない。

![](https://storage.googleapis.com/zenn-user-upload/1f88984ea5aba496969a7ed1.png)

引用元: https://zenn.dev/hsaki/books/golang-context/viewer/definition

処理が複数個の goroutine にまたがると、「**情報伝達全般」が難しくなる。

基本的に、Go では「**異なる goroutine 間での情報共有は、ロックを使ってメモリを共有するよりも、チャネルを使った伝達を使うべし**」という考え方を取っている。 並行に動いている複数の goroutine 上から、メモリ上に存在する 1 つのデータにそれぞれが「安全に」アクセスできることを担保するのはとても難しいためである。

「**複数の goroutine 間で安全に、そして簡単に情報伝達を行いたい**」という要望は、チャネルによる伝達だけで実現しようとすると難しい。

context では、goroutine 間での情報伝達のうち、特に需要が多い、

- **処理の締め切りを伝達**
- **キャンセル信号の伝播**
- **リクエストスコープ値の伝達**

の 3 つについて、「goroutine 上で起動される関数の第一引数に、`context.Context`型を 1 つ渡す」だけで簡単に実現できるようになっている。

## context の定義

- `context.Context` 型

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key interface{}) interface{}
}
```

## References

- [(2021/06, Zenn) Go での並行処理を徹底解剖！](https://zenn.dev/hsaki/books/golang-concurrency)
- [(2021, Zenn) よくわかる context の使い方](https://zenn.dev/hsaki/books/golang-context)
