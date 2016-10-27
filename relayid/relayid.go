package relayid

import (
	"time"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

const (
	TypeTwitter  = "twitter"
	TypeFacebook = "facebook"
	TypePointer  = "pointer"
)

const (
	TypeRelayIdProjectId = "ProjectId"
	TypeRelayIdName      = "Name"
	TypeRelayIdId        = "Id"
	TypeRelayIdType      = "Type"
	TypeRelayIdUserName  = "UserName"
	TypeRelayIdInfo      = "Info"
	TypeRelayIdUpdate    = "Update"
	TypeRelayIdSign      = "Sign"
)

type GaeRelayIdItem struct {
	ProjectId    string
	IdentifyName string
	IdentifyId   string
	Type         string
	UserName     string
	Info         string
	Update       time.Time
	Sign         string
}

type RelayIdManagerConfig struct {
	Kind      string
	ProjectId string
}

type RelayId struct {
	gaeObj *GaeRelayIdItem
	gaeKey *datastore.Key
	kind   string
}

type RelayIdManager struct {
	kind      string
	projectId string
}

func NewRelayIdManager(config RelayIdManagerConfig) *RelayIdManager {
	return &RelayIdManager{
		kind:      config.Kind,
		projectId: config.ProjectId,
	}
}

func (obj *RelayId) ToJson() []byte {
	propObj := miniprop.NewMiniProp()
	propObj.SetString(TypeRelayIdProjectId, obj.gaeObj.ProjectId)
	propObj.SetString(TypeRelayIdName, obj.gaeObj.IdentifyName)
	propObj.SetString(TypeRelayIdId, obj.gaeObj.IdentifyId)
	propObj.SetString(TypeRelayIdUserName, obj.gaeObj.UserName)
	propObj.SetString(TypeRelayIdInfo, obj.gaeObj.Info)
	propObj.SetTime(TypeRelayIdUpdate, obj.gaeObj.Update)
	propObj.SetString(TypeRelayIdSign, obj.gaeObj.Sign)
	return propObj.ToJson()
}

func (obj *RelayId) SetValueFromJson(data []byte) {
	propObj := miniprop.NewMiniPropFromJson(data)
	obj.gaeObj.ProjectId = propObj.GetString(TypeRelayIdProjectId, "")
	obj.gaeObj.IdentifyName = propObj.GetString(TypeRelayIdName, "")
	obj.gaeObj.IdentifyId = propObj.GetString(TypeRelayIdId, "")
	obj.gaeObj.UserName = propObj.GetString(TypeRelayIdUserName, "")
	obj.gaeObj.Info = propObj.GetString(TypeRelayIdInfo, "")
	obj.gaeObj.Update = propObj.GetTime(TypeRelayIdUpdate, time.Now())
	obj.gaeObj.Sign = propObj.GetString(TypeRelayIdSign, "")
}

func (obj *RelayId) UpdateMemcache(ctx context.Context) {
	userObjMemSource := obj.ToJson()
	userObjMem := &memcache.Item{
		Key:   obj.gaeKey.StringID(),
		Value: []byte(userObjMemSource), //
	}
	memcache.Set(ctx, userObjMem)
}

func (obj *RelayId) GetName() string {
	return obj.gaeObj.IdentifyName
}

func (obj *RelayId) GetId() string {
	return obj.gaeObj.IdentifyId
}
func (obj *RelayId) GetType() string {
	return obj.gaeObj.Type
}

func (obj *RelayId) GetUserName() string {
	return obj.gaeObj.UserName
}

func (obj *RelayId) SetUserName(v string) {
	obj.gaeObj.UserName = v
}

func (obj *RelayId) GetSign() string {
	return obj.gaeObj.Sign
}

func (obj *RelayId) SetSign(v string) {
	obj.gaeObj.Sign = v
}

func (obj *RelayId) GetInfo() string {
	return obj.gaeObj.Info
}

func (obj *RelayId) GetUpdate() time.Time {
	return obj.gaeObj.Update
}

func (obj *RelayId) Save(ctx context.Context) error {
	_, err := datastore.Put(ctx, obj.gaeKey, obj.gaeObj)
	if err == nil {
		obj.UpdateMemcache(ctx)
	}
	return err
}

func (obj *RelayIdManager) GetRelayIdForTwitter(ctx context.Context, screenName string, userId string, oauthToken string) *RelayId {
	return obj.GetRelayIdWithNew(ctx, screenName, userId, TypeTwitter, map[string]string{"token": oauthToken})
}
func (obj *RelayIdManager) GetRelayIdAsPointer(ctx context.Context, userName string) *RelayId {
	//	Debug(ctx, "GetRelayIdAsPointer ==> "+userName)
	return obj.GetRelayIdWithNew(ctx, userName, userName, TypePointer, map[string]string{})
}

func (obj *RelayIdManager) NewRelayId(ctx context.Context, screenName string, //
	userId string, identifyType string, infos map[string]string) *RelayId {
	gaeKey := obj.NewRelayIdGaeKey(ctx, userId, identifyType)
	gaeObj := GaeRelayIdItem{
		IdentifyName: screenName,
		IdentifyId:   userId,
		Type:         TypeTwitter,
		ProjectId:    obj.projectId,
	}
	propObj := miniprop.NewMiniPropFromJson([]byte(gaeObj.Info))
	for k, v := range infos {
		propObj.SetString(k, v)
	}
	gaeObj.Info = string(propObj.ToJson())
	gaeObj.Update = time.Now()
	return &RelayId{
		gaeObj: &gaeObj,
		gaeKey: gaeKey,
		kind:   obj.kind,
	}
}

func (obj *RelayIdManager) GetRelayId(ctx context.Context, identify string, identifyType string) (*RelayId, error) {
	gaeKey := obj.NewRelayIdGaeKey(ctx, identify, identifyType)
	gaeObj := GaeRelayIdItem{}

	//
	// mem
	memItemObj, errMemObj := memcache.Get(ctx, obj.MakeRelayIdStringId(identify, identifyType))
	if errMemObj == nil {
		ret := &RelayId{
			gaeObj: &gaeObj,
			gaeKey: gaeKey,
			kind:   obj.kind,
		}
		ret.SetValueFromJson(memItemObj.Value)
		return ret, nil
	}
	//
	// db
	err := datastore.Get(ctx, gaeKey, &gaeObj)
	if err != nil {
		return nil, err
	}
	ret := &RelayId{
		gaeObj: &gaeObj,
		gaeKey: gaeKey,
		kind:   obj.kind,
	}
	//
	//
	ret.UpdateMemcache(ctx)
	return ret, nil
}

func (obj *RelayIdManager) GetRelayIdWithNew(ctx context.Context, screenName string, userId string, userIdType string, infos map[string]string) *RelayId {
	relayObj, err := obj.GetRelayId(ctx, userId, userIdType)
	if err != nil {
		relayObj = obj.NewRelayId(ctx, screenName, userId, userIdType, infos)
	}
	//
	propObj := miniprop.NewMiniPropFromJson([]byte(relayObj.gaeObj.Info))
	for k, v := range infos {
		propObj.SetString(k, v)
	}
	relayObj.gaeObj.Info = string(propObj.ToJson())
	relayObj.gaeObj.Update = time.Now()
	return relayObj
}

func (obj *RelayIdManager) NewRelayIdGaeKey(ctx context.Context, identify string, identifyType string) *datastore.Key {
	return datastore.NewKey(ctx, obj.kind, obj.MakeRelayIdStringId(identify, identifyType), 0, nil)
}

func (obj *RelayIdManager) MakeRelayIdStringId(identify string, identifyType string) string {
	return obj.kind + ":" + obj.projectId + ":" + identifyType + ":" + identify
}
