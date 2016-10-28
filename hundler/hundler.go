package hundler

import (
	"net/http"

	//	"strings"

	"io/ioutil"

	"github.com/firefirestyle/go.miniarticle/article"
	"github.com/firefirestyle/go.miniblob"
	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
)

const (
	ErrorCodeFailedToSave                = 2001
	ErrorCodeFailedToCheckAboutGetCalled = 2002
	ErrorCodeNotFoundArticleId           = 2003
)

type ArticleHandler struct {
	projectId   string
	articleKind string
	blobKind    string
	artMana     *article.ArticleManager
	blobHundler *miniblob.BlobHandler
	onEvents    ArticleHandlerOnEvent
}

type ArticleHandlerManagerConfig struct {
	ProjectId       string
	ArticleKind     string
	BlobKind        string
	BlobCallbackUrl string
	BlobSign        string
}

type ArticleHandlerOnEvent struct {
	OnGetCalled  func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error
	OnGetFailed  func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp)
	OnGetSuccess func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error
}

func (obj *ArticleHandler) GetManager() *article.ArticleManager {
	return obj.artMana
}

func NewArtHandler(config ArticleHandlerManagerConfig, onEvents ArticleHandlerOnEvent) *ArticleHandler {
	blobHundler := miniblob.NewBlobHandler(config.BlobCallbackUrl, config.BlobSign,
		miniblob.BlobManagerConfig{
			ProjectId:   config.ProjectId,
			Kind:        config.BlobKind,
			CallbackUrl: config.BlobCallbackUrl,
		}, miniblob.BlobHandlerOnEvent{})
	artMana := article.NewArticleManager(config.ProjectId, config.ArticleKind, "art-", 10)
	//
	//
	if onEvents.OnGetCalled == nil {
		onEvents.OnGetCalled = func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error {
			return nil
		}
	}
	if onEvents.OnGetFailed == nil {
		onEvents.OnGetFailed = func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) {
			return
		}
	}
	if onEvents.OnGetSuccess == nil {
		onEvents.OnGetSuccess = func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error {
			return nil
		}
	}
	return &ArticleHandler{
		projectId:   config.ProjectId,
		articleKind: config.ArticleKind,
		blobKind:    config.BlobKind,
		artMana:     artMana,
		blobHundler: blobHundler,
		onEvents:    onEvents,
	}
}

func HandleError(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, errorCode int, errorMessage string) {
	//
	//
	if errorCode != 0 {
		outputProp.SetInt("errorCode", errorCode)
	}
	if errorMessage != "" {
		outputProp.SetString("errorMessage", errorMessage)
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write(outputProp.ToJson())
}

func (obj *ArticleHandler) GetInputProp(w http.ResponseWriter, r *http.Request) *miniprop.MiniProp {
	v, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return miniprop.NewMiniProp()
	} else {
		return miniprop.NewMiniPropFromJson(v)
	}
}

func (obj *ArticleHandler) HandleNew(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	propObj := miniprop.NewMiniProp()
	//
	// load param from json
	inputProp := obj.GetInputProp(w, r)
	ownerName := inputProp.GetString("ownerName", "")
	title := inputProp.GetString("title", "")
	target := inputProp.GetString("target", "")
	content := inputProp.GetString("content", "")
	//
	//
	outputProp := miniprop.NewMiniProp()
	{
		err := obj.onEvents.OnGetCalled(w, r, obj, inputProp, outputProp)
		if nil != err {
			HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, err.Error())
			return
		}
	}
	artObj := obj.GetManager().NewArticle(ctx, ownerName, "")
	artObj.SetTitle(title)
	artObj.SetTarget(target)
	artObj.SetCont(content)
	//
	//
	errSave := obj.GetManager().SaveOnDB(ctx, artObj)
	if errSave != nil {
		HandleError(w, r, outputProp, ErrorCodeFailedToSave, errSave.Error())
		return
	} else {
		propObj.SetPropString("", "articleId", artObj.GetArticleId())
		errOnSc := obj.onEvents.OnGetSuccess(w, r, obj, inputProp, outputProp)
		if nil != errOnSc {
			if nil != obj.GetManager().DeleteFromArticleId(ctx, artObj.GetArticleId(), "") {
				Debug(ctx, "<GOMIDATA>articleId="+artObj.GetArticleId())
			}
			HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, errOnSc.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(propObj.ToJson())
	}
}

