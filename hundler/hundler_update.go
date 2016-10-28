package hundler

import (
	"net/http"

	"github.com/firefirestyle/go.miniprop"
	"google.golang.org/appengine"
)

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
	if articleId == "" {
		obj.onEvents.OnUpdateArtFailed(w, r, obj, inputProp, outputProp)
		HandleError(w, r, outputProp, ErrorCodeNotFoundArticleId, "Not Found Article")
		return
	}

	errOnGe := obj.onEvents.OnUpdateRequest(w, r, obj, inputProp, outputProp)
	if nil != errOnGe {
		obj.onEvents.OnUpdateArtFailed(w, r, obj, inputProp, outputProp)
		HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, errOnGe.Error())
		return
	}

	artObj, errGetArt := obj.GetManager().GetArticleFromArticleIdOnQuery(ctx, articleId)
	if errGetArt != nil {
		obj.onEvents.OnUpdateArtFailed(w, r, obj, inputProp, outputProp)
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
		obj.onEvents.OnUpdateArtFailed(w, r, obj, inputProp, outputProp)
		HandleError(w, r, outputProp, ErrorCodeFailedToSave, errSave.Error())
		return
	} else {
		propObj.SetPropString("", "articleId", artObj.GetArticleId())
		errOnSc := obj.onEvents.OnUpdateArtSuccess(w, r, obj, inputProp, outputProp)
		if nil != errOnSc {
			obj.onEvents.OnUpdateArtFailed(w, r, obj, inputProp, outputProp)
			HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, errOnSc.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(propObj.ToJson())
	}
}
