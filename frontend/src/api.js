import axios from 'axios';
import {getToken} from "./utils";

const API = axios.create({ baseURL: 'http://localhost:8080' });

API.interceptors.request.use((config) => {
    const token = getToken()
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

export const login = (username, password) =>
    API.post('/login', JSON.stringify({ username, password }), {
        headers: { 'Content-Type': 'application/json' },
    });

export const register = (username, password) =>
    API.post('/register', JSON.stringify({ username, password }), {
        headers: { 'Content-Type': 'application/json' },
    });

export const uploadPDF = (formData) =>
    API.post('/submit', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
    });
