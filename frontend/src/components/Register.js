import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { register } from '../api';

const Register = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await register(username, password);
            alert('Register successfully!');
            navigate('/login');
        } catch (error) {
            if (error.response && error.response.data && error.response.data.error) {
                alert(`Register failed: ${error.response.data.error}`);
            } else {
                alert('Register failed: An unknown error occurred.');
            }
        }
    };
    const handleLoginRedirect = () => {
        navigate('/login');
    };
    return (<div>
    <form onSubmit={handleSubmit}>
        <h2>Register</h2>
        <input type="text" placeholder="username" value={username} onChange={(e) => setUsername(e.target.value)}
               required/>
        <input type="password" placeholder="password" value={password} onChange={(e) => setPassword(e.target.value)}
               required/>
        <button type="submit">Register</button>
    </form>
    <div style={{marginTop: '10px'}}>
        <p>Already registered?</p>
        <button onClick={handleLoginRedirect}>Sign in</button>
    </div>
        </div>
);
};

export default Register;
