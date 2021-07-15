import React, { useRef, useEffect, useState } from 'react';
import mapboxgl from '!mapbox-gl'; // eslint-disable-line import/no-webpack-loader-syntax
import MapPopup from './MapPopUp/MapPopUp'
import {NotificationManager} from 'react-notifications';

import './Home.css'

const HeartBeatThreshold = 60

mapboxgl.accessToken = 'pk.eyJ1IjoiYWhtZWRhZmlmaSIsImEiOiJja3F6YzJibGUwNXEyMnNsZ2U2N2lod2xqIn0.U2YYTWHCYqkCUBaAFd9MfA';

const numDeltas = 100;

function Home() {
    const mapContainer = useRef(null);
    const map = useRef(null);
    const [units, setUnits] = useState({})
    const [selectedUnit, setSelectedUnit] = useState(null);
    const [eventSource, setEventSource] = useState(new EventSource("http://localhost:3170/events"))

    var onTimeout = (unitIP) => {
        var el = document.createElement('div');
        el.className = 'map-unit-inactive';
        window.$(el).bind("click", () => {
            setSelectedUnit(() => {
                return { ip: unitIP }
            })
        });

        units[unitIP].marker.remove()
        units[unitIP].marker = new mapboxgl.Marker(el)
                                            .setLngLat([units[unitIP].lng, units[unitIP].lat])
                                            .addTo(map.current);
        NotificationManager.error(unitIP + " is inactive for 2 minutes!");
    }

    var onDanger = (unitIP) => {
        var el = document.createElement('div');
        el.className = 'map-unit-danger';
        window.$(el).bind("click", () => {
            setSelectedUnit(() => {
                return { ip: unitIP }
            })
        });

        units[unitIP].marker.remove()
        units[unitIP].marker = new mapboxgl.Marker(el)
                                            .setLngLat([units[unitIP].lng, units[unitIP].lat])
                                            .addTo(map.current);
        NotificationManager.error(unitIP + " is in danger!!");
    }

    var onDataChange = (newData) => {
        setUnits(() => {
            var unitsCopy = JSON.parse(JSON.stringify(units));
            if (newData.src in unitsCopy) {
                var oldPosition = [unitsCopy[newData.src].lng, unitsCopy[newData.src].lat]
                var newPosition = [newData.loc_x, newData.loc_y]
                onPositionChange(oldPosition, newPosition, unitsCopy[newData.src].marker)

                clearTimeout(unitsCopy[newData.src].timerID);
            } else {
                var el = document.createElement('div');
                el.className = 'map-unit' + unitsCopy[newData.src].groupID;
                window.$(el).bind("click", () => {
                    setSelectedUnit(() => {
                        return { ip: newData.src }
                    })
                });
                unitsCopy[newData.src].marker = new mapboxgl.Marker(el)
                    .setLngLat([newData.loc_x, newData.loc_y])
                    .addTo(map.current);
                map.current.fitBounds(getBounds());

            }

            if(unitsCopy[newData.src].heartbeat < HeartBeatThreshold)
                onDanger(newData.src)
                
            unitsCopy[newData.src].lng = newData.loc_x
            unitsCopy[newData.src].lat = newData.loc_y
            unitsCopy[newData.src].heartbeat = newData.heartbeat
            unitsCopy[newData.src].timerID = setTimeout(() => { onTimeout(newData.src) }, 2 * 60 * 1000)
        })
    }

    var getBounds = () => {
        var coordinates = []
        setUnits(() => {
            for (const [key, value] of Object.entries(units)) {
                coordinates.push([value.loc_x, value.loc_y])
            }
            return units
        })

        var lngB = [Number.MAX_SAFE_INTEGER, Number.MIN_SAFE_INTEGER]
        var latB = [Number.MAX_SAFE_INTEGER, Number.MIN_SAFE_INTEGER]
        coordinates.map((c) => {
            if (c[0] < lngB[0]) lngB[0] = c[0];
            if (c[0] > lngB[1]) lngB[1] = c[0];
            if (c[1] < latB[0]) latB[0] = c[1];
            if (c[1] > latB[1]) latB[1] = c[1];
        });

        return [[lngB[0] - 2, latB[0] - 2], [lngB[1] + 2, latB[1] + 2]]
    }

    var moveMarker = (marker, steps, delta) => {
        marker.setLngLat([marker.getLngLat().lng + delta[0], marker.getLngLat().lat + delta[1]])

        if (steps !== numDeltas) {
            steps++;
            setTimeout(() => moveMarker(marker, steps, delta), 100);
        }
    }

    var onPositionChange = (origin, destination, marker) => {
        var delta = [];
        delta[0] = (destination[0] - origin[0]) / numDeltas;
        delta[1] = (destination[1] - origin[1]) / numDeltas;
        moveMarker(marker, delta, 0)
    }

    useEffect(() => {
        eventSource.addEventListener("sensors-data", ev => {
            onDataChange(JSON.parse(ev.data))
        })

        if (map.current) return; // initialize map only once
        map.current = new mapboxgl.Map({
            container: mapContainer.current,
            style: 'mapbox://styles/ahmedafifi/ckr0krxez6p641ao9vl2p71vf',
        });

        map.current.addControl(new mapboxgl.FullscreenControl());
        map.current.addControl(new mapboxgl.NavigationControl());
    });

    return (
        <>
            <MapPopup selectedUnit={units[selectedUnit]} />
            <div>
                <div ref={mapContainer} className="map-container" />
            </div>
        </>
    );
}

export default Home;