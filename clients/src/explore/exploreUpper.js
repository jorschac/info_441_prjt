import React, { Component } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faHeart } from '@fortawesome/free-solid-svg-icons';
//import * as convert from 'color-convert';
//import { FilterContainer } from './exploreFilter.js';



export class UpperContainer extends Component {
    render() {
        /*let searchBoxProp = {handleAddFilter: this.props.propList.handleAddFilter, handleSearch: this.props.propList.handleSearch, 
            searchQuery: this.props.propList.searchQuery, handleLock: this.props.propList.handleLock, selectedPalette: this.props.propList.selectedPalette, 
            filterList: this.props.propList.filterList, handleError: this.props.propList.handleError}*/
        
        /*let filterContainerProp = {filterList: this.props.propList.filterList, handleLock: this.props.propList.handleLock, 
            handleRemoveFilter: this.props.propList.handleRemoveFilter, selectedPalette: this.props.propList.selectedPalette}*/
        let UsrProp = {UsrName: this.props.propList.UsrName, UsrDesc: this.props.propList.UsrDesc, followingNum: this.props.propList.followingNum, followersNum: this.props.propList.followersNum}

        let UserImgProp = {UserImg: this.props.propList.UserImg}
            return (
            <section id="uppercontainer">
                <UserPicture propList={UserImgProp} />
                <UserInformation propList={UsrProp} />
            </section>
        );
    }
}

class UserInformation extends Component {
    render() {
        return (
            <section id="UserInfo">
                <div id="NameAndDesc">
                    <p id="UserName">{this.props.propList.UsrName}</p>
                    <p id="Describtion">{this.props.propList.UsrDesc}</p>
                </div>
                <div id="FollowStatus">
                    <div id="following">
                        <FontAwesomeIcon icon={faHeart} className='far fa-heart' aria-hidden="true" />
                        <p>{'following ' + this.props.propList.followingNum}</p>
                    </div>
                    <div id="followers">
                        <FontAwesomeIcon icon={faHeart} className='far fas-heart' aria-hidden="true" />
                        <p>{'followers ' + this.props.propList.followersNum}</p>
                    </div>
                </div>
            </section>    
        );
    }
}

class UserPicture extends Component {
    render() {
        return (
            <img src={this.props.propList.UserImg} alt="avatar" id="avatar"/>    
        );
    }
}

