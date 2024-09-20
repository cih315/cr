package model

import (
	"github.com/cloudreve/Cloudreve/v3/pkg/util"
	"github.com/jinzhu/gorm"
	"time"
)

// MagnetUser 定义用户模型
type MagnetUser struct {
	UID           int       // 用户ID
	Email         string    // 用户邮箱
	RegisterTime  time.Time // 注册时间
	FileBan       int       // 文件禁用状态
	Points        int       // 用户积分
	MaxView1      int       // 最大查看次数1
	MaxView2      int       // 最大查看次数2
	VipMaxView1   int       // VIP最大查看次数1
	VipMaxView2   int       // VIP最大查看次数2
	VipEndTime    int       // VIP结束时间
	TotalPayment  float64   // 总支付金额
	Tiyan         int       // 体验次数
	FavoriteSize  string    // 收藏大小
	RegUUID       string    // 注册UUID
	RegIP         string    // 注册IP
	LoginUUID     string    // 登录UUID
	LoginIP       string    // 登录IP
	LastActivity  time.Time // 最后活动时间
	RegChannel    string    // 注册渠道
	Remark        string    // 备注
	Token         string    // 用户Token
}

// NewMagnetUser 返回一个新的空 MagnetUser
func NewMagnetUser() MagnetUser {
	return MagnetUser{}
}

// TableName 设置表名
func (MagnetUser) TableName() string {
    return "user" // 设置为您想要的表名
}


// GetMagnetUserByID 用UID获取MagnetUser
func GetMagnetUserByID(uid int) (MagnetUser, error) {
	var magnetUser MagnetUser
	result := DB2.Where("uid = ?", uid).First(&magnetUser)
	return magnetUser, result.Error
}

// UpdateMagnetUser 更新MagnetUser信息
func (magnetUser *MagnetUser) UpdateMagnetUser(val map[string]interface{}) error {
	return DB2.Model(magnetUser).Updates(val).Error
}

// CheckEmail 检查邮箱是否已被使用
func CheckEmail(email string) (bool, error) {
	  // 将布尔值 chk 转换为字符串
	  util.Log().Info("check email....") // 使用 fmt.Sprintf

	var count int
	result := DB2.Model(&MagnetUser{}).Where("email = ?", email).Count(&count)
	return count > 0, result.Error
}

// GetMagnetUserByEmail 根据邮箱查找并返回 MagnetUser
func GetMagnetUserByEmail(email string) (MagnetUser, error) {
    var magnetUser MagnetUser
    result := DB2.Debug().Where("email = ?", email).First(&magnetUser)

    // 检查是否找到用户
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            return MagnetUser{}, nil // 没有找到用户，返回空的 MagnetUser
        }
        return MagnetUser{}, result.Error // 发生错误
    }

    return magnetUser, nil // 找到用户，返回用户信息
}
