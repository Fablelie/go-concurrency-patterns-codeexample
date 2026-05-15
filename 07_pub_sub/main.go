package main

import (
	"fmt"
	"sync"
	"time"
)

// Broker ทำหน้าที่เป็นตัวกลางในการจัดการคลังผู้รับข่าวสาร
type Broker struct {
	mu          sync.Mutex
	subscribers []chan string
}

// ฟังก์ชันสำหรับฝั่งผู้รับมาลงชื่อขอฟังข่าวสาร
func (b *Broker) Subscribe() chan string {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string, 1)
	b.subscribers = append(b.subscribers, ch)
	return ch
}

// ฟังก์ชันสำหรับผู้ส่งเพื่อกระจายข่าวให้ทุกคนที่รออยู่
func (b *Broker) Publish(msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, ch := range b.subscribers {
		// ใช้ select เพื่อไม่ให้ชะงักหากช่องรับเต็ม
		select {
		case ch <- msg:
		default:
			fmt.Println("Buffer full, skipping subscriber")
		}
	}
}

func prepareShipping(orderID string) {
	fmt.Printf("[Delivery Dept] Packing items for %s...\n", orderID)
	// Code สำหรับจองขนส่ง หรือตัดสต็อกสินค้าจริง
}

func createInvoice(orderID string) {
	fmt.Printf("[Accounting Dept] Generating invoice for %s...\n", orderID)
	// Code สำหรับคำนวณเงิน และส่ง Email ใบเสร็จให้ลูกค้า
}

func main() {
	fmt.Println("Starting Pub/Sub Demo...")
	broker := &Broker{}

	// 1. สร้างผู้รับคนที่ 1 (เช่น แผนกจัดส่งของ)
	sub1 := broker.Subscribe()
	go func() {
		for msg := range sub1 {
			fmt.Printf("[Delivery Dept] Received alert: %s\n", msg)
			// นำ msg ไปทำงานต่อ
			prepareShipping(msg)
		}
	}()

	// 2. สร้างผู้รับคนที่ 2 (เช่น แผนกบัญชี)
	sub2 := broker.Subscribe()
	go func() {
		for msg := range sub2 {
			fmt.Printf("[Accounting Dept] Received alert: %s\n", msg)
			// นำ msg ไปทำงานต่อ
			createInvoice(msg)
		}
	}()

	// 3. ฝั่งผู้ส่งทำการส่งข่าวสารเข้าระบบกลาง
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Publishing news: 'New Order #1001 Created'")
	broker.Publish("Order #1001")

	// รอให้ Goroutines ทำงานแสดงผลเสร็จสิ้น
	time.Sleep(200 * time.Millisecond)
}
