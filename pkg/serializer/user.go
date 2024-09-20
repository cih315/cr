package serializer

import (
	"fmt"
	"strconv"

	"github.com/cloudreve/Cloudreve/v3/pkg/util"
	model "github.com/cloudreve/Cloudreve/v3/models"
	"github.com/cloudreve/Cloudreve/v3/pkg/hashid"
	"github.com/duo-labs/webauthn/webauthn"
	"time"
)

// CheckLogin 检查登录
func CheckLogin() Response {
	return Response{
		Code: CodeCheckLogin,
		Msg:  "Login required",
	}
}

// MagnetUser 序列化器
type MagnetUser struct {
    UID           int       `json:"uid"`
    Email         string    `json:"email"`
    RegisterTime  time.Time `json:"register_time"`
    FileBan       int       `json:"file_ban"`
    Points        int       `json:"points"`
    MaxView1      int       `json:"max_view_1"`
    MaxView2      int       `json:"max_view_2"`
    VipMaxView1   int       `json:"vip_max_view_1"`
    VipMaxView2   int       `json:"vip_max_view_2"`
    VipEndTime    int       `json:"vip_end_time"`
    TotalPayment  float64   `json:"total_payment"`
    Tiyan         int       `json:"tiyan"`
    FavoriteSize  string    `json:"favorite_size"`
    RegUUID       string    `json:"reg_uuid"`
    RegIP         string    `json:"reg_ip"`
    LoginUUID     string    `json:"login_uuid"`
    LoginIP       string    `json:"login_ip"`
    LastActivity  time.Time `json:"last_activity_time"`
    RegChannel    string    `json:"reg_channel"`
    Remark        string    `json:"remark"`
    Token         string    `json:"token"`
}


// User 用户序列化器
type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"user_name"`
	Nickname       string    `json:"nickname"`
	Status         int       `json:"status"`
	Avatar         string    `json:"avatar"`
	CreatedAt      time.Time `json:"created_at"`
	PreferredTheme string    `json:"preferred_theme"`
	Anonymous      bool      `json:"anonymous"`
	Group          group     `json:"group"`
	Tags           []tag     `json:"tags"`
}

type group struct {
	ID                   uint   `json:"id"`
	Name                 string `json:"name"`
	AllowShare           bool   `json:"allowShare"`
	AllowRemoteDownload  bool   `json:"allowRemoteDownload"`
	AllowArchiveDownload bool   `json:"allowArchiveDownload"`
	ShareDownload        bool   `json:"shareDownload"`
	CompressEnabled      bool   `json:"compress"`
	WebDAVEnabled        bool   `json:"webdav"`
	SourceBatchSize      int    `json:"sourceBatch"`
	AdvanceDelete        bool   `json:"advanceDelete"`
	AllowWebDAVProxy     bool   `json:"allowWebDAVProxy"`
}

type tag struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	Color      string `json:"color"`
	Type       int    `json:"type"`
	Expression string `json:"expression"`
}

type storage struct {
	Used  uint64 `json:"used"`
	Free  uint64 `json:"free"`
	Total uint64 `json:"total"`
}

// WebAuthnCredentials 外部验证器凭证
type WebAuthnCredentials struct {
	ID          []byte `json:"id"`
	FingerPrint string `json:"fingerprint"`
}

// BuildWebAuthnList 构建设置页面凭证列表
func BuildWebAuthnList(credentials []webauthn.Credential) []WebAuthnCredentials {
	res := make([]WebAuthnCredentials, 0, len(credentials))
	for _, v := range credentials {
		credential := WebAuthnCredentials{
			ID:          v.ID,
			FingerPrint: fmt.Sprintf("% X", v.Authenticator.AAGUID),
		}
		res = append(res, credential)
	}

	return res
}

// BuildUser 序列化用户
func BuildUser(user model.User) User {
	tags, _ := model.GetTagsByUID(user.ID)
	return User{
		ID:             hashid.HashID(user.ID, hashid.UserID),
		Email:          user.Email,
		Nickname:       user.Nick,
		Status:         user.Status,
		Avatar:         user.Avatar,
		CreatedAt:      user.CreatedAt,
		PreferredTheme: user.OptionsSerialized.PreferredTheme,
		Anonymous:      user.IsAnonymous(),
		Group: group{
			ID:                   user.GroupID,
			Name:                 user.Group.Name,
			AllowShare:           user.Group.ShareEnabled,
			AllowRemoteDownload:  user.Group.OptionsSerialized.Aria2,
			AllowArchiveDownload: user.Group.OptionsSerialized.ArchiveDownload,
			ShareDownload:        user.Group.OptionsSerialized.ShareDownload,
			CompressEnabled:      user.Group.OptionsSerialized.ArchiveTask,
			WebDAVEnabled:        user.Group.WebDAVEnabled,
			AllowWebDAVProxy:     user.Group.OptionsSerialized.WebDAVProxy,
			SourceBatchSize:      user.Group.OptionsSerialized.SourceBatchSize,
			AdvanceDelete:        user.Group.OptionsSerialized.AdvanceDelete,
		},
		Tags: buildTagRes(tags),
	}
}

// BuildUserResponse 序列化用户响应
func BuildUserResponse(user model.User) Response {
	return Response{
		Data: BuildUser(user),
	}
}

// BuildUserStorageResponse 序列化用户存储概况响应
func BuildUserStorageResponse(user model.User) Response {
	total := user.Group.MaxStorage
	storageResp := storage{
		Used:  user.Storage,
		Free:  total - user.Storage,
		Total: total,
	}

	if total < user.Storage {
		storageResp.Free = 0
	}

	return Response{
		Data: storageResp,
	}
}

func BuildUserSpaceResponse(user model.User) Response {

	util.Log().Info(user.Email)

//	chk, _:= model.CheckEmail(user.Email)
	uinfo, _:= model.GetMagnetUserByEmail(user.Email)

    // 将布尔值 chk 转换为字符串
    util.Log().Info(fmt.Sprintf("UID result: %d", uinfo.UID)) // 使用 fmt.Sprintf
	util.Log().Info(fmt.Sprintf("Email check result: %s", uinfo.Email)) // 使用 fmt.Sprintf
	util.Log().Info(fmt.Sprintf("FavoriteSize: %s", uinfo.FavoriteSize)) // 使用 fmt.Sprintf
	size, _ := model.GetTotalSizeByUID(uinfo.UID)
	util.Log().Info(fmt.Sprintf("Use Size: %f", size)) // 使用 fmt.Sprintf

	 // 将 FavoriteSize 转换为 uint64
	 allSize, err := strconv.ParseUint(uinfo.FavoriteSize, 10, 64)
	 if err != nil {
		 // 处理错误，例如返回默认值或记录日志
		 allSize = 0 // 或者根据需要处理
	 }
 
	 useSize := uint64(size) // size 已经是 uint64 类型
 
	 storageResp := storage{
		 Used:  useSize,
		 Free:  allSize - useSize,
		 Total: allSize,
	 }

	return Response{
		Data: storageResp,
	}
}

// buildTagRes 构建标签列表
func buildTagRes(tags []model.Tag) []tag {
	res := make([]tag, 0, len(tags))
	for i := 0; i < len(tags); i++ {
		newTag := tag{
			ID:    hashid.HashID(tags[i].ID, hashid.TagID),
			Name:  tags[i].Name,
			Icon:  tags[i].Icon,
			Color: tags[i].Color,
			Type:  tags[i].Type,
		}
		if newTag.Type != 0 {
			newTag.Expression = tags[i].Expression

		}
		res = append(res, newTag)
	}

	return res
}
