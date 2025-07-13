import { useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";
import { RefreshCw } from "lucide-react";
import ENDPOINT_URL from "../utils/config.js";
import AddFriend from "./AddFriendButton.jsx";
import NewGroup from "./NewGroupButton.jsx";

function SideBar() {
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [conversations, setConversations] = useState([]);
    const navigate = useNavigate();

    useEffect(() => {
        if(loading){
            fetch("/api/conversations", {
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
                    setConversations(data.conversations)
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
    }, [loading])

    if(loading) 
        return (
            <>
                <div className="flex items-center justify-center h-full bg-blue-200 my-1 rounded-md">
                    <div className="w-16 h-16 border-4 border-dashed rounded-full border-blue-500 animate-spin"></div>
                </div>
            </>
        );

    if(error) 
        return (
            <>
                <div className="flex items-center justify-center h-full bg-blue-200 my-1 rounded-md">
                    <div className="items-center justify-center flex flex-col w-full h-fit">
                        <button className="w-fit p-2 h-fit bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors my-3"
                            onClick={(e) => {e.preventDefault(); setLoading(true)}}>
                            <RefreshCw />
                        </button>
                    </div>
                </div>
            </>
        );

    return (
        <>
            <div className="h-full w-full overflow-y-auto bg-blue-200 my-1 rounded-md p-1">
                <div className="flex flex-row w-full h-fit justify-end p-1">
                    <NewGroup returnData={(convo) => {setConversations([...conversations, convo])}} />
                    <AddFriend returnData={(convo) => {setConversations([...conversations, convo])}} />
                </div>
                {loading ? (
                    <div className="flex items-center m-2 justify-center min-h-40">
                        <div className="w-16 h-16 border-4 border-dashed rounded-full border-blue-500 animate-spin m-1"></div>
                    </div>
                ) : (
                    conversations.length != 0 ? conversations.map((convo, index) => (
                        <button className="grid grid-cols-1 w-full p-1 h-fit bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors my-2"
                            onClick={(e) => {e.stopPropagation(); navigate(`/Messanger/${convo.Id}`)}}>
                            <div className="items-center justify-center">
                                <p className="text-1xl text-gray-100">{convo.Name}</p>
                            </div>
                        {
                            // <div className="items-center justify-center">
                            //     <p className="text-sm text-gray-200">{convo.CreateDate}</p>
                            // </div>
                        }
                        </button>
                    )) : (
                        <div className="items-center justify-center w-full">
                            <p className="text-1xl text-gray-700 text-center">You have no chats.</p>
                        </div>
                    )
                )}
            </div>
        </>
    );
    
}

export default SideBar;
