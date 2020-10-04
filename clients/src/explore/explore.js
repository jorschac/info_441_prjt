import React, { Component } from 'react';
import { UpperContainer } from './exploreUpper.js';
import {CardContainer } from './explorePalettes.js';
import { Spinner } from '../common/spinner.js';
import * as convert from 'color-convert'; // for converting color values
//import firebase from "firebase/app";
import 'firebase/database';
import './explore.css';
//import { strict } from 'assert'
export class Explore extends Component {
    constructor(props) {
        super(props);
        this.state = {palettes: [], filteredPalettes: [], nFiltered: 0, error: '', 
            dataLoaded: false, filterList: [], lockStatus: [false, false, false, false, false], searchQuery: ''
            ,UsrName: '', UsrDesc: '', currUsrFollowingNum: 0, currUsrFollowersNum: 0, Usrpic: "", 
            WidgetCoverPic: ["http://localhost:3000/img/samplecover1.jpg", "http://localhost:3000/img/samplecover2.jpeg", "http://localhost:3000/img/samplecover3.jpg"], FeaturedWigets: [],  
            RecentTracksWigets: [], SpotifyWigets: []}
        this.handleUsrInfo = this.handleUsrInfo.bind(this)
        this.handleWidgetInfo = this.handleWidgetInfo.bind(this)
    }

    handleWidgetInfo = () => {
        console.log("place holder")
    }

    handleUsrInfo = () => {
        
        console.log(this.props.propList.usrID)
        //debugger;
        let url = 'https://api.cahillaw.me/v1/profile/'+this.props.propList.usrID
        var xhttp = new XMLHttpRequest();
        xhttp.open('GET', url, true);
        xhttp.setRequestHeader('Content-Type', 'Authorization');
        console.log(this.props.propList.usrAuth)
        xhttp.setRequestHeader('Authorization', this.props.propList.usrAuth);
        xhttp.onload = () => {
            //debugger;
           console.log("user ptofile request sent sucessfully")
           console.log(xhttp.status)
           if(xhttp.status !== 200) { //needs to be changed later
               alert(xhttp.response);
           } else {
            let responseJson = JSON.parse(xhttp.response)
            console.log(responseJson)
            if (responseJson.user.id !== null) {
            //when a user struct got back, reset State
            this.setState({UsrName: responseJson.user.userName, UsrDesc: responseJson.user.description
                            , currUsrFollowersNum: responseJson.numFollowers, currUsrFollowingNum: responseJson.numFollowing, Usrpic: responseJson.user.photoUrl})
            var i;
            var len = responseJson.featuredMusicWidget.length
            var featuredWidgetInfo = []
            for (i = 0; i < len; i ++) {
                console.log(responseJson.featuredMusicWidget[i].description)
                featuredWidgetInfo.push({WidgetCoverPic: this.state.WidgetCoverPic[0], WigetTitle: responseJson.featuredMusicWidget[i].musicName, WigetDesc: responseJson.featuredMusicWidget[i].description})
            }
            this.setState({FeaturedWigets: featuredWidgetInfo})
            var o;
            var recentLen = responseJson.recentTracksWidgets.length
            var recentTracksWidgetInfo = []
            for (o = 0; o < recentLen; o ++) {
                console.log(responseJson.recentTracksWidgets[i].description)
                recentTracksWidgetInfo.push({baseInfo: responseJson.recentTracksWidgets[i].baseInfo, WidgetCoverPic: this.state.WidgetCoverPic[0], numTracks: responseJson.recentTracksWidgets[i].numTracks, WigetDesc: responseJson.recentTracksWidgets[i].description})
            }
            this.setState({RecentTracksWigets: recentTracksWidgetInfo})
            var p;
            var spotifyLen = responseJson.spotifyPlaylistWidgets.length
            var spotifyWidgetInfo = []
            for (p = 0; p < spotifyLen; p ++) {
                console.log(responseJson.spotifyPlaylistWidgets[i].description)
                spotifyWidgetInfo.push({baseInfo: responseJson.spotifyPlaylistWidgets[i].baseInfo, WidgetCoverPic: this.state.WidgetCoverPic[0], numTracks: responseJson.spotifyPlaylistWidgets[i].numTracks, WigetDesc: responseJson.spotifyPlaylistWidgets[i].description})
            }
            this.setState({SpotifyWigets: spotifyWidgetInfo})
            } else {
                alert('specific user id recieved is null');
            }
        }
        }
        xhttp.onprogress = () => {
            debugger;
            console.log('loading', xhttp.status)
        }
        xhttp.send();
        console.log("specific user profile req sent, waiting for results...")
    }

