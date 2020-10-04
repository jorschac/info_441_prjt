package widgetsrc

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type SQLStore struct {
	DB *sql.DB
}

func (sqls *SQLStore) CreateDefaultWidget(dw *DefaultWidgetInfo) (*DefaultWidgetInfo, error) {
	dw.CreatedAt = time.Now()
	dw.UpdatedAt = time.Now()

	insq := "insert into widget(user_id, created_at, updated_at) values(?,?,?)"

	res, errExec := sqls.DB.Exec(insq, dw.UserID, dw.CreatedAt, dw.UpdatedAt)
	if errExec != nil {
		return nil, errExec
	}

	wid, idErr := res.LastInsertId()
	if idErr != nil {
		return nil, idErr
	}
	dw.WidgetID = wid

	insq = "insert into widget_location(widget_id, location) values(?,?)"
	for _, loc := range dw.Location {
		res, errExec = sqls.DB.Exec(insq, wid, loc.Location)
		if errExec != nil {
			return nil, errExec
		}
	}

	return dw, nil
}

func (sqls *SQLStore) GetDefaultWidget(wid int) (*DefaultWidgetInfo, error) {
	dw := &DefaultWidgetInfo{}
	insq := "select w.widget_id, wt.widget_type_name, w.user_id, w.created_at, w.updated_at from widget w join widget_type wt on w.widget_type_id = wt.widget_type_id where w.widget_id = ?"
	errQuery := sqls.DB.QueryRow(insq, wid).Scan(&dw.WidgetID, &dw.WidgetType, &dw.UserID, &dw.CreatedAt, &dw.UpdatedAt)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errQuery
	}

	locs := make([]*WidgetLocation, 0)
	insq = "select location from widget_location where widget_id = ?"

	res, errQuery := sqls.DB.Query(insq, dw.WidgetID)
	if errQuery != nil {
		return nil, errQuery
	}

	for res.Next() {
		loc := &WidgetLocation{}
		errScan := res.Scan(&loc.Location)
		if errScan != nil {
			return nil, errScan
		}
		locs = append(locs, loc)
	}
	dw.Location = locs
	return dw, nil
}

func (sqls *SQLStore) EditDefaultWidget(dw *DefaultWidgetInfo) error {
	insq := "update widget set updated_at = ? where widget_id = ?"
	_, errExec := sqls.DB.Exec(insq, dw.UpdatedAt, dw.WidgetID)
	if errExec != nil {
		return errExec
	}

	insq = "delete from widget_location where widget_id = ?"
	_, errExec = sqls.DB.Exec(insq, dw.WidgetID)
	if errExec != nil {
		return errExec
	}

	insq = "insert into widget_location(widget_id, location) values(?,?)"
	for _, loc := range dw.Location {
		_, errExec = sqls.DB.Exec(insq, dw.WidgetID, loc.Location)
		if errExec != nil {
			return errExec
		}
	}
	return nil
}

func (sqls *SQLStore) DeleteDefaultWidget(wid int) error {
	insq := "delete from widget_location where widget_id = ?"
	_, errExec := sqls.DB.Exec(insq, wid)
	if errExec != nil {
		fmt.Println("error default")
		return errExec
	}

	insq = "delete from widget_like where widget_id = ?"

	_, errExec = sqls.DB.Exec(insq, wid)
	if errExec != nil {
		fmt.Println("error default")
		return errExec
	}

	insq = "delete from widget_comment where widget_id = ?"

	_, errExec = sqls.DB.Exec(insq, wid)
	if errExec != nil {
		fmt.Println("error default")
		return errExec
	}

	insq = "delete from widget where widget_id = ?"

	_, errExec = sqls.DB.Exec(insq, wid)
	if errExec != nil {
		fmt.Println("error default")
		return errExec
	}

	return nil
}

func (sqls *SQLStore) GetTextBoxWidget(wid int) (*TextBoxWidget, error) {
	tbw := &TextBoxWidget{}
	insq := "select `text` from text_box_widget where widget_id = ?"
	errQuery := sqls.DB.QueryRow(insq, wid).Scan(&tbw.Text)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errQuery
	}
	return tbw, nil
}

