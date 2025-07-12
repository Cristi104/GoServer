import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { useEffect } from 'react';
import ENDPOINT_URL from "./utils/config.js"
import './index.css';
import Layout from './pages/Layout.jsx';
import Conversation from './pages/Conversation.jsx';
import Profile from './pages/Profile.jsx';
import PublicHomePage from './pages/PublicHomePage.jsx';
import SignIn from './pages/SignIn.jsx';
import SignUp from './pages/SignUp.jsx';
import Dashboard from './pages/Dashboard.jsx';
import getCSRFToken from './utils/csrf.js'


function App() {
    useEffect(() => {
        fetch(ENDPOINT_URL + "/api/csrf", {
            method: "GET",
            credentials: "include",
        })
        .then(response => {
            localStorage.setItem("csrfToken", response.headers.get("X-CSRF-Token"))
        })
    }, [])

    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<PublicHomePage />} />
                <Route path="/SignIn" element={<SignIn />} />
                <Route path="/SignUp" element={<SignUp />} />
                <Route path="/Messanger" element={<Layout />}>
                    <Route index element={<Dashboard />} />
                    <Route path=":id" element={<Conversation />} />
                    <Route path="Profile" element={<Profile />} />
                </Route>
                <Route path="*" element={<Navigate to="/" />} />
            </Routes>
        </BrowserRouter>
    );
}

export default App
