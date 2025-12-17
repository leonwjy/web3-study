package homework

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// task 1:
// 函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10
func Task1(n *int) {
	*n += 10
}

// task 2:
// 实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
func Task2(nums *[]int) {
	for i := range *nums {
		(*nums)[i] *= 2
	}
}

// task 3:
// 编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
func Task3() {
	go func() {
		for i := 1; i <= 10; i += 2 {
			fmt.Println(i)
		}
	}()
	go func() {
		for i := 2; i <= 10; i += 2 {
			fmt.Println(i)
		}
	}()
}

// task 4:
// 设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
func Task4(tasks []func()) {
	for _, task := range tasks {
		go func() {
			start := time.Now()
			task()
			elapsed := time.Since(start)
			fmt.Printf("Task executed in %s\n", elapsed)
		}()
	}
}

// task 5:
// 定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。
// 然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
// 在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

type Circle struct {
	Radius float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func Task5() {
	rect := Rectangle{Width: 10, Height: 5}
	circle := Circle{Radius: 3}
	fmt.Println(rect.Area())
	fmt.Println(rect.Perimeter())
	fmt.Println(circle.Area())
	fmt.Println(circle.Perimeter())
}

// task 6:
// 使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，
// 再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。
// 为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	Person
	EmployeeID string
}

func (e Employee) PrintInfo() {
	fmt.Printf("Name: %s, Age: %d, EmployeeID: %s\n", e.Name, e.Age, e.EmployeeID)
}

func Task6() {
	employee := Employee{Person: Person{Name: "John", Age: 30}, EmployeeID: "123456"}
	employee.PrintInfo()
}

// task 7:
// 编写一个程序，使用通道实现两个协程之间的通信。
// 一个协程生成从1到10的整数，并将这些整数发送到通道中，
// 另一个协程从通道中接收这些整数并打印出来。
func Task7() {
	ch := make(chan int)
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- i
		}
	}()
	go func() {
		for i := range ch {
			fmt.Println(i)
		}
	}()
}

// task 8
// 实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
func Task8() {
	ch := make(chan int, 10)
	go func() {
		for i := 0; i < 100; i++ {
			ch <- i
		}
	}()
	go func() {
		for i := range ch {
			fmt.Println(i)
		}
	}()
}

// task 9
// 编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
// 启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
func Task() {
	count := 0
	var mutex sync.Mutex
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				mutex.Lock()
				count++
				mutex.Unlock()
			}
		}()
	}
	wg.Wait()
	fmt.Println(count)
}

// task 10
// 使用原子操作（ sync/atomic 包）实现一个无锁的计数器。
// 启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
func Task10() {
	var count int64
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&count, 1)
			}
		}()
	}
	wg.Wait()
	fmt.Println(count)
}
