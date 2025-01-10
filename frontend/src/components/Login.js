import React, {useEffect, useState} from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../api';
import { saveToken } from '../utils';

const Login = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [turnstileToken, setTurnstileToken] = useState('');
    const navigate = useNavigate();

    window.handleTurnstileCallback = (token) => {
        console.log('Turnstile token:', token); // 调试用
        setTurnstileToken(token);
    };

    useEffect(() => {
        if (window.turnstile) {
            const instance = window.turnstile.render('.cf-turnstile');
            return () => {
                if (instance) {
                    window.turnstile.remove(instance);
                }
            };
        }
    }, []);


    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!turnstileToken) {
            alert('Please complete the Turnstile verification.');
            return;
        }
        try {
            const { data } = await login(username, password, turnstileToken);
            saveToken(data.token);
            navigate(0);
        } catch (error) {
            if (error.response && error.response.data && error.response.data.error) {
                alert(`Login failed: ${error.response.data.error}`);
            } else {
                alert('Login failed: An unknown error occurred.');
            }
        }
    };

    return (
        <div>
            <form onSubmit={handleSubmit}>
                <h2>Sign In</h2>
                <input
                    type="text"
                    placeholder="username"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}
                    required
                />
                <input
                    type="password"
                    placeholder="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                />
                <div
                    className="cf-turnstile"
                    data-sitekey="0x4AAAAAAA4m5v6QQl0Eov1I"
                    data-callback="handleTurnstileCallback"
                ></div>
                <button type="submit">Sign In</button>
            </form>
            <div style={{ marginTop: '10px' }}>
                <p>No account yet?</p>
                <button onClick={() => navigate('/register')}>Register</button>
            </div>
        </div>
    );
};

export default Login;
