import React, { Component } from "react";
import { Redirect } from 'react-router-dom'
import './css/signin.css'



export  class SignUp extends Component {
    //this function handles api request send part, it will read all values from form,
    //create NewUser json object, send a POST request to 'v1/users' api
    //when it loads the response User Json object, the ID will be returned 
    //back to App.js for further use
    constructor(props) {
        super(props);
        console.log(this.props)
        this.state = {redirect: false /*callback: this.props.propList.callback*/}
        this.renderRedirect = this.renderRedirect.bind(this)
        this.handleSend = this.handleSend.bind(this)
    }
    renderRedirect = () => {
        if (this.state.redirect) {
          return <Redirect to='/signUpSuccess' />
        }
      }

    handleSend = (e) => {
        e.preventDefault();
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
       let fname = document.getElementsByName("fname")[0].value;
       let lname = document.getElementsByName("lname")[0].value;
       let email = document.getElementsByName("email")[0].value;
       let passwrd = document.getElementsByName("passwrd")[0].value;
       let usrname = fname;
       let desc = "wonderful me";
       let newUsr = {"Email" : email,   
       "Password"  : passwrd,   
       "PasswordConf": passwrd,
       "UserName" : usrname,
       "FirstName" : fname,
       "LastName" : lname, 
       "Description" : desc};
    //send the request
    console.log('fuck me')
    if(filled === 4) {
     var xhttp = new XMLHttpRequest();
     xhttp.open('POST', 'https://api.cahillaw.me/v1/users', true);
     xhttp.setRequestHeader('Content-Type', 'application/json');
     //return ID to App.js when sent sucessfully.
     //console.log(this.state.callback)
     //let fuckitplzwork = this.state.callback
     xhttp.onload = () => {

        console.log("request sent sucessfully")
         let responseJson = JSON.parse(xhttp.response)
         console.log('retrieved ID: ', responseJson.id)
         //fuckitplzwork(responseJson.id)
         this.setState({redirect: true})
     }
     xhttp.onprogress = function () {
         debugger;
         console.log('loading', xhttp.status)
     }
     xhttp.send(JSON.stringify(newUsr));
     console.log("req sent, waiting for results...")
                      } else {
                          alert('all form must be filled')
                      }
    }

  

    render() {
        return (
            <div className="auth-wrapper">
            <div className="auth-inner">
            <form name="frm">
                <h3>Sign Up</h3>
                {this.renderRedirect()}
                <div className="form-group">
                    <label>First name</label>
                    <input type="text" name="fname" className="form-control" placeholder="First name" required/>
                </div>

                <div className="form-group">
                    <label>Last name</label>
                    <input type="text" name="lname" className="form-control" placeholder="Last name" required/>
                </div>

                <div className="form-group">
                    <label>Email address</label>
                    <input type="email" name="email" className="form-control" placeholder="Enter email" required/>
                </div>

                <div className="form-group">
                    <label>Password</label>
                    <input type="password" name="passwrd" className="form-control" placeholder="Enter password" required/>
                </div>

                <button type="submit" className="btn btn-primary btn-block" onClick={this.handleSend}>Sign Up</button>
                <p className="forgot-password text-right">
                    Already registered <a href="/LogInUp">sign in?</a>
                </p>
            </form>
            </div>
            </div>
        );
    }
}