package article

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"crypto/rand"

	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

func NewArticleManager(projectId string, kindArticle string, prefixOfId string, limitOfFinding int) *ArticleManager {
	ret := new(ArticleManager)
	ret.projectId = projectId
	ret.prefixOfId = prefixOfId
	ret.kindArticle = kindArticle
	ret.limitOfFinding = limitOfFinding
	return ret
}

func (obj *ArticleManager) NewArticleFromArticleId(ctx context.Context, articleId string, sign string) (*Article, error) {
	return obj.NewArticleFromGaeObjectKey(ctx, obj.NewGaeObjectKey(ctx, articleId, sign, ""))
}

func (obj *ArticleManager) NewArticleFromGaeObjectKey(ctx context.Context, key *datastore.Key) (*Article, error) {
	k := key
	//
	//
	artObjFMem, errNewFMem := obj.NewArticleFromMemcache(ctx, k.StringID())
	if errNewFMem == nil {
		log.Infof(ctx, ">>>> new article Obj from memcache")
		return artObjFMem, nil
	}
	//
	//
	var gaeObj GaeObjectArticle
	err := datastore.Get(ctx, k, &gaeObj)
	if err != nil {
		return nil, err
	}
	//
	//
	return obj.NewArticleFromGaeObject(ctx, k, &gaeObj), nil
}

func (obj *ArticleManager) NewArticleFromMemcache(ctx context.Context, stringId string) (*Article, error) {
	ret := new(Article)
	ret.gaeObject = new(GaeObjectArticle)
	idInfo := obj.ExtractInfoFromStringId(stringId)
	ret.gaeObjectKey = obj.NewGaeObjectKey(ctx, idInfo.ArticleId, idInfo.Sign, "")
	ret.kind = obj.kindArticle
	artObjSource, errGetFMem := memcache.Get(ctx, ret.gaeObjectKey.StringID())
	if errGetFMem != nil {
		return nil, errGetFMem
	}
	errSetParam := ret.SetParamFromsJson(ctx, string(artObjSource.Value))

	return ret, errSetParam
}

func (obj *ArticleManager) NewArticleFromGaeObject(ctx context.Context, gaeKey *datastore.Key, gaeObj *GaeObjectArticle) *Article {
	ret := new(Article)
	ret.gaeObject = gaeObj
	ret.gaeObjectKey = gaeKey
	ret.kind = obj.kindArticle
	//
	//

	return ret
}

func (obj *ArticleManager) NewArticle(ctx context.Context, userName string, sign string) *Article {
	created := time.Now()
	var secretKey string
	var articleId string
	var key *datastore.Key
	var art GaeObjectArticle
	for {
		secretKey = obj.makeRandomId() + obj.makeRandomId()
		articleId = obj.makeArticleId(userName, created, secretKey)
		//stringId := obj.makeStringId(articleId, sign)
		//
		Debug(ctx, "<NewArticle>"+articleId+"::"+sign)
		key = obj.NewGaeObjectKey(ctx, articleId, sign, "")
		err := datastore.Get(ctx, key, &art)
		if err != nil {
			break
		}
	}
	//
	ret := new(Article)
	ret.kind = obj.kindArticle
	ret.gaeObject = &art
	ret.gaeObjectKey = key
	ret.gaeObject.ProjectId = obj.projectId
	ret.gaeObject.UserName = userName
	ret.gaeObject.Sign = sign
	ret.gaeObject.Created = created
	ret.gaeObject.Updated = created
	ret.gaeObject.SecretKey = secretKey
	ret.gaeObject.ArticleId = articleId
	//
	//
	//

	return ret
}

func (obj *ArticleManager) NewGaeObjectKey(ctx context.Context, articleId string, sign string, kind string) *datastore.Key {
	if kind == "" {
		kind = obj.kindArticle
	}
	return datastore.NewKey(ctx, kind, obj.makeStringId(articleId, sign), 0, nil)
}

func (obj *ArticleManager) GetKind() string {
	return obj.kindArticle
}

func (obj *ArticleManager) makeArticleId(userName string, created time.Time, secretKey string) string {
	hashKey := obj.hashStr(fmt.Sprintf("%sv1e%s%s%s%d", obj.prefixOfId, secretKey, userName, userName, created.UnixNano()))
	userName64 := base64.StdEncoding.EncodeToString([]byte(userName))
	return "" + obj.prefixOfId + "v1e" + hashKey + userName64
}

func (obj *ArticleManager) makeStringId(articleId string, sign string) string {
	propObj := miniprop.NewMiniProp()
	propObj.SetString("i", articleId)
	propObj.SetString("s", sign)
	return string(propObj.ToJson())
}

type StringIdInfo struct {
	ArticleId string
	Sign      string
}

func (obj *ArticleManager) ExtractInfoFromStringId(stringId string) *StringIdInfo {
	propObj := miniprop.NewMiniPropFromJson([]byte(stringId))
	return &StringIdInfo{
		ArticleId: propObj.GetString("i", ""),
		Sign:      propObj.GetString("s", ""),
	}
}

func (obj *ArticleManager) hash(v string) string {
	sha1Obj := sha1.New()
	sha1Obj.Write([]byte(v))
	return string(sha1Obj.Sum(nil))
}

func (obj *ArticleManager) hashStr(v string) string {
	sha1Obj := sha1.New()
	sha1Obj.Write([]byte(v))
	return string(base64.StdEncoding.EncodeToString(sha1Obj.Sum(nil)))
}

func (obj *ArticleManager) makeRandomId() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}

func (obj *ArticleManager) SaveOnDB(ctx context.Context, artObj *Article) error {
	return artObj.saveOnDB(ctx)
}

func (obj *ArticleManager) SaveUsrWithImmutable(ctx context.Context, artObj *Article) error {
	sign := strconv.Itoa(time.Now().Nanosecond())
	nextArObj, nextArtErr := obj.NewArticleFromArticleId(ctx, artObj.GetArticleId(), sign)
	if nextArtErr != nil {
		return nextArtErr
	}
	usrObjData := artObj.ToMap()
	usrObjData[TypeSign] = sign
	usrObjData[TypeUpdated] = artObj.GetUpdated().UnixNano()
	nextArObj.SetParamFromsMap(usrObjData)
	saveErr := nextArObj.saveOnDB(ctx)
	if saveErr != nil {
		return saveErr
	}
	obj.DeleteFromArticleId(ctx, artObj.GetArticleId(), artObj.GetParentId())
	return nil
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
