package biz

// IBiz 定义了 Biz 层需要实现的方法.
type IBiz interface {
	// TODO: Add your biz interfaces here
}

// biz 是 IBiz 的一个具体实现.
type biz struct {
	// TODO: Add your dependencies here
}

// 确保 biz 实现了 IBiz 接口.
var _ IBiz = (*biz)(nil)

// NewBiz 创建一个 IBiz 类型的实例.
func NewBiz() *biz {
	return &biz{}
}

// TODO: Implement your biz methods here