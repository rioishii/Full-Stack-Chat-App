import React, { useState } from 'react';
import { Button, Form, FormGroup, Label, Input } from "reactstrap";
import { Redirect, NavLink } from 'react-router-dom';
import { AuthContext } from "../App";
import '../App.css';


const Signup = () => {
    const { state, dispatch } = React.useContext(AuthContext);
    const initialState = {
        email: "",
        password: "",
        passconf: "",
        username: "",
        firstname: "",
        lastname: ""
    };
    const [data, setData] = useState(initialState)
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
        event.preventDefault();
        setData({
            ...data,
        });
        const url = new URL("https://api.rioishii.me/v1/users")
        fetch(url.href, {
            method: "POST",
            headers: { 
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                Email: data.email,
                Password: data.password,
                PasswordConf: data.passconf,
                UserName: data.username,
                FirstName: data.firstname,
                LastName: data.lastname
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
        <div className="signup">
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
                <FormGroup>
                    <Label for="passconf">Confirm Password</Label>
                    <Input type="password" name="passconf" id="passconf" onChange={handleInputChange} />
                </FormGroup>
                <FormGroup>
                    <Label for="username">Username</Label>
                    <Input name="username" id="username" onChange={handleInputChange} />
                </FormGroup>
                <FormGroup>
                    <Label for="firstname">First Name</Label>
                    <Input name="firstname" id="firstname" onChange={handleInputChange}/>
                </FormGroup>
                <FormGroup>
                    <Label for="lastname">Last Name</Label>
                    <Input name="lastname" id="lastname" onChange={handleInputChange} />
                </FormGroup>
                <Button color="danger" type="submit">Register</Button>
                <NavLink to="/login">
                    <Button>Log in</Button>
                </NavLink>
            </Form>
        </div>
        
    )
}

export default Signup;