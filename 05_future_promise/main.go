package main

import (
	"fmt"
	"time"
)

// โครงสร้าง Future ที่ใช้ส่งกลับไปทันที ภายในมี Channel รอรับผลลัพธ์
type FuturePrice struct {
	result <-chan string
}

// ฟังก์ชันดึงราคาสินค้าที่ทำงานแบบดึงผลลัพธ์ในอนาคต
func fetchPriceFuture(product string) FuturePrice {
	c := make(chan string, 1)

	// สั่งทำงานเบื้องหลังทันที
	go func() {
		defer close(c)
		time.Sleep(1000 * time.Millisecond) // จำลองว่ารอโหลดข้อมูลนาน 1 วินาที
		c <- fmt.Sprintf("Price of %s is $99", product)
	}()

	// ส่งโครงสร้างที่สัญญาว่าจะให้ผลลัพธ์กลับไปให้ทันทีโดยไม่บล็อกโค้ดหลัก
	return FuturePrice{result: c}
}

func main() {
	fmt.Println("Starting Future/Promise Demo...")
	startTime := time.Now()

	// 1. สั่งเริ่มดึงข้อมูล (โค้ดจะวิ่งผ่านบรรทัดนี้ไปทันที ไม่ต้องรอ 1 วินาที)
	future := fetchPriceFuture("Go Book")

	// 2. ระหว่างที่ระบบหลังบ้านกำลังโหลดข้อมูล เราสามารถประมวลผลงานอื่นไปพร้อมกันได้
	fmt.Println("Doing other useful tasks in main thread...")
	time.Sleep(400 * time.Millisecond)
	fmt.Println("Main thread tasks finished.")

	// 3. เมื่อถึงเวลาที่ต้องใช้ข้อมูลจริง ๆ ค่อยเปิดอ่านค่าจาก Future (ถ้ายังไม่เสร็จ โค้ดจะหยุดรอตรงนี้)
	fmt.Println("Now waiting for the future result...")
	price := <-future.result

	fmt.Printf("Result: %s (Total time: %v)\n", price, time.Since(startTime).Round(time.Millisecond))
}
