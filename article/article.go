package article

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	//	"google.golang.org/appengine/log"
	//"google.golang.org/appengine/blobstore"
	"errors"

	"github.com/firefirestyle/go.miniprop"
	"google.golang.org/appengine/memcache"
)

const (
	ArticleStatePublic  = "public"
	ArticleStatePrivate = "private"
	ArticleStateAll     = ""
)

type GaeObjectArticle struct {
	ProjectId string
	UserName  string
	Title     string `datastore:",noindex"`
	Tag       string `datastore:",noindex"`
	Cont      string `datastore:",noindex"`
	Info      string `datastore:",noindex"`
	Type      string
	Sign      string
	ArticleId string
	Created   time.Time
	Updated   time.Time
	SecretKey string `datastore:",noindex"`
	Target    string
}

type Article struct {
	gaeObjectKey *datastore.Key
	gaeObject    *GaeObjectArticle
	kind         string
}

type ArticleManager struct {
	projectId      string
	prefixOfId     string
	kindArticle    string
	limitOfFinding int
}

const (
	TypeProjectId = "ProjectId"
	TypeUserName  = "UserName"
	TypeTitle     = "Title"
	TypeTag       = "Tag"
	TypeCont      = "Cont"
	TypeInfo      = "Info"
	TypeType      = "Type"
	TypeSign      = "Sign"
	TypeArticleId = "ArticleId"
	TypeCreated   = "Created"
	TypeUpdated   = "Updated"
	TypeSecretKey = "SecretKey"
	TypeTarget    = "Target"
)

func (obj *Article) updateMemcache(ctx context.Context) error {
	userObjMemSource := obj.ToJson()
	userObjMem := &memcache.Item{
		Key:   obj.gaeObjectKey.StringID(),
		Value: []byte(userObjMemSource), //
	}
	memcache.Set(ctx, userObjMem)
	return nil
}

//
func getStringFromProp(requestPropery map[string]interface{}, key string, defaultValue string) string {
	v := requestPropery[key]
	if v == nil {
		return defaultValue
	} else {
		return v.(string)
	}
}

