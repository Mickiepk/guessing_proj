import React, { useState, useContext } from 'react';
import { AuthContext } from './AuthContext';
import axios from 'axios';

const Guess = () => {
    const [guess, setGuess] = useState('');
    const [result, setResult] = useState('');
    const { token } = useContext(AuthContext);

    const handleGuess = async () => {
        try {
            console.log('Sending guess to backend:', guess);
            const response = await axios.post('http://localhost:8080/guess', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                },
                body: JSON.stringify({ guess })
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const data = await response.json();
            console.log('Received response from backend:', data);
            setResult(data.message);
        } catch (error) {
            console.error('Failed to fetch:', error);
            setResult('Failed to fetch data from backend');
        }
    };

    return (
        <div>
            <input type="text" value={guess} onChange={(e) => setGuess(e.target.value)} placeholder="Your guess" />
            <button onClick={handleGuess}>Guess</button>
            <p>{result}</p>
        </div>
    );
};

export default Guess;
