package hundler

import (
	"net/http"

	"google.golang.org/appengine"
)

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
