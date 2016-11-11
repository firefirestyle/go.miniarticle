package template

import (
	"net/http"

	"errors"

	"github.com/firefirestyle/go.miniarticle/article"
	arthundler "github.com/firefirestyle/go.miniarticle/handler"

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
	"sync"
)

const (
	UrlArtNew             = "/api/v1/art/new"
	UrlArtUpdate          = "/api/v1/art/update"
	UrlArtFind            = "/api/v1/art/find"
	UrlArtFindMe          = "/api/v1/art/find_with_token"
	UrlArtGet             = "/api/v1/art/get"
	UrlArtBlobGet         = "/api/v1/art/getblob"
	UrlArtRequestBlobUrl  = "/api/v1/art/requestbloburl"
	UrlArtCallbackBlobUrl = "/api/v1/art/callbackbloburl"
	UrlArtDelete          = "/api/v1/art/delete"
)

type ArtTemplateConfig struct {
	GroupName                  string
	KindBaseName               string
	PrivateKey                 string
	MemcachedOnlyInBlobPointer bool
}

type ArtTemplate struct {
	config         ArtTemplateConfig
	artHandlerObj  *arthundler.ArticleHandler
	getUserHundler func(context.Context) *userHandler.UserHandler
	initOpt        func(context.Context)
	once           *sync.Once
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
		once:           new(sync.Once),
	}
}

func (tmpObj *ArtTemplate) SetInitFunc(f func(ctx context.Context)) {
	tmpObj.initOpt = f
}

func (tmpObj *ArtTemplate) InitalizeTemplate(ctx context.Context) {
	tmpObj.GetArtHundlerObj(ctx)
	tmpObj.getUserHundler(ctx)
	if tmpObj != nil {
		tmpObj.initOpt(ctx)
	}
}

func (tmpObj *ArtTemplate) CheckLogin(r *http.Request, token string) minisession.CheckLoginIdResult {

	return tmpObj.getUserHundler(appengine.NewContext(r)).CheckLogin(r, token)
}

func (tmpObj *ArtTemplate) GetArtHundlerObj(ctx context.Context) *arthundler.ArticleHandler {
	if tmpObj.artHandlerObj == nil {
		tmpObj.artHandlerObj = arthundler.NewArtHandler(
			arthundler.ArticleHandlerConfig{
				RootGroup:       tmpObj.config.GroupName,
				ArticleKind:     tmpObj.config.KindBaseName,
				BlobCallbackUrl: UrlArtCallbackBlobUrl,
				BlobSign:        appengine.VersionID(ctx),
				MemcachedOnly:   tmpObj.config.MemcachedOnlyInBlobPointer,
			})
		tmpObj.artHandlerObj.AddOnNewBeforeSave(func(w http.ResponseWriter, r *http.Request, handler *arthundler.ArticleHandler, artObj *article.Article, input *miniprop.MiniProp, output *miniprop.MiniProp) error {
			ret := tmpObj.CheckLogin(r, input.GetString("token", ""))
			if ret.IsLogin == false {
				return errors.New("Failed in token check")
			} else {
				artObj.SetUserName(ret.AccessTokenObj.GetUserName())
				return nil
			}
		})
		tmpObj.artHandlerObj.AddOnUpdateRequest(func(w http.ResponseWriter, r *http.Request, handler *arthundler.ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error {
			ret := tmpObj.CheckLogin(r, input.GetString("token", ""))
			if ret.IsLogin == false {
				return errors.New("Failed in token check")
			} else {
				return nil
			}
		})
		tmpObj.artHandlerObj.AddOnGetArtSuccess(func(w http.ResponseWriter, r *http.Request, h *arthundler.ArticleHandler, i *article.Article, o *miniprop.MiniProp) error {
			pointerObj, pointerErr := tmpObj.getUserHundler(appengine.NewContext(r)).GetManager().GetPointerFromUserName(appengine.NewContext(r), i.GetUserName())
			if pointerErr == nil {
				o.SetString("userSign", pointerObj.GetSign())
			}
			return nil
		})
		tmpObj.artHandlerObj.GetBlobHandler().AddOnBlobRequest(func(w http.ResponseWriter, r *http.Request, input *miniprop.MiniProp, output *miniprop.MiniProp, h *blobhandler.BlobHandler) (map[string]string, error) {
			ret := tmpObj.CheckLogin(r, input.GetString("token", ""))
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
		tmpObj.once.Do(func() {
			tmpObj.InitalizeTemplate(ctx)
		})
		tmpObj.GetArtHundlerObj(ctx).HandleNew(w, r)
	})

	// art
	// UrlArtNew
	http.HandleFunc(UrlArtUpdate, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.once.Do(func() {
			tmpObj.InitalizeTemplate(ctx)
		})
		tmpObj.GetArtHundlerObj(ctx).HandleUpdate(w, r)
	})

	http.HandleFunc(UrlArtFind, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.once.Do(func() {
			tmpObj.InitalizeTemplate(ctx)
		})
		tmpObj.GetArtHundlerObj(ctx).HandleFind(w, r)
	})

	http.HandleFunc(UrlArtFindMe, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.once.Do(func() {
			tmpObj.InitalizeTemplate(ctx)
		})
		propObj := miniprop.NewMiniPropFromJsonReader(r.Body)
		loginInfo := tmpObj.CheckLogin(r, propObj.GetString("token", ""))
		if loginInfo.IsLogin == false {
			tmpObj.GetArtHundlerObj(ctx).HandleError(w, r, nil, 4001, "failed to login")
		} else {
			tmpObj.GetArtHundlerObj(ctx).HandleFindBase(w, r, //
				propObj.GetString("cursor", ""), loginInfo.AccessTokenObj.GetUserName(), "", "")
		}
	})

	http.HandleFunc(UrlArtGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.once.Do(func() {
			tmpObj.InitalizeTemplate(ctx)
		})
		tmpObj.GetArtHundlerObj(ctx).HandleGet(w, r)
	})
	//UrlArtGet

	http.HandleFunc(UrlArtRequestBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.once.Do(func() {
			tmpObj.InitalizeTemplate(ctx)
		})
		tmpObj.GetArtHundlerObj(ctx).HandleBlobRequestToken(w, r)
	})

	http.HandleFunc(UrlArtCallbackBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.once.Do(func() {
			tmpObj.InitalizeTemplate(ctx)
		})
		tmpObj.GetArtHundlerObj(ctx).HandleBlobUpdated(w, r)
	})

	http.HandleFunc(UrlArtBlobGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.once.Do(func() {
			tmpObj.InitalizeTemplate(ctx)
		})
		tmpObj.GetArtHundlerObj(ctx).HandleBlobGet(w, r)
	})

	http.HandleFunc(UrlArtDelete, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		tmpObj.once.Do(func() {
			tmpObj.InitalizeTemplate(ctx)
		})
		propObj := miniprop.NewMiniPropFromJsonReader(r.Body)
		loginInfo := tmpObj.CheckLogin(r, propObj.GetString("token", ""))
		if loginInfo.IsLogin == false {
			tmpObj.GetArtHundlerObj(ctx).HandleError(w, r, nil, 4001, "failed to login")
		} else {
			tmpObj.GetArtHundlerObj(ctx).HandleDeleteBaseWithFile(w, r, propObj.GetString("articleId", ""))
		}
	})
	//
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