func (sqls *SQLStore) CreateTextBoxWidget(tbw *TextBoxWidget) error {
	insq := "insert into text_box_widget(widget_id, `text`) values(?,?)"

	_, errExec := sqls.DB.Exec(insq, tbw.BaseInfo.WidgetID, tbw.Text)
	if errExec != nil {
		return errExec
	}

	insq = "update widget set widget_type_id = 1 where widget_id = ?"

	_, errExec = sqls.DB.Exec(insq, tbw.BaseInfo.WidgetID)
	if errExec != nil {
		return errExec
	}

	return nil
}

func (sqls *SQLStore) EditTextBoxWidget(tbw *TextBoxWidget) error {
	insq := "update text_box_widget set `text` = ? where widget_id = ?"
	_, errExec := sqls.DB.Exec(insq, tbw.Text, tbw.BaseInfo.WidgetID)
	if errExec != nil {
		return errExec
	}
	return nil
}

func (sqls *SQLStore) DeleteTextBoxWidget(wid int) error {
	insq := "delete from text_box_widget where widget_id = ?"

	_, errExec := sqls.DB.Exec(insq, wid)
	if errExec != nil {
		fmt.Println("error tb")
		return errExec
	}

	return nil
}

func (sqls *SQLStore) CreateRecentTracksWidget(rtw *RecentTracksWidget) error {
	insq := "insert into recent_tracks_widget(widget_id, num_tracks, lastfm, description) values(?,?,?,?)"

	_, errExec := sqls.DB.Exec(insq, rtw.BaseInfo.WidgetID, rtw.NumTracks, rtw.Lastfm, rtw.Description)
	if errExec != nil {
		return errExec
	}

	insq = "update widget set widget_type_id = 2 where widget_id = ?"

	_, errExec = sqls.DB.Exec(insq, rtw.BaseInfo.WidgetID)
	if errExec != nil {
		return errExec
	}

	return nil
}

func (sqls *SQLStore) GetRecentTracksWidget(wid int) (*RecentTracksWidget, error) {
	rtw := &RecentTracksWidget{}
	insq := "select num_tracks, lastfm, description from recent_tracks_widget where widget_id = ?"
	errQuery := sqls.DB.QueryRow(insq, wid).Scan(&rtw.NumTracks, &rtw.Lastfm, &rtw.Description)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errQuery
	}
	return rtw, nil
}

func (sqls *SQLStore) EditRecentTracksWidget(rtw *RecentTracksWidget) error {
	insq := "update recent_tracks_widget set num_tracks = ?, lastfm = ?, description = ? where widget_id = ?"
	_, errExec := sqls.DB.Exec(insq, rtw.NumTracks, rtw.Lastfm, rtw.Description, rtw.BaseInfo.WidgetID)
	if errExec != nil {
		return errExec
	}
	return nil
}

func (sqls *SQLStore) DeleteRecentTracksWidget(wid int) error {
	insq := "delete from recent_tracks_widget where widget_id = ?"

	_, errExec := sqls.DB.Exec(insq, wid)
	if errExec != nil {
		fmt.Println("error tb")
		return errExec
	}
	return nil
}

func (sqls *SQLStore) CreateTopMusicWidget(tmw *TopMusicWidget) error {
	insq := "insert into top_music_widget(widget_id, num_tracks, lastfm, description, type, time_period) values(?,?,?,?,?,?)"

	_, errExec := sqls.DB.Exec(insq, tmw.BaseInfo.WidgetID, tmw.NumTracks, tmw.Lastfm, tmw.Description, tmw.Type, tmw.TimePeriod)
	if errExec != nil {
		return errExec
	}

	insq = "update widget set widget_type_id = 3 where widget_id = ?"

	_, errExec = sqls.DB.Exec(insq, tmw.BaseInfo.WidgetID)
	if errExec != nil {
		return errExec
	}

	return nil
}

func (sqls *SQLStore) GetTopMusicWidget(wid int) (*TopMusicWidget, error) {
	tmw := &TopMusicWidget{}
	insq := "select num_tracks, lastfm, description, type, time_period from top_music_widget where widget_id = ?"
	errQuery := sqls.DB.QueryRow(insq, wid).Scan(&tmw.NumTracks, &tmw.Lastfm, &tmw.Description, &tmw.Type, &tmw.TimePeriod)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errQuery
	}
	return tmw, nil
}

func (sqls *SQLStore) EditTopMusicWidget(tmw *TopMusicWidget) error {
	insq := "update top_music_widget set num_tracks = ?, lastfm = ?, description = ?, type = ?, time_period = ? where widget_id = ?"
	_, errExec := sqls.DB.Exec(insq, tmw.NumTracks, tmw.Lastfm, tmw.Description, tmw.Type, tmw.TimePeriod, tmw.BaseInfo.WidgetID)
	if errExec != nil {
		return errExec
	}
	return nil
}

