package timingwheel

// Task 延时任务
type Task struct {
	key   string      //任务的唯一标识
	round int         //当前任务轮转到第几轮执行 round = 0时当前可执行
	info  interface{} //任务具体内容
	delay int         //当前任务延时多少秒
}
