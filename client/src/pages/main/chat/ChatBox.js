import React, { useState, useEffect } from "react";
import axios from "axios";
import { getToken } from "../../../general/token";
import Message from "../message/Message";

function ChatBox(props) {
  const [messages, setMessages] = useState([]);
  const [inputValue, setInputValue] = useState("");

  const fetchMessages = () => {
    axios
      .get("http://localhost:8080/chat?token=" + getToken())
      .then((res) => {
        setMessages(res.data);
      })
      .catch((err) => {
        alert(err);
      });
      const chatContainer = document.getElementById("chatbox-messages");
      chatContainer.scrollTop = chatContainer.scrollHeight;
  };

  const sendMessage = () => {
    const token = getToken();
    if (!inputValue) return;
    axios
      .post("http://localhost:8080/chat/send", {
        Token: token,
        Message: inputValue,
      })
      .then((res) => {
        console.log(res.data);
        setInputValue(""); // Reset input Value to empty string
      })
      .catch((err) => {
        alert(err);
      });
  };

  useEffect(() => {
    fetchMessages();
  }, []);

  useEffect(() => {
    const interval = setInterval(fetchMessages, 100);

    return () => {
      clearInterval(interval);
    };
  }, []);

  const handleInputChange = (event) => {
    setInputValue(event.target.value);
  };

  const handleKeyDown = (event) => {
    if (event.key === "Enter") {
      event.preventDefault();
      sendMessage();
    }
  };

  return (
    <section id="info_chatbox" >
      <div className="chatbox-container">
        <h1>Live Chat</h1>
        <ul id="chatbox-messages" style={{ scrollBehavior: "smooth", height: "500px", overflowY: "scroll" }}>
          {(!messages)? "" :messages.map((message, index) => (
            <Message
              key={index}
              userId={message.id_sender}
              setPage={props.setPage}
              setBody={props.setBody}
              message={message.message}
            />
          ))}
        </ul>
        <input
          type="text"
          value={inputValue}
          onChange={handleInputChange}
          onKeyDown={handleKeyDown}
        />
        <button onClick={sendMessage}>Send</button>
      </div>
    </section>
  );
}

export default ChatBox;
