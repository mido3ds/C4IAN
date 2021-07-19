import React from 'react';
import './Stream.css'
import ReactHlsPlayer from 'react-hls-player';
// import video from './v.mp4'

const baseURL = "http://localhost:3170/api/";
const M3U8Name   = "index.m3u8";

const HLSConfig = {
    maxLoadingDelay: 4,
    minAutoBitrate: 0,
    lowLatencyMode: true,
}

class Stream extends React.Component {
    constructor(props) {
        super(props)
        this.setState({ 
            video: null
        })
    }

    componentDidMount() {
    }

    render() {
        const playerRef = React.useRef();

        React.useEffect(() => {
          function fireOnVideoStart() {
            // Do some stuff when the video starts/resumes playing
          }
      
          playerRef.current.addEventListener('play', fireOnVideoStart);
      
          return playerRef.current.removeEventListener('play', fireOnVideoStart);
        }, []);
      
        React.useEffect(() => {
          function fireOnVideoEnd() {
            // Do some stuff when the video ends
          }
      
          playerRef.current.addEventListener('ended', fireOnVideoEnd);
      
          return playerRef.current.removeEventListener('ended', fireOnVideoEnd);
        }, []);

        return (
            <ReactHlsPlayer
            playerRef={playerRef}
            src={baseURL + "stream/" + ip + "/" + id + "/" + M3U8Name}
            autoPlay={true}
            controls={true}
            width={auto}
            height={auto}
            hlsConfig={HLSConfig}  
            />
        );
    }

} export default Stream;
