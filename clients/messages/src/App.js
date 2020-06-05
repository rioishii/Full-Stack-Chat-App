import React, { useReducer } from 'react';
import './App.css';
import { BrowserRouter, Switch, Route } from 'react-router-dom';
import Signup from './components/Signup';
import Login from './components/Login'
import Dashboard from './components/Dashboard';
import Update from './components/Update';
import Search from './components/Search';

export const AuthContext = React.createContext();

const initialState = {
  isAuthenticated: false,
  user: null
};

const reducer = (state, action) => {
  switch (action.type) {
    case "LOGIN":
      localStorage.setItem("user", JSON.stringify(action.payload.user));
      return {
        ...state,
        isAuthenticated: true,

        user: {
          username: action.payload.userName,
          firstname: action.payload.firstName,
          lastname: action.payload.lastName,
          photoURL: action.payload.photoURL
        }
      };
    case "LOGOUT":
      return {
        ...state,
        isAuthenticated: false,
        user: null
      };
    case "UPDATE":
      localStorage.setItem("user", JSON.stringify(action.payload.user));
      return {
        ...state,
        isAuthenticated: true,

        user: {
          username: action.payload.userName,
          firstname: action.payload.firstName,
          lastname: action.payload.lastName,
          photoURL: action.payload.photoURL
        }
      };
    default:
      return state;
  }
};

function App() {
  const [state, dispatch] = useReducer(reducer, initialState);

  return (
    <AuthContext.Provider
      value={{
        state,
        dispatch
      }}
    >
      <BrowserRouter>
        <div className="App">
          <Switch>
            <Route exact path='/' component={Dashboard} />
            <Route path='/login' component={Login} />
            <Route path='/signup' component={Signup} />
            <Route path='/update' component={Update} />
            <Route path='/search' component={Search} />
          </Switch>
        </div>
      </BrowserRouter>
    </AuthContext.Provider>
  );
}

export default App;