func (sqls *SQLStore) DeleteTopMusicWidget(wid int) error {
	insq := "delete from top_music_widget where widget_id = ?"

	_, errExec := sqls.DB.Exec(insq, wid)
	if errExec != nil {
		fmt.Println("error tb")
		return errExec
	}
	return nil
}

func (sqls *SQLStore) CreateSpotifyPlaylistWidget(spw *SpotifyPlaylistWidget) error {
	insq := "insert into spotify_playlist_widget(widget_id, num_tracks, description, spotify_uri, playlist_order) values(?,?,?,?,?)"

	_, errExec := sqls.DB.Exec(insq, spw.BaseInfo.WidgetID, spw.NumTracks, spw.Description, spw.SpotifyURI, spw.PlaylistOrder)
	if errExec != nil {
		return errExec
	}

	insq = "update widget set widget_type_id = 4 where widget_id = ?"

	_, errExec = sqls.DB.Exec(insq, spw.BaseInfo.WidgetID)
	if errExec != nil {
		return errExec
	}

	return nil
}

func (sqls *SQLStore) GetSpotifyPlaylistWidget(wid int) (*SpotifyPlaylistWidget, error) {
	spw := &SpotifyPlaylistWidget{}
	insq := "select num_tracks, description, spotify_uri, playlist_order from spotify_playlist_widget where widget_id = ?"
	errQuery := sqls.DB.QueryRow(insq, wid).Scan(&spw.NumTracks, &spw.Description, &spw.SpotifyURI, &spw.PlaylistOrder)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errQuery
	}
	return spw, nil
}

func (sqls *SQLStore) EditSpotifyPlaylistWidget(spw *SpotifyPlaylistWidget) error {
	insq := "update spotify_playlist_widget set num_tracks = ?, description = ?, spotify_uri = ?, playlist_order = ? where widget_id = ?"
	_, errExec := sqls.DB.Exec(insq, spw.NumTracks, spw.Description, spw.SpotifyURI, spw.PlaylistOrder, spw.BaseInfo.WidgetID)
	if errExec != nil {
		return errExec
	}
	return nil
}

func (sqls *SQLStore) DeleteSpotifyPlaylistWidget(wid int) error {
	insq := "delete from spotify_playlist_widget where widget_id = ?"

	_, errExec := sqls.DB.Exec(insq, wid)
	if errExec != nil {
		fmt.Println("error tb")
		return errExec
	}
	return nil
}

func (sqls *SQLStore) CreateFeaturedMusicWidget(fmw *FeaturedMusicWidget) error {
	insq := "insert into featured_music_widget(widget_id, description, type, music_name) values(?,?,?,?)"

	_, errExec := sqls.DB.Exec(insq, fmw.BaseInfo.WidgetID, fmw.Description, fmw.Type, fmw.MusicName)
	if errExec != nil {
		return errExec
	}

	insq = "update widget set widget_type_id = 5 where widget_id = ?"

	_, errExec = sqls.DB.Exec(insq, fmw.BaseInfo.WidgetID)
	if errExec != nil {
		return errExec
	}

	return nil
}

func (sqls *SQLStore) GetFeaturedMusicWidget(wid int) (*FeaturedMusicWidget, error) {
	fmw := &FeaturedMusicWidget{}
	insq := "select description, type, music_name from featured_music_widget where widget_id = ?"
	errQuery := sqls.DB.QueryRow(insq, wid).Scan(&fmw.Description, &fmw.Type, &fmw.MusicName)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errQuery
	}
	return fmw, nil
}

func (sqls *SQLStore) EditFeaturedMusicWidget(fmw *FeaturedMusicWidget) error {
	insq := "update featured_music_widget set description = ?, type = ?, music_name = ? where widget_id = ?"
	_, errExec := sqls.DB.Exec(insq, fmw.Description, fmw.Type, fmw.MusicName)
	if errExec != nil {
		return errExec
	}
	return nil
}

func (sqls *SQLStore) DeleteFeaturedMusicWidget(wid int) error {
	insq := "delete from featured_music_widget where widget_id = ?"

	_, errExec := sqls.DB.Exec(insq, wid)
	if errExec != nil {
		fmt.Println("error tb")
		return errExec
	}
	return nil
}

