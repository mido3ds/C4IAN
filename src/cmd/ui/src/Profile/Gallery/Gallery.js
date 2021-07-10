import './Gallery.css';
import React from 'react';

function Gallery() {
    return (
        <div className="gallery-container">
            <div className="gallery-header">

            </div>
            <div className="items-container"> 
                <div data-augmented-ui="border" className="gallery-item"> </div>
                <div data-augmented-ui="border" className="gallery-item"> </div>
                <div data-augmented-ui="border" className="gallery-item"> </div>
                <div data-augmented-ui="border" className="gallery-item"> </div>           
            </div>
            <div className="pagination"></div>
        </div>
    );
}
export default Gallery;
