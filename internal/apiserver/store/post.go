package store

import (
	"context"

	"github.com/onexstack/onexstack/pkg/store/where"
	"github.com/wshadm/miniblog/internal/apiserver/model"
	"github.com/wshadm/miniblog/internal/pkg/log"
)

type PostStore interface {
	Create(ctx context.Context, obj *model.PostM) error
	Update(ctx context.Context, obj *model.PostM) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.PostM, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.PostM, error)
	PostExpansion
}

// PostExpansion 定义了帖子操作的附加方法.
type PostExpansion interface{}

// postStore 是 PostStore 接口的实现.
type postStore struct {
	store *datastore
}

// 确保 postStore 实现了 PostStore 接口.
var _ PostStore = (*postStore)(nil)

// newPostStore 创建 postStore 的实例.
func newPostStore(store *datastore) *postStore {
	return &postStore{store: store}
}

// Create 插入一条帖子记录.
func (s *postStore) Create(ctx context.Context, obj *model.PostM) error {
	if err := s.store.DB(ctx).Create(obj).Error; err != nil {
		log.Errorw("Failed to insert post into database", "err", err, "post", obj)
		return err
	}
	return nil
}

// Update 更新帖子数据库记录.
func (s *postStore) Update(ctx context.Context, obj *model.PostM) error {
	if err := s.store.DB(ctx).Save(obj).Error; err != nil {
		log.Errorw("Failed to insert post into database", "err", err, "post", obj)
		return err
	}
	return nil
}

// Delete 根据条件删除帖子记录.
func (s *postStore) Delete(ctx context.Context, opts *where.Options) error {
	err := s.store.DB(ctx, opts).Delete(&model.PostM{}).Error
	if err != nil {
		log.Errorw("Failed to delete post from database", "err", err, "conditions", opts)
		return err
	}
	return nil
}

// Get 根据条件查询帖子记录.
func (s *postStore) Get(ctx context.Context, opts *where.Options) (*model.PostM, error) {
	var obj model.PostM //目的是作为返回实例使用
	if err := s.store.DB(ctx, opts).First(&obj).Error; err != nil {
		log.Errorw("Failed to retrieve post from database", "err", err, "conditions", opts)
		return nil, err
	}
	return &obj, nil
}

// List 返回帖子列表和总数.
func (s *postStore) List(ctx context.Context, opts *where.Options) (count int64, ret []*model.PostM, err error) {
	//直接使用返回值中的变量ret，也相当于直接修改它的值。由于在返回值中已经预先声明了变量。所以err要用“=”赋值
	//只使用return，是因为返回参数中已经命名了，会自动返回这些变量。
	err = s.store.DB(ctx, opts).Order("id desc").Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		log.Errorw("Failed to list posts from database", "err", err, "conditions", opts)
		return 0, nil, err
	}
	return
}
