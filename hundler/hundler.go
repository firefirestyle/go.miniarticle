package hundler

import (
	"net/http"

	//	"strings"

	"io/ioutil"

	"github.com/firefirestyle/go.miniarticle/article"
	miniblob "github.com/firefirestyle/go.miniblob/blob"
	blobhandler "github.com/firefirestyle/go.miniblob/handler"
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
	pointerKind string
	artMana     *article.ArticleManager
	blobHundler *blobhandler.BlobHandler
	onEvents    ArticleHandlerOnEvent
}

type ArticleHandlerManagerConfig struct {
	ProjectId       string
	ArticleKind     string
	PointerKind     string
	BlobKind        string
	BlobCallbackUrl string
	BlobSign        string
}

type ArticleHandlerOnEvent struct {
	OnNewRequest    func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error
	OnNewBeforeSave func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, artObj *article.Article, input *miniprop.MiniProp, output *miniprop.MiniProp) error
	OnNewArtFailed  func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp)
	OnNewArtSuccess func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error
	//
	OnUpdateRequest    func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error
	OnUpdateArtFailed  func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp)
	OnUpdateArtSuccess func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error
	//
	blobOnEvent blobhandler.BlobHandlerOnEvent
}

func NewArtHandler(config ArticleHandlerManagerConfig, onEvents ArticleHandlerOnEvent) *ArticleHandler {

	artMana := article.NewArticleManager(config.ProjectId, config.ArticleKind, config.PointerKind, "art", 10)
	//
	//
	if onEvents.OnNewRequest == nil {
		onEvents.OnNewRequest = func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error {
			return nil
		}
	}
	if onEvents.OnNewArtFailed == nil {
		onEvents.OnNewArtFailed = func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) {
			return
		}
	}
	if onEvents.OnNewArtSuccess == nil {
		onEvents.OnNewArtSuccess = func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error {
			return nil
		}
	}
	if onEvents.OnNewBeforeSave == nil {
		onEvents.OnNewBeforeSave = func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, artObj *article.Article, input *miniprop.MiniProp, output *miniprop.MiniProp) error {
			return nil
		}
	}

	//
	//
	if onEvents.OnUpdateRequest == nil {
		onEvents.OnUpdateRequest = func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error {
			return nil
		}
	}
	if onEvents.OnUpdateArtFailed == nil {
		onEvents.OnUpdateArtFailed = func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) {
			return
		}
	}
	if onEvents.OnUpdateArtSuccess == nil {
		onEvents.OnUpdateArtSuccess = func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error {
			return nil
		}
	}

	//
	//
	//
	artHandlerObj := &ArticleHandler{
		projectId:   config.ProjectId,
		articleKind: config.ArticleKind,
		blobKind:    config.BlobKind,
		artMana:     artMana,
		onEvents:    onEvents,
	}
	completeFunc := onEvents.blobOnEvent.OnBlobComplete
	onEvents.blobOnEvent.OnBlobComplete = func(w http.ResponseWriter, r *http.Request, o *miniprop.MiniProp, hh *blobhandler.BlobHandler, i *miniblob.BlobItem) error {
		dirSrc := r.URL.Query().Get("dir")
		articlId := artHandlerObj.GetArticleIdFromDir(dirSrc)
		dir := artHandlerObj.GetDirFromDir(dirSrc)
		//
		//
		ctx := appengine.NewContext(r)
		Debug(ctx, "OnBlobComplete "+articlId+":"+dir)
		_, errGet := artHandlerObj.GetManager().GetArticleFromPointer(ctx, articlId)
		if errGet != nil {
			return errGet
		}
		//
		if completeFunc != nil {
			return completeFunc(w, r, o, hh, i)
		} else {
			return nil
		}
	}
	//
	artHandlerObj.blobHundler = blobhandler.NewBlobHandler(config.BlobCallbackUrl, config.BlobSign,
		miniblob.BlobManagerConfig{
			ProjectId:   config.ProjectId,
			Kind:        config.BlobKind,
			CallbackUrl: config.BlobCallbackUrl,
		}, onEvents.blobOnEvent)

	return artHandlerObj
}

func (obj *ArticleHandler) GetManager() *article.ArticleManager {
	return obj.artMana
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

//
//
//

// HandleBlobRequestTokenFromParams

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
