package handler

import (
	"net/http"

	//	"strings"

	"github.com/firefirestyle/go.miniarticle/article"
	"github.com/firefirestyle/go.miniprop"
	//	"github.com/firefirestyle/go.minitag/tag"
	"google.golang.org/appengine"
)

func (obj *ArticleHandler) HandleFind(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	cursor := values.Get("cursor")
	userName := values.Get("userName")
	target := values.Get("target")
	tag := values.Get("tag")
	obj.HandleFindBase(w, r, cursor, userName, target, tag)
}

func (obj *ArticleHandler) HandleFindBase(w http.ResponseWriter, r *http.Request, cursor, userName, target, tag string) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	var foundObj *article.FoundArticles
	//if tag != "" {
	//	obj.HandleFindTagBase(w, r, cursor, tag)
	//} else {
	Debug(ctx, ">>>>>>>>>>>>target ="+target)
	if tag != "" {
		foundObj = obj.GetManager().FindArticleFromTag(ctx, []string{tag}, cursor, true)
	} else if userName != "" {
		foundObj = obj.GetManager().FindArticleFromUserName(ctx, userName, cursor, true)
	} else if target != "" {
		foundObj = obj.GetManager().FindArticleFromProp(ctx, "target", target, cursor, true)
	} else {
		foundObj = obj.GetManager().FindArticleWithNewOrder(ctx, cursor, true)
	}
	propObj.SetPropStringList("", "keys", foundObj.ArticleIds)
	propObj.SetPropString("", "cursorOne", foundObj.CursorOne)
	propObj.SetPropString("", "cursorNext", foundObj.CursorNext)
	w.Write(propObj.ToJson())
	//}
}
