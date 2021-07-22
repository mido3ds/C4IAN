import './GalleryItem.css';
import React, { useRef, useState, useEffect } from 'react';
import Moment from 'react-moment';
import PlayAudio from '../../../PlayAudio/PlayAudio'
import PlayVideo from '../../../PlayVideo/PlayVideo'

function GalleryItem({ type, data, time, name }) {
    const playAudioRef = useRef(null);
    const playVideoRef = useRef(null);
    const [audioData, setAudioData] = useState(null)
    const [videoData, setVideoData] = useState(null)

    useEffect(()=> {
        if (type === "audio") {
            setAudioData(data.body)
        } else if (type === "video") {
            setVideoData(data)
        }
    }, [data, type])


    var playMedia = () => {
        if (type === "audio") {
            playAudioRef.current.open()
        } else if (type === "video") {
            playVideoRef.current.open()
        }
    }

    return (
        <>
            {type === "video" ?
                <PlayVideo videoData={videoData} ref={playVideoRef}></PlayVideo>
                : <PlayAudio name={name} audio={audioData} ref={playAudioRef}></PlayAudio>
            }
            <div data-augmented-ui="border" className="gallery-item">
                <i onClick={playMedia} className="fas fa-play-circle fa-3x gallery-item-play-icon"></i>
                <Moment className="gallery-item-time" format="wo MMM hh:mm" unix>{time/ (1000*1000)}</Moment>
            </div>
        </>
    );
}
export default GalleryItem;
