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
		obj.OnUpdateArtFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, ErrorCodeNotFoundArticleId, "Not Found Article")
		return
	}

	errOnGe := obj.OnUpdateRequest(w, r, obj, inputProp, outputProp)
	if nil != errOnGe {
		obj.OnUpdateArtFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, errOnGe.Error())
		return
	}

	artObj, _, errGetArt := obj.GetManager().GetArticleFromPointer(ctx, articleId)
	if errGetArt != nil {
		obj.OnUpdateArtFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, ErrorCodeNotFoundArticleId, "Not Found Article")
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
	obj.tagMana.DeleteTagsFromOwner(appengine.NewContext(r), articleId)
	obj.tagMana.AddBasicTags(ctx, tags, "art://"+artObj.GetArticleId()+"@"+artObj.GetSign(), artObj.GetArticleId(), "")

	if errSave != nil {
		obj.OnUpdateArtFailed(w, r, obj, inputProp, outputProp)
		obj.HandleError(w, r, outputProp, ErrorCodeFailedToSave, errSave.Error())
		return
	} else {
		propObj.SetPropString("", "articleId", artObj.GetArticleId())
		errOnSc := obj.OnUpdateArtSuccess(w, r, obj, inputProp, outputProp)
		if nil != errOnSc {
			obj.OnUpdateArtFailed(w, r, obj, inputProp, outputProp)
			obj.HandleError(w, r, outputProp, ErrorCodeFailedToCheckAboutGetCalled, errOnSc.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(propObj.ToJson())
	}
}
