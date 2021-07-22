import React, { useRef, useEffect } from 'react';
import mapboxgl from '!mapbox-gl'; // eslint-disable-line import/no-webpack-loader-syntax
import { getSensorsData } from '../../Api/Api'
import './Map.css'

mapboxgl.accessToken = 'pk.eyJ1IjoiYWhtZWRhZmlmaSIsImEiOiJja3F6YzJibGUwNXEyMnNsZ2U2N2lod2xqIn0.U2YYTWHCYqkCUBaAFd9MfA';

function Map({ unit, port }) {
    const profileMapContainer = useRef(null);
    const map = useRef(null);

    var getBounds = coordinates => {
        var lngB = [Number.MAX_SAFE_INTEGER, Number.MIN_SAFE_INTEGER]
        var latB = [Number.MAX_SAFE_INTEGER, Number.MIN_SAFE_INTEGER]
        coordinates.map((coordinate) => {
            if (coordinate[0] < lngB[0]) lngB[0] = coordinate[0];
            if (coordinate[0] > lngB[1]) lngB[1] = coordinate[0];
            if (coordinate[1] < latB[0]) latB[0] = coordinate[1];
            if (coordinate[1] > latB[1]) latB[1] = coordinate[1];
            return coordinates
        });

        return [[lngB[0] - 2, latB[0] - 2], [lngB[1] + 2, latB[1] + 2]]
    }

    useEffect(() => {
        if(!unit || !port) return
        var coordinates = []
        getSensorsData(unit.ip, port).then(sensorData => {
            if (!sensorData || !sensorData.length) {
                if (map.current) {
                    map.current.remove()
                    map.current = null;
                }
                return 
            }

            sensorData.forEach((item, index) => {
                coordinates.push([item.lon, item.lat])
            })

            if (coordinates === null) {
                if (map.current)  {
                    map.current.remove()
                    map.current = null;
                }
                return 
            }

            if (map.current) {
                if (map.current) {
                    map.current.remove() 
                    map.current = null;
                }
            }

            var center = [...coordinates[coordinates.length - 1]]
            center[0] -= 0.005

            map.current = new mapboxgl.Map({
                container: profileMapContainer.current,
                style: 'mapbox://styles/ahmedafifi/ckr3eqazg5ndn18p3pgmuffc1',
                center: center,
                zoom: 7
            });

            map.current.addControl(new mapboxgl.FullscreenControl());
            map.current.addControl(new mapboxgl.NavigationControl());

            map.current.fitBounds(getBounds(coordinates));
    
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
            new mapboxgl.Marker({ color: 'black' })
                .setLngLat(coordinates[coordinates.length - 1])
                .addTo(map.current);
    
        })
    },[unit, port])

    return (
        <>  
            <div className="no-data-gallery-msg"> 
                <p> No data to be previewed </p> 
            </div>: 
            <div className="profile-map-wrapper">
                <div ref={profileMapContainer} className="profile-map-container" />
            </div>
        </>
    );
}

export default Map;