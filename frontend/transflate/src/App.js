import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Login from './components/Login';
import Register from './components/Register';
import Translate from './components/Translate';
import { isAuthenticated } from './utils';
import './App.css'
const App = () => {
  return (
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/register" element={<Register />} />
          <Route
              path="/"
              element={isAuthenticated() ? <Translate /> : <Navigate to="/login" />}
          />
        </Routes>
      </Router>
  );
};

export default App;
