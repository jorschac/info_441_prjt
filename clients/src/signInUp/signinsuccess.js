import React, { Component } from "react";
import {NavLink} from 'react-router-dom';
import './css/signin.css'


export  class SignInSuccess extends Component {
 
    
    render() {
        return (
            <div>
                <p style={{color: "black"}}>very good, click to transfer to your home page</p>
                <NavLink exact to='/Home' activeClassName='active' className='redirectButton'><p style={{color: "blue"}}>transfer</p></NavLink>
            </div>
        );
    }
}