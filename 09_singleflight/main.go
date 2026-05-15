package main

import (
	"fmt"
	"sync"
	"time"

	// นำเข้าแพ็กเกจมาตรฐานของ Go สำหรับทำ Singleflight
	"golang.org/x/sync/singleflight"
)

// จำลองระบบ Cache และ Database
var (
	cacheMutex sync.RWMutex
	mockCache  = make(map[string]string)
)

// ฟังก์ชันจำลองการดึงข้อมูลจาก Cache
func getFromCache(key string) (string, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	val, exists := mockCache[key]
	return val, exists
}

// ฟังก์ชันจำลองการดึงข้อมูลจาก Database (ใช้เวลาโหลดนานและกินทรัพยากรสูง)
func getFromDB(productID string) string {
	fmt.Printf("[Database] Executing heavy SQL Query for product: %s...\n", productID)
	time.Sleep(1000 * time.Millisecond) // จำลองว่าดึงข้อมูลจาก DB นาน 1 วินาที
	return fmt.Sprintf("Detail of %s (Fetched at %s)", productID, time.Now().Format("15:04:05"))
}

func main() {
	fmt.Println("Starting Production Singleflight Demo (Cache Stampede Protection)...")
	startTime := time.Now()

	// สร้าง Group ของ singleflight (ปกติจะประกาศเป็น Global หรืออยู่ใน Struct ของ Service)
	var sfGroup singleflight.Group
	var wg sync.WaitGroup

	productKey := "product_123"
	// ณ ตอนนี้ mockCache ว่างเปล่า เปรียบเสมือน Cache หมดอายุ (Cache Miss)

	// จำลองสถานการณ์: มีลูกค้า 5 คน กดดูสินค้าชิ้นเดียวกันเข้ามาพร้อมกันในเสี้ยววินาที
	totalUsers := 5
	for userID := 1; userID <= totalUsers; userID++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			// ขั้นตอนที่ 1: ตรวจสอบใน Cache ก่อน
			if data, found := getFromCache(productKey); found {
				fmt.Printf("User %d -> [Cache Hit] Data: %s\n", id, data)
				return
			}

			// ขั้นตอนที่ 2: ถ้า Cache Miss ให้ใช้ singleflight คุมการยิงไปที่ DB
			// ฟังก์ชัน sfGroup.Do จะบล็อกคำขอที่ซ้ำกันไว้ แล้วแชร์ผลลัพธ์ร่วมกัน
			// v: ผลลัพธ์ที่ได้คืนกลับมา, err: ข้อผิดพลาด, shared: บอกว่ามีคนอื่นแชร์ผลลัพธ์นี้ร่วมกับเราไหม
			v, _, shared := sfGroup.Do(productKey, func() (interface{} /* ข้อมูลที่คืนค่า */, error) {
				// โค้ดในบล็อกนี้จะทำงานเพียง "ครั้งเดียว" สำหรับ Request ที่เข้ามาพร้อมกัน
				dbData := getFromDB(productKey)

				// นำข้อมูลที่ได้จาก DB ไปอัปเดตลง Cache เพื่อให้คนรอบหน้าอ่านจาก Cache ได้เลย
				cacheMutex.Lock()
				mockCache[productKey] = dbData
				cacheMutex.Unlock()

				return dbData, nil
			})

			// ขั้นตอนที่ 3: แสดงผลลัพธ์ที่ทุกคนได้รับ
			fmt.Printf("User %d -> Received Data: %v (Shared result: %t)\n", id, v, shared)
		}(userID)
	}

	wg.Wait()
	fmt.Printf("All users served. Total time elapsed: %v\n", time.Since(startTime).Round(time.Millisecond))
}
