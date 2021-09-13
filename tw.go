package goTimeWheel

import (
	"time"
)

// TimeWheel Struct
type TimeWheel struct {
	interval time.Duration // ticker run interval

	ticker *time.Ticker

	slots [][]*Task

	keyPosMap map[interface{}]int // keep each timer's position

	slotNum int
	currPos int // timewheel current position

	addChannel    chan Task        // channel to add Task
	removeChannel chan interface{} // channel to remove Task
	stopChannel   chan bool        // stop signal
}

// Task Struct
type Task struct {
	key interface{} // Timer Task ID

	delay  time.Duration // Run after delay
	circle int           // when circle equal 0 will trigger

	fn     func(interface{}) // custom function
	params interface{}       // custom params
}

// NewTimeWheel Func: Generate TimeWheel with ticker and slotNum
func NewTimeWheel(interval time.Duration, slotNum int) *TimeWheel {

	if interval <= 0 || slotNum <= 0 {
		return nil
	}

	tw := &TimeWheel{
		interval:      interval,
		slots:         make([][]*Task, slotNum),
		keyPosMap:     make(map[interface{}]int),
		currPos:       0,
		slotNum:       slotNum,
		addChannel:    make(chan Task),
		removeChannel: make(chan interface{}),
		stopChannel:   make(chan bool),
	}

	for i := 0; i < slotNum; i++ {
		tw.slots[i] = make([]*Task, 16)
	}

	return tw
}

// Start Func: start ticker and monitor channel
func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.start()
}

func (tw *TimeWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.handle()
		case task := <-tw.addChannel:
			tw.addTask(&task)
		case key := <-tw.removeChannel:
			tw.removeTask(key)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimeWheel) Stop() {
	tw.stopChannel <- true
}

func (tw *TimeWheel) AddTimer(delay time.Duration, key interface{}, fn func(interface{}), params interface{}) {
	if delay < 0 {
		return
	}
	tw.addChannel <- Task{delay: delay, key: key, fn: fn, params: params}
}

func (tw *TimeWheel) RemoveTimer(key interface{}) {
	if key == nil {
		return
	}
	tw.removeChannel <- key
}

// handle Func: Do currPosition slots Task
func (tw *TimeWheel) handle() {
	currentSlice := tw.slots[tw.currPos]
	newSlice := make([]*Task, 16)

	for _, task := range currentSlice {
		if task == nil {
			continue
		}
		if task.circle > 0 {
			task.circle--
			newSlice = append(newSlice, task)
			continue
		}
		go task.fn(task.params)
		if task.key != nil {
			delete(tw.keyPosMap, task.key)
		}
	}

	tw.slots[tw.currPos] = newSlice
	tw.currPos = (tw.currPos + 1) % tw.slotNum
}

// getSlotNumAndCircle Func: parse duration by interval to get slotNum and circle.
func (tw *TimeWheel) getSlotNumAndCircle(d time.Duration) (slotNum int, circle int) {
	// circle represents how many iterations of the slots to wait before executing the task
	circle = int(d.Seconds()) / int(tw.interval.Seconds()) / tw.slotNum
	slotNum = (tw.currPos + int(d.Seconds())/int(tw.interval.Seconds())) % tw.slotNum
	return
}

func (tw *TimeWheel) addTask(task *Task) {
	slotNum, circle := tw.getSlotNumAndCircle(task.delay)
	task.circle = circle

	tw.slots[slotNum] = append(tw.slots[slotNum], task)

	if task.key != nil {
		tw.keyPosMap[task.key] = slotNum
	}
}

func (tw *TimeWheel) removeTask(key interface{}) {
	slotNum, ok := tw.keyPosMap[key]
	if !ok {
		return
	}

	slotSlice := tw.slots[slotNum]

	for taskIdx, task := range slotSlice {
		if task.key == key {
			slotSlice[taskIdx] = nil
			delete(tw.keyPosMap, task.key)
		}
	}
}
