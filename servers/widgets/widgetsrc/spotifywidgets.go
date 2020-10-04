package widgetsrc

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

func (ctx *WidgetContext) SpotifyPlaylistHandlerCreate(w http.ResponseWriter, r *http.Request) {
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

		spw := &SpotifyPlaylistWidget{}
		errDecode = json.NewDecoder(r.Body).Decode(&spw)
		if errDecode != nil {
			http.Error(w, "Bad input", http.StatusBadRequest)
			return
		}

		spw.BaseInfo.UserID = au.ID
		dw, errDB := ctx.WStore.CreateDefaultWidget(spw.BaseInfo)
		if errDB != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}
		spw.BaseInfo = dw
		spw.BaseInfo.WidgetType = "Spotify Playlist"

		errDB = ctx.WStore.CreateSpotifyPlaylistWidget(spw)
		if errDB != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}

		encoded, errEncode := json.Marshal(spw)
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

func (ctx *WidgetContext) SpotifyPlaylistHandler(w http.ResponseWriter, r *http.Request) {
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

	spw, errDB := ctx.WStore.GetSpotifyPlaylistWidget(wid)
	if errDB != nil {
		http.Error(w, errDB.Error(), http.StatusInternalServerError)
		return
	}
	if spw == nil {
		http.Error(w, "Wrong widget type", http.StatusBadRequest)
		return
	}
	spw.BaseInfo = dw

	if r.Method == http.MethodGet {
		encoded, errEncode := json.Marshal(spw)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encoded)

	} else if r.Method == http.MethodPatch {
		spwUpdates := &SpotifyPlaylistWidget{}
		errDecode = json.NewDecoder(r.Body).Decode(&spwUpdates)
		if errDecode != nil {
			http.Error(w, errDecode.Error(), http.StatusInternalServerError)
			return
		}

		if spwUpdates.BaseInfo.Location != nil {
			spw.BaseInfo.Location = spwUpdates.BaseInfo.Location
		}
		if spwUpdates.NumTracks != 0 {
			spw.NumTracks = spwUpdates.NumTracks
		}
		if spwUpdates.Description != "" {
			spw.Description = spwUpdates.Description
		}
		if spwUpdates.SpotifyURI != "" {
			spw.SpotifyURI = spwUpdates.SpotifyURI
		}
		if spwUpdates.PlaylistOrder != spw.PlaylistOrder {
			spw.PlaylistOrder = spwUpdates.PlaylistOrder
		}
		spw.BaseInfo.UpdatedAt = time.Now()

		spwUpdates.BaseInfo.UserID = au.ID
		spwUpdates.BaseInfo.WidgetID = dw.WidgetID
		spwUpdates.BaseInfo.UpdatedAt = spw.BaseInfo.UpdatedAt
		errEditDefault := ctx.WStore.EditDefaultWidget(spwUpdates.BaseInfo)
		if errEditDefault != nil {
			http.Error(w, errEditDefault.Error(), http.StatusInternalServerError)
			return
		}

		errDB = ctx.WStore.EditSpotifyPlaylistWidget(spwUpdates)
		if errDB != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}
		encoded, errEncode := json.Marshal(spw)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encoded)

		//delete
	} else {
		errDelete := ctx.WStore.DeleteSpotifyPlaylistWidget(wid)
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

func (ctx *WidgetContext) FeaturedMusicHandlerCreate(w http.ResponseWriter, r *http.Request) {
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

		fmw := &FeaturedMusicWidget{}
		errDecode = json.NewDecoder(r.Body).Decode(&fmw)
		if errDecode != nil {
			http.Error(w, "Bad input", http.StatusBadRequest)
			return
		}

		fmw.BaseInfo.UserID = au.ID
		dw, errDB := ctx.WStore.CreateDefaultWidget(fmw.BaseInfo)
		if errDB != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}
		fmw.BaseInfo = dw
		fmw.BaseInfo.WidgetType = "Featured Music"

		errDB = ctx.WStore.CreateFeaturedMusicWidget(fmw)
		if errDB != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}

		encoded, errEncode := json.Marshal(fmw)
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

func (ctx *WidgetContext) FeaturedMusicHandler(w http.ResponseWriter, r *http.Request) {
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

	fmw, errDB := ctx.WStore.GetFeaturedMusicWidget(wid)
	if errDB != nil {
		http.Error(w, errDB.Error(), http.StatusInternalServerError)
		return
	}
	if fmw == nil {
		http.Error(w, "Not a featured music widget", http.StatusBadRequest)
		return
	}
	fmw.BaseInfo = dw

	if r.Method == http.MethodGet {
		encoded, errEncode := json.Marshal(fmw)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encoded)

	} else if r.Method == http.MethodPatch {
		fmwUpdates := &FeaturedMusicWidget{}
		errDecode = json.NewDecoder(r.Body).Decode(&fmwUpdates)
		if errDecode != nil {
			http.Error(w, errDecode.Error(), http.StatusInternalServerError)
			return
		}

		if fmwUpdates.BaseInfo.Location != nil {
			fmw.BaseInfo.Location = fmwUpdates.BaseInfo.Location
		}

		if fmwUpdates.Description != "" {
			fmw.Description = fmwUpdates.Description
		}

		if fmwUpdates.Type != "" {
			fmw.Type = fmwUpdates.Type
		}

		if fmwUpdates.MusicName != "" {
			fmw.MusicName = fmwUpdates.MusicName
		}
		fmw.BaseInfo.UpdatedAt = time.Now()

		fmwUpdates.BaseInfo.UserID = au.ID
		fmwUpdates.BaseInfo.WidgetID = dw.WidgetID
		fmwUpdates.BaseInfo.UpdatedAt = fmw.BaseInfo.UpdatedAt
		errEditDefault := ctx.WStore.EditDefaultWidget(fmwUpdates.BaseInfo)
		if errEditDefault != nil {
			http.Error(w, errEditDefault.Error(), http.StatusInternalServerError)
			return
		}

		errDB = ctx.WStore.EditFeaturedMusicWidget(fmwUpdates)
		if errDB != nil {
			http.Error(w, errDB.Error(), http.StatusInternalServerError)
			return
		}
		encoded, errEncode := json.Marshal(fmw)
		if errEncode != nil {
			http.Error(w, "Error encoding user to JSON", http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(encoded)

		//delete
	} else {
		errDelete := ctx.WStore.DeleteFeaturedMusicWidget(wid)
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
