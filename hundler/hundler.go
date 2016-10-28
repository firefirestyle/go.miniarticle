package hundler

import (
	"net/http"

	//	"strings"

	"io/ioutil"

	"github.com/firefirestyle/go.miniarticle/article"
	"github.com/firefirestyle/go.miniblob"
	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
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
	OnNewArtCalled  func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error
	OnNewArtFailed  func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp)
	OnNewArtSuccess func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error
	//
	OnUpdateArtCalled  func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error
	OnUpdateArtFailed  func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp)
	OnUpdateArtSuccess func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error
	//
	OnBlobRequest func(http.ResponseWriter, *http.Request, *miniprop.MiniProp, *miniblob.BlobHandler) (string, map[string]string, error)
}

func NewArtHandler(config ArticleHandlerManagerConfig, onEvents ArticleHandlerOnEvent) *ArticleHandler {

	artMana := article.NewArticleManager(config.ProjectId, config.ArticleKind, "art-", 10)
	//
	//
	if onEvents.OnNewArtCalled == nil {
		onEvents.OnNewArtCalled = func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error {
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
	//
	if onEvents.OnUpdateArtCalled == nil {
		onEvents.OnUpdateArtCalled = func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error {
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
	if onEvents.OnBlobRequest == nil {
		onEvents.OnBlobRequest = func(http.ResponseWriter, *http.Request, *miniprop.MiniProp, *miniblob.BlobHandler) (string, map[string]string, error) {
			return "dummy", map[string]string{}, nil
		}
	}
	blobHundler := miniblob.NewBlobHandler(config.BlobCallbackUrl, config.BlobSign,
		miniblob.BlobManagerConfig{
			ProjectId:   config.ProjectId,
			Kind:        config.BlobKind,
			CallbackUrl: config.BlobCallbackUrl,
		}, miniblob.BlobHandlerOnEvent{
			OnRequest: func(w http.ResponseWriter, r *http.Request, input *miniprop.MiniProp, blob *miniblob.BlobHandler) (string, map[string]string, error) {
				return onEvents.OnBlobRequest(w, r, input, blob)
			},
		})
	return &ArticleHandler{
		projectId:   config.ProjectId,
		articleKind: config.ArticleKind,
		blobKind:    config.BlobKind,
		artMana:     artMana,
		blobHundler: blobHundler,
		onEvents:    onEvents,
	}
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
