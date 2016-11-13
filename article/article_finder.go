package article

import (
	"github.com/firefirestyle/go.miniprop"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

/*
https://cloud.google.com/appengine/docs/go/config/indexconfig#updating_indexes

*/

type FoundArticles struct {
	Articles   []*Article
	ArticleIds []string
	CursorOne  string
	CursorNext string
}

func (obj *ArticleManager) FindArticleFromUserName(ctx context.Context, userName string, cursorSrc string, keyOnly bool) *FoundArticles {
	q := datastore.NewQuery(obj.kindArticle)
	q = q.Filter("RootGroup =", obj.projectId)
	q = q.Filter("UserName =", userName) ////
	q = q.Order("-Updated").Limit(obj.limitOfFinding)
	return obj.FindArticleFromQuery(ctx, q, cursorSrc, keyOnly)
}

func (obj *ArticleManager) FindArticleFromTag(ctx context.Context, tags []string, cursorSrc string, keyOnly bool) *FoundArticles {
	q := datastore.NewQuery(obj.kindArticle)
	q = q.Filter("RootGroup =", obj.projectId)
	for _, tag := range tags {
		q = q.Filter("Tags.Tag =", tag) ////
	}
	q = q.Order("-Updated").Limit(obj.limitOfFinding)
	return obj.FindArticleFromQuery(ctx, q, cursorSrc, keyOnly)
}

func (obj *ArticleManager) FindArticleFromProp(ctx context.Context, name string, value string, cursorSrc string, keyOnly bool) *FoundArticles {
	Debug(ctx, "======> Find Article target")
	p := miniprop.NewMiniProp()
	p.SetString(name, value)
	v := string(p.ToJson())
	q := datastore.NewQuery(obj.kindArticle)
	q = q.Filter("RootGroup =", obj.projectId)
	q = q.Filter("Props.Value =", v) ////
	q = q.Order("-Updated").Limit(obj.limitOfFinding)
	return obj.FindArticleFromQuery(ctx, q, cursorSrc, keyOnly)
}

func (obj *ArticleManager) FindArticleWithNewOrder(ctx context.Context, cursorSrc string, keyOnly bool) *FoundArticles {
	q := datastore.NewQuery(obj.kindArticle)
	q = q.Filter("RootGroup =", obj.projectId)
	//	q = q.Order("-Updated").Limit(obj.limitOfFinding)

	return obj.FindArticleFromQuery(ctx, q, cursorSrc, keyOnly)
}

func (obj *ArticleManager) FindArticleFromQuery(ctx context.Context, q *datastore.Query, cursorSrc string, keyOnly bool) *FoundArticles {
	cursor := obj.newCursorFromSrc(cursorSrc)
	if cursor != nil {
		q = q.Start(*cursor)
	}
	q = q.KeysOnly()
	founds := q.Run(ctx)

	var retUser []*Article
	var articleIds []string = make([]string, 0)

	var cursorNext string = ""
	var cursorOne string = ""
	for i := 0; ; i++ {
		key, err := founds.Next(nil)

		if err != nil || err == datastore.Done {
			break
		} else {
			articleIds = append(articleIds, key.StringID())
			if keyOnly == false {
				userObj, errNewUserObj := obj.NewArticleFromGaeObjectKey(ctx, key)
				if errNewUserObj == nil {
					retUser = append(retUser, userObj)
				}
			}
		}
		if i == 0 {
			cursorOne = obj.makeCursorSrc(founds)
		}
	}
	cursorNext = obj.makeCursorSrc(founds)
	return &FoundArticles{
		Articles:   retUser,
		ArticleIds: articleIds,
		CursorNext: cursorNext,
		CursorOne:  cursorOne,
	}
}