//
//
//
func (obj *ArticleHandler) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	propObj := miniprop.NewMiniProp()
	//
	// load param from json
	inputProp := obj.GetInputProp(w, r)
	articleId := inputProp.GetString("articleId", "")
	//	ownerName := inputProp.GetString("ownerName", "")
	title := inputProp.GetString("title", "")
	target := inputProp.GetString("target", "")
	content := inputProp.GetString("content", "")
	tags := inputProp.GetPropStringList("", "tags", make([]string, 0))
	//
	//
	outputProp := miniprop.NewMiniProp()
	//
	//
	if articleId == "" {
		HandleError(w, r, outputProp, ErrorCodeNotFoundArticleId, "Not Found Article")
		return
	}

	errOnGe := obj.onEvents.OnGetCalled(w, r, obj, inputProp, outputProp)
	if nil != errOnGe {
		HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, errOnGe.Error())
		return
	}

	artObj, errGetArt := obj.GetManager().GetArticleFromArticleIdOnQuery(ctx, articleId)
	if errGetArt != nil {
		HandleError(w, r, outputProp, ErrorCodeNotFoundArticleId, "Not Found Article")
		return
	}
	//

	artObj.SetTitle(title)
	artObj.SetTarget(target)
	artObj.SetCont(content)
	artObj.SetTags(tags)
	//
	//
	errSave := obj.GetManager().SaveUsrWithImmutable(ctx, artObj)
	if errSave != nil {
		HandleError(w, r, outputProp, ErrorCodeFailedToSave, errSave.Error())
		return
	} else {
		propObj.SetPropString("", "articleId", artObj.GetArticleId())
		errOnSc := obj.onEvents.OnGetSuccess(w, r, obj, inputProp, outputProp)
		if nil != errOnSc {
			HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, errOnSc.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(propObj.ToJson())
	}
}

func (obj *ArticleHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	values := r.URL.Query()
	key := values.Get("key")
	articleId := values.Get("articleId")
	sign := values.Get("sign")
	mode := values.Get("m")
	//
	if key != "" {
		keyInfo := obj.GetManager().ExtractInfoFromStringId(key)
		articleId = keyInfo.ArticleId
		sign = keyInfo.Sign
	}
	var artObj *article.Article
	var err error
	if mode != "q" {
		artObj, err = obj.GetManager().GetArticleFromArticleId(ctx, articleId, sign)
	} else {
		artObj, err = obj.GetManager().GetArticleFromArticleIdOnQuery(ctx, articleId)
	}
	if err != nil {
		HandleError(w, r, propObj, ErrorCodeNotFoundArticleId, "not found article")
		return
	}
	if mode != "q" {
		w.Header().Set("Cache-Control", "public, max-age=2592000")
	}
	w.Write(artObj.ToJsonPublicOnly())
}

func (obj *ArticleHandler) HandleFind(w http.ResponseWriter, r *http.Request) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	values := r.URL.Query()
	cursor := values.Get("cursor")
	foundObj := obj.GetManager().FindArticleWithNewOrder(ctx, cursor, true)

	propObj.SetPropStringList("", "keys", foundObj.ArticleIds)
	propObj.SetPropString("", "cursorOne", foundObj.CursorOne)
	propObj.SetPropString("", "cursorOne", foundObj.CursorNext)
	w.Write(propObj.ToJson())
}

// HandleBlobRequestTokenFromParams
func (obj *ArticleHandler) HandleBlobRequestToken(w http.ResponseWriter, r *http.Request) {
	//
	// load param from json
	articleId := r.URL.Query().Get("articleId")
	dir := r.URL.Query().Get("dir")
	name := r.URL.Query().Get("file")
	//
	// todo check articleId

	//
	//
	obj.blobHundler.HandleBlobRequestTokenFromParams(w, r, "/art/"+articleId+"/"+dir, name)
}

func (obj *ArticleHandler) HandleBlobUpdated(w http.ResponseWriter, r *http.Request) {
	//
	ctx := appengine.NewContext(r)
	Debug(ctx, "callbeck AAAA")
	obj.blobHundler.HandleUploaded(w, r)
}
func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
