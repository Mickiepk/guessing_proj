import React, { useState, useContext } from 'react';
import { AuthContext } from './AuthContext';

const Guess = () => {
    const [guess, setGuess] = useState('');
    const [result, setResult] = useState('');
    const { token } = useContext(AuthContext);

    const handleGuess = async () => {
        const response = await fetch('/guess', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ guess })
        });
        const data = await response.json();
        setResult(data.message);
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