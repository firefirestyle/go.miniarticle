package handler

import (
	"net/http"

	"github.com/firefirestyle/go.miniarticle/article"
	"github.com/firefirestyle/go.miniprop"
)

//
func (obj *ArticleHandler) AddOnNewRequest(f func(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error) {
	obj.onEvents.OnNewRequestList = append(obj.onEvents.OnNewRequestList, f)
}
func (obj *ArticleHandler) OnNewRequest(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnNewRequestList {
		e := f(w, r, h, i, o)
		if e != nil {
			return e
		}
	}
	return nil
}

//
func (obj *ArticleHandler) AddOnNewBeforeSave(f func(w http.ResponseWriter, r *http.Request, h *ArticleHandler, a *article.Article, i *miniprop.MiniProp, o *miniprop.MiniProp) error) {
	obj.onEvents.OnNewBeforeSaveList = append(obj.onEvents.OnNewBeforeSaveList, f)
}

func (obj *ArticleHandler) OnNewBeforeSave(w http.ResponseWriter, r *http.Request, h *ArticleHandler, a *article.Article, i *miniprop.MiniProp, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnNewBeforeSaveList {
		e := f(w, r, h, a, i, o)
		if e != nil {
			return e
		}
	}
	return nil
}

//
func (obj *ArticleHandler) AddOnNewArtFailed(f func(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp)) {
	obj.onEvents.OnNewArtFailedList = append(obj.onEvents.OnNewArtFailedList, f)
}
func (obj *ArticleHandler) OnNewArtFailed(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) {
	for _, f := range obj.onEvents.OnNewArtFailedList {
		f(w, r, h, i, o)
	}
}

//
func (obj *ArticleHandler) AddOnNewArtSucces(f func(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error) {
	obj.onEvents.OnNewArtSuccessList = append(obj.onEvents.OnNewArtSuccessList, f)
}

func (obj *ArticleHandler) OnNewArtSuccess(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnNewArtSuccessList {
		e := f(w, r, h, i, o)
		if e != nil {
			return e
		}
	}
	return nil
}

//
func (obj *ArticleHandler) AddOnUpdateRequest(f func(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error) {
	obj.onEvents.OnUpdateRequestList = append(obj.onEvents.OnUpdateRequestList, f)
}
func (obj *ArticleHandler) OnUpdateRequest(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnUpdateRequestList {
		e := f(w, r, h, i, o)
		if e != nil {
			return e
		}
	}
	return nil
}

//
func (obj *ArticleHandler) AddOnUpdateArtFailed(f func(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp)) {
	obj.onEvents.OnUpdateArtFailedList = append(obj.onEvents.OnUpdateArtFailedList, f)
}
func (obj *ArticleHandler) OnUpdateArtFailed(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) {
	for _, f := range obj.onEvents.OnUpdateArtFailedList {
		f(w, r, h, i, o)
	}
}

//
func (obj *ArticleHandler) AddOnUpdateArtSuccess(f func(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error) {
	obj.onEvents.OnUpdateArtSuccessList = append(obj.onEvents.OnUpdateArtSuccessList, f)
}

func (obj *ArticleHandler) OnUpdateArtSuccess(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnUpdateArtSuccessList {
		e := f(w, r, h, i, o)
		if e != nil {
			return e
		}
	}
	return nil
}

//
// GET
//
func (obj *ArticleHandler) AddOnGetArtRequest(f func(w http.ResponseWriter, r *http.Request, h *ArticleHandler, o *miniprop.MiniProp) error) {
	obj.onEvents.OnGetArtRequestList = append(obj.onEvents.OnGetArtRequestList, f)
}

func (obj *ArticleHandler) OnGetArtRequest(w http.ResponseWriter, r *http.Request, h *ArticleHandler, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnGetArtRequestList {
		e := f(w, r, h, o)
		if e != nil {
			return e
		}
	}
	return nil
}

//
func (obj *ArticleHandler) AddOnGetArtFailed(f func(w http.ResponseWriter, r *http.Request, h *ArticleHandler, o *miniprop.MiniProp)) {
	obj.onEvents.OnGetArtFailedList = append(obj.onEvents.OnGetArtFailedList, f)
}

func (obj *ArticleHandler) OnGetArtFailed(w http.ResponseWriter, r *http.Request, h *ArticleHandler, o *miniprop.MiniProp) {
	for _, f := range obj.onEvents.OnGetArtFailedList {
		f(w, r, h, o)
	}
}

//
func (obj *ArticleHandler) AddOnGetArtSuccess(f func(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *article.Article, o *miniprop.MiniProp) error) {
	obj.onEvents.OnGetArtSuccessList = append(obj.onEvents.OnGetArtSuccessList, f)
}
func (obj *ArticleHandler) OnGetArtSuccess(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *article.Article, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnGetArtSuccessList {
		e := f(w, r, h, i, o)
		if e != nil {
			return e
		}
	}
	return nil
}

//
// DELETE
//
func (obj *ArticleHandler) AddOnDeleteArtRequest(f func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp) error) {
	obj.onEvents.OnDeleteArtRequestList = append(obj.onEvents.OnDeleteArtRequestList, f)
}

func (obj *ArticleHandler) OnDeleteArtRequest(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnDeleteArtRequestList {
		e := f(w, r, h, i, o)
		if e != nil {
			return e
		}
	}
	return nil
}

//
func (obj *ArticleHandler) AddOnDeleteArtFailed(f func(w http.ResponseWriter, r *http.Request, handler *ArticleHandler, input *miniprop.MiniProp, output *miniprop.MiniProp)) {
	obj.onEvents.OnDeleteArtFailedList = append(obj.onEvents.OnDeleteArtFailedList, f)
}

func (obj *ArticleHandler) OnDeleteArtFailed(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) {
	for _, f := range obj.onEvents.OnDeleteArtFailedList {
		f(w, r, h, o, i)
	}
}

//
func (obj *ArticleHandler) AddOnDeleteArtSuccess(f func(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error) {
	obj.onEvents.OnDeleteArtSuccessList = append(obj.onEvents.OnDeleteArtSuccessList, f)
}
func (obj *ArticleHandler) OnDeleteArtSuccess(w http.ResponseWriter, r *http.Request, h *ArticleHandler, i *miniprop.MiniProp, o *miniprop.MiniProp) error {
	for _, f := range obj.onEvents.OnDeleteArtSuccessList {
		e := f(w, r, h, i, o)
		if e != nil {
			return e
		}
	}
	return nil
}
