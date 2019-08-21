package blockchain

// MedianTimeSource 提供一种机制，通过添加几个时间样本来确定中间时间
// 用于表示本地时钟的偏移量
// provides a mechanism to add several time samples which are
// used to determine a median time which is then used as an offset to the local
// clock.
type MedianTimeSource interface{}
