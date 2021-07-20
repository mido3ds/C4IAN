import './GalleryItem.css';
import React, { useRef, useState, useEffect } from 'react';
import Moment from 'react-moment';
import PlayAudio from '../../../PlayAudio/PlayAudio'

function GalleryItem({ type, data, time, name }) {
    const playAudioRef = useRef(null);
    const [audioData, setAudioData] = useState(null)

    useEffect(()=> {
        if (type === "audio") {
            setAudioData(data.body)
        }
    }, [data, type])


    var playMedia = () => {
        if (type === "audio") {
            playAudioRef.current.open()
        }
    }

    return (
        <>
            {type === "audio" ?
                <PlayAudio name={name} audio={audioData} ref={playAudioRef}></PlayAudio> :
                <></>
            }
            <div data-augmented-ui="border" className="gallery-item">
                <i onClick={playMedia} className="fas fa-play-circle fa-3x gallery-item-play-icon"></i>
                <Moment className="gallery-item-time" format="wo MMM hh:mm" unix>{time}</Moment>
            </div>
        </>
    );
}
export default GalleryItem;
