import React, { useState } from 'react';
import { NavLink } from 'react-router-dom';
import { Button, Input, InputGroup, InputGroupAddon,  } from "reactstrap";
import '../App.css';

const Search = () => {

    const [results, setResults] = useState([]);
    const [query, setQuery] = useState("");

    function handleError(resp) {
        if (!resp.ok) {
            throw Error(resp.statusText);
        }
        return resp.json()
    }

    function handleSearch() {
        const url = new URL('https://api.rioishii.me/v1/users')
        url.searchParams.append('q', query)
        fetch(url.href, {
            method: "GET",
            headers: { 
                "Content-Type": "application/json",
                "Authorization": localStorage.getItem("authToken")
            }
        }).then(handleError)
        .then(json => {
            const users = [];
            json.forEach(user => {
                users.push({
                   userName: user.userName,
                   firstName: user.firstName,
                   lastName: user.lastName,
                   photoURL: user.photoURL
                })
            });
            setResults(users)
        }).catch(err => {
            alert(err)
        })
    }

    const handleInputChange = event => {
        setQuery(event.target.value);
    };

    let displayResults = results.map((user, index) => (
        <div id="search-item" key={index}>
            <img src={user.photoURL} alt="gravatar"/>
            <h2>{user.userName}</h2>
            <h2>{user.firstName} {user.lastName}</h2>
        </div>
    ));

    return (
        <div id="search">
            <div className="container">
                <h1>Search Users</h1>
                <NavLink to="/">
                    <Button>Back</Button>
                </NavLink>
                <InputGroup>
                    <Input onChange={handleInputChange} />
                    <InputGroupAddon addonType="append"><Button color="secondary" onClick={handleSearch}>Search</Button></InputGroupAddon>
                </InputGroup>
                {displayResults}
            </div>
        </div>
    )
}

export default Search;