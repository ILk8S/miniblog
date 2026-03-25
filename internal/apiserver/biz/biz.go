package biz

import (
	postv1 "github.com/wshadm/miniblog/internal/apiserver/biz/v1/post"
	userv1 "github.com/wshadm/miniblog/internal/apiserver/biz/v1/user"
	"github.com/wshadm/miniblog/internal/apiserver/store"
)

// IBiz 定义了业务层需要实现的方法.
type IBiz interface {
	// 获取用户业务接口.
	Userv1() userv1.UserBiz
	// 获取帖子业务接口.
	Postv1() postv1.PostBiz
}

// biz 是 IBiz 的一个具体实现.
type biz struct {
	store store.IStore
}

// 确保 biz 实现了 IBiz 接口.
var _ IBiz = (*biz)(nil)

// NewBiz 创建一个 IBiz 类型的实例.
func NewBiz(s store.IStore) *biz {
	return &biz{
		store: s,
	}
}

// Userv1 返回一个实现了 UserBiz 接口的实例.
func (b *biz) Userv1() userv1.UserBiz {
	return userv1.NewuserBiz(b.store)
}

// Postv1 返回一个实现了 PostBiz 接口的实例.
func (b *biz) Postv1() postv1.PostBiz {
	return postv1.NewPostBiz(b.store)
}
