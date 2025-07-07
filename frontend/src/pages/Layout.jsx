import getCookie from "../utils/cookies.js";
import { Outlet, Navigate } from "react-router-dom";
import { useState, useEffect } from 'react';
import AllConversations from './AllConversations.jsx';
import ProfilePreview from './ProfilePreview.jsx';

function Layout() {
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [profile, setProfile] = useState("");

    console.log(getCookie("auth"))
    // if (getCookie("auth") == "") {
    //     return(<Navigate to="/" />);
    // }   

    return (
        <>
            <div className="grid grid-cols-2 bg-blue-100 mx-1">
                <div className="flex flex-col w-48 h-screen">
                    <ProfilePreview />
                    <AllConversations />
                </div>
                <div>
                    <Outlet />
                </div>
            </div>
        </>
    );
}

export default Layout
