package model

import (
    "time"
)

// Magnet app Users 用户模型
type Users struct {
    UID               int       `gorm:"primary_key"`
    Email             string    `gorm:"type:varchar(100);not null"`
    Password          string    `gorm:"type:varchar(100);not null"`
    RegisterTime      time.Time `gorm:"type:timestamp;default:current_timestamp()"`
    Fileban           int       `gorm:"type:int(11);default:0;comment:'文件屏蔽'"`
    Points            int       `gorm:"type:int(10);default:0;comment:'邀请积分'"`
    MaxView1          int       `gorm:"type:int(11);default:0"`
    MaxView2          int       `gorm:"type:int(11);default:0"`
    VipMaxView1       int       `gorm:"type:int(11);default:0"`
    VipMaxView2       int       `gorm:"type:int(11);default:0"`
    VipEndTime        int       `gorm:"type:int(11);default:0"`
    TotalPayment      float64   `gorm:"type:double;default:0;comment:'用户一共交了多少钱'"`
    Tiyan             int       `gorm:"type:int(11);default:0;comment:'体验会员买了几次'"`
    FavoriteSize      string    `gorm:"type:varchar(50);default:'536870912000';comment:'收藏夹容量默认10G'"`
    RegUUID           string    `gorm:"type:varchar(50);default:'0';comment:'注册uuid'"`
    RegIP             string    `gorm:"type:varchar(200);not null;comment:'注册IP'"`
    LoginUUID         string    `gorm:"type:varchar(50);default:'0';comment:'登陆uuid'"`
    LoginIP           string    `gorm:"type:varchar(200);default:null;comment:'登陆IP'"`
    LastActivityTime  time.Time `gorm:"type:timestamp;default:current_timestamp();comment:'最后登陆时间'"`
    RegChannel        string    `gorm:"type:varchar(50);not null;comment:'注册渠道'"`
    Remark            string    `gorm:"type:text;default:null;comment:'用户备注信息'"`
    Token             string    `gorm:"type:varchar(100);default:null"`
}

// TableName 返回自定义表名
func (Users) TableName() string {
    return "user" // 强制指定表名
}
