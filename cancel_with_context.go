package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var wgCtx sync.WaitGroup

// キャンセルされるまで num をひたすら送信し続けるチャネルを生成
func numGeneratorWithContext(ctx context.Context, num int) <-chan int {
	gen := make(chan int)

	// キャンセルされるまで num をひたすら送信し続ける
	go func() {
		defer wgCtx.Done()

	LOOP:
		for {
			// chan による情報伝達：ループごとに select で ctx.Done() チャネルを監視
			select {
			case <-ctx.Done(): // ctx.Done() チャネルが close されたら break が実行され、LOOP ラベルのついた for ループを抜ける
				fmt.Println("[numGeneratorWithContext] done")
				break LOOP
			case gen <- num: // キャンセルされてなければ num を gen チャネルに送信
				fmt.Println("[numGeneratorWithContext] send num:", num)
			}

			// increment num
			num++

			// wait 0.5 sec
			time.Sleep(500 * time.Millisecond)
		}

		close(gen)
		fmt.Println("[numGeneratorWithContext] closed channel `gen`")
	}()

	return gen
}

func CancelGoroutineWithContext() {
	fmt.Println("----- start CancelGoroutineWithContext -----")

	ctx, cancel := context.WithCancel(context.Background())
	gen := numGeneratorWithContext(ctx, 96)

	wgCtx.Add(1)

	// gen から 5 回数値を受信する
	for i := 0; i < 5; i++ {
		num := <-gen
		fmt.Println("[main] receive a number from channel `gen`:", num)
	}

	// 5 回 gen から cancel する。cancel は ctx.Done() によって監視できる。
	cancel()
	fmt.Println("[main] canceled context")

	wgCtx.Wait()
}
