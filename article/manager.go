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

type ArticleManagerConfig struct {
	RootGroup      string
	KindArticle    string
	KindPointer    string
	PrefixOfId     string
	LimitOfFinding int
}
type ArticleManager struct {
	pointerMgr *minipointer.PointerManager
	config     ArticleManagerConfig
}

func NewArticleManager(config ArticleManagerConfig) *ArticleManager {
	if config.RootGroup == "" {
		config.RootGroup = "FFArt"
	}
	if config.KindArticle == "" {
		config.KindArticle = "FFArt"
	}
	if config.KindPointer == "" {
		config.KindPointer = config.KindArticle + "pointer"
	}
	if config.PrefixOfId == "" {
		config.PrefixOfId = "ffart"
	}
	if config.LimitOfFinding <= 0 {
		config.LimitOfFinding = 20
	}

	ret := new(ArticleManager)
	ret.config = config
	ret.pointerMgr = minipointer.NewPointerManager(minipointer.PointerManagerConfig{
		RootGroup: config.RootGroup,
		Kind:      config.KindPointer,
	})
	return ret
}

func (obj *ArticleManager) GetKind() string {
	return obj.config.KindArticle
}

func (obj *ArticleManager) GetPointerMgr() *minipointer.PointerManager {
	return obj.pointerMgr
}

func (obj *ArticleManager) makeArticleId(created time.Time, secretKey string) string {
	hashKey := obj.hashStr(fmt.Sprintf("p:%s;s:%s;c:%d;", obj.config.PrefixOfId, secretKey, created.UnixNano()))
	return "" + obj.config.PrefixOfId + "-v1e-" + hashKey
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
	pointerObj := obj.pointerMgr.GetPointerWithNewForRelayId(ctx, artObj.GetArticleId())
	pointerObj.SetValue(nextArObj.GetArticleId())
	pointerObj.SetSign(nextArObj.gaeObject.Sign)
	pointerObj.SetOwner(artObj.GetArticleId())
	Debug(ctx, ".>>>>>>>> AAA > "+pointerObj.GetOwner())
	savePointerErr := obj.pointerMgr.Save(ctx, pointerObj)
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
