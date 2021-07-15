import './GalleryItem.css';
import React, { useRef } from 'react';
import Moment from 'react-moment';
import PlayAudio from '../../../PlayAudio/PlayAudio'

function GalleryItem({type, data, time, name}) {
    const playAudioRef = useRef(null);
    var playAudio = (audio) => {
        playAudioRef.current.openModal()
    }
    
    return (
        <>
        <PlayAudio name={name} audioBolb={data} ref={playAudioRef}></PlayAudio>
        <div  data-augmented-ui="border" className="gallery-item">
            <i onClick={playAudio} className="fas fa-play-circle fa-3x gallery-item-play-icon"></i>
            <Moment className="gallery-item-time" format="wo MMM hh:mm" unix>{time}</Moment>
        </div>
        </>
    );
}
export default GalleryItem;
