package handler

import (
	"net/http"

	//	"github.com/firefirestyle/go.miniarticle/article"
	"github.com/firefirestyle/go.miniprop"
	"google.golang.org/appengine"
)

//
// you must to delete file before call this method, if there are many articleid's file.
//
func (obj *ArticleHandler) HandleDeleteBaseWithFile(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	outputObj := miniprop.NewMiniProp()
	inputObj := miniprop.NewMiniPropFromJsonReader(r.Body)
	articleId := inputObj.GetString("articleId", "")
	reqCheckErr := obj.OnDeleteArtRequest(w, r, obj, inputObj, outputObj)
	if reqCheckErr != nil {
		obj.OnDeleteArtFailed(w, r, obj, inputObj, outputObj)
		obj.HandleError(w, r, outputObj, 2001, reqCheckErr.Error())
		return
	}
	deleteFileErr := obj.GetBlobHandler().GetManager().DeleteBlobItemsWithPointerAtRecursiveMode(appengine.NewContext(r), obj.MakePath(articleId, ""))
	if deleteFileErr != nil {
		obj.OnDeleteArtFailed(w, r, obj, inputObj, outputObj)
		obj.HandleError(w, r, outputObj, 2002, deleteFileErr.Error())
		return
	}
	//	obj.tagMana.DeleteTagsFromOwner(appengine.NewContext(r), articleId)
	err := obj.GetManager().DeleteFromArticleIdWithPointer(ctx, articleId)
	if err != nil {
		obj.HandleError(w, r, outputObj, ErrorCodeNotFoundArticleId, "not found article")
		return
	}
	obj.OnDeleteArtSuccess(w, r, obj, inputObj, outputObj)
	w.WriteHeader(http.StatusOK)
}
