package hundler

import (
	"net/http"

	"github.com/firefirestyle/go.miniprop"
	"google.golang.org/appengine"
)

func (obj *ArticleHandler) HandleNew(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	propObj := miniprop.NewMiniProp()
	//
	// load param from json
	inputProp := obj.GetInputProp(w, r)
	title := inputProp.GetString("title", "")
	target := inputProp.GetString("target", "")
	content := inputProp.GetString("content", "")
	ownerName := inputProp.GetString("ownerName", "")

	//
	//
	outputProp := miniprop.NewMiniProp()
	errCall := obj.onEvents.OnNewRequest(w, r, obj, inputProp, outputProp)
	if nil != errCall {
		obj.onEvents.OnNewArtFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, errCall.Error())
		return
	}
	//
	artObj := obj.GetManager().NewArticle(ctx)
	artObj.SetTitle(title)
	artObj.SetTarget(target)
	artObj.SetCont(content)
	artObj.SetUserName(ownerName)
	//
	errNew := obj.onEvents.OnNewBeforeSave(w, r, obj, artObj, inputProp, outputProp)
	if nil != errNew {
		obj.onEvents.OnNewArtFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, errNew.Error())
		return
	}
	//
	errSave := obj.GetManager().SaveUsrWithImmutable(ctx, artObj)
	if errSave != nil {
		obj.onEvents.OnNewArtFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, ErrorCodeFailedToSave, errSave.Error())
		return
	} else {
		propObj.SetPropString("", "articleId", artObj.GetArticleId())
		errOnSc := obj.onEvents.OnNewArtSuccess(w, r, obj, inputProp, outputProp)
		if nil != errOnSc {
			if nil != obj.GetManager().DeleteFromArticleId(ctx, artObj.GetArticleId(), "") {
				Debug(ctx, "<GOMIDATA>articleId="+artObj.GetArticleId())
			}
			obj.onEvents.OnNewArtFailed(w, r, obj, inputProp, outputProp)
			obj.HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, errOnSc.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(propObj.ToJson())
	}
}
