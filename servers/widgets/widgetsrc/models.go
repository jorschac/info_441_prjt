package widgetsrc

import "time"

type WidgetContext struct {
	WStore WidgetStore
}

type AuthUser struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	UserName    string `json:"userName"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhotoURL    string `json:"photoUrl"`
	Description string `json:"description"`
}

type Artist struct {
	Playcount string `json:"playcount"`
	Name      string `json:"name"`
	Url       string `json:"url"`
}

type ArtistChart struct {
	Artists []*Artist `json:"artist"`
}

type Response struct {
	Weeklyartistchart *ArtistChart `json:"weeklyartistchart"`
}

type WidgetLocation struct {
	Location int64 `json:"location"`
}

type DefaultWidgetInfo struct {
	WidgetID   int64             `json:"widgetId"`
	WidgetType string            `json:"widgetType"`
	UserID     int64             `json:"userId"`
	CreatedAt  time.Time         `json:"createdAt"`
	UpdatedAt  time.Time         `json:"updatedAt"`
	Location   []*WidgetLocation `json:"locations"`
}

type TextBoxWidget struct {
	BaseInfo *DefaultWidgetInfo `json:"baseInfo"`
	Text     string             `json:"text"`
}

type RecentTracksWidget struct {
	BaseInfo    *DefaultWidgetInfo `json:"baseInfo"`
	NumTracks   int64              `json:"numTracks"`
	Lastfm      string             `json:"lastfm"`
	Description string             `json:"description"`
}

type TopMusicWidget struct {
	BaseInfo    *DefaultWidgetInfo `json:"baseInfo"`
	NumTracks   int64              `json:"numTracks"`
	Lastfm      string             `json:"lastfm"`
	Description string             `json:"description"`
	Type        string             `json:"type"`
	TimePeriod  int64              `json:"timePeriod"`
}

type SpotifyPlaylistWidget struct {
	BaseInfo      *DefaultWidgetInfo `json:"baseInfo"`
	NumTracks     int64              `json:"numTracks"`
	Description   string             `json:"description"`
	SpotifyURI    string             `json:"spotifyUri"`
	PlaylistOrder bool               `json:"playlistOrder"`
}

type FeaturedMusicWidget struct {
	BaseInfo    *DefaultWidgetInfo `json:"baseInfo"`
	Description string             `json:"description"`
	Type        string             `json:"type"`
	MusicName   string             `json:"musicName"`
}

type WidgetLike struct {
	WidgetLikeID int64     `json:"likeId"`
	WidgetID     int64     `json:"widgetId"`
	UserID       int64     `json:"userId"`
	CreatedAt    time.Time `json:"createdAt"`
}

type WidgetComment struct {
	WidgetCommentID int64     `json:"commentId"`
	WidgetID        int64     `json:"widgetId"`
	UserID          int64     `json:"userId"`
	Comment         string    `json:"comment"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	Liked           bool      `json:"liked"`
	NumLikes        int64     `json:"numLikes"`
}

type WidgetCommentLike struct {
	WidgetCommentLikeID int64     `json:"commentLikeId"`
	WidgetCommentID     int64     `json:"commentId"`
	UserID              int64     `json:"userId"`
	CreatedAt           time.Time `json:"createdAt"`
}

type WidgetSocial struct {
	WidgetID int64            `json:"widgetId"`
	Liked    bool             `json:"like"`
	NumLikes int              `json:"numLikes"`
	Comments []*WidgetComment `json:"comments"`
}

type UserPage struct {
	User             *AuthUser                `json:"user"`
	IsMe             bool                     `json:"isMe"`
	IsFollowing      bool                     `json:"isFollowing"`
	NumFollowers     int                      `json:"numFollowers"`
	NumFollowing     int                      `json:"numFollowing"`
	TBW              []*TextBoxWidget         `json:"textBoxWidgets"`
	RTW              []*RecentTracksWidget    `json:"recentTracksWidgets"`
	TMW              []*TopMusicWidget        `json:"topMusicWidgets"`
	SPW              []*SpotifyPlaylistWidget `json:"spotifyPlaylistWidgets"`
	FMW              []*FeaturedMusicWidget   `json:"featuredMusicWidget"`
	WidgetSocialInfo []*WidgetSocial          `json:"widgetSocialInfo"`
}
