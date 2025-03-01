import React from 'react';
import { Theme } from '@radix-ui/themes';
import '@radix-ui/themes/styles.css';
import { Yukon } from './components/Yukon/Yukon';

function App() {
  return (
    <Theme>
      <Yukon initialPath="/" />
    </Theme>
  );
}

export default App; 