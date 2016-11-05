package template

import (
	"net/http"

	"errors"

	"github.com/firefirestyle/go.miniarticle/article"
	arthundler "github.com/firefirestyle/go.miniarticle/hundler"

	blobhandler "github.com/firefirestyle/go.miniblob/handler"
	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.minisession"
	userHandler "github.com/firefirestyle/go.miniuser/handler"

	//
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	//

	//"io/ioutil"
)

const (
	UrlArtNew             = "/api/v1/art/new"
	UrlArtUpdate          = "/api/v1/art/update"
	UrlArtFind            = "/api/v1/art/find"
	UrlArtGet             = "/api/v1/art/get"
	UrlArtBlobGet         = "/api/v1/art/getblob"
	UrlArtRequestBlobUrl  = "/api/v1/art/requestbloburl"
	UrlArtCallbackBlobUrl = "/api/v1/art/callbackbloburl"
)

//
//

type ArtTemplateConfig struct {
	GroupName    string
	KindBaseName string
	PrivateKey   string
}

type ArtTemplate struct {
	config         ArtTemplateConfig
	artHandlerObj  *arthundler.ArticleHandler
	getUserHundler func(context.Context) *userHandler.UserHandler
}

func NewArtTemplate(config ArtTemplateConfig, getUserHundler func(context.Context) *userHandler.UserHandler) *ArtTemplate {
	if config.GroupName == "" {
		config.GroupName = "FFS"
	}
	if config.KindBaseName == "" {
		config.KindBaseName = "FFSArt"
	}

	return &ArtTemplate{
		config:         config,
		getUserHundler: getUserHundler,
	}
}

//
//

func (tmpObj *ArtTemplate) CheckLogin(r *http.Request, input *miniprop.MiniProp) minisession.CheckLoginIdResult {
	token := input.GetString("token", "")
	return tmpObj.getUserHundler(appengine.NewContext(r)).CheckLogin(r, token)
}

func (tmpObj *ArtTemplate) GetArtHundlerObj(ctx context.Context) *arthundler.ArticleHandler {
	if tmpObj.artHandlerObj == nil {
		tmpObj.artHandlerObj = arthundler.NewArtHandler(
			arthundler.ArticleHandlerManagerConfig{
				RootGroup:       tmpObj.config.GroupName,
				ArticleKind:     tmpObj.config.KindBaseName,
				BlobCallbackUrl: UrlArtCallbackBlobUrl,
				BlobSign:        appengine.VersionID(ctx),
			}, //
			arthundler.ArticleHandlerOnEvent{
				OnNewBeforeSave: func(w http.ResponseWriter, r *http.Request, handler *arthundler.ArticleHandler, artObj *article.Article, input *miniprop.MiniProp, output *miniprop.MiniProp) error {
					ret := tmpObj.CheckLogin(r, input)
					if ret.IsLogin == false {
						return errors.New("Failed in token check")
					} else {
						artObj.SetUserName(ret.AccessTokenObj.GetUserName())
						return nil
					}
				},
				OnUpdateRequest: func(w http.ResponseWriter, r *http.Request, handler *arthundler.ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error {
					ret := tmpObj.CheckLogin(r, input)
					if ret.IsLogin == false {
						return errors.New("Failed in token check")
					} else {
						return nil
					}
				},
			})
		tmpObj.artHandlerObj.GetBlobHandler().AddOnBlobRequest(func(w http.ResponseWriter, r *http.Request, input *miniprop.MiniProp, output *miniprop.MiniProp, h *blobhandler.BlobHandler) (map[string]string, error) {
			ret := tmpObj.CheckLogin(r, input)
			if ret.IsLogin == false {
				return map[string]string{}, errors.New("Failed in token check")
			}
			return map[string]string{"sst": ret.AccessTokenObj.GetLoginId()}, nil
		})
	}
	return tmpObj.artHandlerObj
}

func (tmpObj *ArtTemplate) InitArtApi() {

	// art
	// UrlArtNew
	http.HandleFunc(UrlArtNew, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.GetArtHundlerObj(ctx).HandleNew(w, r)
	})

	// art
	// UrlArtNew
	http.HandleFunc(UrlArtUpdate, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.GetArtHundlerObj(ctx).HandleUpdate(w, r)
	})

	http.HandleFunc(UrlArtFind, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.GetArtHundlerObj(ctx).HandleFind(w, r)
	})

	http.HandleFunc(UrlArtGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.GetArtHundlerObj(ctx).HandleGet(w, r)
	})
	//UrlArtGet

	http.HandleFunc(UrlArtRequestBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.GetArtHundlerObj(ctx).HandleBlobRequestToken(w, r)
	})

	http.HandleFunc(UrlArtCallbackBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		Debug(ctx, "asdfasdfasdf")

		tmpObj.GetArtHundlerObj(ctx).HandleBlobUpdated(w, r)
	})

	http.HandleFunc(UrlArtBlobGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		Debug(ctx, "asdfasdfasdf")

		tmpObj.GetArtHundlerObj(ctx).HandleBlobGet(w, r)
	})

}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
