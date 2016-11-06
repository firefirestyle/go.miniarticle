package hundler

import (
	"net/http"

	//	"github.com/firefirestyle/go.miniarticle/article"
	"github.com/firefirestyle/go.miniprop"
	"google.golang.org/appengine"
)

func (obj *ArticleHandler) HandleDeleteBase(w http.ResponseWriter, r *http.Request, articleId string) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	err := obj.GetManager().DeleteFromArticleIdWithPointer(ctx, articleId)
	if err != nil {
		obj.HandleError(w, r, propObj, ErrorCodeNotFoundArticleId, "not found article")
		return
	}
	w.WriteHeader(http.StatusOK)
}

//
// you must to delete file before call this method, if there are many articleid's file.
//
func (obj *ArticleHandler) HandleDeleteBaseWithFile(w http.ResponseWriter, r *http.Request, articleId string) {
	obj.GetBlobHandler().GetManager().DeleteBlobItemsFormOnwer(appengine.NewContext(r), articleId)
	obj.HandleDeleteBase(w, r, articleId)
}
