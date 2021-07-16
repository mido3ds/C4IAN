import React from 'react';
import ReactPlayer from 'react-player'
import './Stream.css'
import video from './v.mp4'

class Stream extends React.Component {
    constructor(props) {
        super(props)
        this.setState({ 
            video: null 
        })
    }

    componentDidMount() {
        /*import('./v.mp4')
        .then(module => this.setState({ video: module.default }))*/
    }

    render() {
        return (
            <ReactPlayer width="100%" height="100%" controls url={video} />
        );
    }

} export default Stream;
