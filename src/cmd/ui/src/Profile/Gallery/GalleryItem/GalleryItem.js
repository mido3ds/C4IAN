import './GalleryItem.css';
import React, { useRef } from 'react';
import Moment from 'react-moment';
import PlayAudio from '../../../PlayAudio/PlayAudio'
import PlayVideo from '../../../PlayVideo/PlayVideo'

function GalleryItem({type, data, time, name}) {
    const playAudioRef = useRef(null);
    const playVideoRef = useRef(null);

    var playMedia = () => {
        if(type === "audio") {
            playAudioRef.current.openModal()
        } else if(type === "video") {
            playVideoRef.current.openModal()
        }
    }
    
    return (
        <>
        <PlayVideo videoUrl={data} ref={playVideoRef}></PlayVideo>
        <PlayAudio name={name} audioBolb={data} ref={playAudioRef}></PlayAudio>
        <div  data-augmented-ui="border" className="gallery-item">
            <i onClick={playMedia} className="fas fa-play-circle fa-3x gallery-item-play-icon"></i>
            <Moment className="gallery-item-time" format="wo MMM hh:mm" unix>{time}</Moment>
        </div>
        </>
    );
}
export default GalleryItem;
