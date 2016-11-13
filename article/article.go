package article

import (
	//	"encoding/json"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	//	"google.golang.org/appengine/log"
	//"google.golang.org/appengine/blobstore"
	//	"errors"

	"github.com/firefirestyle/go.minipointer"
	"github.com/firefirestyle/go.miniprop"
	"google.golang.org/appengine/memcache"
)

type GaeObjectArticle struct {
	RootGroup   string
	UserName    string
	Title       string    `datastore:",noindex"`
	Tags        []string  `datastore:"Tags.Tag"`
	PointNames  []string  `datastore:"Points.Name"`
	PointValues []float64 `datastore:"Points.Value"`
	PropNames   []string  `datastore:"Props.Name"`
	PropValues  []string  `datastore:"Props.Value"`
	Cont        string    `datastore:",noindex"`
	Info        string    `datastore:",noindex"`
	Sign        string    `datastore:",noindex"`
	ArticleId   string
	Created     time.Time
	Updated     time.Time
	SecretKey   string `datastore:",noindex"`
	IconUrl     string `datastore:",noindex"`
}

type Article struct {
	gaeObjectKey *datastore.Key
	gaeObject    *GaeObjectArticle
	kind         string
}

const (
	TypeRootGroup   = "RootGroup"
	TypeUserName    = "UserName"
	TypeTitle       = "Title"
	TypeTag         = "Tag"
	TypePointNames  = "PointNames"
	TypePointValues = "PointValues"
	TypePropNames   = "PropNames"
	TypePropValues  = "PropValues"
	TypeCont        = "Cont"
	TypeInfo        = "Info"
	TypeType        = "Type"
	TypeSign        = "Sign"
	TypeArticleId   = "ArticleId"
	TypeCreated     = "Created"
	TypeUpdated     = "Updated"
	TypeSecretKey   = "SecretKey"
	TypeTarget      = "Target"
	TypeIconUrl     = "IconUrl"
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

func (obj *Article) GetSign() string {
	return obj.gaeObject.Sign
}

func (obj *Article) GetIconUrl() string {
	return obj.gaeObject.IconUrl
}

func (obj *Article) SetIconUrl(v string) {
	obj.gaeObject.IconUrl = v
}

func (obj *Article) GetInfo() string {
	return obj.gaeObject.Info
}

func (obj *Article) SetInfo(v string) {
	obj.gaeObject.Info = v
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
	ret := make([]string, 0)
	for _, v := range obj.gaeObject.Tags {
		//		ret = append(ret, v.Tag)
		ret = append(ret, v)
	}
	return ret
}

func (obj *Article) SetTags(vs []string) {
	//	obj.gaeObject.Tags = make([]Tag, 0)
	obj.gaeObject.Tags = make([]string, 0)
	for _, v := range vs {
		//		obj.gaeObject.Tags = append(obj.gaeObject.Tags, Tag{Tag: v})
		obj.gaeObject.Tags = append(obj.gaeObject.Tags, v)
	}
}

func (obj *Article) GetCont() string {
	return obj.gaeObject.Cont
}

func (obj *Article) SetCont(v string) {
	obj.gaeObject.Cont = v
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
func (obj *Article) GetPoint(name string) float64 {
	index := -1
	for i, v := range obj.gaeObject.PointNames {
		if v == name {
			index = i
			break
		}
	}
	if index < 0 {
		return 0
	}
	return obj.gaeObject.PointValues[index]
}

func (obj *Article) SetPoint(name string, v float64) {
	index := -1
	for i, iv := range obj.gaeObject.PointNames {
		if iv == name {
			index = i
			break
		}
	}
	if index == -1 {
		obj.gaeObject.PointValues = append(obj.gaeObject.PointValues, v)
		obj.gaeObject.PointNames = append(obj.gaeObject.PointNames, name)
	} else {
		obj.gaeObject.PointValues[index] = v
	}
}

func (obj *Article) GetProp(name string) string {
	index := -1
	for i, v := range obj.gaeObject.PropNames {
		if v == name {
			index = i
			break
		}
	}
	if index < 0 {
		return ""
	}
	p := miniprop.NewMiniPropFromJson([]byte(obj.gaeObject.PropValues[index]))
	return p.GetString(name, "")
}

func (obj *Article) SetProp(name, v string) {
	index := -1
	p := miniprop.NewMiniProp()
	p.SetString(name, v)
	v = string(p.ToJson())
	for i, iv := range obj.gaeObject.PropNames {
		if iv == name {
			index = i
			break
		}
	}
	if index == -1 {
		obj.gaeObject.PropValues = append(obj.gaeObject.PropValues, v)
		obj.gaeObject.PropNames = append(obj.gaeObject.PropNames, name)
	} else {
		obj.gaeObject.PropValues[index] = v
	}
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

func (mgrObj *ArticleManager) DeleteFromArticleIdWithPointer(ctx context.Context, articleId string) error {
	artObj, pointerObj, _ := mgrObj.GetArticleFromPointer(ctx, articleId)
	if artObj != nil {
		Debug(ctx, "===> art DEL")
		deleteErr := mgrObj.DeleteFromArticleId(ctx, articleId, pointerObj.GetSign())
		if deleteErr != nil {
			return deleteErr
		}
	}
	if pointerObj != nil {
		Debug(ctx, "===> pointer DEL")
		return mgrObj.pointerMgr.DeletePointerFromObj(ctx, pointerObj)
	}
	return nil
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

func (obj *ArticleManager) GetArticleFromPointer(ctx context.Context, articleId string) (*Article, *minipointer.Pointer, error) {
	pointerObj, pointerErr := obj.pointerMgr.GetPointer(ctx, articleId, minipointer.TypePointer)
	if pointerErr != nil {
		Debug(ctx, "---> pointer")
		return nil, nil, pointerErr
	}
	pointerArticleId := pointerObj.GetValue()
	pointerSign := pointerObj.GetSign()
	Debug(ctx, "---> pointer "+":"+pointerArticleId+":"+pointerSign+":")

	artObj, artErr := obj.GetArticleFromArticleId(ctx, pointerArticleId, pointerSign)
	return artObj, pointerObj, artErr
}

func (obj *ArticleManager) GetPointerFromArticleId(ctx context.Context, articleId string) (*minipointer.Pointer, error) {
	return obj.pointerMgr.GetPointer(ctx, articleId, minipointer.TypePointer)
}

/*
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
*/
