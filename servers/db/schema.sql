create table if not exists user (
    id int not null auto_increment primary key,
    email varchar(320) not null,
    passhash binary(64) not null,
    username varchar(255) not null,
    first_name varchar(128) not null,
    last_name varchar(128) not null,
    photourl varchar(256),
    description varchar(256),
    UNIQUE KEY unique_email(email),
    UNIQUE KEY unique_username(username)
);

create table if not exists logins (
    loginid int not null auto_increment primary key,
    id int not null,
    login_date date not null,
    ip_address varchar(128) not null,
    foreign key (id) REFERENCES user(id)
);

create table if not exists follow (
    follow_id int not null auto_increment primary key,
    user_following int not null,
    user_followed int not null,
    date_followed DATETIME not null,
    foreign key (user_following) REFERENCES user(id),
    foreign key (user_following) REFERENCES user(id)
);

create table if not exists widget_type (
    widget_type_id int not null auto_increment primary key,
    widget_type_name varchar(125) not null
);

create table if not exists widget (
    widget_id int not null auto_increment primary key,
    widget_type_id int,
    user_id int not null,
    created_at DATETIME not null,
    updated_at DATETIME not null,
    foreign key (widget_type_id) REFERENCES widget_type(widget_type_id),
    foreign key (user_id) REFERENCES user(id)
);

create table if not exists widget_location (
    wl_id int not null auto_increment primary key,
    widget_id int not null,
    location int not null,
    foreign key (widget_id) REFERENCES widget(widget_id)
);

create table if not exists widget_comment (
    wc_id int not null auto_increment primary key,
    widget_id int not null,
    user_id int not null,
    comment varchar(128) not null,
    created_at DATETIME not null,
    updated_at DATETIME not null,
    foreign key (widget_id) REFERENCES widget(widget_id),
    foreign key (user_id) REFERENCES user(id)
);

create table if not exists widget_like (
    wl_id int not null auto_increment primary key,
    widget_id int not null,
    user_id int not null,
    created_at DATETIME not null,
    foreign key (widget_id) REFERENCES widget(widget_id),
    foreign key (user_id) REFERENCES user(id)
);

create table if not exists comment_like (
    cl_id int not null auto_increment primary key,
    wc_id int not null,
    user_id int not null,
    created_at DATETIME not null,
    foreign key (wc_id) REFERENCES widget_comment(wc_id),
    foreign key (user_id) REFERENCES user(id)
);

create table if not exists text_box_widget (
    widget_id int not null primary key REFERENCES widget(widget_id),
    `text` varchar(256) not null
);

create table if not exists recent_tracks_widget (
    widget_id int not null primary key REFERENCES widget(widget_id),
    num_tracks int not null,
    lastfm varchar(125) not null,
    description varchar(256)
);

create table if not exists top_music_widget (
    widget_id int not null primary key REFERENCES widget(widget_id),
    num_tracks int not null,
    lastfm varchar(125) not null,
    description varchar(256),
    type varchar(125) not null,
    time_period int not null
);

create table if not exists spotify_playlist_widget (
    widget_id int not null primary key REFERENCES widget(widget_id),
    num_tracks int not null,
    description varchar(256),
    spotify_uri varchar(256) not null,
    playlist_order boolean not null
);

create table if not exists featured_music_widget (
    widget_id int not null primary key REFERENCES widget(widget_id),
    description varchar(256),
    type varchar(125) not null,
    music_name varchar(256) not null
);

insert into user(email, passhash, username, first_name, last_name, photourl, description)
values('cahillawx@gmail.com', 2321321, 'cahillaw', 'Andy','Cahill','url','I like music lol');

insert into widget_type(widget_type_name)
values('Text Box'),('Recent Tracks'),('Top Music'),('Spotify Playlist'),('Featured Music')