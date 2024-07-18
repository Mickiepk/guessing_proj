import React, { useState, useContext } from 'react';
import { AuthContext } from './AuthContext';
import axios from 'axios';

const Login = () => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const { setToken } = useContext(AuthContext);

    const handleLogin = async () => {
        try {
            console.log('Sending login request to backend:', { username, password });
            const response = await axios.post('http://localhost:8080/login', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ username, password })
            });
    
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
    
            const data = await response.json();
            console.log('Received response from backend:', data);
            if (data.token) {
                setToken(data.token);
                window.location.href = '/guess';
            } else {
                console.error('Login failed: ', data);
            }
        } catch (error) {
            console.error('Failed to fetch: ', error);
        }
    };

    return (
        <div>
            <input type="text" value={username} onChange={(e) => setUsername(e.target.value)} placeholder="Username" />
            <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Password" />
            <button onClick={handleLogin}>Login</button>
            
        </div>
    );
}

export default Login;
