package model

import (
    "time"
	"fmt"
	"github.com/cloudreve/Cloudreve/v3/pkg/util"
	//"database/sql"
)

// Favorite 用户收藏模型
type Favorite struct {
    ID     int       // 主键
    UID    int       // 用户ID
    DHash  string    // 文件包hash
    Hash   string    // 文件hash
    DName  string    // 种子名-》目录名
    Name   string    // 文件名
    Size   uint64    // 文件大小
    AddTime time.Time // 添加时间
    AddIP  string    // 添加IP
}

// TableName 设置表名
func (Favorite) TableName() string {
    return "favorite" // 设置为您想要的表名
}


// GetTotalSizeByUID 根据用户ID统计已使用的字节数
/*
func GetTotalSizeByUID(uid int) (uint64, error) {
    var totalSize uint64
    result := DB2.Debug().Model(&Favorite{}).Where("uid = ?", uid).Select("SUM(size)").Row().Scan(&totalSize)

    if result != nil {
        return 0, result // 返回错误
    }

    return totalSize, nil // 返回已使用的字节数
}*/


// GetTotalSizeByUID 根据用户ID统计已使用的字节数
/*
func GetTotalSizeByUID(uid int) (float64, error) {
    var totalSize float64 // 使用 sql.NullInt64 来处理 NULL 值
    err := DB.Debug().Model(&Favorite{}).Where("uid = ?", uid).Select("SUM(size)").Scan(&totalSize).Error
    return totalSize, err
}*/

/*
func GetTotalSizeByUID(uid int) (uint64, error) {
    var totalSize uint64
    err := DB2.Model(&Favorite{}).Where("uid = ?", uid).Pluck("SUM(size)", &totalSize).Error
    return totalSize, err
}
*/

func GetTotalSizeByUID(uid int) (float64, error) {
    type Result struct {
        TotalSize float64
    }

    var result Result
    err := DB2.Debug().Model(&Favorite{}).Where("uid = ?", uid).Select("SUM(size) as total_size").Scan(&result).Error

	util.Log().Info(fmt.Sprintf("Use Size 111: %f",result.TotalSize)) // 使用 fmt.Sprintf
    return result.TotalSize, err
}