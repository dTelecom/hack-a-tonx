import 'regenerator-runtime/runtime'
import React from 'react'
import {Route, Routes} from 'react-router-dom'
import Home from './pages/Home/Home'
import Call from './pages/Call/Call'
import {JoinModeSelect} from "./pages/JoinModeSelect/JoinModeSelect"
import Join from "./pages/JoinParticipant/JoinParticipant"
import {JoinViewer} from "./pages/JoinViewer/JoinViewer"
import './App.scss'

const App = () => {
  return (
    <Routes>
      <Route
        index
        element={<Home/>}
      />
      <Route
        path={'/call'}
        element={<Call/>}
      />
      <Route
        path={'/join/:sid'}
        element={<JoinModeSelect/>}
      />
      <Route
        path={'/join/viewer/:sid'}
        element={<JoinViewer/>}
      />
      <Route
        path={'/join/participant/:sid'}
        element={<Join/>}
      />
      <Route
        path={'/call/:sid'}
        element={<Call/>}
      />
    </Routes>
  )
}

export default App
