package hundler

import (
	"net/http"

	"github.com/firefirestyle/go.miniarticle/article"
	"github.com/firefirestyle/go.miniprop"
	"google.golang.org/appengine"
)

func (obj *ArticleHandler) HandleFind(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	cursor := values.Get("cursor")
	userName := values.Get("userName")
	target := values.Get("target")
	obj.HandleFindBase(w, r, cursor, userName, target)
}

func (obj *ArticleHandler) HandleFindBase(w http.ResponseWriter, r *http.Request, cursor, userName, target string) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	var foundObj *article.FoundArticles
	if userName != "" {
		foundObj = obj.GetManager().FindArticleFromUserName(ctx, userName, cursor, true)
	} else if target != "" {
		foundObj = obj.GetManager().FindArticleFromTarget(ctx, target, cursor, true)
	} else {
		foundObj = obj.GetManager().FindArticleWithNewOrder(ctx, cursor, true)
	}
	propObj.SetPropStringList("", "keys", foundObj.ArticleIds)
	propObj.SetPropString("", "cursorOne", foundObj.CursorOne)
	propObj.SetPropString("", "cursorNext", foundObj.CursorNext)
	w.Write(propObj.ToJson())
}
