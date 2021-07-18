import './Streams.css';
import React from 'react';
import Stream from './Stream/Stream'

class Streams extends React.Component {

    constructor(props) {
        super(props)

        this.state = {
            Streams: []
        }

        var eventSource = new EventSource("http://localhost:3170/events")
        eventSource.addEventListener("video-fragment", ev => {
            this.onReceiveFragment(JSON.parse(ev.data))
        })
    }

    onReceiveFragment(data) {
        if (this.state.Streams.some(e => e.ID === data.ID)) {
            // Add Fragment
        } else {
            this.state.Streams.push({ ID: data.ID })
        }
    }

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
                        <Stream> </Stream>
                    </div>
                    {this.state.Streams.map((value, index) => {
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


