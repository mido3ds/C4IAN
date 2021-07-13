import './Gallery.css';
import React from 'react';
import GalleryItem from './GalleryItem/GalleryItem'
import Pagination from './Pagination/Pagination'

function Gallery({type}) {
    return (
        <div className="gallery-container">
            <div className="items-container">
                <GalleryItem />
                <GalleryItem />
                <GalleryItem />
                <GalleryItem />
            </div> 
            <Pagination
                paginate={5}
                hasNext={true}
            />
        </div>
    );
}
export default Gallery;
