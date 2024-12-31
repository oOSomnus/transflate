import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../api';
import { saveToken } from '../utils';

const Login = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const { data } = await login(username, password);
            saveToken(data.token);
            navigate('/');
        } catch (error) {
            alert('登录失败');
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <h2>登录</h2>
            <input type="text" placeholder="用户名" value={username} onChange={(e) => setUsername(e.target.value)} required />
            <input type="password" placeholder="密码" value={password} onChange={(e) => setPassword(e.target.value)} required />
            <button type="submit">登录</button>
        </form>
    );
};

export default Login;
