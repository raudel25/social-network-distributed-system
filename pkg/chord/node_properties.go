package chord

func (n *Node) setPredecessorProp(predecessor *Node) {
	n.predLock.Lock()
	defer n.predLock.Unlock()

	n.predecessor = predecessor
}

func (n *Node) getPredecessorProp() *Node {
	n.predLock.RLock()
	defer n.predLock.RUnlock()

	return n.predecessor
}

func (n *Node) successorsPushBack(successor *Node) {
	n.sucLock.Lock()
	defer n.sucLock.Unlock()

	n.successors.PushBack(successor)
}

func (n *Node) successorsPushFront(successor *Node) {
	n.sucLock.Lock()
	defer n.sucLock.Unlock()

	n.successors.PushFront(successor)
}

func (n *Node) successorsPopBack() {
	n.sucLock.Lock()
	defer n.sucLock.Unlock()

	n.successors.PopBack()
}

func (n *Node) successorsPopFront() {
	n.sucLock.Lock()
	defer n.sucLock.Unlock()

	n.successors.PopFront()
}

func (n *Node) successorsFront() *Node {
	n.sucLock.RLock()
	defer n.sucLock.RUnlock()

	return n.successors.Front()
}

func (n *Node) successorsBack() *Node {
	n.sucLock.RLock()
	defer n.sucLock.RUnlock()

	return n.successors.Back()
}

func (n *Node) successorsLen() int {
	n.sucLock.RLock()
	defer n.sucLock.RUnlock()

	return n.successors.Len()
}
