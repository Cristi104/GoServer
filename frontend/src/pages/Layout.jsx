import { Outlet, Navigate } from "react-router-dom";
import { useState, useEffect } from 'react';
import ProfilePreview from './ProfilePreview.jsx';
import SideBar from "./SideBar.jsx"

function Layout() {
    return (
        <>
            <div className="flex flex-row bg-blue-100">
                <div className="flex flex-col w-48 h-screen">
                    <ProfilePreview />
                    <SideBar />
                </div>
                <div className="flex flex-col w-full h-screen mx-1 py-1">
                    <Outlet />
                </div>
            </div>
        </>
    );
}

export default Layout
