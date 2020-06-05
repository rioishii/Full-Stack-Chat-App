import React, { useState, useContext } from 'react';
import { Button, Form, FormGroup, Label, Input } from "reactstrap";
import { NavLink } from 'react-router-dom';
import { AuthContext } from "../App";

const Update = (props) => {
    const { dispatch } = useContext(AuthContext);
    const initialState = {
        firstname: "",
        lastname: "",
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
            throw Error(resp.statusText);
        }
        return resp.json()
    }

    function handleSubmit(event) {
        event.preventDefault()
        setData({
            ...data,
        });
        const url = new URL("https://api.rioishii.me/v1/users/me")
        fetch(url.href, {
            method: "PATCH",
            headers: { 
                "Content-Type": "application/json",
                "Authorization": localStorage.getItem("authToken")
            },
            body: JSON.stringify({
                FirstName: data.firstname,
                LastName: data.lastname,
            })
        }).then(handleError)
        .then(json => {
            dispatch({
                type: "UPDATE",
                payload: json
            });
            props.history.push('/')
        }).catch(err => {
            alert(err)
        });
    }

    return (
        <div className="container">
            <Form onSubmit={handleSubmit}>
                <FormGroup>
                    <Label for="firstname">First Name</Label>
                    <Input name="firstname" id="firstname" onChange={handleInputChange}/>
                </FormGroup>
                <FormGroup>
                    <Label for="lastname">Last Name</Label>
                    <Input name="lastname" id="lastname" onChange={handleInputChange}/>
                </FormGroup>
                <Button color="primary" type="submit">Update</Button>
                <NavLink to="/">
                    <Button>Cancel</Button>
                </NavLink>
            </Form>
        </div>
    )
}

export default Update;