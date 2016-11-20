package handler

import (
	"net/http"

	"strings"

	"github.com/firefirestyle/go.miniarticle/article"
	"github.com/firefirestyle/go.miniprop"
	//	"github.com/firefirestyle/go.minitag/tag"
	"google.golang.org/appengine"
)

func (obj *ArticleHandler) HandleFind(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	cursor := values.Get("cursor")
	userName := values.Get("userName")
	tag := []string{}
	props := map[string]string{}
	for k, v := range values {
		if strings.HasPrefix(k, "p-") {
			key := strings.Replace(k, "p-", "", 1)
			props[key] = v[0]
		} else if strings.HasPrefix(k, "t-") {
			tag = append(tag, v[0])
		}
	}
	obj.HandleFindBase(w, r, cursor, userName, props, tag)
}

func (obj *ArticleHandler) HandleFindBase(w http.ResponseWriter, r *http.Request, cursor, userName string, props map[string]string, tags []string) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	var foundObj *article.FoundArticles
	//if tag != "" {
	//	obj.HandleFindTagBase(w, r, cursor, tag)
	//} else {
	///Debug(ctx, ">>>>>>>>>>>>target ="+target)
	if len(tags) > 0 {
		foundObj = obj.GetManager().FindArticleFromTag(ctx, tags, cursor, true)
	} else if userName != "" {
		foundObj = obj.GetManager().FindArticleFromUserName(ctx, userName, cursor, true)
	} else if len(props) > 0 {
		foundObj = obj.GetManager().FindArticleFromProp(ctx, props, cursor, true)
	} else {
		foundObj = obj.GetManager().FindArticleWithNewOrder(ctx, cursor, true)
	}
	propObj.SetPropStringList("", "keys", foundObj.ArticleIds)
	propObj.SetPropString("", "cursorOne", foundObj.CursorOne)
	propObj.SetPropString("", "cursorNext", foundObj.CursorNext)
	w.Write(propObj.ToJson())
	//}
}
