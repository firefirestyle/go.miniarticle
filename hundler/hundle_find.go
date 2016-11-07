package hundler

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
	target := values.Get("target")
	tag := values.Get("tag")
	obj.HandleFindBase(w, r, cursor, userName, target, tag)
}

func (obj *ArticleHandler) HandleFindBase(w http.ResponseWriter, r *http.Request, cursor, userName, target, tag string) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	var foundObj *article.FoundArticles
	if tag != "" {
		obj.HandleFindTagBase(w, r, cursor, tag)
	} else {
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
}

func (obj *ArticleHandler) HandleFindTagBase(w http.ResponseWriter, r *http.Request, cursor, tag string) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	foundObj := obj.tagMana.FindTags(ctx, tag, "", cursor)
	keys := make([]string, 0)
	for _, v := range foundObj.Keys {
		keyinfo := obj.tagMana.GetKeyInfoFromStringId(v)
		if strings.HasPrefix(keyinfo.Value, "art://") {
			stringId := strings.Replace(keyinfo.Value, "art://", "", 1)
			keys = append(keys, stringId)
		}
	}
	propObj.SetPropStringList("", "keys", keys)
	propObj.SetPropString("", "cursorOne", foundObj.CursorOne)
	propObj.SetPropString("", "cursorNext", foundObj.CursorNext)
	w.Write(propObj.ToJson())
}
