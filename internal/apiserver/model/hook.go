package model

import (
	"github.com/wshadm/miniblog/internal/pkg/rid"
	"github.com/wshadm/miniblog/pkg/auth"
	"gorm.io/gorm"
)

// BeforeCreate 在创建数据库记录之前加密明文密码.
func (m *UserM) BeforeCreate(tx *gorm.DB) error {
	var err error
	m.Password, err = auth.Encrypt(m.Password)
	if err != nil {
		return err
	}
	return nil
}

// AfterCreate 在创建数据库记录之后生成postID
func (m *PostM) AfterCreate(tx *gorm.DB) error {
	m.PostID = rid.PostID.New(uint64(m.ID))
	return tx.Save(m).Error
}

// AfterCreate 在创建数据库记录之后生成 userID.
func (m *UserM) AfterCreate(tx *gorm.DB) error {
	m.UserID = rid.UserID.New(uint64(m.ID))
	return tx.Save(m).Error
}
