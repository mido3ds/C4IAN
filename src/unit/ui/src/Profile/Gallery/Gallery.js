import './Gallery.css';
import React, { useState, useEffect } from 'react';
import GalleryItem from './GalleryItem/GalleryItem'
import Pagination from './Pagination/Pagination'
import { AudiosDB } from '../../db'

 function Gallery({type, audios}) {
    const [data, setData] = useState(null);
    const [startItem, setStartItem] = useState(0)

    var paginate = (pageNumber) => {
        setStartItem(4 * (pageNumber - 1));
    }

    useEffect(() => {
        if (audios) {
            setData(audios)
        }
    },[audios])

    return (
        <div className="gallery-container">
            {!data || !data.length ?
            <div className="no-data-gallery-msg"> 
                <p> No data to be previewed </p> 
            </div>: 
            <> 
                <div className="items-container">
                    {
                        [...Array(Math.min(data.length - startItem, 4))].map((x, i) => 
                            <GalleryItem name= {"Command Center"} type={type} data={data[i]} key={i} />
                        )
                    }
                </div>
                <Pagination
                    dataSize={data.length}
                    paginate={paginate}
                />
            </>
        }
        </div>
    );
}
export default Gallery;
