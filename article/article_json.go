package article

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"

	"github.com/firefirestyle/go.miniprop"
)

//
func (obj *Article) ToMap() map[string]interface{} {
	return map[string]interface{}{
		TypeRootGroup:   obj.gaeObject.RootGroup,
		TypeUserName:    obj.gaeObject.UserName, //
		TypeTitle:       obj.gaeObject.Title,    //
		TypeTag:         obj.GetTags(),          //
		TypeCont:        obj.gaeObject.Cont,
		TypeInfo:        obj.gaeObject.Info,
		TypeType:        obj.gaeObject.Type,
		TypePointNames:  obj.gaeObject.PointNames,  //
		TypePointValues: obj.gaeObject.PointValues, //
		TypePropNames:   obj.gaeObject.PropNames,   //
		TypePropValues:  obj.gaeObject.PropValues,  //
		TypeSign:        obj.gaeObject.Sign,
		TypeArticleId:   obj.gaeObject.ArticleId,
		TypeCreated:     obj.gaeObject.Created.UnixNano(),
		TypeUpdated:     obj.gaeObject.Updated.UnixNano(),
		TypeSecretKey:   obj.gaeObject.SecretKey,
		TypeTarget:      obj.gaeObject.Target,
		TypeIconUrl:     obj.gaeObject.IconUrl,
	}
}
func (obj *Article) ToMapPublicOnly() map[string]interface{} {
	v := obj.ToMap()
	delete(v, TypeSecretKey)
	return v
}
func (obj *Article) ToJson() []byte {
	vv, _ := json.Marshal(obj.ToMap())
	return vv
}

func (obj *Article) ToJsonPublicOnly() []byte {
	v := obj.ToMap()
	delete(v, TypeSecretKey)
	vv, _ := json.Marshal(v)
	return vv
}

func (userObj *Article) SetParamFromsMap(v map[string]interface{}) error {
	propObj := miniprop.NewMiniPropFromMap(v)
	//
	userObj.gaeObject.RootGroup = propObj.GetString(TypeRootGroup, "")
	userObj.gaeObject.UserName = propObj.GetString(TypeUserName, "")
	userObj.gaeObject.Title = propObj.GetString(TypeTitle, "")
	userObj.SetTags(propObj.GetPropStringList("", TypeTag, make([]string, 0)))
	userObj.gaeObject.PointValues = propObj.GetPropFloatList("", TypePointValues, []float64{})
	userObj.gaeObject.PointNames = propObj.GetPropStringList("", TypePointNames, []string{})
	userObj.gaeObject.PropValues = propObj.GetPropStringList("", TypePropValues, []string{})
	userObj.gaeObject.PropNames = propObj.GetPropStringList("", TypePropNames, []string{})
	userObj.gaeObject.Cont = propObj.GetString(TypeCont, "")
	userObj.gaeObject.Info = propObj.GetString(TypeInfo, "")
	userObj.gaeObject.Type = propObj.GetString(TypeType, "")
	userObj.gaeObject.Sign = propObj.GetString(TypeSign, "")
	userObj.gaeObject.ArticleId = propObj.GetString(TypeArticleId, "")
	userObj.gaeObject.Created = propObj.GetTime(TypeCreated, time.Now()) //srcCreated
	userObj.gaeObject.Updated = propObj.GetTime(TypeUpdated, time.Now()) //srcLogin
	userObj.gaeObject.SecretKey = propObj.GetString(TypeSecretKey, "")
	userObj.gaeObject.Target = propObj.GetString(TypeTarget, "")
	userObj.gaeObject.IconUrl = propObj.GetString(TypeIconUrl, "")
	return nil
}

func (userObj *Article) SetParamFromsJson(ctx context.Context, source string) error {
	v := make(map[string]interface{})
	e := json.Unmarshal([]byte(source), &v)
	if e != nil {
		return e
	}
	//
	userObj.SetParamFromsMap(v)

	return nil
}
