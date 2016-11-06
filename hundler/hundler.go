package hundler

import (
	"io/ioutil"
	"net/http"

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
	RootGroup       string
	ArticleKind     string
	PointerKind     string
	BlobKind        string
	BlobPointerKind string
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
}

func NewArtHandler(config ArticleHandlerManagerConfig, onEvents ArticleHandlerOnEvent) *ArticleHandler {
	if config.RootGroup == "" {
		config.RootGroup = "ffstyle"
	}
	if config.ArticleKind == "" {
		config.ArticleKind = "ffart"
	}
	if config.PointerKind == "" {
		config.PointerKind = config.ArticleKind + "-pointer"
	}
	if config.BlobKind == "" {
		config.BlobKind = config.ArticleKind + "-blob-pointer"
	}
	if config.BlobPointerKind == "" {
		config.BlobPointerKind = config.ArticleKind + "-blob-pointer"
	}

	artMana := article.NewArticleManager(config.RootGroup, config.ArticleKind, config.PointerKind, "art", 10)
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
		projectId:   config.RootGroup,
		articleKind: config.ArticleKind,
		blobKind:    config.BlobKind,
		artMana:     artMana,
		onEvents:    onEvents,
	}

	//
	artHandlerObj.blobHundler = blobhandler.NewBlobHandler(config.BlobCallbackUrl, config.BlobSign,
		miniblob.BlobManagerConfig{
			RootGroup:   config.RootGroup,
			Kind:        config.BlobKind,
			CallbackUrl: config.BlobCallbackUrl,
			PointerKind: config.PointerKind,
		})
	artHandlerObj.blobHundler.AddOnBlobComplete(func(w http.ResponseWriter, r *http.Request, o *miniprop.MiniProp, hh *blobhandler.BlobHandler, i *miniblob.BlobItem) error {
		dirSrc := r.URL.Query().Get("dir")
		articlId := artHandlerObj.GetArticleIdFromDir(dirSrc)
		dir := artHandlerObj.GetDirFromDir(dirSrc)
		fileName := r.URL.Query().Get("file")
		//
		//
		ctx := appengine.NewContext(r)
		Debug(ctx, "OnBlobComplete ::"+articlId+"::"+dir+"::"+fileName+"::")
		artObj, errGet := artHandlerObj.GetManager().GetArticleFromPointer(ctx, articlId)
		if errGet != nil {
			Debug(ctx, "From Pointer GEt ER "+articlId)
			return errGet
		}
		if dir == "/" && fileName == "icon" {
			artObj.SetIconUrl("key://" + i.GetBlobKey())
			errSave := artHandlerObj.GetManager().SaveUsrWithImmutable(ctx, artObj)
			if errSave != nil {
				return errSave
			}
		}
		return nil
	})
	return artHandlerObj
}

func (obj *ArticleHandler) GetManager() *article.ArticleManager {
	return obj.artMana
}

func (obj *ArticleHandler) GetBlobHandler() *blobhandler.BlobHandler {
	return obj.blobHundler
}

func (obj *ArticleHandler) HandleError(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, errorCode int, errorMessage string) {
	//
	//
	if outputProp == nil {
		outputProp = miniprop.NewMiniProp()
	}
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

///
//
