package hundler

import (
	"net/http"

	"strings"

	"google.golang.org/appengine"
)

func (obj *ArticleHandler) GetArticleIdFromDir(dir string) string {
	if false == strings.HasPrefix(dir, "/art/") {
		return ""
	}
	t1 := strings.Replace(dir, "/art/", "", 1)
	t2 := strings.Index(t1, "/")
	if t2 == -1 {
		t2 = len(t1)
	}

	return t1[0:t2]
}

func (obj *ArticleHandler) GetDirFromDir(dir string) string {
	if false == strings.HasPrefix(dir, "/art/") {
		return ""
	}
	t1 := strings.Replace(dir, "/art/", "", 1)
	t2 := strings.Index(t1, "/")
	if t2 == -1 {
		t2 = 0
	}

	return t1[t2:len(t1)]
}

func (obj *ArticleHandler) HandleBlobRequestToken(w http.ResponseWriter, r *http.Request) {
	//
	// load param from json
	articleId := r.URL.Query().Get("articleId")
	dir := r.URL.Query().Get("dir")
	name := r.URL.Query().Get("file")
	//
	// todo check articleId

	//
	//
	obj.blobHundler.HandleBlobRequestTokenFromParams(w, r, "/art/"+articleId+"/"+dir, name)
}

func (obj *ArticleHandler) HandleBlobUpdated(w http.ResponseWriter, r *http.Request) {
	//
	ctx := appengine.NewContext(r)
	Debug(ctx, "callbeck AAAA")
	obj.blobHundler.HandleUploaded(w, r)
}
