import './GalleryItem.css';
import React, { useRef } from 'react';
import Moment from 'react-moment';
import PlayAudio from '../../../PlayAudio/PlayAudio'

function GalleryItem({type, data, time, name}) {
    const playAudioRef = useRef(null);
    var playMedia = (audio) => {
        if(type === "audio") {
            playAudioRef.current.openModal()
        } else if(type === "audio") {
            
        }
    }
    
    return (
        <>
        <PlayAudio name={name} audioBolb={data} ref={playAudioRef}></PlayAudio>
        <div  data-augmented-ui="border" className="gallery-item">
            <i onClick={playMedia} className="fas fa-play-circle fa-3x gallery-item-play-icon"></i>
            <Moment className="gallery-item-time" format="wo MMM hh:mm" unix>{time}</Moment>
        </div>
        </>
    );
}
export default GalleryItem;
