package widgetsrc

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
)

func (ctx *WidgetContext) WidgetLikeHandler(w http.ResponseWriter, r *http.Request) {
	allowedMethods := http.MethodPost + http.MethodDelete
	if !strings.Contains(allowedMethods, r.Method) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	isAuth := r.Header.Get("X-User")
	if len(isAuth) == 0 {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	byteAuth := []byte(isAuth)
	au := &AuthUser{}
	errDecode := json.Unmarshal(byteAuth, &au)
	if errDecode != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	path := r.URL.Path
	split := strings.Split(path, "/")
	wid, errConv := strconv.Atoi(split[3])
	if errConv != nil {
		http.Error(w, "WidgetID not an integer", http.StatusBadRequest)
		return
	}

	wl := &WidgetLike{}
	wl.WidgetID = int64(wid)
	wl.UserID = au.ID

	if r.Method == http.MethodPost {
		wl, errLike := ctx.WStore.CreateLike(wl)
		if errLike != nil {
			http.Error(w, "Error getting user", http.StatusInternalServerError)
			return
		}

		encoded, errEncode := json.Marshal(wl)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(encoded)
		//delete
	} else {
		errLike := ctx.WStore.DeleteLike(wl)
		if errLike != nil {
			http.Error(w, "Error removing like", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Unliked"))
	}
}
func (ctx *WidgetContext) WidgetCommentHandlerCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		isAuth := r.Header.Get("X-User")
		if len(isAuth) == 0 {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}
		byteAuth := []byte(isAuth)
		au := &AuthUser{}
		errDecode := json.Unmarshal(byteAuth, &au)
		if errDecode != nil {
			http.Error(w, "Error getting user", http.StatusInternalServerError)
			return
		}

		path := r.URL.Path
		split := strings.Split(path, "/")
		wid, errConv := strconv.Atoi(split[3])
		if errConv != nil {
			http.Error(w, "commentID is not an integer", http.StatusBadRequest)
			return
		}

		wc := &WidgetComment{}
		errDecode = json.NewDecoder(r.Body).Decode(&wc)
		if errDecode != nil {
			http.Error(w, "Bad input", http.StatusBadRequest)
			return
		}

		wc.WidgetID = int64(wid)
		wc.UserID = au.ID

		wc, errComment := ctx.WStore.CreateComment(wc)
		if errComment != nil {
			http.Error(w, "Error getting user", http.StatusInternalServerError)
			return
		}

		encoded, errEncode := json.Marshal(wc)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(encoded)

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func (ctx *WidgetContext) WidgetCommentHandler(w http.ResponseWriter, r *http.Request) {
	allowedMethods := http.MethodGet + http.MethodPatch + http.MethodDelete
	if !strings.Contains(allowedMethods, r.Method) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	isAuth := r.Header.Get("X-User")
	if len(isAuth) == 0 {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	byteAuth := []byte(isAuth)
	au := &AuthUser{}
	errDecode := json.Unmarshal(byteAuth, &au)
	if errDecode != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	pathid := path.Base(r.URL.Path)
	split := strings.Split(pathid, "&")
	pathid = split[0]
	cid, errConv := strconv.Atoi(pathid)
	if errConv != nil {
		http.Error(w, "commentID is not an integer", http.StatusBadRequest)
		return
	}

	wc, errGet := ctx.WStore.GetComment(cid)
	if errGet != nil {
		http.Error(w, "Error getting comment", http.StatusInternalServerError)
		return
	}
	if wc == nil {
		http.Error(w, "CommentID does not exist", http.StatusBadRequest)
		return
	}
	if wc.UserID != au.ID {
		http.Error(w, "Cannot access other user's comments", http.StatusForbidden)
		return
	}
	if r.Method == http.MethodGet {
		encoded, errEncode := json.Marshal(wc)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encoded)
	} else if r.Method == http.MethodPatch {
		wcUpdates := &WidgetComment{}
		errDecode = json.NewDecoder(r.Body).Decode(&wcUpdates)
		if errDecode != nil {
			http.Error(w, errDecode.Error(), http.StatusInternalServerError)
			return
		}

		wc.Comment = wcUpdates.Comment
		errComment := ctx.WStore.EditComment(wc)
		if errComment != nil {
			http.Error(w, errComment.Error(), http.StatusInternalServerError)
			return
		}

		encoded, errEncode := json.Marshal(wc)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encoded)
		//delete
	} else {
		errComment := ctx.WStore.DeleteComment(cid)
		if errComment != nil {
			http.Error(w, "Error deleting comment", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Comment removed."))
	}
}

func (ctx *WidgetContext) WidgetCommentLikeHandler(w http.ResponseWriter, r *http.Request) {
	allowedMethods := http.MethodPost + http.MethodDelete
	if !strings.Contains(allowedMethods, r.Method) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	isAuth := r.Header.Get("X-User")
	if len(isAuth) == 0 {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	byteAuth := []byte(isAuth)
	au := &AuthUser{}
	errDecode := json.Unmarshal(byteAuth, &au)
	if errDecode != nil {
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	path := r.URL.Path
	split := strings.Split(path, "/")
	cid, errConv := strconv.Atoi(split[3])
	if errConv != nil {
		http.Error(w, "WidgetID not an integer", http.StatusBadRequest)
		return
	}

	wcl := &WidgetCommentLike{}
	wcl.WidgetCommentID = int64(cid)
	wcl.UserID = au.ID

	if r.Method == http.MethodPost {
		wcl, errLike := ctx.WStore.CreateCommentLike(wcl)
		if errLike != nil {
			http.Error(w, "Error getting user", http.StatusInternalServerError)
			return
		}

		encoded, errEncode := json.Marshal(wcl)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write(encoded)
		//delete
	} else {
		errLike := ctx.WStore.DeleteCommentLike(wcl)
		if errLike != nil {
			http.Error(w, "Error removing like", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Unliked"))
	}

}
