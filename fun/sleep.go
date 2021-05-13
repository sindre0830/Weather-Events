package fun

import "time"

/**
* HookSleep
* Function that gets passed an integer, and sleeps the appropriate amount of time based on it
**/
func HookSleep(number int64) {
	nextTime := time.Now().Truncate(time.Second)
	nextTime = nextTime.Add(time.Duration(number) * time.Second) // Change to hour!
	time.Sleep(time.Until(nextTime))
}