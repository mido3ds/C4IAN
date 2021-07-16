import './Streams.css';
import React from 'react';
import Stream from './Stream/Stream'

class Streams extends React.Component {

    onReceiveFragment(data) {
        if(this.state.Streams.some(e => e.ID === data.ID)) {
            // Add Fragment
        } else {
            Streams.push({ID: Stream.ID})
        }
    }

    componentDidMount() {
        var number = window.$('.stream').length;
        var columns = Math.ceil(Math.sqrt(number) - 0.1);
        var rows = Math.ceil(number / columns);
        window.$('.stream').css("width", `${100 / columns}%`);
        window.$('.stream').css("height", `${100 / rows}%`);

        this.state = {
            eventSource: new EventSource("http://localhost:3170/events"),
            Streams : []
        }

        this.state.eventSource.addEventListener("video-fragment", ev => {
            this.onReceiveFragment(JSON.parse(ev.data))
        })
    }

    render() {
        return (
            <div className="video-root">
                <div className="streams-container">
                    {Streams.map((value, index) => {
                        return <div className="stream">
                                     <Stream streamID={value.ID}> </Stream>
                                </div>
                    })}
                </div>
            </div>
        );
    }
}
export default Streams;


