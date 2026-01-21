package base62

//62进制转换
//0123456789abcdefjhijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
//0-9: 0-9
//a-z:10-35
//A-Z:36-61

// 定义一个转换的字符串
// 为了避免被人恶意请求，我们把字符串打乱(已经定义在配置文件)
var (
	base62Str string
)

// MustInit 要使用base62 必须调用该函数完成初始化
func MustInit(bs string) {
	if len(bs) == 0 {
		panic("need base string")
	}
	base62Str = bs
}

// Int2String 十进制转换为62进制
func Int2String(seq uint64) string {
	if seq == 0 {
		return string(base62Str[0])
	}
	// 1.预分配一个足够大的字节数组 (11位足够存 uint64 的最大值)
	// 使用数组而不是切片，分配在栈上，速度会更快
	var buf [11]byte

	// 2. 定义一个指针或索引，从数组的末尾开始往前填
	i := len(buf)

	for seq > 0 {
		mod := seq % 62
		seq = seq / 62
		// 索引前移
		i--
		// 直接赋值，不使用 append
		buf[i] = base62Str[mod]
	}
	// 3. 返回切片：从我们填入的第一个字符开始切到最后
	// string(buf[i:]) 会发生一次内存拷贝，但这无法避免
	return string(buf[i:])
}
