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
            alert('注册成功，请登录');
            navigate('/login');
        } catch (error) {
            alert('注册失败');
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <h2>注册</h2>
            <input type="text" placeholder="用户名" value={username} onChange={(e) => setUsername(e.target.value)} required />
            <input type="password" placeholder="密码" value={password} onChange={(e) => setPassword(e.target.value)} required />
            <button type="submit">注册</button>
        </form>
    );
};

export default Register;
