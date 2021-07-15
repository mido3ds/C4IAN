import React, { useRef, useEffect, useState } from 'react';
import mapboxgl from '!mapbox-gl'; // eslint-disable-line import/no-webpack-loader-syntax
import './Map.css'
import { getSensorsData } from '../../Api/Api'

mapboxgl.accessToken = 'pk.eyJ1IjoiYWhtZWRhZmlmaSIsImEiOiJja3F6YzJibGUwNXEyMnNsZ2U2N2lod2xqIn0.U2YYTWHCYqkCUBaAFd9MfA';

function Map({unit}) {
    const profileMapContainer = useRef(null);
    const map = useRef(null);

    var getCoordinates = () => {
        var sensorData = getSensorsData(unit.ip)
        var coordinates = []
        sensorData.forEach((item, index) => {
            coordinates.push([item.lat, item.lon])
        })
        return coordinates
    }

    useEffect(() => {
        const coordinates = getCoordinates()
        if(coordinates.length === 0) return
        if (map.current) return; 
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
        <>  <div className="no-data-gallery-msg"> 
                <p> No data to be previewed </p> 
            </div>
            <div className="profile-map-wrapper">
                <div ref={profileMapContainer} className="profile-map-container" />
            </div>
        </>
    );
}

export default Map;