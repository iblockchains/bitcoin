package mempool

import "sync"

type FeeEstimator struct {
	// The last known height
	lastKnownHeight int32

	mtx sync.RWMutex
}

// LastKnownHeight 返回最后注册的高度
// returns the height of the last block which was registered.
func (ef *FeeEstimator) LastKnownHeight() int32 {
	ef.mtx.Lock()
	defer ef.mtx.Unlock()
	return ef.lastKnownHeight
}
