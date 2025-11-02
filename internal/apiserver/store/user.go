package store

import (
	"context"
	"errors"

	"github.com/onexstack/onexstack/pkg/store/where"
	"github.com/wshadm/miniblog/internal/apiserver/model"
	"github.com/wshadm/miniblog/internal/pkg/log"
	"gorm.io/gorm"
)

// UserStore 定义了 user 模块在 store 层所实现的方法.
type UserStore interface {
	Create(ctx context.Context, obj *model.UserM) error
	Update(ctx context.Context, obj *model.UserM) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.UserM, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.UserM, error)
	UserExpansion
}

// UserExpansion 定义了用户操作的附加方法.
type UserExpansion interface{}

// userStore 是 UserStore 接口的实现.
type userStore struct {
	store *datastore
}

// newUserStore 创建 userStore 的实例.
func newUserStore(store *datastore) *userStore {
	return &userStore{store}
}

// 确保 userStore 实现了 UserStore 接口.
var _ UserStore = (*userStore)(nil)

// Create 插入一条用户记录.
func (s *userStore) Create(ctx context.Context, obj *model.UserM) error {
	if err := s.store.DB(ctx).Create(obj).Error; err != nil {
		log.Errorw("Failed to insert user into database", "err", err, "user", obj)
		return err
	}
	return nil
}

// Update 更新用户数据库记录.
func (s *userStore) Update(ctx context.Context, obj *model.UserM) error {
	//save：如果记录有主键，那就执行update。如果记录没有主键就执行inster
	if err := s.store.DB(ctx).Save(obj).Error; err != nil {
		log.Errorw("Failed to update user in database", "err", err, "user", obj)
		return err
	}
	return nil
}

// Delete 根据条件删除用户记录.
func (s *userStore) Delete(ctx context.Context, opts *where.Options) error {
	if err := s.store.DB(ctx, opts).Delete(new(model.UserM)).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorw("Failed to delete user from database", "err", err, "conditions", opts)
		return err
	}
	return nil
}

// Get 根据条件查询用户记录.
func (s *userStore) Get(ctx context.Context, opts *where.Options) (*model.UserM, error) {
	var obj model.UserM
	if err := s.store.DB(ctx, opts).First(&obj).Error; err != nil && errors.As(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &obj, nil
}
func (s *userStore) List(ctx context.Context, opts *where.Options) (count int64, ret []*model.UserM, err error) {
	err = s.store.DB(ctx, opts).Order("id desc").Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		log.Errorw("Failed to list users from database", "err", err, "conditions", opts)
		return
	}
	return
}
