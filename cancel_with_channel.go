package main

import (
	"fmt"
	"sync"
	"time"
)

var wgChan sync.WaitGroup

// キャンセルされるまで num をひたすら送信し続けるチャネルを生成
func numGeneratorWithChannel(done chan struct{}, num int) <-chan int {
	gen := make(chan int)

	// キャンセルされるまで num をひたすら送信し続ける
	go func() {
		defer wgChan.Done()

	LOOP:
		for {
			// chan による情報伝達：ループごとに select で done チャネルを監視
			select {
			case <-done: // done チャネルが close されたら break が実行され、LOOP ラベルのついた for ループを抜ける
				fmt.Println("[numGeneratorWithChannel] done")
				break LOOP
			case gen <- num: // キャンセルされてなければ num を gen チャネルに送信
				fmt.Println("[numGeneratorWithChannel] send num:", num)
			}

			// increment num
			num++

			// wait 0.5 sec
			time.Sleep(500 * time.Millisecond)
		}

		close(gen)
		fmt.Println("[numGeneratorWithChannel] closed channel `gen`")
	}()

	return gen
}

func CancelGoroutineWithChannel() {
	fmt.Println("----- start CancelGoroutineWithChannel -----")

	done := make(chan struct{})
	gen := numGeneratorWithChannel(done, 1)

	wgChan.Add(1)

	// gen から 5 回数値を受信する
	for i := 0; i < 5; i++ {
		num := <-gen
		fmt.Println("[main] receive a number from channel `gen`:", num)
	}

	// 5 回 gen から受信したら done チャネルを close する
	close(done)
	fmt.Println("[main] closed channel `done`")

	wgChan.Wait()
}
