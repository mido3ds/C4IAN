import './Streams.css';
import React, { useEffect } from 'react';
import Stream from './Stream/Stream'


function Streams({streams, onEndStream, port}) {

    useEffect(() => {
        var number = window.$('.stream').length;
        var columns = Math.ceil(Math.sqrt(number) - 0.1);
        var rows = Math.ceil(number / columns);
        window.$('.stream').css("width", `${100 / columns}%`);
        window.$('.stream').css("height", `${100 / rows}%`);
    }, [streams])

    return (
        <div className="video-root">
            <div className="streams-container">
            {!streams || !streams.length ?
                <p className="no-streams-msg"> No streams available </p> :
                streams.map((value, index) => {
                    return <div className="stream">
                        <Stream onEndStream={onEndStream} stream={value} port={port}> </Stream>
                    </div>
                })}
            </div>
        </div>
    );
    
}
export default Streams;


