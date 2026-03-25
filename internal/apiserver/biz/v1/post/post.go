package post

import (
	"context"

	"github.com/wshadm/miniblog/internal/apiserver/store"
	apiv1 "github.com/wshadm/miniblog/pkg/api/apiserver/v1"
)

// PostBiz 定义帖子相关业务方法.
type PostBiz interface {
	Create(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error)
	Get(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error)
	List(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error)
}

type postBiz struct {
	store store.IStore
}

var _ PostBiz = (*postBiz)(nil)

// NewPostBiz 创建 PostBiz 实现.
func NewPostBiz(s store.IStore) *postBiz {
	return &postBiz{store: s}
}

func (p *postBiz) Create(ctx context.Context, rq *apiv1.CreatePostRequest) (*apiv1.CreatePostResponse, error) {
	panic("implement me")
}

func (p *postBiz) Update(ctx context.Context, rq *apiv1.UpdatePostRequest) (*apiv1.UpdatePostResponse, error) {
	panic("implement me")
}

func (p *postBiz) Delete(ctx context.Context, rq *apiv1.DeletePostRequest) (*apiv1.DeletePostResponse, error) {
	panic("implement me")
}

func (p *postBiz) Get(ctx context.Context, rq *apiv1.GetPostRequest) (*apiv1.GetPostResponse, error) {
	panic("implement me")
}

func (p *postBiz) List(ctx context.Context, rq *apiv1.ListPostRequest) (*apiv1.ListPostResponse, error) {
	panic("implement me")
}
