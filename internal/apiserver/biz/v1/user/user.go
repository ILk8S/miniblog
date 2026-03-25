package user

import (
	"context"
	"errors"
	"time"

	"github.com/jinzhu/copier"
	"github.com/onexstack/onexstack/pkg/store/where"
	"github.com/wshadm/miniblog/internal/apiserver/model"
	"github.com/wshadm/miniblog/internal/apiserver/store"
	"github.com/wshadm/miniblog/internal/pkg/log"
	apiv1 "github.com/wshadm/miniblog/pkg/api/apiserver/v1"
	"github.com/wshadm/miniblog/pkg/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserBiz interface {
	Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error)
	Get(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error)
	List(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error)

	UserExpansion
}

// UserExpansion 定义用户操作的扩展方法.
type UserExpansion interface {
	Login(ctx context.Context, rq *apiv1.LoginRequest) (*apiv1.LoginResponse, error)
	RefreshToken(ctx context.Context, rq *apiv1.RefreshTokenRequest) (*apiv1.RefreshTokenResponse, error)
	ChangePassword(ctx context.Context, rq *apiv1.ChangePasswordRequest) (*apiv1.ChangePasswordResponse, error)
	ListWithBadPerformance(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error)
}

// userBiz 是 UserBiz 接口的实现.
type userBiz struct {
	store store.IStore
}

// 确保 userBiz 实现了 UserBiz 接口.
var _ UserBiz = (*userBiz)(nil)

func NewuserBiz(store store.IStore) *userBiz {
	return &userBiz{
		store: store,
	}
}

// Create 实现 UserBiz 接口中的 Create 方法.
func (u *userBiz) Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error) {
	var userM model.UserM
	_ = copier.Copy(&userM, rq)
	err := u.store.User().Create(ctx, &userM)
	if err != nil {
		return nil, err
	}
	return &apiv1.CreateUserResponse{UserID: userM.UserID}, nil
}

func (u *userBiz) Update(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userBiz) Delete(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userBiz) Get(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userBiz) List(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userBiz) Login(ctx context.Context, rq *apiv1.LoginRequest) (*apiv1.LoginResponse, error) {
	whr := where.F("username", rq.GetUsername())
	userM, err := u.store.User().Get(ctx, whr)
	if err != nil {
		return nil, errors.New("User Not Found")
	}
	// 对比传入的明文密码和数据库中已加密过的密码是否匹配
	if err := auth.Compare(userM.Password, rq.GetPassword()); err != nil {
		log.W(ctx).Errorw("Failed to compare password", "err", err)
		return nil, errors.New("username or password Invalid")
	}
	// TODO：实现 Token 签发逻辑
	return &apiv1.LoginResponse{Token: "<placeholder>",
		ExpireAt: timestamppb.New(time.Now().Add(2 * time.Hour))}, nil
}

func (u *userBiz) RefreshToken(ctx context.Context, rq *apiv1.RefreshTokenRequest) (*apiv1.RefreshTokenResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangePassword 实现 UserBiz 接口中的 ChangePassword 方法.
func (u *userBiz) ChangePassword(ctx context.Context, rq *apiv1.ChangePasswordRequest) (*apiv1.ChangePasswordResponse, error) {
	userM, err := u.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}
	//对比库里存储的密码和发起修改请求的密码是否一致。确认用户是否知道当前的密码
	//一个安全验证的步骤
	if err := auth.Compare(userM.Password, rq.GetOldPassword()); err != nil {
		log.W(ctx).Errorw("Failed to compare password", "err", err)
		return nil, err
	}
	userM.Password, _ = auth.Encrypt(rq.GetNewPassword())
	if err := u.store.User().Create(ctx, userM); err != nil {
		return nil, err
	}
	return &apiv1.ChangePasswordResponse{}, nil

}

func (u *userBiz) ListWithBadPerformance(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error) {
	//TODO implement me
	panic("implement me")
}
