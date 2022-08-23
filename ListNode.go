package timingwheel

// ListNode 链表
type ListNode struct {
	val   int
	next  *ListNode
	tasks map[string]*Task //当前节点的任务集合 任务的唯一key作为集合键提高任务的删除效率
}
