import { Navigate, Link, useNavigate } from "react-router-dom";
import { useState, useEffect } from 'react';

function PublicHomePage() {

    const navigate = useNavigate();

    return (
        <>
            <div className="h-screen flex items-center justify-center bg-blue-100">
                <div className="w-1/2 mx-auto bg-blue-200 p-2 rounded-lg">
                    <h1 className="text-center text-4xl pb-3 font-bold">GoServer</h1>
                    <h2 className="text-xl pb-2">Welcome!</h2>
                    <p className = "p-1">GoServer is a simple instant messaging app similar to WhatsApp</p>
                    <p className = "p-1">This app has not been created for actual usage, I only created it to improve my own dev skills and capabilities</p>
                    <p className = "p-1">The application's stack consists of a Postgresql database, the backend written in Golang, and the frontend is in React.js and Tailwindcss</p>
                    <p className = "p-1">To continue exploring the app create an account</p>
                    <div className="grid grid-cols-2 gap-4 mt-3">
                        <button className="bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition-colors"
                            onClick={(e) => {e.preventDefault(); navigate("/SignUp")}}>
                            SignUp
                        </button>
                        <button className="bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition-colors"
                            onClick={(e) => {e.preventDefault(); navigate("/SignIn")}}>
                            SignIn
                        </button>
                    </div>
                </div>
            </div>
        </>
    );
}

export default PublicHomePage
