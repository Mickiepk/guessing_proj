// App.js
import React from 'react';
import { AuthProvider } from './AuthContext';
import Login from './Login';
import Guess from './Guess';

const App = () => {
  return (
    <AuthProvider>
      <div className="App">
        <Login />
        <Guess />
      </div>
    </AuthProvider>
  );
}

export default App;
