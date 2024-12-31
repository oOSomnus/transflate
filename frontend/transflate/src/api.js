import axios from 'axios';

const API = axios.create({ baseURL: 'http://localhost:8080' });

API.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export const login = (username, password) => API.post('/login', { username, password });
export const register = (username, password) => API.post('/register', { username, password });
export const uploadPDF = (formData) => API.post('/upload', formData);