//
func (obj *Article) ToMap() map[string]interface{} {
	return map[string]interface{}{
		TypeProjectId: obj.gaeObject.ProjectId,
		TypeUserName:  obj.gaeObject.UserName, //
		TypeTitle:     obj.gaeObject.Title,    //
		TypeTag:       obj.gaeObject.Tag,      //
		TypeCont:      obj.gaeObject.Cont,
		TypeInfo:      obj.gaeObject.Info,
		TypeType:      obj.gaeObject.Type,
		TypeSign:      obj.gaeObject.Sign,
		TypeArticleId: obj.gaeObject.ArticleId,
		TypeCreated:   obj.gaeObject.Created.UnixNano(),
		TypeUpdated:   obj.gaeObject.Updated.UnixNano(),
		TypeSecretKey: obj.gaeObject.SecretKey,
		TypeTarget:    obj.gaeObject.Target,
	}
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

//
// func (userObj *User) SetUserFromsMap(ctx context.Context, v map[string]interface{}) {
//	propObj := miniprop.NewMiniPropFromMap(v)
//
func (userObj *Article) SetParamFromsMap(v map[string]interface{}) error {
	propObj := miniprop.NewMiniPropFromMap(v)
	//
	userObj.gaeObject.ProjectId = propObj.GetString(TypeProjectId, "")
	userObj.gaeObject.UserName = propObj.GetString(TypeUserName, "")
	userObj.gaeObject.Title = propObj.GetString(TypeTitle, "")
	userObj.gaeObject.Tag = propObj.GetString(TypeTag, "")
	userObj.gaeObject.Cont = propObj.GetString(TypeCont, "")
	userObj.gaeObject.Info = propObj.GetString(TypeInfo, "")
	userObj.gaeObject.Type = propObj.GetString(TypeType, "")
	userObj.gaeObject.Sign = propObj.GetString(TypeSign, "")
	userObj.gaeObject.ArticleId = propObj.GetString(TypeArticleId, "")
	userObj.gaeObject.Created = propObj.GetTime(TypeCreated, time.Now()) //srcCreated
	userObj.gaeObject.Updated = propObj.GetTime(TypeUpdated, time.Now()) //srcLogin
	userObj.gaeObject.SecretKey = propObj.GetString(TypeSecretKey, "")
	userObj.gaeObject.Target = propObj.GetString(TypeTarget, "")

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

//
//
//
func (obj *Article) GetGaeObjectKind() string {
	return obj.kind
}

func (obj *Article) GetGaeObjectKey() *datastore.Key {
	return obj.gaeObjectKey
}

func (obj *Article) GetUserName() string {
	return obj.gaeObject.UserName
}

func (obj *Article) GetInfo() string {
	return obj.gaeObject.Info
}

func (obj *Article) SetInfo(v string) {
	obj.gaeObject.Info = v
}

func (obj *Article) GetTarget() string {
	return obj.gaeObject.Target
}

func (obj *Article) SetTarget(v string) {
	obj.gaeObject.Target = v
}

func (obj *Article) SetUserName(v string) {
	obj.gaeObject.UserName = v
}

func (obj *Article) GetTitle() string {
	return obj.gaeObject.Title
}

func (obj *Article) SetTitle(v string) {
	obj.gaeObject.Title = v
}

func (obj *Article) GetTags() []string {
	var tags []string
	json.Unmarshal([]byte(obj.gaeObject.Tag), &tags)
	return tags
}

func (obj *Article) SetTags(v []string) {
	if v == nil || len(v) == 0 {
		obj.gaeObject.Tag = ""
	} else {
		b, _ := json.Marshal(v)
		obj.gaeObject.Tag = string(b)
	}
}

func (obj *Article) GetCont() string {
	return obj.gaeObject.Cont
}

func (obj *Article) SetCont(v string) {
	obj.gaeObject.Cont = v
}

func (obj *Article) GetState() string {
	return obj.gaeObject.Type
}

func (obj *Article) SetState(v string) {
	obj.gaeObject.Type = v
}

func (obj *Article) GetParentId() string {
	return obj.gaeObject.Sign
}

func (obj *Article) SetParentId(v string) {
	obj.gaeObject.Sign = v
}

func (obj *Article) GetArticleId() string {
	return obj.gaeObject.ArticleId
}

func (obj *Article) GetCreated() time.Time {
	return obj.gaeObject.Created
}

func (obj *Article) GetUpdated() time.Time {
	return obj.gaeObject.Updated
}

func (obj *Article) SetUpdated(v time.Time) {
	obj.gaeObject.Updated = v
}

//
//
//
// ArticleManager

func (obj *Article) saveOnDB(ctx context.Context) error {
	Debug(ctx, "<saveOnDB> "+obj.gaeObjectKey.StringID())
	_, err := datastore.Put(ctx, obj.gaeObjectKey, obj.gaeObject)
	obj.updateMemcache(ctx)
	return err
}

func (mgrObj *ArticleManager) SaveOnOtherDB(ctx context.Context, obj *Article, kind string) error {
	_, err := datastore.Put(ctx, mgrObj.NewGaeObjectKey(ctx, obj.GetArticleId(), obj.gaeObject.Sign, kind), obj.gaeObject)
	return err
}

func (mgrObj *ArticleManager) DeleteFromArticleId(ctx context.Context, articleId string, sign string) error {
	key := mgrObj.NewGaeObjectKey(ctx, articleId, sign, mgrObj.GetKind())
	memcache.Delete(ctx, key.StringID())
	return datastore.Delete(ctx, mgrObj.NewGaeObjectKey(ctx, articleId, sign, mgrObj.GetKind()))
}

func (obj *ArticleManager) newCursorFromSrc(cursorSrc string) *datastore.Cursor {
	c1, e := datastore.DecodeCursor(cursorSrc)
	if e != nil {
		return nil
	} else {
		return &c1
	}
}

func (obj *ArticleManager) makeCursorSrc(founds *datastore.Iterator) string {
	c, e := founds.Cursor()
	if e == nil {
		return c.String()
	} else {
		return ""
	}
}

func (obj *ArticleManager) GetArticleFromArticleIdOnQuery(ctx context.Context, articleId string) (*Article, error) {
	q := datastore.NewQuery(obj.kindArticle)
	q = q.Filter("ProjectId =", obj.projectId)
	q = q.Filter("ArticleId =", articleId)
	arts := obj.FindArticleFromQuery(ctx, q, "", false)
	if len(arts.Articles) > 0 {
		Debug(ctx, "======> A")
		return arts.Articles[0], nil
	} else {
		Debug(ctx, "======> B : "+obj.projectId+" : "+articleId)
		return nil, errors.New("--")
	}
}
