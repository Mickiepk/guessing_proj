
import React, { useState, useContext } from 'react';
import { AuthContext } from './AuthContext';
import axios from 'axios';
import './Login.css';
import {
    useNavigate
} from 'react-router-dom';
const Login = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const { setToken } = useContext(AuthContext);
    const navigate = useNavigate();

    

    const handleLogin = async (e) => {
        e.preventDefault();
        
        try {
            console.log('Sending login request to backend:', { username, password });
            const response = await axios.post('http://localhost:8080/login', {
                username: username,
                password: password
            });

            console.log('Login successful:', response.data);
            setToken(response.data.token);
            navigate('/guess');
        } catch (error) {
            console.error('Error logging in:', error.response ? error.response.data : error.message);
        }
    };

    return (
        <form onSubmit={handleLogin}>
        <div className="login-container">
            <input 
                type="text" 
                value={username} 
                onChange={(e) => setUsername(e.target.value)} 
                placeholder="Username" 
                className="login-input"
            />
            <input 
                type="password" 
                value={password} 
                onChange={(e) => setPassword(e.target.value)} 
                placeholder="Password" 
                className="login-input"
            />
            <button type='submit' className="login-button">Login</button>
        </div>
        </form>
    );
}

export default Login;
