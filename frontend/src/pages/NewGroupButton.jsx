import { useState, useEffect } from "react";
import { SquarePlus, UserPlus, Search, SquareX } from "lucide-react";
import ENDPOINT_URL from "../utils/config.js";

function NewGroup({ returnData }) {
    const [isOpen, setIsOpen] = useState(false);
    const [formData, setFormData] = useState({ name: "" });
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [friends, setFriends] = useState([]);
    const [addedFriends, setAddedFriends] = useState([])

    useEffect(() => {
        setLoading(true);
        if(isOpen){
            fetch("/api/profiles/friends", {
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
                    setFriends(data.friends)
                    setLoading(false)
                } else {
                    setError(data.error)
                    setLoading(false)
                }
            })
            .catch((error) => {
                setError(error.message)
                setLoading(false)
            })
        }
    }, [isOpen])

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
        fetch(ENDPOINT_URL + "/api/conversations", {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-type": "application/json",
                "X-CSRF-Token": localStorage.getItem("csrfToken"),
            },
            credentials: "same-origin",
            body: JSON.stringify({Action: "createGroup", Name: formData.name, Users: addedFriends}),
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

    function handleChange(e) {
        e.preventDefault();
        const { name, value } = e.target;
        setFormData((prevData) => ({
            ...prevData,
            [name]: value,
        }));
    }

    function addFriend(id) {
        setAddedFriends([...addedFriends, id])
    }

    function removeFriend(id) {
        setAddedFriends(addedFriends.filter((item) => item != id))
    }

    return (
        <div>
            <button className="w-fit p-2 h-fit bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors mx-1"
                onClick={(e) => setIsOpen(true)}
                title="Create group">
                <SquarePlus />
            </button>

            {isOpen && (
                <div className="fixed inset-0 z-50 bg-black/30 bg-opacity-80 flex items-center justify-center"
                    onClick={(e) => setIsOpen(false)}>
                    <div className="bg-blue-100 rounded-2xl p-8 w-1/3 text-center h-fit"
                        onClick={(e) => e.stopPropagation()}>
                        <form onSubmit={handleSubmit} className="w-full h-fit">
                            <div className="flex flex-row">
                                <input autocomplete="off"
                                    placeholder="Group name"
                                    type="text"
                                    name="name"
                                    value={formData.name}
                                    onChange={handleChange}
                                    className="w-full placeholder-gray-600 border border-blue-300 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                />
                                <button type="submit" className="w-fit p-2 h-fit bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors mx-1"
                                    onClick={(e) => setIsOpen(true)}>
                                    <SquarePlus />
                                </button>
                            </div>
                            <div className="flex w-full my-1">
                                {friends.length == 0 ? "" : (friends.map((friend, index) => (
                                    addedFriends.includes(friend.id) ? (
                                        <div className="w-fit h-fit bg-blue-200 rounded-lg transition-colors my-1 flex flex-row pl-1 m-1">
                                            <p className="text-xs p-1">{friend.nickname}</p>
                                            <button className="w-fit h-fit p-1 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors ml-auto"
                                                onClick={(e) => removeFriend(friend.id)}>
                                                <SquareX />
                                            </button>
                                        </div>
                                    ) : ""
                                )))}
                            </div>
                        </form>
                        {loading ? (
                            <div className="flex items-center m-2 justify-center min-h-40">
                                <div className="w-4 h-4 border-4 border-dashed rounded-full border-blue-500 animate-spin m-1"></div>
                            </div>
                        ) : (
                            <div className="flex flex-col items-center my-2 justify-center min-h-40 overflow-y-auto max-h-120">
                                {friends.length == 0 ? (
                                    <div className="w-full h-4 items-center justify-center">
                                        <p className="text-1xl text-gray-700 text-center">{error}</p>
                                    </div>
                                ) : (friends.map((friend, index) => (
                                    <div className="w-full h-fit bg-blue-200 rounded-lg transition-colors my-1 flex flex-row pl-1">
                                        <p className="text-md py-2 px-1 p-1">{friend.nickname}</p>
                                        <p className="text-xs py-2 px-1 text-gray-600">@{friend.username}</p>
                                        {!addedFriends.includes(friend.id) ? (
                                            <button className="w-fit p-2 h-fit bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors ml-auto"
                                                onClick={(e) => addFriend(friend.id)}>
                                                <SquarePlus />
                                            </button>
                                        ) : ""}
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

export default NewGroup;
