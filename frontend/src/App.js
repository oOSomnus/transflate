import React from 'react';
import {BrowserRouter as Router, Navigate, Route, Routes} from 'react-router-dom';
import Login from './components/Login';
import Register from './components/Register';
import Translate from './components/Translate';
import {isAuthenticated} from './utils';
import './App.css'
import Header from "./components/Header";
import TaskPage from "./components/TaskPage";

const App = () => {
  return (
      <Router>
          <Header />
        <Routes>
            <Route
                path="/login"
                element={isAuthenticated() ? <Navigate to="/" /> : <Login />}
            />
          <Route path="/register" element={<Register />} />
          <Route
              path="/"
              element={isAuthenticated() ? <Translate /> : <Navigate to="/login" />}
          />
            <Route path={"/tasks"} element={isAuthenticated() ? <TaskPage/> : <Navigate to="/login"/>}/>
        </Routes>
      </Router>
  );
};

export default App;
