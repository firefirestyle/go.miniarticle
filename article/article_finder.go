package article

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	//	"google.golang.org/appengine/log"
	//"google.golang.org/appengine/blobstore"
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

func (obj *ArticleManager) FindArticleFromUserName(ctx context.Context, userName string, parentId string, state string, cursorSrc string, keyOnly bool) *FoundArticles {
	q := datastore.NewQuery(obj.kindArticle)
	q = q.Filter("ProjectId =", obj.projectId)
	q = q.Filter("UserName =", userName) ////
	q = q.Filter("ParentId =", parentId)
	if state != "" {
		q = q.Filter("State =", ArticleStatePublic) //
	}
	q = q.Order("-Updated").Limit(obj.limitOfFinding)
	return obj.FindArticleFromQuery(ctx, q, cursorSrc, keyOnly)
}

func (obj *ArticleManager) FindArticleFromTarget(ctx context.Context, targetName string, parentId string, state string, cursorSrc string, keyOnly bool) *FoundArticles {
	q := datastore.NewQuery(obj.kindArticle)
	q = q.Filter("ProjectId =", obj.projectId)
	q = q.Filter("Target =", targetName) ////
	q = q.Filter("ParentId =", parentId)
	if state != "" {
		q = q.Filter("State =", ArticleStatePublic) //
	}
	q = q.Order("-Updated").Limit(obj.limitOfFinding)
	return obj.FindArticleFromQuery(ctx, q, cursorSrc, keyOnly)
}

func (obj *ArticleManager) FindArticleWithNewOrder(ctx context.Context, cursorSrc string, keyOnly bool) *FoundArticles {
	q := datastore.NewQuery(obj.kindArticle)
	q = q.Filter("ProjectId =", obj.projectId)
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
			if keyOnly != false {
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
