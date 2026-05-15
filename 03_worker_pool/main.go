package main

import (
	"fmt"
	"sync"
	"time"
)

// worker คือฟังก์ชันคนงานที่จะคอยดึงงานจากคิวไปทำ
// id: หมายเลขประจำตัวคนงาน
// jobs: Channel ขาเข้า (ไว้อ่านงานอย่างเดียว)
// results: Channel ขาออก (ไว้ส่งผลลัพธ์อย่างเดียว)
func worker(id int, jobs <-chan int, results chan<- int, waitGroup *sync.WaitGroup) {
	// บอก WaitGroup ว่าคนงานคนนี้ทำงานเสร็จแล้วเมื่อฟังก์ชันนี้จบ (เพราะคิวงานปิด)
	defer waitGroup.Done()

	// คนงานจะต่อคิวรอรับงานจากช่อง jobs ไปเรื่อยๆ จนกว่าช่อง jobs จะถูกสั่ง close
	for job := range jobs {
		fmt.Printf("Worker %d started job %d\n", id, job)

		// จำลองว่างานแต่ละชิ้นใช้เวลาทำ 500 มิลลิวินาที
		time.Sleep(500 * time.Millisecond)

		fmt.Printf("Worker %d finished job %d\n", id, job)

		// ส่งผลลัพธ์ของงานกลับไปในช่อง results
		results <- job * 2
	}
}

func main() {
	fmt.Println("Starting Worker Pool Demo...")
	startTime := time.Now()

	numJobs := 6    // จำนวนงานทั้งหมดในระบบ
	numWorkers := 3 // จำกัดจำนวนคนงานให้ทำงานพร้อมกันได้แค่ 3 ตัวเท่านั้น

	// สร้างคิวสำหรับส่งงานและคิวสำหรับรับผลลัพธ์ (แบบมี Buffer เพื่อไม่ให้บล็อกการส่ง)
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	var wg sync.WaitGroup

	// 1. สร้าง Worker (คนงาน) ขึ้นมาทำงานตามจำนวนที่กำหนดไว้
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		// แตก Goroutine ให้คนงานแต่ละตัวไปสแตนด์บายรองานในระบบ
		go worker(w, jobs, results, &wg)
	}

	// 2. ส่งงานจำนวน 6 ชิ้นเข้าไปในคิวกลาง
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	// ปิดคิวงานทันทีเมื่อส่งครบ เพื่อบอกให้คนงานรู้ว่า "ไม่มีงานใหม่มาเพิ่มแล้วนะ"
	close(jobs)

	// 3. แตก Goroutine เบื้องหลังมารอให้คนงานทุกคนทำงานเสร็จ แล้วปิดคิวผลลัพธ์
	go func() {
		wg.Wait()
		close(results)
	}()

	// 4. ฝั่งรับดึงผลลัพธ์จากคิว results มาแสดงผล
	for res := range results {
		fmt.Printf("Result received: %d\n", res)
	}

	fmt.Printf("All jobs finished. Total time: %v\n", time.Since(startTime).Round(time.Millisecond))
}
