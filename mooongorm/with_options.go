// Package mooongorm
// Wrote by yijian on 2025/03/13
package mooongorm

import (
    "gorm.io/gorm"
)

// WithOmit 选项函数
func WithOmit(fields ...string) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Omit(fields...)
    }
}

// WithSelect 选项函数
func WithSelect(fields ...string) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Select(fields)
    }
}