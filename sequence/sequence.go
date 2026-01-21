package sequence

// Sequence取号器接口 这样我们不在乎底层是Mysql还是Redis实现，我们只在乎Next方法是否存在并且能使用
type Sequence interface {
	Next() (uint64, error)
}
