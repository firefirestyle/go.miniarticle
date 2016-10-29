package hundler

import (
	"net/http"

	"github.com/firefirestyle/go.miniarticle/article"
	"github.com/firefirestyle/go.miniprop"
	"google.golang.org/appengine"
)

func (obj *ArticleHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	values := r.URL.Query()
	key := values.Get("key")
	articleId := values.Get("articleId")
	sign := values.Get("sign")
	mode := values.Get("m")
	//
	if key != "" {
		keyInfo := obj.GetManager().ExtractInfoFromStringId(key)
		articleId = keyInfo.ArticleId
		sign = keyInfo.Sign
	}
	var artObj *article.Article
	var err error
	if mode != "q" {
		artObj, err = obj.GetManager().GetArticleFromArticleId(ctx, articleId, sign)
	} else {
		artObj, err = obj.GetManager().GetArticleFromPointer(ctx, articleId)
	}
	if err != nil {
		HandleError(w, r, propObj, ErrorCodeNotFoundArticleId, "not found article")
		return
	}
	if mode != "q" {
		w.Header().Set("Cache-Control", "public, max-age=2592000")
	}
	w.Write(artObj.ToJsonPublicOnly())
}
