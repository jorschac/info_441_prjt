package widgetsrc

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
)

func (ctx *WidgetContext) UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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
	userid, errConv := strconv.Atoi(pathid)
	if errConv != nil {
		http.Error(w, "ChannelID not an integer", http.StatusBadRequest)
		return
	}

	usertoload, errGetUser := ctx.WStore.GetUser(userid)
	if errGetUser != nil {
		http.Error(w, errGetUser.Error(), http.StatusInternalServerError)
		return
	}

	up := &UserPage{}
	up.User = usertoload
	up.TBW = make([]*TextBoxWidget, 0)
	up.RTW = make([]*RecentTracksWidget, 0)
	up.TMW = make([]*TopMusicWidget, 0)
	up.SPW = make([]*SpotifyPlaylistWidget, 0)
	up.WidgetSocialInfo = make([]*WidgetSocial, 0)
	if userid == int(au.ID) {
		up.IsMe = true
	} else {
		up.IsMe = false
	}

	isFollowing, errCheck := ctx.WStore.CheckIfFollowing(int(au.ID), int(usertoload.ID))
	if errCheck != nil {
		http.Error(w, errCheck.Error(), http.StatusInternalServerError)
		return
	}
	up.IsFollowing = isFollowing

	followerCount, errGetFollowers := ctx.WStore.GetFollowerCount(userid)
	if errGetFollowers != nil {
		http.Error(w, errGetFollowers.Error(), http.StatusInternalServerError)
		return
	}
	up.NumFollowers = followerCount

	followingCount, errGetFollowing := ctx.WStore.GetFollowingCount(userid)
	if errGetFollowing != nil {
		http.Error(w, errGetFollowing.Error(), http.StatusInternalServerError)
		return
	}
	up.NumFollowing = followingCount

	dws, errGetWidgets := ctx.WStore.GetUserWidgets(int(usertoload.ID))
	if errGetWidgets != nil {
		http.Error(w, errGetWidgets.Error(), http.StatusInternalServerError)
		return
	}
	//user has no widgets
	if dws == nil {
		dws = make([]*DefaultWidgetInfo, 0)
	}

	wsi := make([]*WidgetSocial, 0)
	for _, dw := range dws {
		if dw.WidgetType == "Text Box" {
			tbw, errGet := ctx.WStore.GetTextBoxWidget(int(dw.WidgetID))
			if errGet != nil {
				http.Error(w, errGet.Error(), http.StatusInternalServerError)
				return
			}
			tbw.BaseInfo = dw
			up.TBW = append(up.TBW, tbw)
		} else if dw.WidgetType == "Recent Tracks" {
			rtw, errGet := ctx.WStore.GetRecentTracksWidget(int(dw.WidgetID))
			if errGet != nil {
				http.Error(w, errGet.Error(), http.StatusInternalServerError)
				return
			}
			rtw.BaseInfo = dw
			up.RTW = append(up.RTW, rtw)
		} else if dw.WidgetType == "Top Music" {
			tmw, errGet := ctx.WStore.GetTopMusicWidget(int(dw.WidgetID))
			if errGet != nil {
				http.Error(w, errGet.Error(), http.StatusInternalServerError)
				return
			}
			tmw.BaseInfo = dw
			up.TMW = append(up.TMW, tmw)
		} else if dw.WidgetType == "Spotify Playlist" {
			spw, errGet := ctx.WStore.GetSpotifyPlaylistWidget(int(dw.WidgetID))
			if errGet != nil {
				http.Error(w, errGet.Error(), http.StatusInternalServerError)
				return
			}
			spw.BaseInfo = dw
			up.SPW = append(up.SPW, spw)
		} else if dw.WidgetType == "Featured Music" {
			fmw, errGet := ctx.WStore.GetFeaturedMusicWidget(int(dw.WidgetID))
			if errGet != nil {
				http.Error(w, errGet.Error(), http.StatusInternalServerError)
				return
			}
			fmw.BaseInfo = dw
			up.FMW = append(up.FMW, fmw)
		}

		ws, errGetWS := ctx.WStore.GetWidgetSocialInfo(int(dw.WidgetID), int(usertoload.ID))
		if errGetWS != nil {
			http.Error(w, errGetWS.Error(), http.StatusInternalServerError)
			return
		}
		wsi = append(wsi, ws)
	}
	up.WidgetSocialInfo = wsi

	encoded, errEncode := json.Marshal(up)
	if errEncode != nil {
		http.Error(w, errEncode.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encoded)
}
