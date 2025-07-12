import { Outlet, Navigate } from "react-router-dom";
import { useState, useEffect } from 'react';
import AllConversations from './AllConversations.jsx';
import ProfilePreview from './ProfilePreview.jsx';

function Layout() {
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [profile, setProfile] = useState("");

    return (
        <>
            <div className="flex flex-row bg-blue-100">
                <div className="flex flex-col w-48 h-screen">
                    <ProfilePreview />
                    <AllConversations />
                </div>
                <div className="flex flex-col w-full h-screen mx-1 py-1">
                    <Outlet />
                </div>
            </div>
        </>
    );
}

export default Layout
