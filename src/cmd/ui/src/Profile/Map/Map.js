import React, { useRef, useEffect, useState } from 'react';
import mapboxgl from '!mapbox-gl'; // eslint-disable-line import/no-webpack-loader-syntax
import './Map.css'

mapboxgl.accessToken = 'pk.eyJ1IjoiYWhtZWRhZmlmaSIsImEiOiJja3F6YzJibGUwNXEyMnNsZ2U2N2lod2xqIn0.U2YYTWHCYqkCUBaAFd9MfA';

const coordinates = [
    [-122.483696, 37.833818],
    [-122.483482, 37.833174],
    [-122.483396, 37.8327],
    [-122.483568, 37.832056],
    [-122.48404, 37.831141],
    [-122.48404, 37.830497],
    [-122.483482, 37.82992],
    [-122.483568, 37.829548],
    [-122.48507, 37.829446],
    [-122.4861, 37.828802],
    [-122.486958, 37.82931],
    [-122.487001, 37.830802],
    [-122.487516, 37.831683],
    [-122.488031, 37.832158],
    [-122.488889, 37.832971],
    [-122.489876, 37.832632],
    [-122.490434, 37.832937],
    [-122.49125, 37.832429],
    [-122.491636, 37.832564],
    [-122.492237, 37.833378],
    [-122.493782, 37.833683]
]

function Map() {
    const profileMapContainer = useRef(null);
    const map = useRef(null);

    useEffect(() => {
        if (map.current) return; // initialize map only once
        var center = [...coordinates[Math.ceil(coordinates.length / 2)]]
        center[0] -= 0.005
        console.log(center)
        map.current = new mapboxgl.Map({
            container: profileMapContainer.current,
            style: 'mapbox://styles/ahmedafifi/ckr3eqazg5ndn18p3pgmuffc1',
            center: center,
            zoom: 15
        });

        map.current.on('load', function () {
            map.current.addSource('route', {
                'type': 'geojson',
                'data': {
                    'type': 'Feature',
                    'properties': {},
                    'geometry': {
                        'type': 'LineString',
                        'coordinates': coordinates
                    }
                }
            });

            map.current.addLayer({
                'id': 'route',
                'type': 'line',
                'source': 'route',
                'layout': {
                    'line-join': 'round',
                    'line-cap': 'round'
                },
                'paint': {
                    'line-color': '#888',
                    'line-width': 8
                }
            });
        })
        var marker = new mapboxgl.Marker({color: 'black'})
                .setLngLat(coordinates[coordinates.length - 1])
                .addTo(map.current);

    })

    return (
        <>
            <div className="profile-map-wrapper">
                <div ref={profileMapContainer} className="profile-map-container" />
            </div>
        </>
    );
}

export default Map;