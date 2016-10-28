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
	ownerName := inputProp.GetString("ownerName", "")
	title := inputProp.GetString("title", "")
	target := inputProp.GetString("target", "")
	content := inputProp.GetString("content", "")
	//
	//
	outputProp := miniprop.NewMiniProp()
	{
		err := obj.onEvents.OnNewArtCalled(w, r, obj, inputProp, outputProp)
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
		obj.onEvents.OnNewArtFailed(w, r, obj, inputProp, outputProp)
		HandleError(w, r, outputProp, ErrorCodeFailedToSave, errSave.Error())
		return
	} else {
		propObj.SetPropString("", "articleId", artObj.GetArticleId())
		errOnSc := obj.onEvents.OnNewArtSuccess(w, r, obj, inputProp, outputProp)
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
