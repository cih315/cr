package share

import (
	"net/url"
	"time"
	"github.com/cloudreve/Cloudreve/v3/pkg/util"
	"fmt"
	"strconv"

	model "github.com/cloudreve/Cloudreve/v3/models"
	"github.com/cloudreve/Cloudreve/v3/pkg/hashid"
	"github.com/cloudreve/Cloudreve/v3/pkg/serializer"
	"github.com/gin-gonic/gin"
)

// ShareCreateService 创建新分享服务
type ShareCreateService struct {
	SourceID        string `json:"id" binding:"required"`
	IsDir           bool   `json:"is_dir"`
	Password        string `json:"password" binding:"max=255"`
	RemainDownloads int    `json:"downloads"`
	Expire          int    `json:"expire"`
	Preview         bool   `json:"preview"`
}

// ShareUpdateService 分享更新服务
type ShareUpdateService struct {
	Prop  string `json:"prop" binding:"required,eq=password|eq=preview_enabled"`
	Value string `json:"value" binding:"max=255"`
}

// Delete 删除分享
func (service *Service) Delete(c *gin.Context, user *model.User) serializer.Response {
	share := model.GetShareByHashID(c.Param("id"))
	if share == nil || share.Creator().ID != user.ID {
		return serializer.Err(serializer.CodeShareLinkNotFound, "", nil)
	}

	if err := share.Delete(); err != nil {
		return serializer.DBErr("Failed to delete share record", err)
	}

	return serializer.Response{}
}



// Update 更新分享属性
func (service *ShareUpdateService) Update(c *gin.Context) serializer.Response {
	shareCtx, _ := c.Get("share")
	share := shareCtx.(*model.Share)

	switch service.Prop {
	case "password":
		err := share.Update(map[string]interface{}{"password": service.Value})
		if err != nil {
			return serializer.DBErr("Failed to update share record", err)
		}
	case "preview_enabled":
		value := service.Value == "true"
		err := share.Update(map[string]interface{}{"preview_enabled": value})
		if err != nil {
			return serializer.DBErr("Failed to update share record", err)
		}
		return serializer.Response{
			Data: value,
		}
	}
	return serializer.Response{
		Data: service.Value,
	}
}

func processSourceID(sourceID interface{}) {
	switch v := sourceID.(type) {
	case string:
		fmt.Println("SourceID is a string:", v)
	case int:
		fmt.Println("SourceID is an int:", v)
	default:
		fmt.Println("Unknown type for SourceID")
	}
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}


