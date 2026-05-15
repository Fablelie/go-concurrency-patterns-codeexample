package main

import (
	"fmt"
	"time"
)

// Stage 1: ฟังก์ชันสำหรับเตรียมตัวเลข (Generator)
// รับ Array ของตัวเลขเข้ามา แล้วส่งออกไปทีละตัวผ่าน Channel ขาออก
func generator(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, n := range nums {
			out <- n
		}
	}()
	return out
}

// Stage 2: ฟังก์ชันสำหรับคำนวณกำลังสอง (Transformer)
// รับ Channel ขาเข้า นำตัวเลขมาคูณตัวเอง แล้วส่งออกไปที่ Channel ขาออก
func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			// จำลองว่าขั้นตอนนี้ใช้เวลาประมวลผล 300 มิลลิวินาที
			time.Sleep(300 * time.Millisecond)
			out <- n * n
		}
	}()
	return out
}

// Stage 3: ฟังก์ชันสำหรับบวกเลขเพิ่ม (Transformer)
// รับค่าจากสเตจกำลังสองมาบวกเพิ่มอีก 10 แล้วส่งออกไปที่ Channel ขาออกขั้นสุดท้าย
func addTen(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for n := range in {
			// จำลองว่าขั้นตอนนี้ใช้เวลาประมวลผล 200 มิลลิวินาที
			time.Sleep(200 * time.Millisecond)
			out <- n + 10
		}
	}()
	return out
}

func main() {
	fmt.Println("Starting Pipeline Demo...")
	startTime := time.Now()

	// ข้อมูลดิบเริ่มต้นที่ต้องการนำเข้าสายพาน
	inputs := []int{2, 3, 4}

	// เชื่อมต่อแต่ละ Stage เข้าด้วยกันเป็นสายพาน
	// [inputs] -> generator -> [chan1] -> square -> [chan2] -> addTen -> [resultChan]
	chan1 := generator(inputs)
	chan2 := square(chan1)
	resultChan := addTen(chan2)

	// Stage สุดท้าย: ดึงข้อมูลปลายสายพานมาแสดงผล (Consumer)
	for result := range resultChan {
		fmt.Printf("Result: %d (Time elapsed: %v)\n", result, time.Since(startTime).Round(time.Millisecond))
	}

	fmt.Println("Pipeline processing finished.")
}
