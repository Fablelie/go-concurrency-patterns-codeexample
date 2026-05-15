package main

import (
	"fmt"
	"time"
)

// orderGenerator คือฟังก์ชันผู้ผลิตข้อมูล โดยจะส่ง Channel ที่เอาไว้อ่านข้อมูลกลับไป
func orderGenerator(prefix string, count int) <-chan string {
	// สร้าง Channel สำหรับส่งข้อมูลตัวอักษรออกไป
	out := make(chan string)

	// แตก Goroutine ออกไปทำงานเบื้องหลังทันที เพื่อไม่ให้บล็อกฟังก์ชันหลัก
	go func() {
		// เมื่อทำงานในลูปเสร็จหมดแล้ว ให้ปิด Channel เพื่อบอกฝั่งรับว่าข้อมูลหมดแล้ว
		defer close(out)

		for i := 1; i <= count; i++ {
			// จำลองเวลาในการผลิตข้อมูล (เช่น รอผลจาก Database หรือ API)
			time.Sleep(500 * time.Millisecond)

			// ส่งข้อมูลที่ผลิตเสร็จแล้วเข้าไปใน Channel
			out <- fmt.Sprintf("%s-%04d", prefix, i)
		}
	}()

	// ส่ง Channel กลับไปให้ฝั่งรับใช้งานทันที โดยไม่ต้องรอให้ Goroutine ข้างบนรันเสร็จ
	return out
}

func main() {
	fmt.Println("Starting Order Generator Demo...")
	startTime := time.Now()

	// เรียกใช้งาน Generator เพื่อขอรับ Channel ในการดึงข้อมูล Order
	orderChannel := orderGenerator("ORD", 5)

	// ฝั่งรับใช้ลูป range เพื่อดึงข้อมูลออกจาก Channel มาทีละตัว
	// ลูปนี้จะรอข้อมูลโดยอัตโนมัติ และจะจบการทำงานเองเมื่อ Channel ถูกสั่ง close
	for order := range orderChannel {
		fmt.Printf("Received: %s at %v\n", order, time.Since(startTime).Round(time.Millisecond))
	}

	fmt.Println("All orders processed successfully.")
}
