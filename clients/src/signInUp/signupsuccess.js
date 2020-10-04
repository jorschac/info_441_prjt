import React, { Component } from "react";
//import {NavLink} from 'react-router-dom';


export  class SignUpSuccess extends Component {

    render() {
        return (
            <div>
                <p>very good, now Sign In</p>
                <a href='/LogInUp' id='redirectButton'>: Sign In</a>            
            </div>
        );
    }
}
