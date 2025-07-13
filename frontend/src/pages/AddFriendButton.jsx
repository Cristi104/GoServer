import { useState, useEffect } from 'react';
import { SquarePlus, UserPlus, Search, SquareX } from "lucide-react";
import ENDPOINT_URL from "../utils/config.js";

function AddFriend({ returnData }) {
    const [isOpen, setIsOpen] = useState(false);
    const [formData, setFormData] = useState({ username: "" });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState("");
    const [profiles, setProfiles] = useState([]);

    useEffect(() => {
        const handleKeyDown = (event) => {
            if (event.key === 'Escape') {
                setIsOpen(false);
            }
        };

        if (isOpen) {
            document.addEventListener('keydown', handleKeyDown);
        }

        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, [isOpen]);

    function handleSubmit(e) {
        e.preventDefault();
        setLoading(true);
        fetch(ENDPOINT_URL + `/api/profiles?username=${formData.username}`, {
            method: "GET",
            headers: {
                Accept: "application/json",
                "Content-type": "application/json",
                "X-CSRF-Token": localStorage.getItem("csrfToken"),
            },
            credentials: "same-origin",
        })
        .then(response => response.json())
        .then(data => {
            if(data.success){
                setProfiles(data.profiles)
                setError("");
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
    }

    function handleChange(e) {
        e.preventDefault();
        const { name, value } = e.target;
        setFormData((prevData) => ({
            ...prevData,
            [name]: value,
        }));
    }

    function addFriend(id) {
        setLoading(true)
        fetch(ENDPOINT_URL + "/api/profiles/friends", {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-type": "application/json",
                "X-CSRF-Token": localStorage.getItem("csrfToken"),
            },
            credentials: "same-origin",
            body: JSON.stringify({Action: "addFriend", Id: id}),
        })
        .then(response => response.json())
        .then(data => {
            if(data.success){
                returnData(data.conversation);
                setError("");
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
    }

    return (
        <div>
            <button className="w-fit p-2 h-fit bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors mx-1"
                onClick={(e) => setIsOpen(true)}
                title="New direct message">
                <UserPlus />
            </button>

            {isOpen && (
                <div className="fixed inset-0 z-50 bg-black/30 bg-opacity-80 flex items-center justify-center"
                    onClick={(e) => setIsOpen(false)}>
                    <div className="bg-blue-100 rounded-2xl p-8 w-1/3 text-center h-fit"
                        onClick={(e) => e.stopPropagation()}>
                        <form onSubmit={handleSubmit} className="w-full h-fit">
                            <div className="flex flex-row">
                                <input autocomplete="off"
                                    placeholder="username"
                                    type="text"
                                    name="username"
                                    value={formData.username}
                                    onChange={handleChange}
                                    className="w-full placeholder-gray-600 border border-blue-300 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                />
                                <button type="submit" className="w-fit p-2 h-fit bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors mx-1"
                                    onClick={(e) => setIsOpen(true)}>
                                    <Search />
                                </button>
                            </div>
                        </form>
                        {loading ? (
                            <div className="flex items-center m-2 justify-center min-h-40">
                                <div className="w-4 h-4 border-4 border-dashed rounded-full border-blue-500 animate-spin m-1"></div>
                            </div>
                        ) : (
                            <div className="flex flex-col items-center m-2 justify-center overflow-y-auto max-h-120 min-h-40">
                                {profiles.length == 0 ? (
                                    <div className="w-full h-4 items-center justify-center">
                                        <p className="text-1xl text-gray-700 text-center">{error}</p>
                                    </div>
                                ) : (profiles.map((profile, index) => (
                                    <div className="w-full h-fit bg-blue-200 rounded-lg transition-colors my-1 flex flex-row pl-1">
                                        <p className="text-md py-2 px-1 p-1">{profile.nickname}</p>
                                        <p className="text-xs py-2 px-1 text-gray-600">@{profile.username}</p>
                                        <button className="w-fit p-2 h-fit bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors ml-auto"
                                            onClick={(e) => addFriend(profile.id)}>
                                            <SquarePlus />
                                        </button>
                                    </div>
                                )))}
                            </div>
                        )}
                        <button
                            onClick={() => setIsOpen(false)}
                            className="px-2 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded">
                            <SquareX />
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
}

export default AddFriend;
