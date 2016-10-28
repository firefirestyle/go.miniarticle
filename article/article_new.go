package article

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

func (obj *ArticleManager) GetArticleFromArticleId(ctx context.Context, articleId string, sign string) (*Article, error) {
	return obj.NewArticleFromGaeObjectKey(ctx, obj.NewGaeObjectKey(ctx, articleId, sign, ""))
}

func (obj *ArticleManager) NewArticleFromGaeObjectKey(ctx context.Context, key *datastore.Key) (*Article, error) {
	k := key
	//
	//
	artObjFMem, errNewFMem := obj.NewArticleFromMemcache(ctx, k.StringID())
	if errNewFMem == nil {
		log.Infof(ctx, ">>>> new article Obj from memcache")
		return artObjFMem, nil
	}
	//
	//
	var gaeObj GaeObjectArticle
	err := datastore.Get(ctx, k, &gaeObj)
	if err != nil {
		return nil, err
	}
	//
	//
	return obj.NewArticleFromGaeObject(ctx, k, &gaeObj), nil
}

func (obj *ArticleManager) NewArticleFromMemcache(ctx context.Context, stringId string) (*Article, error) {
	ret := new(Article)
	ret.gaeObject = new(GaeObjectArticle)
	idInfo := obj.ExtractInfoFromStringId(stringId)
	ret.gaeObjectKey = obj.NewGaeObjectKey(ctx, idInfo.ArticleId, idInfo.Sign, "")
	ret.kind = obj.kindArticle
	artObjSource, errGetFMem := memcache.Get(ctx, ret.gaeObjectKey.StringID())
	if errGetFMem != nil {
		return nil, errGetFMem
	}
	errSetParam := ret.SetParamFromsJson(ctx, string(artObjSource.Value))

	return ret, errSetParam
}

func (obj *ArticleManager) NewArticleFromGaeObject(ctx context.Context, gaeKey *datastore.Key, gaeObj *GaeObjectArticle) *Article {
	ret := new(Article)
	ret.gaeObject = gaeObj
	ret.gaeObjectKey = gaeKey
	ret.kind = obj.kindArticle
	//
	//

	return ret
}

func (obj *ArticleManager) NewArticleFromArticle(ctx context.Context, baseArtObj *Article, sign string) *Article {
	//
	ret := new(Article)
	ret.kind = obj.kindArticle
	ret.gaeObject = &GaeObjectArticle{}
	ret.gaeObjectKey = obj.NewGaeObjectKey(ctx, baseArtObj.GetArticleId(), sign, "")

	//
	baseArtData := baseArtObj.ToMap()
	baseArtData[TypeSign] = sign
	ret.SetParamFromsMap(baseArtData)
	return ret
}

func (obj *ArticleManager) NewArticle(ctx context.Context, userName string, sign string) *Article {
	created := time.Now()
	var secretKey string
	var articleId string
	var key *datastore.Key
	var art GaeObjectArticle
	for {
		secretKey = obj.makeRandomId() + obj.makeRandomId()
		articleId = obj.makeArticleId(userName, created, secretKey)
		//stringId := obj.makeStringId(articleId, sign)
		//
		Debug(ctx, "<NewArticle>"+articleId+"::"+sign)
		key = obj.NewGaeObjectKey(ctx, articleId, sign, "")
		err := datastore.Get(ctx, key, &art)
		if err != nil {
			break
		}
	}
	//
	ret := new(Article)
	ret.kind = obj.kindArticle
	ret.gaeObject = &art
	ret.gaeObjectKey = key
	ret.gaeObject.ProjectId = obj.projectId
	ret.gaeObject.UserName = userName
	ret.gaeObject.Sign = sign
	ret.gaeObject.Created = created
	ret.gaeObject.Updated = created
	ret.gaeObject.SecretKey = secretKey
	ret.gaeObject.ArticleId = articleId
	//
	//
	//

	return ret
}

func (obj *ArticleManager) NewGaeObjectKey(ctx context.Context, articleId string, sign string, kind string) *datastore.Key {
	if kind == "" {
		kind = obj.kindArticle
	}
	return datastore.NewKey(ctx, kind, obj.makeStringId(articleId, sign), 0, nil)
}
