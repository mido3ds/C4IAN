import React, { useEffect, useState } from 'react';
import uImage from '../../images/unit.png';
import { getMsgs } from '../../Api/Api'
import moment from 'moment';
import { receivedCodes } from '../../codes'
import './ChatBox.css'

const sentCodes = { 2: "Start Streming", 3: "End Streaming", 4: "Attack", 5: "Defense", 6: "Escape", 7: "Regroup" }

function ChatBox({ unit, port }) {
    const [msgs, setMsgs] = useState(null)

    useEffect(() => {
        if (unit && port) {
            getMsgs(unit.ip, port).then(msgsData => {
                setMsgs(msgsData)
                if (!msgsData || !msgsData.length) {
                    window.$('.chat-container').css('overflow-y', 'hidden')
                } else {
                    window.$('.chat-container').css('overflow-y', 'scroll')
                    var element = document.querySelector(".chat-container");
                    if(element) element.scrollTop = element.scrollHeight;
                }
                return
            })
        }
    },[unit, port])

    return (
        <div className="content-wrapper">
            <div className="chat-container">
                {msgs && msgs.length ?
                    <ul className="chat-box chatContainerScroll">
                        {msgs.map((msg, index) => {
                            return <>
                                {msg.sent ?
                                    <li className="chat-right">
                                        <div className="chat-hour"> {moment.unix(msg.time / (1000*1000)).format('hh:mm')} </div>
                                        <div className="chat-text"> {sentCodes[msg.code]} </div>
                                        <div className="chat-avatar">
                                            <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                                            <div className="chat-name"> Me </div>
                                        </div>
                                    </li> :
                                    <li className="chat-left">
                                        <div className="chat-avatar">
                                            <img className="unit-item-profile-image" alt="unit" src={uImage}></img>
                                            <div className="chat-name"> {unit.name.substr(0, unit.name.indexOf(' '))} </div>
                                        </div>
                                        <div className="chat-text"> {receivedCodes.hasOwnProperty(msg.code) ? receivedCodes[msg.code] : (msg.code).toString(10)
} </div>
                                        <div className="chat-hour"> {moment.unix(msg.time / (1000*1000)).format('hh:mm')} </div>
                                    </li>
                                }
                            </>
                        })}
                    </ul> :
                    <div className="no-data-chat-msg">
                        <p> No data to be previewed </p>
                    </div>}
            </div>
        </div>
    )

} export default ChatBox;