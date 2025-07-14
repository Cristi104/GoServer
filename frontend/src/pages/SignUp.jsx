import { Navigate, Link, useNavigate } from "react-router-dom";
import { useState, useEffect } from 'react';
import ENDPOINT_URL from "../utils/config.js";

function SignUp() {
    const [formData, setFormData] = useState({ username: "", email: "", password: "" });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");

    const navigate = useNavigate();
    function handleSubmit(e) {
        e.preventDefault();
        setLoading(true);
        fetch(ENDPOINT_URL + "/api/auth/signup", {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-type": "application/json",
                "X-CSRF-Token": localStorage.getItem("csrfToken"),
            },
            credentials: "same-origin",
            body: JSON.stringify(formData),
        })
        .then(response => response.json())
        .then(data => {
            if(data.success){
                localStorage.setItem("userId", data.userId)
                navigate("/Messanger")
                // setError("");
                // setLoading(false);
            } else {
                setError(data.error);
                setLoading(false);
            }
        })
        .catch((error) => {
            setError(error.message);
            setLoading(false);
        })
    }

    function handleChange(e) {
        e.preventDefault();
        const { name, value } = e.target;
        setFormData((prevData) => ({
            ...prevData,
            [name]: value,
        }));
    }

    return (
        <>
            <div className="w-screen h-screen bg-blue-100 justify-center items-center flex">
                <div className="p-8 rounded-2xl shadow-lg w-full max-w-md mx-auto h-fit bg-blue-200">
                    <h2 className="text-2xl font-bold mb-6 text-center">Sign Up</h2>
                    <form onSubmit={handleSubmit} className="space-y-4">
                        <div>
                            <label className="block text-gray-700 font-medium mb-1">Username</label>
                            <input autocomplete="off"
                                type="text"
                                name="username"
                                value={formData.username}
                                onChange={handleChange}
                                className="w-full border border-blue-300 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                        <div>
                            <label className="block text-gray-700 font-medium mb-1">Email</label>
                            <input autocomplete="off"
                                type="text"
                                name="email"
                                value={formData.email}
                                onChange={handleChange}
                                className="w-full border border-blue-300 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                        <div>
                            <label className="block text-gray-700 font-medium mb-1">
                                Password
                            </label>
                            <input autocomplete="off"
                                type="password"
                                name="password"
                                value={formData.password}
                                onChange={handleChange}
                                className="w-full border border-blue-300 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                            />
                        </div>
                        <p className="text-1xl text-red-600 mb-3">{error}</p>
                        {loading ? (
                            <div className="flex items-center justify-center">
                                <div className="w-4 h-4 border-4 border-dashed rounded-full border-blue-500 animate-spin"></div>
                            </div>
                        ) : (
                            <div className="flex items-center justify-center">
                                <div className="w-4 h-4"></div>
                            </div>
                        )}
                        <button type="submit" className="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition-colors">
                            Sign Up
                        </button>
                    </form>
                    <button className="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition-colors my-3"
                        onClick={(e) => {e.preventDefault(); navigate("/SignIn")}}>
                        Allready have an account?
                    </button>
                </div>
            </div>
        </>
    );
}

export default SignUp
