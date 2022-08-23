package timingwheel

import (
	"time"
)

// TimingWheel 时间轮结构体
type TimingWheel struct {
	interval       time.Duration //指针每隔多久往前移动一格
	ticker         *time.Ticker  //定时器
	slots          *ListNode     //时间轮槽
	slotNum        int           //槽的数量
	handler        Handler       //任务处理回调函数
	addTaskChannel chan Task     //新增任务channel
	delTaskChannel chan string   //删除任务channel
	stopChannel    chan bool     //停止定时器channel
}

// Handler 任务回调处理函数
type Handler func(interface{})

// CreateTimingWheel 创建时间轮
// interval 指针每隔多久往前移动一格
// slotTotal 槽的数量
// handler 任务回调处理函数
func CreateTimingWheel(interval time.Duration, slotTotal int, handler Handler) *TimingWheel {
	//创建时间槽环形链表
	var slotsListNode *ListNode
	var tail *ListNode

	for i := 0; i < slotTotal; i++ {
		if slotsListNode == nil {
			slotsListNode = &ListNode{val: i}
			tail = slotsListNode
			continue
		}
		tail.next = &ListNode{val: i}
		tail = tail.next
	}
	tail.next = slotsListNode

	return &TimingWheel{
		interval:       interval,
		slots:          slotsListNode,
		slotNum:        slotTotal,
		handler:        handler,
		addTaskChannel: make(chan Task),
		delTaskChannel: make(chan string),
		stopChannel:    make(chan bool),
	}
}

//转动时间轮 响应处理各类时间
func (tw *TimingWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			//获取当前时间槽的任务集合
			tasks := tw.slots.tasks
			//任务处理的回调方法
			tw.tickHandler(tasks)
			//自动转下一个时间槽
			tw.slots = tw.slots.next
		case task := <-tw.addTaskChannel:
			//添加任务
			tw.addTask(task)
		case key := <-tw.delTaskChannel:
			//删除任务
			tw.delTask(key)
		case <-tw.stopChannel:
			//停止定时器
			tw.ticker.Stop()
			return
		}
	}
}

//任务处理的回调方法
func (tw *TimingWheel) tickHandler(tasks map[string]*Task) {
	for key, task := range tasks {
		if task.round > 0 {
			//当前任务还不到能执行的时候
			task.round -= 1
			continue
		}
		delete(tasks, key)
		go tw.handler(task)
	}
}

//添加任务
func (tw *TimingWheel) addTask(taskInfo Task) {
	//获取时间槽位置 和圈数
	delaySlots := taskInfo.delay / int(tw.interval.Seconds())
	round := delaySlots / tw.slotNum
	step := delaySlots % tw.slotNum //时间轮当前指针位置再往前走几格
	taskInfo.round = round
	curSlots := tw.slots
	for curSlots != nil {
		if step <= 0 {
			if curSlots.tasks == nil {
				curSlots.tasks = make(map[string]*Task)
			}
			curSlots.tasks[taskInfo.key] = &taskInfo
			break
		}
		curSlots = curSlots.next
		step--
	}
}

//删除任务
func (tw *TimingWheel) delTask(key string) {
	curSlots := tw.slots
	curVal := curSlots.val
	for curSlots != nil {
		if _, ok := curSlots.tasks[key]; ok {
			delete(curSlots.tasks, key)
		}
		curSlots = curSlots.next
		if curVal == curSlots.val { //已经遍历了一轮
			break
		}
	}
}

// Running 启动时间轮
func (tw *TimingWheel) Running() {
	//创建定时器
	tw.ticker = time.NewTicker(tw.interval)
	//开启协程 开始转动时间轮
	go tw.start()
}

// AddTask 添加延时任务到时间轮
func (tw *TimingWheel) AddTask(taskInfo Task) {
	tw.addTaskChannel <- taskInfo
}

// DelTask 删除延时任务
func (tw *TimingWheel) DelTask(key string) {
	tw.delTaskChannel <- key
}

// Stop 停止时间轮
func (tw *TimingWheel) Stop() {
	tw.stopChannel <- true
}
