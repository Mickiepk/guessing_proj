import React, { useState, useContext, useEffect } from 'react';
import { AuthContext } from './AuthContext';
import axios from 'axios';
import './Guess.css';
import { useNavigate } from 'react-router-dom';
const Guess = () => {
    const [guess, setGuess] = useState('');
    const [result, setResult] = useState('');
    const { token } = useContext(AuthContext);
    const navigator = useNavigate();
    useEffect(() => {
        if (!token) {

            navigator('/login');


        }
    }, [token])



    const onGuestChangeHandler = (e) => {
        e.preventDefault();
        setGuess(prevGuess => e.target.value);
    }

    const handleGuess = async () => {
        try {
            console.log('Sending guess to backend:', guess);
            const response = await axios.post('http://localhost:8080/guess',
                { guess },
                {
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${token}`
                    }
                }
            );

            console.log('Received response from backend:', response.data);
            setResult(prevResult => response.data.message);
        } catch (error) {
            console.error('Failed to fetch:', error.response ? error.response.data : error.message);
            setResult('Failed to fetch data from backend');
        }
    };

    return (
        <div>
            <input type="text" value={guess} onChange={onGuestChangeHandler} placeholder="Your guess" />
            <button onClick={handleGuess}>Guess</button>
            <p>{result}</p>
        </div>
    );
};

export default Guess;
