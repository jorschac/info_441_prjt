package widgetsrc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

func (ctx *WidgetContext) TopTracksHandler(w http.ResponseWriter, r *http.Request) {
	username := "Xergoyf"
	apikey := "32fab5088d004b7dadb721e0f381ca71"
	from := 1587772800
	to := time.Now().Unix()
	base := "http://ws.audioscrobbler.com/2.0/?method=user.getweeklyartistchart&user=" + username
	base2 := "&from" + string(from)
	base3 := "&to" + string(to)
	base4 := "&limit=50&api_key=" + apikey + "&format=json"

	concat := base + base2 + base3 + base4

	response, errGet := http.Get(concat)
	if errGet != nil {
		http.Error(w, errGet.Error(), http.StatusInternalServerError)
		return
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)
	swag := responseObject.Weeklyartistchart.Artists[2].Name
	toReturn := responseObject.Weeklyartistchart.Artists
	fmt.Println(swag)

	encoded, errEncode := json.Marshal(toReturn)
	if errEncode != nil {
		http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(encoded)
}

func (ctx *WidgetContext) RecentTracksHandlerCreate(w http.ResponseWriter, r *http.Request) {
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

		rtw := &RecentTracksWidget{}
		errDecode = json.NewDecoder(r.Body).Decode(&rtw)
		if errDecode != nil {
			http.Error(w, "Bad input", http.StatusBadRequest)
			return
		}

		rtw.BaseInfo.UserID = au.ID
		dw, errDB := ctx.WStore.CreateDefaultWidget(rtw.BaseInfo)
		if errDB != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}
		rtw.BaseInfo = dw
		rtw.BaseInfo.WidgetType = "Recent Tracks"

		errDB = ctx.WStore.CreateRecentTracksWidget(rtw)
		if errDB != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}

		encoded, errEncode := json.Marshal(rtw)
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

func (ctx *WidgetContext) RecentTracksHandler(w http.ResponseWriter, r *http.Request) {
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

	rtw, errDB := ctx.WStore.GetRecentTracksWidget(wid)
	if errDB != nil {
		http.Error(w, errDB.Error(), http.StatusInternalServerError)
		return
	}
	if rtw == nil {
		http.Error(w, "Not a recent tracks widget", http.StatusBadRequest)
		return
	}
	rtw.BaseInfo = dw

	if r.Method == http.MethodGet {
		encoded, errEncode := json.Marshal(rtw)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encoded)

	} else if r.Method == http.MethodPatch {
		rtwUpdates := &RecentTracksWidget{}
		errDecode = json.NewDecoder(r.Body).Decode(&rtwUpdates)
		if errDecode != nil {
			http.Error(w, errDecode.Error(), http.StatusInternalServerError)
			return
		}

		if rtwUpdates.BaseInfo.Location != nil {
			rtw.BaseInfo.Location = rtwUpdates.BaseInfo.Location
		}
		if rtwUpdates.NumTracks != 0 {
			rtw.NumTracks = rtwUpdates.NumTracks
		}
		if rtwUpdates.Lastfm != "" {
			rtw.Lastfm = rtwUpdates.Lastfm
		}
		if rtwUpdates.Description != "" {
			rtw.Description = rtwUpdates.Description
		}
		rtw.BaseInfo.UpdatedAt = time.Now()

		rtwUpdates.BaseInfo.UserID = au.ID
		rtwUpdates.BaseInfo.WidgetID = dw.WidgetID
		rtwUpdates.BaseInfo.UpdatedAt = rtw.BaseInfo.UpdatedAt
		errEditDefault := ctx.WStore.EditDefaultWidget(rtwUpdates.BaseInfo)
		if errEditDefault != nil {
			http.Error(w, errEditDefault.Error(), http.StatusInternalServerError)
			return
		}

		errDB = ctx.WStore.EditRecentTracksWidget(rtwUpdates)
		if errDB != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}

		encoded, errEncode := json.Marshal(rtw)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encoded)

		//delete
	} else {
		errDelete := ctx.WStore.DeleteRecentTracksWidget(wid)
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
	}
}

func (ctx *WidgetContext) TopMusicHandlerCreate(w http.ResponseWriter, r *http.Request) {
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

		tmw := &TopMusicWidget{}
		errDecode = json.NewDecoder(r.Body).Decode(&tmw)
		if errDecode != nil {
			http.Error(w, "Bad input", http.StatusBadRequest)
			return
		}

		tmw.BaseInfo.UserID = au.ID
		dw, errDB := ctx.WStore.CreateDefaultWidget(tmw.BaseInfo)
		if errDB != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}
		tmw.BaseInfo = dw
		tmw.BaseInfo.WidgetType = "Top Music"

		errDB = ctx.WStore.CreateTopMusicWidget(tmw)
		if errDB != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}

		encoded, errEncode := json.Marshal(tmw)
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

func (ctx *WidgetContext) TopMusicHandler(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "error getting user", http.StatusInternalServerError)
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

	tmw, errDB := ctx.WStore.GetTopMusicWidget(wid)
	if errDB != nil {
		http.Error(w, errDB.Error(), http.StatusInternalServerError)
		return
	}
	if tmw == nil {
		http.Error(w, "Not a recent tracks widget", http.StatusBadRequest)
		return
	}
	tmw.BaseInfo = dw

	if r.Method == http.MethodGet {
		encoded, errEncode := json.Marshal(tmw)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encoded)

	} else if r.Method == http.MethodPatch {
		tmwUpdates := &TopMusicWidget{}
		errDecode = json.NewDecoder(r.Body).Decode(&tmwUpdates)
		if errDecode != nil {
			http.Error(w, errDecode.Error(), http.StatusInternalServerError)
			return
		}

		if tmwUpdates.BaseInfo.Location != nil {
			tmw.BaseInfo.Location = tmwUpdates.BaseInfo.Location
		}
		if tmwUpdates.NumTracks != 0 {
			tmw.NumTracks = tmwUpdates.NumTracks
		}
		if tmwUpdates.Lastfm != "" {
			tmw.Lastfm = tmwUpdates.Lastfm
		}
		if tmwUpdates.Description != "" {
			tmw.Description = tmwUpdates.Description
		}
		if tmwUpdates.Type != "" {
			tmw.Type = tmwUpdates.Type
		}
		if tmw.TimePeriod != 0 {
			tmw.TimePeriod = tmwUpdates.TimePeriod
		}
		tmw.BaseInfo.UpdatedAt = time.Now()

		tmwUpdates.BaseInfo.UserID = au.ID
		tmwUpdates.BaseInfo.WidgetID = dw.WidgetID
		tmwUpdates.BaseInfo.UpdatedAt = tmw.BaseInfo.UpdatedAt
		errEditDefault := ctx.WStore.EditDefaultWidget(tmwUpdates.BaseInfo)
		if errEditDefault != nil {
			http.Error(w, errEditDefault.Error(), http.StatusInternalServerError)
			return
		}

		errDB = ctx.WStore.EditTopMusicWidget(tmwUpdates)
		if errDB != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}

		encoded, errEncode := json.Marshal(tmw)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encoded)

		//delte
	} else {
		errDelete := ctx.WStore.DeleteTopMusicWidget(wid)
		if errDelete != nil {
			http.Error(w, errDelete.Error(), http.StatusInternalServerError)
			return
		}

		errDelete = ctx.WStore.DeleteDefaultWidget(wid)
		if errDelete != nil {
			http.Error(w, errDelete.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Widget deleted successfully"))
	}
}
