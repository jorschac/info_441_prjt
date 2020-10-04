package widgetsrc

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

func (ctx *WidgetContext) TextBoxHandlerCreate(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, "Incorrect formating on line 20", http.StatusInternalServerError)
			return
		}

		tbw := &TextBoxWidget{}
		errDecode = json.NewDecoder(r.Body).Decode(&tbw)
		if errDecode != nil {
			http.Error(w, "Incorrect formmating on line 27", http.StatusInternalServerError)
			return
		}

		dw, errDB := ctx.WStore.CreateDefaultWidget(tbw.BaseInfo)
		if errDecode != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}
		tbw.BaseInfo = dw
		tbw.BaseInfo.WidgetType = "Text Box"

		errDB = ctx.WStore.CreateTextBoxWidget(tbw)
		if errDB != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}

		encoded, errEncode := json.Marshal(tbw)
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

func (ctx *WidgetContext) TextBoxHandler(w http.ResponseWriter, r *http.Request) {
	isAuth := r.Header.Get("X-User")
	if len(isAuth) == 0 {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}
	byteAuth := []byte(isAuth)
	au := &AuthUser{}
	errDecode := json.Unmarshal(byteAuth, &au)
	if errDecode != nil {
		http.Error(w, "Incorrect formating on line 20", http.StatusInternalServerError)
		return
	}

	pathid := path.Base(r.URL.Path)
	split := strings.Split(pathid, "&")
	pathid = split[0]
	wid, errConv := strconv.Atoi(pathid)
	if errConv != nil {
		http.Error(w, "ChannelID not an integer", http.StatusBadRequest)
		return
	}

	dw, errDB := ctx.WStore.GetDefaultWidget(wid)
	if errDecode != nil {
		http.Error(w, errDB.Error(), http.StatusInternalServerError)
		return
	}
	if dw == nil {
		http.Error(w, "WidgetID does not exist", http.StatusBadRequest)
		return
	}
	if dw.UserID != au.ID {
		http.Error(w, "Cannot access widget information for other users.", http.StatusForbidden)
		return
	}
	tbw, errDB := ctx.WStore.GetTextBoxWidget(wid)
	if errDB != nil {
		http.Error(w, errDB.Error(), http.StatusInternalServerError)
		return
	}
	if tbw == nil {
		http.Error(w, "Not a textbox widget", http.StatusBadRequest)
		return
	}
	tbw.BaseInfo = dw

	if r.Method == http.MethodGet {
		encoded, errEncode := json.Marshal(tbw)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encoded)

	} else if r.Method == http.MethodPatch {
		tbwUpdates := &TextBoxWidget{}
		errDecode = json.NewDecoder(r.Body).Decode(&tbwUpdates)
		if errDecode != nil {
			http.Error(w, errDecode.Error(), http.StatusInternalServerError)
			return
		}

		if tbwUpdates.BaseInfo.Location != nil {
			tbw.BaseInfo.Location = tbwUpdates.BaseInfo.Location
		}
		if tbwUpdates.Text != "" {
			tbw.Text = tbwUpdates.Text
		}
		tbw.BaseInfo.UpdatedAt = time.Now()

		tbwUpdates.BaseInfo.UserID = au.ID
		tbwUpdates.BaseInfo.WidgetID = dw.WidgetID
		tbwUpdates.BaseInfo.UpdatedAt = tbw.BaseInfo.UpdatedAt
		errEditDefault := ctx.WStore.EditDefaultWidget(tbwUpdates.BaseInfo)
		if errEditDefault != nil {
			http.Error(w, errEditDefault.Error(), http.StatusInternalServerError)
			return
		}

		errDB = ctx.WStore.EditTextBoxWidget(tbwUpdates)
		if errDB != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}

		encoded, errEncode := json.Marshal(tbw)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encoded)

	} else if r.Method == http.MethodDelete {
		errDelete := ctx.WStore.DeleteTextBoxWidget(wid)
		if errDelete != nil {
			http.Error(w, errDelete.Error(), http.StatusInternalServerError)
			return
		}

		errDelete = ctx.WStore.DeleteDefaultWidget(wid)
		if errDelete != nil {
			http.Error(w, errDelete.Error(), http.StatusInternalServerError)
			return
		}
		//will need to delete comments and likes as well when I implement that
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Widget deleted successfully"))

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