func (sqls *SQLStore) CreateLike(wl *WidgetLike) (*WidgetLike, error) {
	insq := "insert into widget_like(widget_id, user_id, created_at) values(?,?,?)"

	wl.CreatedAt = time.Now()
	res, errExec := sqls.DB.Exec(insq, wl.WidgetID, wl.UserID, wl.CreatedAt)
	if errExec != nil {
		return nil, errExec
	}

	wlid, idErr := res.LastInsertId()
	if idErr != nil {
		return nil, idErr
	}
	wl.WidgetLikeID = wlid

	return wl, nil
}

func (sqls *SQLStore) DeleteLike(wl *WidgetLike) error {
	insq := "delete from widget_like where widget_id = ? and user_id = ?"

	_, errExec := sqls.DB.Exec(insq, wl.WidgetID, wl.UserID)
	if errExec != nil {
		fmt.Println("error tb")
		return errExec
	}
	return nil
}

func (sqls *SQLStore) CreateComment(wc *WidgetComment) (*WidgetComment, error) {
	insq := "insert into widget_comment(widget_id, user_id, comment, created_at, updated_at) values(?,?,?,?,?)"

	wc.CreatedAt = time.Now()
	wc.UpdatedAt = time.Now()
	res, errExec := sqls.DB.Exec(insq, wc.WidgetID, wc.UserID, wc.Comment, wc.CreatedAt, wc.UpdatedAt)
	if errExec != nil {
		return nil, errExec
	}

	wcid, idErr := res.LastInsertId()
	if idErr != nil {
		return nil, idErr
	}
	wc.WidgetCommentID = wcid

	return wc, nil
}

func (sqls *SQLStore) GetComment(cid int) (*WidgetComment, error) {
	wc := &WidgetComment{}
	insq := "select * from widget_comment where wc_id = ?"

	errQuery := sqls.DB.QueryRow(insq, cid).Scan(&wc.WidgetCommentID, &wc.WidgetID, &wc.UserID, &wc.Comment, &wc.CreatedAt, &wc.UpdatedAt)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errQuery
	}
	return wc, nil
}

func (sqls *SQLStore) EditComment(wc *WidgetComment) error {
	wc.UpdatedAt = time.Now()
	insq := "update widget_comment set comment = ?, updated_at = ? where wc_id = ?"
	_, errExec := sqls.DB.Exec(insq, wc.Comment, wc.UpdatedAt, wc.WidgetID, wc.UserID)
	if errExec != nil {
		return errExec
	}
	return nil
}

func (sqls *SQLStore) DeleteComment(cid int) error {
	insq := "delete from widget_comment where wc_id = ?"

	_, errExec := sqls.DB.Exec(insq, cid)
	if errExec != nil {
		fmt.Println("error tb")
		return errExec
	}
	return nil
}

func (sqls *SQLStore) CreateCommentLike(wcl *WidgetCommentLike) (*WidgetCommentLike, error) {
	insq := "insert into comment_like(wc_id, user_id, created_at) values(?,?,?)"

	wcl.CreatedAt = time.Now()
	res, errExec := sqls.DB.Exec(insq, wcl.WidgetCommentID, wcl.UserID, wcl.CreatedAt)
	if errExec != nil {
		return nil, errExec
	}

	wclid, idErr := res.LastInsertId()
	if idErr != nil {
		return nil, idErr
	}
	wcl.WidgetCommentLikeID = wclid

	return wcl, nil
}

func (sqls *SQLStore) DeleteCommentLike(wcl *WidgetCommentLike) error {
	insq := "delete from comment_like where wc_id = ? and user_id = ?"

	_, errExec := sqls.DB.Exec(insq, wcl.WidgetCommentID, wcl.UserID)
	if errExec != nil {
		fmt.Println("error tb")
		return errExec
	}
	return nil
}

func (sqls *SQLStore) GetUser(userid int) (*AuthUser, error) {
	au := &AuthUser{}
	insq := "select id, username, photourl, description from user where id = ?"
	errQuery := sqls.DB.QueryRow(insq, userid).Scan(&au.ID, &au.UserName, &au.PhotoURL, &au.Description)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errQuery
	}
	return au, nil
}

