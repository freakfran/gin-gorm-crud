package models

import "gorm.io/gorm"

// Post 文章
type Post struct {
	gorm.Model
	Title string // 标题
	Body  string // 正文
}
