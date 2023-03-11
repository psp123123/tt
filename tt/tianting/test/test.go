// package main

// import (
// 	"fmt"

// 	"runtime"

// )

// func main() {
// 	sysType:=runtime.GOOS
// 	if sysType=="linux" {
// 		fmt.Printf("当前系统是：%v\n",sysType)
// 	} else {
// 		fmt.Printf("当前系统是：%v\n",sysType)
// 	}
// }

// package main

// import "fmt"

//	func main() {
//		for i := 0; i < 5; i++ {
//			go func() {
//				fmt.Println(i)
//			}()
//		}
//	}
package main

import (
	"fmt"
	"time"
)

// import (
// 	"fmt"
// 	"time"
// )

// func f2(ch chan int) {
// 	for v := range ch {
// 		fmt.Printf("----:%v", v)
// 		time.Sleep(time.Second * 2)
// 		fmt.Println("end rsv")
// 	}
// }

// func main() {
// 	ch := make(chan int, 2)
// 	for i := 1; i < 100; i++ {
// 		ch <- i
// 		fmt.Println("发送i")
// 		time.Sleep(time.Second)
// 		f2(ch)
// 	}

// 	close(ch)

// }

// func main() {
// 	ch := make(chan int, 1)
// 	for i := 0; i <= 10; i++ {
// 		select {
// 		case x := <-ch:
// 			fmt.Println(x)
// 		case y := <-ch:
// 			fmt.Println(y, "this is y")
// 		case ch <- i:
// 		}
// 	}
// }

// demo1 通道误用导致的bug
// func demo1() {
// 	wg := sync.WaitGroup{}

// 	ch := make(chan int, 10)
// 	for i := 0; i < 10; i++ {
// 		ch <- i
// 	}
// 	close(ch)

// 	wg.Add(3)
// 	for j := 0; j < 3; j++ {
// 		go func() {
// 			fmt.Printf("this is xiechng:%v\n", j)
// 			for {

// 				task, ok := <-ch
// 				if !ok {
// 					fmt.Println("!ok")
// 					break
// 				}
// 				// 这里假设对接收的数据执行某些操作

// 				time.Sleep(time.Second)
// 				fmt.Printf("zhi------:%v\n", task)

// 			}
// 			wg.Done()
// 		}()
// 	}
// 	wg.Wait()
// }

// demo2 通道误用导致的bug
func demo2() {
	ch := make(chan string)
	go func() {
		// 这里假设执行一些耗时的操作
		time.Sleep(3 * time.Second)
		ch <- "job result"
	}()

	select {
	case result := <-ch:
		fmt.Println(result)
	case <-time.After(time.Second): // 较小的超时时间
		return
	}
}
func main() {
	demo2()
}
