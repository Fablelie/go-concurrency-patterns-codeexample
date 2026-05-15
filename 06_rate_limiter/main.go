package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting Rate Limiter Demo...")
	startTime := time.Now()

	// จำลองงานส่งคำขอ 5 งาน
	requests := []int{1, 2, 3, 4, 5}

	// สร้างนาฬิกาจับเวลา (Ticker) ให้ส่งสัญญาณออกมาทุก ๆ 200 มิลลิวินาที
	// เปรียบเสมือนใบอนุญาต (Token) ที่จะปล่อยให้งานวิ่งผ่านไปได้ทีละหนึ่งชิ้นตามเวลาที่กำหนด
	limiter := time.NewTicker(200 * time.Millisecond)
	defer limiter.Stop()

	for _, req := range requests {
		// โค้ดจะหยุดรอตรงนี้ จนกว่าตัว Limiter จะส่งสัญญาณเวลาถัดไปออกมา
		<-limiter.C

		// ประมวลผลงานเมื่อได้รับอนุญาตตามรอบเวลา
		fmt.Printf("Processing request %d at %v\n", req, time.Since(startTime).Round(time.Millisecond))
	}

	fmt.Println("All requests rate-limited and processed successfully.")
}