func (sqls *SQLStore) GetFollowerCount(userid int) (int, error) {
	followers := 0
	insq := "select count(*) from follow where user_followed = ?"
	errQuery := sqls.DB.QueryRow(insq, userid).Scan(&followers)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			return 0, nil
		}
		return 0, errQuery
	}

	return followers, nil
}

func (sqls *SQLStore) GetFollowingCount(userid int) (int, error) {
	following := 0
	insq := "select count(*) from follow where user_following = ?"
	errQuery := sqls.DB.QueryRow(insq, userid).Scan(&following)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			return 0, nil
		}
		return 0, errQuery
	}
	return following, nil
}

func (sqls *SQLStore) GetUserWidgets(userid int) ([]*DefaultWidgetInfo, error) {
	dws := make([]*DefaultWidgetInfo, 0)
	insq := "select w.widget_id, wt.widget_type_name, w.user_id, w.created_at, w.updated_at from widget w join widget_type wt on w.widget_type_id = wt.widget_type_id where w.user_id = ?"
	res, errQuery := sqls.DB.Query(insq, userid)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errQuery
	}

	for res.Next() {
		dw := &DefaultWidgetInfo{}
		errScan := res.Scan(&dw.WidgetID, &dw.WidgetType, &dw.UserID, &dw.CreatedAt, &dw.UpdatedAt)
		if errScan != nil {
			return nil, errScan
		}

		locs := make([]*WidgetLocation, 0)
		insq = "select location from widget_location where widget_id = ?"

		res, errQuery := sqls.DB.Query(insq, dw.WidgetID)
		if errQuery != nil {
			return nil, errQuery
		}

		for res.Next() {
			loc := &WidgetLocation{}
			errScan := res.Scan(&loc.Location)
			if errScan != nil {
				return nil, errScan
			}
			locs = append(locs, loc)
		}
		dw.Location = locs
		dws = append(dws, dw)
	}

	return dws, nil
}

func (sqls *SQLStore) GetWidgetSocialInfo(wid int, curuserid int) (*WidgetSocial, error) {
	ws := &WidgetSocial{}

	insq := "select widget_id from widget_like where widget_id = ? and user_id = ?"
	errQuery := sqls.DB.QueryRow(insq, wid, curuserid).Scan(&ws.WidgetID)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			ws.Liked = false
		} else {
			return nil, errQuery
		}
	} else {
		ws.Liked = true
	}

	insq = "select count(*) from widget_like where widget_id = ?"
	errQuery = sqls.DB.QueryRow(insq, wid).Scan(&ws.NumLikes)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			ws.NumLikes = 0
		} else {
			return nil, errQuery
		}
	}

	ws.WidgetID = int64(wid)
	wcs := make([]*WidgetComment, 0)

	insq = "select * from widget_comment where widget_id = ?"
	res, errQuery := sqls.DB.Query(insq, wid)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			ws.Comments = wcs
			return ws, nil
		} else {
			return nil, errQuery
		}
	}

	for res.Next() {
		wc := &WidgetComment{}
		errScan := res.Scan(&wc.WidgetCommentID, &wc.WidgetID, &wc.UserID, &wc.Comment, &wc.CreatedAt, &wc.UpdatedAt)
		if errScan != nil {
			return nil, errScan
		}

		temp := 0
		insq = "select cl_id from comment_like where wc_id = ? and user_id = ?"
		errQuery := sqls.DB.QueryRow(insq, wc.WidgetCommentID, curuserid).Scan(&temp)
		if errQuery != nil {
			if errQuery == sql.ErrNoRows {
				wc.Liked = false
			} else {
				return nil, errQuery
			}
		} else {
			wc.Liked = true
		}

		insq = "select count(*) from comment_like where wc_id = ?"
		errQuery = sqls.DB.QueryRow(insq, wc.WidgetCommentID).Scan(&wc.NumLikes)
		if errQuery != nil {
			if errQuery == sql.ErrNoRows {
				ws.NumLikes = 0
			} else {
				return nil, errQuery
			}
		}
		wcs = append(wcs, wc)
	}
	ws.Comments = wcs

	return ws, nil
}

func (sqls *SQLStore) CheckIfFollowing(me int, other int) (bool, error) {
	temp := 0
	insq := "select user_following from follow where user_following = ? and user_followed = ?"
	errQuery := sqls.DB.QueryRow(insq, me, other).Scan(&temp)
	if errQuery != nil {
		if errQuery == sql.ErrNoRows {
			return false, nil
		}
		return false, errQuery
	}

	return true, nil
}
