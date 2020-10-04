
import React, { Component } from 'react';
import './index.css';
import { NavBar, MobileNav } from './common/navigate.js';
import { Explore } from './explore/explore.js';
import { Route, Switch, Link, Redirect } from 'react-router-dom';
import { Signin } from './signInUp/signin.js'
import { SignUp } from './signInUp/signup.js'
import {SignUpSuccess} from './signInUp/signupsuccess.js'
import {SignInSuccess} from './signInUp/signinsuccess.js'

//import { FrontPage } from './home/front.js';
// main component
export class App extends Component {
    constructor(props) {
        super(props);
        this.state = {mobileMenuOn: false, selected: false, 
        selectedPalette: ['#ffffff', '#818181', '#ff6f61', '#836e58', '#232326'], currentTheme: ['#ffffff', '#818181', '#ff6f61', '#836e58', '#232326']
        , UserImg: "http://localhost:3000/img/userimg.png", currUsrFollowingNum: 777, currUsrFollowersNum: 666, currUsrName: "SampleName", currUsrDesc: "this is a sample description", 
        WidgetCoverPic: ["http://localhost:3000/img/samplecover1.jpg", "http://localhost:3000/img/samplecover2.jpeg", "http://localhost:3000/img/samplecover3.jpg"], WigetTitle: ["SampleTitle1", "SampleTitle2", "SampleTitle3"], WigetDesc: ["this is a smapleDesc1", "this is a smapleDesc2", "this is a smapleDesc3"], usrID: 0, usrAuth: ""}; //UserImg will use picture sent from go later
        //widget Cover Title Desc needs to be accessed from external database
    }

    // selecte a palette
    handleSelectPalette = (palette) => {
        if (!this.state.selected) {
            this.setState({selected: true});
        }
        this.setState({ selectedPalette: palette });
    }

    // apply selected theme when apply tab is clicked
    handleApplyClick = () => {
        if (this.state.selected) {
            this.setState({ currentTheme: this.state.selectedPalette, selected: false });
        }
    }

    // functionality for mobile menu (hamburger menu)
    handleMobileMenu = () => {
        let status = !this.state.mobileMenuOn;
        this.setState({ mobileMenuOn: status });
    }
    /*
    handleCreateChange = (color, index) => {
        let newSelectedPalette = this.state.selectedPalette;
        newSelectedPalette[index] = color;
        this.setState({selectedPalette: newSelectedPalette, selected: true});
    }*/

    usrIdCallback = (datafromchild) => {
        //debugger;
        this.setState({usrID: datafromchild});
        console.log("id successfully recieved from sign in")
        console.log(this.state.usrID)
    }

    usrAuthCallback = (authfromchild) => {
        this.setState({usrAuth: authfromchild});
        console.log("auth successfully recieved from sign in")
        console.log(this.state.usrAuth)
    }
    

    render() {
        let style = { '--lightShade': this.state.currentTheme[0], '--lightAccent': this.state.currentTheme[1], 
        '--mainColor': this.state.currentTheme[2], '--darkAccent': this.state.currentTheme[3], 
        '--darkShade': this.state.currentTheme[4]};
    
        let mobileNavProp = {handleApply: this.handleApplyClick, mobileMenuOn: this.state.mobileMenuOn, handleMobileMenu: this.handleMobileMenu}
        var i
        var widgetInfo = []
        var len = this.state.WidgetCoverPic.length
        for (i = 0; i < len; i ++) {
            console.log(this.state.WidgetCoverPic[i])
            widgetInfo.push({WidgetCoverPic: this.state.WidgetCoverPic[i], WigetTitle: this.state.WigetTitle[i], WigetDesc: this.state.WigetDesc[i]})
            //widgetInfo += {WidgetCoverPic: this.state.WidgetCoverPic[i], WigetTitle: this.state.WigetTitle[i], WigetDesc: this.state.WigetDesc[i]}
        }


        //let widgetInfo ={WidgetCoverPic: this.state.WidgetCoverPic, WigetTitle: this.state.WigetTitle, WigetDesc: this.state.WigetDesc}

        let exploreProp = {handleSelectPalette: this.handleSelectPalette, handleApplyClick: this.handleApplyClick, selectedPalette: this.state.selectedPalette,
            currentTheme: this.state.currentTheme, selected: this.state.selected, UserImg: this.state.UserImg,
            currUsrFollowingNum: this.state.currUsrFollowingNum, currUsrFollowersNum: this.state.currUsrFollowersNum, 
            currUsrName: this.state.currUsrName, currUsrDesc: this.state.currUsrDesc, WidgetInfo: widgetInfo, usrID: this.state.usrID, usrAuth: this.state.usrAuth};

        let signUpProp = {callback: this.usrIdCallback}
        let signInProp = {callback: this.usrIdCallback, AuthCallBack: this.usrAuthCallback}
        let usrIDProp = {usrID: this.state.usrID, usrAuth: this.state.usrAuth}
        return (

            <div className='appContainer' style={style}>
                <header>
                    <Link exact to='/'><h1>acryline</h1></Link>
                </header>
                <NavBar handleApply={this.handleApplyClick} handleMobileMenu={this.handleMobileMenu}/>
                <MobileNav propList={mobileNavProp} />
                <Switch>
                    <Route exact path='/' render={() => <Explore propList={exploreProp}/>}/>
                    <Route path='/LogInUp' render={() => <Signin propList={signInProp}/>}/>
                    <Route path='/SignUp' render={() => <SignUp propList={signUpProp}/>}/>
                    <Route path='/signUpSuccess' render={() => <SignUpSuccess/>}/>
                    <Route path='/signInSucess' render={() => <SignInSuccess propList={usrIDProp}/>}/>
                    <Route path='/Home' render={() => <Explore propList={exploreProp}/>}/>
                    <Redirect to='/'></Redirect>
                </Switch>
                <Footer />
            </div>
        );
    }
}

export default App;

class Footer extends Component {
    render() {
        return (
            <footer>
                <p>Powered by <a href='https://casesandberg.github.io/react-color/'>React Color</a> | <a href='https://github.com/Qix-/color-convert#readme'>Color-Convert</a></p>
                <p>© 2019 Gunhyung Cho  |  Jiuzhou Zhao</p>
                <address>Contact: <a href='mailto:ghcho@uw.edu'>ghcho@uw.edu</a> |  <a href='mailto:jz73@uw.edu'>jz73@uw.edu</a></address>
            </footer>
        );
    }
}