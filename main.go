package main

import "C"
import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"go_demo/mit"
	"go_demo/picture"
	"sync"
	"time"
	"unsafe"
)

//切片 map是引用传递  结构体是复制

func array_point() {
	array := [3]byte{228, 184, 173} // UTF-8 编码的 "中"
	// 使用 unsafe.Pointer 和 unsafe.Slice 将 array 转换为字符串
	arraySlice := array[:]
	str := (*string)(unsafe.Pointer(&arraySlice))

	fmt.Println(*str)
	fmt.Println(*str)

	s := "中"

	bytes_ := *(*[]byte)(unsafe.Pointer(&s))

	fmt.Println(bytes_)
	fmt.Println(bytes_)
}

func main() {
	//defer Defer()()
	//defer DeferTime(time.Now())
	//fmt.Println("main run")
	//time.Sleep(3 * time.Second)
	//slice2Array()
	//err := errorJoin()
	//if err != nil {
	//	fmt.Printf("%v", err)
	//	return
	//}
	//randT()

	//

	//syncPoolT()

	//time.Sleep(1 * time.Minute)

	picture.Run()

}

//export Run
func Run() {
	mit.Run()
}

func syncPoolT() {

	pool1 := sync.Pool{
		New: func() interface{} {
			return make([]int, 0, 10)
		},
	}
	pool2 := sync.Pool{
		New: func() interface{} {
			return make([]int, 0, 10)
		},
	}

	go func() {
		for {
			array := pool1.Get().([]int)
			array = append(array, 1)
			//fmt.Printf("array1:%v\n", array)
			pool1.Put(array)
			time.Sleep(1 * time.Second)

		}

	}()
	go func() {
		for {
			array := pool2.Get().([]int)
			array = append(array, 2)
			fmt.Printf("array2:%v\n", array)
			pool2.Put(array)
			time.Sleep(1 * time.Second)
		}
	}()

}

func chT() {
	items := []int{1, 2, 3, 4, 5}
	changeSlice(items)
	fmt.Println(items)

	var ch = make(chan struct{})

	go func() {
		i := 0
		for {
			select {
			case _ = <-ch:
				fmt.Println("hello")

			default:
				//fmt.Println("default")
				i++
			}
			if i > 3 {
				close(ch)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("error:%v\n", r)
			}
		}()
		panic("我panic了")
	}()
	time.Sleep(1 * time.Minute)
}
func contextT() {
	ctx, _ := context.WithCancelCause(context.Background())

	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Done")
			fmt.Printf("%v\n", ctx.Err())
			tiem, ok := ctx.Deadline()
			fmt.Printf("%v--%v\n", tiem, ok)

		}

	}()

	time.Sleep(1 * time.Second)
	fmt.Println("Sleep Done")
	///这为什么没有堵塞
	<-ctx.Done()
	fmt.Println("Done Over")
	//cancel(nil)
	time.Sleep(1 * time.Second)

}

func sliceT() {

	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	news := make([]int, 0, len(numbers))
	for _, number := range numbers {

		if number%2 == 0 {
			news = append(news, number)
		}

	}
	fmt.Println(numbers)
	fmt.Println(news)

}

func changeSlice(items []int) {
	items[0] = 10
}

func randT() {
	buf := make([]byte, 8)
	_, err := rand.Read(buf)
	if err != nil {
		return
	}

	fmt.Printf("%x", buf)

}

func errorJoin() error {

	err := errors.New("My error")

	return errors.Join(err, errors.New("My error2"))

}

func slice2Array() {

	s := []int{1, 2, 3, 4, 5}
	///数组 容量可以小于实际，但是实际必须大于等于容量

	a := [4]int(s[0:5])

	s[0] = 999

	fmt.Println(s)
	fmt.Println(a)

}

func MakeEasy() {
	m := make([]int, 0, 5)

	m = append(m, 8)
	fmt.Println(m)
}

func Defer() func() {
	fmt.Println("Defer run")
	return func() {
		fmt.Println("Defer return call")
	}
}
func DeferTime(pre time.Time) {

	fmt.Printf("elapsed: %v\n", time.Since(pre))

}
