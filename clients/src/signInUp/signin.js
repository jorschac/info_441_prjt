import React, { Component } from "react";
import './css/signin.css'
import {NavLink} from 'react-router-dom';
import { Redirect } from 'react-router-dom'



export  class Signin extends Component {
    constructor(props) {
        super(props);
        console.log(this.props)
        this.state = {redirect: false /*callback: this.props.propList.callback*/, usrID: 0}
        this.renderRedirect = this.renderRedirect.bind(this)
        this.handleSend = this.handleSend.bind(this)
        this.callback = this.callback.bind(this)
        this.usrAuthCallback = this.usrAuthCallback.bind(this)
    }

    usrAuthCallback = (auth) => {
        this.props.propList.AuthCallBack(auth)
    }

    callback = (id) => {
        this.props.propList.callback(id)
    }

    renderRedirect = () => {
        if (this.state.redirect) {
            let usrIDProp = {usrID: this.state.usrID}
          return <Redirect to='/signInSucess' propList={usrIDProp}/>
        }
      }

    handleSend = (e) => {
        e.preventDefault()
        //check if all inputs filled
        console.log(document.querySelectorAll("form .form-control"));
        let filled = 0;
       document.querySelectorAll("form .form-control").forEach(element => {
           console.log(element.value)
           if (element.value.length !== 0) {
            console.log('element.value')
               filled += 1;
           }
       });

       let email = document.getElementsByName("email")[0].value;
       let passwrd = document.getElementsByName("passwrd")[0].value;
       let signIngUsr = {"Email" : email,   
       "Password"  : passwrd
    }
    //send the request
    console.log('fuck me')
    if(filled === 2) {
     var xhttp = new XMLHttpRequest();
     xhttp.open('POST', 'https://api.cahillaw.me/v1/sessions', true);
     xhttp.setRequestHeader('Content-Type', 'application/json');
     //return ID to App.js when sent sucessfully.
     //console.log(this.state.callback)
     //let fuckitplzwork = this.state.callback
     xhttp.onload = () => {

        console.log("sign in request sent sucessfully")
         let responseJson = JSON.parse(xhttp.response)
         console.log('retrieved ID: ', responseJson.id)
         let authToken = xhttp.getResponseHeader("Authorization")
         console.log(authToken)
         //fuckitplzwork(responseJson.id)
         if (responseJson.id !== null) {
         this.callback(responseJson.id)
         this.usrAuthCallback(authToken)
         //let c = () => {this.propList.callback}
         this.setState({redirect: true, usrID: responseJson.id})
         } else {
             alert('sign in failed, check your information again');
         }
     }
     xhttp.onprogress = () => {
         debugger;
         console.log('loading', xhttp.status)
     }
     xhttp.send(JSON.stringify(signIngUsr));
     console.log("sign in req sent, waiting for results...")
                      } else {
                          alert('all form must be filled')
                      }
    }

    render() {
        return (
            <div className="auth-wrapper">
            <div className="auth-inner">
            <form>
                <h3>Sign In</h3>
                {this.renderRedirect()}
                <div className="form-group">
                    <label>Email address</label>
                    <input type="email" name="email" className="form-control" placeholder="Enter email" />
                </div>

                <div className="form-group">
                    <label>Password</label>
                    <input type="password" name="passwrd" className="form-control" placeholder="Enter password" />
                </div>

                <div className="form-group">
                    <div className="custom-control custom-checkbox">
                        <input type="checkbox" className="custom-control-input" id="customCheck1" />
                        <label className="custom-control-label" htmlFor="customCheck1">Remember me</label>
                    </div>
                </div>

                <button type="submit" className="btn btn-primary btn-block" onClick={this.handleSend}>Submit</button>
                <p className="forgot-password text-right">
                    Forgot <a href="https://cat-bounce.com/">password?</a>
                </p>
                <div>
                <NavLink exact to='/SignUp' activeClassName='active' className='signUp'>sign Up</NavLink>
                </div>
            </form>
            </div>
            </div>
        );
    }
}