package main

import (
	"fmt"
	"sync"
	"time"
)

// ฟังก์ชันจำลองการดึงราคาจากร้านค้า (ทำงานแบบ Asynchronous)
func fetchPrice(shopName string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		// จำลองว่าแต่ละร้านใช้เวลาตอบกลับไม่เท่ากัน
		if shopName == "Shop-A" {
			time.Sleep(800 * time.Millisecond)
		} else {
			time.Sleep(300 * time.Millisecond)
		}
		out <- fmt.Sprintf("[%s] Price: $100", shopName)
	}()
	return out
}

// Fan-In: ฟังก์ชันสำหรับรวมหลายๆ Channels เข้ามาเป็น Channel เดียว
func merge(channels ...<-chan string) <-chan string {
	var waitGroup sync.WaitGroup
	out := make(chan string)

	// ฟังก์ชันภายในที่จะดึงข้อมูลจาก Channel หนึ่งส่งต่อไปยัง Channel กลาง (out)
	output := func(c <-chan string) {
		defer waitGroup.Done()
		for n := range c {
			out <- n
		}
	}

	// สั่งให้ waitGroup รับรู้จำนวน Channel ทั้งหมดที่ต้องรอ
	waitGroup.Add(len(channels))
	for _, c := range channels {
		// Fan-Out ย่อยๆ: แตก Goroutine ไปคอยดึงข้อมูลจากแต่ละ Channel พร้อมกัน
		go output(c)
	}

	// แตก Goroutine ขนานไปอีกตัวเพื่อคอยปิด Channel กลางเมื่อทุกร้านส่งข้อมูลครบแล้ว
	go func() {
		waitGroup.Wait()
		close(out)
	}()

	return out
}

func main() {
	fmt.Println("Starting Fan-Out & Fan-In Demo...")
	startTime := time.Now()

	// 1. [Fan-Out] กระจายการทำงานออกไปยิง API 3 ร้านพร้อมกันแยกเป็นอิสระ
	chan1 := fetchPrice("Shop-A")
	chan2 := fetchPrice("Shop-B")
	chan3 := fetchPrice("Shop-C")

	// 2. [Fan-In] รวบรวม Channels ของทั้ง 3 ร้านเข้ามาเหลือช่องทางเดียวด้วยฟังก์ชัน merge
	mergedChannel := merge(chan1, chan2, chan3)

	// 3. ฝั่งรับเปิดอ่านข้อมูลจาก Channel กลางที่รวมผลลัพธ์มาให้แล้ว
	for result := range mergedChannel {
		fmt.Printf("Result: %s (Time elapsed: %v)\n", result, time.Since(startTime).Round(time.Millisecond))
	}

	fmt.Println("All shop prices merged completely.")
}
