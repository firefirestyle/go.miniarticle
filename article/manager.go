package article

import (
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"crypto/rand"

	"github.com/firefirestyle/go.minipointer"
	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

type ArticleManager struct {
	projectId      string
	prefixOfId     string
	kindArticle    string
	kindPointer    string
	pointerMgr     *minipointer.PointerManager
	limitOfFinding int
}

func NewArticleManager(projectId string, kindArticle string, kindPointer string, prefixOfId string, limitOfFinding int) *ArticleManager {
	ret := new(ArticleManager)
	ret.projectId = projectId
	ret.prefixOfId = prefixOfId
	ret.kindArticle = kindArticle
	ret.kindPointer = kindPointer
	ret.limitOfFinding = limitOfFinding
	ret.pointerMgr = minipointer.NewPointerManager(minipointer.PointerManagerConfig{
		RootGroup: projectId,
		Kind:      kindPointer,
	})
	return ret
}

func (obj *ArticleManager) GetKind() string {
	return obj.kindArticle
}

func (obj *ArticleManager) GetPointerMgr() *minipointer.PointerManager {
	return obj.pointerMgr
}

func (obj *ArticleManager) makeArticleId(created time.Time, secretKey string) string {
	hashKey := obj.hashStr(fmt.Sprintf("p:%s;s:%s;c:%d;", obj.prefixOfId, secretKey, created.UnixNano()))
	return "" + obj.prefixOfId + "-v1e-" + hashKey
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
	return string(base32.StdEncoding.EncodeToString(sha1Obj.Sum(nil)))
}

func (obj *ArticleManager) makeRandomId() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}

/*
func (obj *ArticleManager) SaveOnDB(ctx context.Context, artObj *Article) error {
	return artObj.saveOnDB(ctx)
}*/

func (obj *ArticleManager) SaveUsrWithImmutable(ctx context.Context, artObj *Article) (*Article, error) {
	sign := strconv.Itoa(time.Now().Nanosecond())
	nextArObj := obj.NewArticleFromArticle(ctx, artObj, sign)
	nextArObj.SetUpdated(time.Now())
	saveErr := nextArObj.saveOnDB(ctx)
	if saveErr != nil {
		Debug(ctx, ".>>>>>>>> AAA")
		return artObj, saveErr
	}
	pointerObj := obj.pointerMgr.GetPointerForRelayId(ctx, artObj.GetArticleId())
	pointerObj.SetValue(nextArObj.GetArticleId())
	pointerObj.SetSign(nextArObj.gaeObject.Sign)
	pointerObj.SetOwner(artObj.GetArticleId())
	Debug(ctx, ".>>>>>>>> AAA > "+pointerObj.GetOwner())
	savePointerErr := pointerObj.Save(ctx)
	if savePointerErr != nil {
		err := obj.DeleteFromArticleId(ctx, nextArObj.GetArticleId(), nextArObj.gaeObject.Sign)
		if err != nil {
			Debug(ctx, "<GOMIDATA>"+nextArObj.GetArticleId()+":"+nextArObj.gaeObject.Sign+"<GOMIDATA>")
		}
		Debug(ctx, ".>>>>>>>> BBB")

		return artObj, savePointerErr
	}
	if artObj.gaeObject.Sign != "0" {
		obj.DeleteFromArticleId(ctx, artObj.GetArticleId(), artObj.GetSign())
	}
	return nextArObj, nil
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
