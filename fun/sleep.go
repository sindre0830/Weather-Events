package fun

import "time"

func HookSleep(number int64) {
	nextTime := time.Now().Truncate(time.Second)
	nextTime = nextTime.Add(time.Duration(number) * time.Second) // Change to hour!
	time.Sleep(time.Until(nextTime))
}