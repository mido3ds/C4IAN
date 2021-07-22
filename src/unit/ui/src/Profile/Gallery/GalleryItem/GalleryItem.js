import './GalleryItem.css';
import React, { useRef, useState, useEffect } from 'react';
import Moment from 'react-moment';
import PlayAudio from '../../../PlayAudio/PlayAudio'


function GalleryItem({ data, name, time }) {
    const playAudioRef = useRef(null);
    const [audioData, setAudioData] = useState(null)

    useEffect(()=> {
        setAudioData(data.body)
    }, [data])


    var playMedia = () => {
        playAudioRef.current.open()
    }

    return (
        <>
            <PlayAudio name={name} audio={audioData} ref={playAudioRef}> </PlayAudio> 
            
            <div data-augmented-ui="border" className="gallery-item">
                <i onClick={playMedia} className="fas fa-play-circle fa-3x gallery-item-play-icon"></i>
                <Moment className="gallery-item-time" format="wo MMM hh:mm" unix>{time / (1000*1000)}</Moment>
            </div>
        </>
    );
}
export default GalleryItem;