// Create 创建新分享
func (service *ShareCreateService) Create(c *gin.Context) serializer.Response {
	userCtx, _ := c.Get("user")
	user := userCtx.(*model.User)

	// 是否拥有权限
	if !user.Group.ShareEnabled {
		return serializer.Err(serializer.CodeGroupNotAllowed, "", nil)
	}
	
	
   // 源对象真实ID
   var (
	sourceID   uint
	sourceName string
	err        error
	)

	// Convert service.SourceID from string to uint64
	//parsedID, err := strconv.ParseUint(service.SourceID, 10, 32) // Convert to uint64
	//util.Log().Info(fmt.Sprintf("parsedID: %s", 
	//service.SourceID))


	if(isNumeric(service.SourceID)){
		fmt.Printf("%s is a number.\n", service.SourceID)

		

		
		SourceID64, err := strconv.ParseUint(service.SourceID, 10, 32)

		SourceID := uint(SourceID64)

		fav, err := model.GetFavoriteByID(int(SourceID))
		sourceName = fav.Name
		
		newShare := model.Share{
			Password:        service.Password,
			IsDir:           service.IsDir,
			IsMagnet:		 true,
			UserID:          user.ID,
			SourceID:        SourceID,
			RemainDownloads: -1,
			PreviewEnabled:  service.Preview,
			SourceName:      sourceName,
		}

		// 如果开启了自动过期
		if service.RemainDownloads > 0 {
			expires := time.Now().Add(time.Duration(service.Expire) * time.Second)
			newShare.RemainDownloads = service.RemainDownloads
			newShare.Expires = &expires
		}

		// 创建分享
		id, err := newShare.Create()
		if err != nil {
			return serializer.DBErr("Failed to create share link record", err)
		}

		// 获取分享的唯一id
		uid := hashid.HashID(id, hashid.ShareID)
		// 最终得到分享链接
		siteURL := model.GetSiteURL()
		sharePath, _ := url.Parse("/s/" + uid)
		shareURL := siteURL.ResolveReference(sharePath)

		return serializer.Response{
			Code: 0,
			Data: shareURL.String(),
		}


	} else {
		fmt.Printf("%s is a string.\n", service.SourceID)

		if service.IsDir {
			sourceID, err = hashid.DecodeHashID(service.SourceID, hashid.FolderID)
		} else {
			sourceID, err = hashid.DecodeHashID(service.SourceID, hashid.FileID)
		}
		if err != nil {
			return serializer.Err(serializer.CodeNotFound, "", nil)
		}

	}

	
/*
	if err != nil {
		return serializer.Err(serializer.CodeNotFound, "", nil)
	}
*/
	// Convert uint64 to uint
	//sourceID = uint(parsedID)


	exist := true

	// Use sourceID in the function call
	file, err := model.GetFilesByIDs([]uint{sourceID}, user.ID) // Now sourceID is of type uint
	/*
	if err != nil {
		return serializer.Err(serializer.CodeNotFound, "", nil)
	}*/

	if err != nil || len(file) == 0 {
		exist = false
		util.Log().Info(fmt.Sprintf("ShareCreateService Create - SourceID: %s, IsDir: %t, not found!", 
		service.SourceID, service.IsDir))

		// 网盘信息库里没找到，就去磁力那边找。

	} else {

		// 在网盘 file 表找到了,则原流程代码处理。

		util.Log().Info("zhao dao le ")


		sourceName = file[0].Name
		util.Log().Info(fmt.Sprintf("ShareCreateService Create - SourceID: %s, IsDir: %t, sourceName: %s!", 
		service.SourceID, service.IsDir, sourceName))



			if service.IsDir {
				sourceID, err = hashid.DecodeHashID(service.SourceID, hashid.FolderID)
			} else {
				sourceID, err = hashid.DecodeHashID(service.SourceID, hashid.FileID)
			}
			if err != nil {
				return serializer.Err(serializer.CodeNotFound, "", nil)
			}


			if service.IsDir {
				folder, err := model.GetFoldersByIDs([]uint{sourceID}, user.ID)
				if err != nil || len(folder) == 0 {
					exist = false
				} else {
					sourceName = folder[0].Name
				}
			} else {
				file, err := model.GetFilesByIDs([]uint{sourceID}, user.ID)
				if err != nil || len(file) == 0 {
					exist = false
				} else {
					sourceName = file[0].Name
				}
			}



			if !exist {
				return serializer.Err(serializer.CodeNotFound, "", nil)
			}

			newShare := model.Share{
				Password:        service.Password,
				IsDir:           service.IsDir,
				UserID:          user.ID,
				SourceID:        sourceID,
				RemainDownloads: -1,
				PreviewEnabled:  service.Preview,
				SourceName:      sourceName,
			}

			// 如果开启了自动过期
			if service.RemainDownloads > 0 {
				expires := time.Now().Add(time.Duration(service.Expire) * time.Second)
				newShare.RemainDownloads = service.RemainDownloads
				newShare.Expires = &expires
			}

			// 创建分享
			id, err := newShare.Create()
			if err != nil {
				return serializer.DBErr("Failed to create share link record", err)
			}

			// 获取分享的唯一id
			uid := hashid.HashID(id, hashid.ShareID)
			// 最终得到分享链接
			siteURL := model.GetSiteURL()
			sharePath, _ := url.Parse("/s/" + uid)
			shareURL := siteURL.ResolveReference(sharePath)

			return serializer.Response{
				Code: 0,
				Data: shareURL.String(),
			}

		
		// 网盘信息库里找到了，还要拿这个  dourceID 和 sourceName 去磁力那边查询。 
		//util.Log().Info(fmt.Sprintf("ShareCreateService Create - SourceID: %s, IsDir: %t, sourceName: %s", 
	//service.SourceID, service.IsDir, sourceName))


	}

	util.Log().Info("default return..")
	return serializer.Response{
		Code: 0,
		Data: "http://www.default.com/t=1",
	}
}
