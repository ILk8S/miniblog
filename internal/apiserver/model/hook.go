package model

import "gorm.io/gorm"

// AfterCreate 在创建数据库记录之后生成postID
func (m *PostM) AfterCreate(tx *gorm.DB) error {
	return nil
}

// AfterCreate 在创建数据库记录之后生成 userID.
func (m *UserM) AfterCreate(tx *gorm.DB) error {
	return nil
}
