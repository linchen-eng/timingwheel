package timingwheel

// ListNode 链表
type ListNode struct {
	val   int
	next  *ListNode
	tasks map[string]*Task //当前节点的任务集合 任务的唯一key作为集合键提高任务的删除效率
}

//创建环形链表
func initRingList(slotTotal int) *ListNode {
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
	return slotsListNode
}
