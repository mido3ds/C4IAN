import './Streams.css';
import React from 'react';
import Stream from './Stream/Stream'

class Streams extends React.Component {

    componentDidMount() {
        var number = window.$('.stream').length;
        var columns = Math.ceil(Math.sqrt(number) - 0.1);
        var rows = Math.ceil(number / columns);
        window.$('.stream').css("width", `${100 / columns}%`);
        window.$('.stream').css("height", `${100 / rows}%`);
    }


    render() {
        return (
            <div className="video-root">
                <div className="streams-container">
                    <div className="stream">
                        <Stream></Stream>
                    </div>
                </div>
            </div>
        );
    }
}
export default Streams;


