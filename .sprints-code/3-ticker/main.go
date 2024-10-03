package main

// Задание 6/6
// Используя Ticker, напишите программу,
// которая десять раз с интервалом в две секунды
// выведет разницу в секундах между текущим временем и временем запуска программы.
// Лучше выводить только целую часть секунд.

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	ticker := time.NewTicker(2 * time.Second)
	for i := 0; i < 10; i++ {
		<-ticker.C
		fmt.Println((time.Since(start).Seconds()))
	}
	ticker.Stop()
}
