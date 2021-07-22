import React, { useEffect, useState } from 'react';
import uImage from '../../images/unit.png';
import { receivedCodes, sentCodes } from '../../codes'
import './ChatBox.css'

function ChatBox({ msgs }) {

    useEffect(() => {
        if (!msgs || !msgs.length) {
            window.$('.chat-box').css('overflow-y', 'hidden')
        } else {
            window.$('.chat-box').css('overflow-y', 'scroll')
            var element = document.querySelector(".chat-box");
            if(element) element.scrollTop = element.scrollHeight;
        }
        return
    },[msgs])

    return (
            <div className="chat-container">
                {msgs && msgs.length ?
                    <ul className="chat-box chatContainerScroll">
                        {msgs.map((msg, _) => {
                            return <>
                                {msg.sent ?
                                    <li className="chat-right">
                                            <div className="chat-text"> {sentCodes[msg.Code]} </div>
                                        <div className="chat-avatar">
                                            <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                                            <div className="chat-name"> unit soldier </div>
                                        </div>
                                    </li> :
                                    <li className="chat-left">
                                        <div className="chat-avatar">
                                            <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                                            <div className="chat-name"> cmd leader</div>
                                        </div>
                                        <div className="chat-text"> {receivedCodes[msg.Body]} </div>
                                    </li>
                                }
                            </>
                        })}
                    </ul> :
                    <div className="no-data-chat-msg"> 
                        <p> No data to be previewed </p> 
                    </div>
                    }
            </div>
    )

} export default ChatBox;