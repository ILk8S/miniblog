package rid

import (
	"github.com/onexstack/onexstack/pkg/id"
)

const defaultABC = "abcdefghijklmnopqrstuvwxyz1234567890"

type ResourceID string

const (
	//UserID定义用户资源标识符
	UserID ResourceID = "user"
	//PostID 定义博文资源标识符
	PostID ResourceID = "post"
)

// String 将资源标识符转换成字符串
func (rid ResourceID) String() string {
	return string(rid)
}

// New创建带前缀的唯一标识符
func (rid ResourceID) New(conter uint64) string {
	//使用自定义选项生成唯一标识符
	uniqueStr := id.NewCode(conter, id.WithCodeChars([]rune(defaultABC)), id.WithCodeL(6), id.WithCodeSalt(Salt()))
	return rid.String() + "-" + uniqueStr
}
