import React, { useRef, useEffect, useState } from 'react';
import mapboxgl from '!mapbox-gl'; // eslint-disable-line import/no-webpack-loader-syntax
import uImage from '../images/unit.png';
import './Home.css'

mapboxgl.accessToken = 'pk.eyJ1IjoiYWhtZWRhZmlmaSIsImEiOiJja3F6YzJibGUwNXEyMnNsZ2U2N2lod2xqIn0.U2YYTWHCYqkCUBaAFd9MfA';

function Home() {
    const mapContainer = useRef(null);
    const map = useRef(null);
    const [lng, setLng] = useState(-70.9);
    const [lat, setLat] = useState(42.35);
    const [zoom, setZoom] = useState(4);

    useEffect(() => {
        if (map.current) return; // initialize map only once
        map.current = new mapboxgl.Map({
            container: mapContainer.current,
            style: 'mapbox://styles/ahmedafifi/ckqzcjls7amo118mk9xl5a4j3',
            center: [lng, lat],
            zoom: zoom
        });
        
         // create a HTML element for each feature
       /* var el = document.createElement('div');
        el.className = 'map-unit';

        // Create a default Marker and add it to the map.
        var marker1 = new mapboxgl.Marker(el)
        .setLngLat([lng, lat])
        .addTo(map.current);
        
        var marker2 = new mapboxgl.Marker(el)
        .setLngLat([-70.9, 42.35])
        .addTo(map.current);*/
        
    });

    return (
        <>
        <div className="map-wrapper">
            <div ref={mapContainer} className="map-container" />
        </div>
        </>
    );
}

export default Home;
