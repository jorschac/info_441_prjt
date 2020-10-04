package widgetsrc

type WidgetStore interface {
	CreateDefaultWidget(dw *DefaultWidgetInfo) (*DefaultWidgetInfo, error)
	GetDefaultWidget(wid int) (*DefaultWidgetInfo, error)
	EditDefaultWidget(dw *DefaultWidgetInfo) error
	DeleteDefaultWidget(wid int) error

	CreateTextBoxWidget(tbw *TextBoxWidget) error
	GetTextBoxWidget(tbid int) (*TextBoxWidget, error)
	EditTextBoxWidget(tbw *TextBoxWidget) error
	DeleteTextBoxWidget(tbid int) error

	CreateRecentTracksWidget(rtw *RecentTracksWidget) error
	GetRecentTracksWidget(wid int) (*RecentTracksWidget, error)
	EditRecentTracksWidget(rtw *RecentTracksWidget) error
	DeleteRecentTracksWidget(wid int) error

	CreateTopMusicWidget(tmw *TopMusicWidget) error
	GetTopMusicWidget(wid int) (*TopMusicWidget, error)
	EditTopMusicWidget(tmw *TopMusicWidget) error
	DeleteTopMusicWidget(wid int) error

	CreateSpotifyPlaylistWidget(spw *SpotifyPlaylistWidget) error
	GetSpotifyPlaylistWidget(wid int) (*SpotifyPlaylistWidget, error)
	EditSpotifyPlaylistWidget(spw *SpotifyPlaylistWidget) error
	DeleteSpotifyPlaylistWidget(wid int) error

	CreateFeaturedMusicWidget(fmw *FeaturedMusicWidget) error
	GetFeaturedMusicWidget(wid int) (*FeaturedMusicWidget, error)
	EditFeaturedMusicWidget(fmw *FeaturedMusicWidget) error
	DeleteFeaturedMusicWidget(wid int) error

	CreateLike(wl *WidgetLike) (*WidgetLike, error)
	DeleteLike(wl *WidgetLike) error

	CreateComment(wc *WidgetComment) (*WidgetComment, error)
	GetComment(cid int) (*WidgetComment, error)
	EditComment(wc *WidgetComment) error
	DeleteComment(cid int) error

	CreateCommentLike(wcl *WidgetCommentLike) (*WidgetCommentLike, error)
	DeleteCommentLike(wcl *WidgetCommentLike) error

	GetUser(userid int) (*AuthUser, error)
	GetFollowerCount(userid int) (int, error)
	GetFollowingCount(userid int) (int, error)

	GetUserWidgets(userid int) ([]*DefaultWidgetInfo, error)

	GetWidgetSocialInfo(wid int, curuserid int) (*WidgetSocial, error)

	CheckIfFollowing(me int, other int) (bool, error)
}
