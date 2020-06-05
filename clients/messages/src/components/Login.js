import React, { useState } from 'react';
import { Button, Form, FormGroup, Label, Input } from "reactstrap";
import { AuthContext } from "../App";
import { Redirect, NavLink } from 'react-router-dom';
import '../App.css';


const Login = () => {
    const { state, dispatch } = React.useContext(AuthContext);
    const initialState = {
        email: "",
        password: "",
    };
    const [data, setData] = useState(initialState);

    const handleInputChange = event => {
        setData({
          ...data,
          [event.target.name]: event.target.value
        });
    };

    function handleError(resp) {
        if (!resp.ok) {
            localStorage.setItem("authToken", "")
            throw Error(resp.statusText);
        }
        localStorage.setItem("authToken", resp.headers.get("Authorization"));
        return resp.json()
    }

    function handleSubmit(event) {
        event.preventDefault()
        setData({
            ...data,
        });
        const url = new URL("https://api.rioishii.me/v1/sessions")
        fetch(url.href, {
            method: "POST",
            headers: { 
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                Email: data.email,
                Password: data.password,
            })
        }).then(handleError)
        .then(json => {
            dispatch({
                type: "LOGIN",
                payload: json
            });
        }).catch(err => {
            alert(err)
        });
    }

    function checkAuth() {
        if (state.isAuthenticated) {
            return <Redirect to='/' />
        }
    }

    return (
        <div className="login">
            {checkAuth()}
            <Form onSubmit={handleSubmit}>
                <FormGroup>
                    <Label for="email">Email</Label>
                    <Input type="email" name="email" id="email" onChange={handleInputChange}/>
                </FormGroup>
                <FormGroup>
                    <Label for="password">Password</Label>
                    <Input type="password" name="password" id="password" onChange={handleInputChange}/>
                </FormGroup>
                <Button color="primary" type="submit">Sign In</Button>
                <NavLink to='/signup'>
                    <Button>Cancel</Button>
                </NavLink>                
            </Form>
        </div>
    )
}

export default Login;