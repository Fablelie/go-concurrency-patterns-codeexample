package main

import (
	"fmt"
	"time"
)

// orFunc ทำหน้าที่รับช่องสัญญาณยกเลิก (done) และช่องสัญญาณข้อมูล (c)
// หากช่องไหนส่งข้อมูลหรือสั่งปิดก่อน ฟังก์ชันนี้จะหยุดทำงานทันที
func orDone(done <-chan struct{}, c <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for {
			select {
			case <-done: // ถ้าระบบสั่งยกเลิก ให้หลุดจากลูปทันที
				return
			case v, ok := <-c:
				if !ok { // ถ้าช่องข้อมูลหลักปิดตัวเอง ให้หยุดทำงาน
					return
				}
				select {
				case out <- v:
				case <-done:
				}
			}
		}
	}()
	return out
}

func main() {
	fmt.Println("Starting Or-Done Channel Demo...")
	done := make(chan struct{})
	dataChan := make(chan string)

	// จำลองการส่งข้อมูลมาเรื่อยๆ ทุก 200 มิลลิวินาที
	go func() {
		defer close(dataChan)
		for i := 1; i <= 5; i++ {
			time.Sleep(200 * time.Millisecond)
			dataChan <- fmt.Sprintf("Data chunk %d", i)
		}
	}()

	// เรียกใช้งาน orDone หุ้มช่องข้อมูลไว้
	safeChan := orDone(done, dataChan)

	// สั่งจำลองเหตุการณ์ยกเลิกงานกลางคันหลังจากผ่านไป 500 มิลลิวินาที
	go func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("--- User clicked CANCEL button! ---")
		close(done) // ส่งสัญญาณยกเลิก
	}()

	// ลูปอ่านข้อมูลจากช่องสัญญาณที่ปลอดภัย
	for data := range safeChan {
		fmt.Printf("Processed: %s\n", data)
	}

	fmt.Println("Worker safely exited. No goroutine leak.")
}
