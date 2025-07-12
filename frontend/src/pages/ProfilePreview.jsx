import { Outlet, Navigate, Link, useNavigate } from "react-router-dom";
import { useState, useEffect } from 'react';
import ENDPOINT_URL from "../utils/config.js";

function ProfilePreview() {
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [profile, setProfile] = useState("");

    useEffect(() => {
        fetch(ENDPOINT_URL + `/api/profiles/${localStorage.getItem("userId")}`, {
            method: "GET",
            headers: {
                Accept: "application/json",
                "Content-type": "application/json",
                "X-CSRF-Token": localStorage.getItem("csrfToken"),
            },
            credentials: "same-origin",
        })
        .then((response) => response.json())
        .then((data) => {
            if(data.success){
                setProfile(data.profile);
                setLoading(false);
            } else {
                setError(data.error);
                setLoading(false);
            }
        })
        .catch((error) => {
            setError(error.message);
            setLoading(false);
        })
    }, [])

    const navigate = useNavigate()

    if(loading) 
        return (
            <>
                <div className="w-full h-24 flex items-center justify-center bg-blue-200 my-1 rounded-md">
                     <div className="w-24 h-24 border-4 border-dashed rounded-full border-blue-500 animate-spin"></div>
                </div>
            </>
        );

    if(error) 
        return (
            <>
                <div className="w-full h-24 flex items-center justify-center bg-blue-200 my-1 rounded-md">
                    <p className="text-1xl text-red-600">Failed to load profile data</p>
                </div>
            </>
        );

    return (
        <>
            <div className="w-full h-24 flex items-center justify-center bg-blue-200 my-1 rounded-md">
                <button className="w-6/7 h-4/5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors my-3"
                    onClick={(e) => {e.preventDefault(); navigate("/Messanger/Profile")}}>
                    <p className="text-xl text-white p-1">{profile.nickname}</p>
                    <p className="text-xs text-gray-300 p-1">{profile.username}</p>
                </button>
            </div>
        </>
    );
}

export default ProfilePreview
