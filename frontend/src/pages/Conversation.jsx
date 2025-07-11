import { Navigate, Link, useParams } from "react-router-dom";
import { useState, useEffect, useRef } from 'react';
import { Send } from "lucide-react";
import ENDPOINT_URL from "../utils/config.js";

function Conversation() {
    const [formData, setFormData] = useState({ body: "" });
    const [messages, setMessages] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const messageBox = useRef(null);
    let { id } = useParams();


    useEffect(() => {
        if (messageBox.current) {
            messageBox.current.scrollTop = messageBox.current.scrollHeight;
        }
    }, [messages])
    

    useEffect(() => {
        const messageSource = new EventSource(ENDPOINT_URL + `/api/conversations/${id}/messages/listener`)
        messageSource.onmessage = (e) => {
            let data = JSON.parse(e.data)
            if(data.success){
                setMessages((m) => [...m, ...data.messages])
                setError("")
            } else {
                setError(data.error)
            }
            setLoading(false)
        }

        return () => {
            messageSource.close();
        }
    }, [])
    
    function handleSubmit(e) {
        e.preventDefault();
        setLoading(true);
        fetch(ENDPOINT_URL + `/api/conversations/${id}/messages`, {
            method: "POST",
            headers: {
                Accept: "application/json",
                "Content-type": "application/json",
            },
            body: JSON.stringify({Action: "createMessage", Body: formData.body, ConversationId: id}),
            credentials: "same-origin",
        })
        .then(response => response.json())
        .then(data => {
            if(data.success){
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

    return (
        <>
            <div className="flex flex-col w-full h-full bg-blue-200 rounded-md m-1 p-1">
                {loading ? (
                    <div className="flex items-center m-2 justify-center min-h-40">
                        <div className="w-16 h-16 border-4 border-dashed rounded-full border-blue-500 animate-spin m-1"></div>
                    </div>
                ) : (
                    <>
                        <div ref={messageBox} className="w-full h-full overflow-y-auto">
                            {messages.length != 0 ? messages.map((message, index) => (
                                <button className="grid grid-cols-1 w-full p-1 h-fit bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors my-2"
                                    onClick={(e) => {e.stopPropagation(); navigate(`/Messanger/${convo.Id}`)}}>
                                    <div className="items-center justify-center">
                                        <p className="text-1xl text-gray-100">{message.SenderId}</p>
                                    </div>
                                    <div className="items-center justify-center">
                                        <p className="text-sm text-gray-200">{message.Body}</p>
                                    </div>
                                </button>
                            )) : (
                                <div className="items-center justify-center w-full">
                                    <p className="text-1xl text-gray-700 text-center">You have no messages.</p>
                                </div>
                            )}
                        </div>
                        <form onSubmit={handleSubmit} className="w-full h-fit">
                            <div className="flex flex-row">
                                <input
                                    placeholder="Group name"
                                    type="text"
                                    name="body"
                                    value={formData.body}
                                    onChange={handleChange}
                                    className="w-full placeholder-gray-600 border border-blue-300 rounded-lg px-4 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
                                />
                                <button type="submit" className="w-fit p-2 h-fit bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors mx-1">
                                    <Send />
                                </button>
                            </div>
                        </form>
                    </>
                )}
            </div>
        </>
    );
}

export default Conversation
