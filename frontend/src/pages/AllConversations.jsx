import { Navigate } from "react-router-dom";
import { useState, useEffect } from 'react';
import { RefreshCw } from "lucide-react";

function AllConversations() {
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const [conversations, setConversations] = useState("");

    useEffect(() => {
        if(loading){
            fetch("/api/conversations", {
                method: "GET"
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
            <div className="w-full overflow-y-auto">
                {conversations != null ? conversations.map((convo, index) => (
                    <div className="grid grid-cols-1 w-max p-1 h-fit bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors my-2">
                        <div className="items-center justify-center">
                            <p className="text-1xl text-gray-100">{convo.Name}</p>
                        </div>
                        <div className="items-center justify-center">
                            <p className="text-sm text-gray-200">{convo.CreateDate}</p>
                        </div>
                    </div>
                )) : ""}
            </div>
        </>
    );
    
}

export default AllConversations
