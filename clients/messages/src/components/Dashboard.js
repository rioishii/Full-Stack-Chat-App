import React, { useContext } from 'react';
import { Redirect, NavLink } from 'react-router-dom';
import { AuthContext } from "../App";
import { Button } from 'reactstrap';

const Dashboard = () => {
    const { state, dispatch } = useContext(AuthContext);

    function handleSignout() {
        const url = new URL("https://api.rioishii.me/v1/sessions/mine")
        fetch(url.href, {
            method: "DELETE",
            headers: { "Authorization":localStorage.getItem('authToken') }
        }).then(resp => {
            localStorage.setItem("authToken", "")
            if (resp.ok) {
                dispatch({
                    type: "LOGOUT",
                });
            }
        }).catch(err => {
            alert(err)
        });
    }

    function renderUserInfo() {
        if (state.isAuthenticated && !state.user != null) {
            return (
                <div>
                    <img src={state.user.photoURL} alt="gravatar"/>
                    <div>{state.user.firstname}</div>
                    <div>{state.user.lastname}</div>
                </div> 
            )
        } else {
            return <Redirect to='/signup' />
        }
    }
    
    return (
        <div id="dashboard">
            {renderUserInfo()}
            <NavLink to="/search">
                <Button color="primary">Search</Button>
            </NavLink>
            <NavLink to="/update">
                <Button color="warning">Update Profile</Button>
            </NavLink>
            <Button color="danger" onClick={() => handleSignout()}>Sign out</Button>
        </div>            
    );
};

export default Dashboard;