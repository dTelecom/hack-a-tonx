import React, {StrictMode} from 'react';
import App from './App';
import {createRoot} from 'react-dom/client';
import {BrowserRouter} from 'react-router-dom';
import {appStore} from './stores/appStore';

// Initializing contract/user
async function initContract() {
  return {};
}

const rootElement = document.getElementById('root');
const root = createRoot(rootElement);

window.initPromise = initContract().then(
  ({currentUser}) => {
    appStore.setCurrentUser(currentUser);
    void appStore.loadStats();

    root.render(
      <StrictMode>
        <BrowserRouter>
          <App
            currentUser={currentUser}
          />
        </BrowserRouter>
      </StrictMode>
    );
  }
);
