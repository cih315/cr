package model

import (
    //"encoding/json"
    "time"
	"fmt"
	"github.com/cloudreve/Cloudreve/v3/pkg/util"
   // "github.com/cloudreve/Cloudreve/v3/pkg/serializer"
    "strconv"
	//"database/sql"
)

type CustomResponse struct {
    Code int                    `json:"code"`
    Data map[string]interface{} `json:"data"`
    Msg  string                 `json:"msg"`
}


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

// FavoriteResponse 用于返回的数据格式
type FavoriteResponse struct {
    Code int             `json:"code"`
    Data FavoriteDataSet `json:"data"`
    Msg  string          `json:"msg"`
}

// FavoriteDataSet 数据集
type FavoriteDataSet struct {
    Objects []FavoriteObject `json:"objects"`
    Parent  int              `json:"parent"`
}

// FavoriteObject 单个文件对象
type FavoriteObject struct {
    ID            string    `json:"id"`
    Name          string    `json:"name"`
    Path          string    `json:"path"`
    Thumb         bool      `json:"thumb"`
    Size          uint64    `json:"size"`
    Type          string    `json:"type"`
    Date          time.Time `json:"date"`
    CreateDate    time.Time `json:"create_date"`
    SourceEnabled bool      `json:"source_enabled"`
}

// GetFavoriteListByUID 根据uid获取收藏列表
func GetFavoriteListByUID(uid int) CustomResponse {
    var favorites []Favorite
    err := DB2.Where("uid = ?", uid).Find(&favorites).Error
    if err != nil {
        return CustomResponse{
            Code: 1,
            Data: map[string]interface{}{ // 将 FavoriteDataSet 转换为 map
                "objects": []FavoriteObject{}, // 空的对象列表
                "parent":  0,                  // 父级ID为0
            },
            Msg:  err.Error(),
        }
    }

    var objects []map[string]interface{}
    for _, fav := range favorites {
        objects = append(objects, map[string]interface{}{
            "id":             strconv.Itoa(fav.ID),
            "name":           fav.Name,
            "path":           "/",
            "thumb":          false,
            "size":           fav.Size,
            "type":           "file",
            "date":           fav.AddTime,
            "create_date":    time.Now(), // Replace with actual create time if available
            "source_enabled": false,
        })
    }

    // 如果 objects 为 nil，设置为默认值
    if objects == nil {
        objects = []map[string]interface{}{} // 设置为一个空的切片
    }
    
    return CustomResponse{
        Code: 0,
        Data: map[string]interface{}{
            "parent":  0,
            "objects": objects,
        },
        Msg: "",
    }

}


func GetTotalSizeByUID(uid int) (float64, error) {
    type Result struct {
        TotalSize float64
    }

    var result Result
    err := DB2.Model(&Favorite{}).Where("uid = ?", uid).Select("SUM(size) as total_size").Scan(&result).Error

	util.Log().Info(fmt.Sprintf("Use Size 111: %f",result.TotalSize)) // 使用 fmt.Sprintf
    return result.TotalSize, err
}

// GetFavoriteByID 根据ID获取指定的收藏数据行
func GetFavoriteByID(id int) (Favorite, error) {
    var favorite Favorite
    err := DB2.Model(&Favorite{}).Where("id = ?", id).First(&favorite).Error

    if err != nil {
        util.Log().Info(fmt.Sprintf("Error retrieving favorite with ID %d: %v", id, err))
        return favorite, err // Return the favorite and the error
    }

    util.Log().Info(fmt.Sprintf("Retrieved favorite: %+v", favorite)) // Log the retrieved favorite
    return favorite, nil // Return the favorite and nil error
}