    componentDidMount() {
        this.handleUsrInfo()
    }

    // adds new filter
    handleAddFilter = (filter) => {
        if (!this.state.filterList.includes(filter)) {

            let filters = this.state.filterList;
            filters.push(filter);
            
            let list = this.state.filteredPalettes.filter((palette) => {
                
                return (palette.username === filter || convert.hex.keyword(palette.light_shade) === filter || 
                convert.hex.keyword(palette.light_accent) === filter || convert.hex.keyword(palette.main) === filter ||
                convert.hex.keyword(palette.dark_accent) === filter || convert.hex.keyword(palette.dark_shade) === filter);
            });
            this.setState({ filteredPalettes: list , filterList: filters, nFiltered: list.length });
        }
    }

    // remove an existing filter
    handleRemoveFilter = (filter) => {
        let list = this.state.filterList.filter((data) => {
            return data !== filter;
        })
        this.setState({ filterList: list }, () => {
            if (this.state.filterList.length === 0) {
                this.setState({ filteredPalettes: this.state.palettes, nFiltered: this.state.palettes.length });
            } else {
                let list = this.state.palettes;
                let filterList = this.state.filterList;
                list = list.filter((data) => {
                    return (filterList.includes(data.username) || filterList.includes(convert.hex.keyword(data.light_shade)) || 
                        filterList.includes(convert.hex.keyword(data.light_accent)) || filterList.includes(convert.hex.keyword(data.main)) ||
                        filterList.includes(convert.hex.keyword(data.dark_accent)) || filterList.includes(convert.hex.keyword(data.dark_shade)));
                });
                
                this.setState({ filteredPalettes: list , nFiltered: list.length});
            }
        });
    }

    // updates the color lock buttons
    handleUpdateLock = (filter, lockId) => {
        let currLockStatus = this.state.lockStatus;
        
        let selectedColorNames = this.props.propList.selectedPalette.map(x => convert.hex.keyword(x));
        let lockColor = selectedColorNames[lockId];
        for (let i = 0; i < 5; i++) {
            if (lockColor === selectedColorNames[i]) {
                currLockStatus[i] = !currLockStatus[i];
            }
        }
        if (currLockStatus[lockId]) {
            this.handleAddFilter(filter);
        } else {
            this.handleRemoveFilter(filter);
        }
        
        this.setState({ lockStatus: currLockStatus });
    }

    // resets the color lock buttons
    handleResetLock = () => {
        this.setState({ lockStatus: [false, false, false, false, false] });
    }
    
    // tracks the input in search box
    handleUpdateQuery = (input) => {
        let cleanedInput = input.toLowerCase().replace(/\s+/g, '');
        this.setState({ searchQuery: cleanedInput });
    }

    // shows the error message
    handleError = (msg) => {
        this.setState({ error: msg });
        setTimeout(() => {
            this.setState({ error: '' });
        }, 3000);
    }

    render() {

        let upperContainerProp = {filterList: this.state.filterList, handleAddFilter: this.handleAddFilter, handleSearch: this.handleUpdateQuery,
            searchQuery: this.state.searchQuery, handleLock: this.handleUpdateLock, handleRemoveFilter: this.handleRemoveFilter,
            selectedPalette: this.props.propList.selectedPalette, handleError: this.handleError, UserImg: this.state.Usrpic
            , followingNum: this.state.currUsrFollowingNum, followersNum: this.state.currUsrFollowersNum,
            UsrName: this.state.UsrName, UsrDesc: this.state.UsrDesc};
        
        let cardContainerProp = {filteredData: this.state.featuredWidgetInfo.slice(0, 16), handleClick: this.props.propList.handleSelectPalette, 
            handleResetLock: this.handleResetLock};
/*
        var i
        var widgetInfo = []
        var len = this.state.WigetTitle.length
        for (i = 0; i < len; i ++) {
            console.log(this.state.WidgetCoverPic[i])
            widgetInfo.push({WidgetCoverPic: this.state.WidgetCoverPic[0], WigetTitle: this.state.WigetTitle[i], WigetDesc: this.state.WigetDesc[i]})
            //widgetInfo += {WidgetCoverPic: this.state.WidgetCoverPic[i], WigetTitle: this.state.WigetTitle[i], WigetDesc: this.state.WigetDesc[i]}
        }*/

        return (
            <main>
                <UpperContainer propList={upperContainerProp} />
                {!this.state.dataLoaded && <Spinner />}
                <CardContainer propList={cardContainerProp} />
            </main>
        );
    }
}