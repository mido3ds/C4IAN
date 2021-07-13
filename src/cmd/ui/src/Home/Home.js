import React, { useRef, useEffect, useState } from 'react';
import mapboxgl from '!mapbox-gl'; // eslint-disable-line import/no-webpack-loader-syntax
import MapPopup from './MapPopUp/MapPopUp'
import { units } from "../units";

import './Home.css'

mapboxgl.accessToken = 'pk.eyJ1IjoiYWhtZWRhZmlmaSIsImEiOiJja3F6YzJibGUwNXEyMnNsZ2U2N2lod2xqIn0.U2YYTWHCYqkCUBaAFd9MfA';

const numDeltas = 100;

function Home() {
    const mapContainer = useRef(null);
    const map = useRef(null);
    const [lng, setLng] = useState(45.234);
    const [lat, setLat] = useState(35.734);
    const [zoom, setZoom] = useState(5.04);
    const [selectedUnit, setSelectedUnit] = useState(null);

    var getBounds = (coordinates) => {
        var lngB = []
        var latB = []
        lngB[0] = Number.MAX_SAFE_INTEGER;
        lngB[1] = Number.MIN_SAFE_INTEGER;
        latB[0] = Number.MAX_SAFE_INTEGER;
        latB[1] = Number.MIN_SAFE_INTEGER;
        coordinates.map((c) => {
            if (c[0] < lngB[0]) lngB[0] = c[0];
            if (c[0] > lngB[1]) lngB[1] = c[0];
            if (c[1] < latB[0]) latB[0] = c[1];
            if (c[1] > latB[1]) latB[1] = c[1];
        });

        return [[lngB[0] - 2, latB[0] - 2], [lngB[1] + 2, latB[1] + 2]]
    }

    var moveMarker = (marker, steps, delta) => {
        marker.setLngLat([marker.getLngLat().lng + delta[lng], marker.getLngLat().lat + delta[lat]])

        if (steps !== numDeltas) {
            steps++;
            setTimeout(() => moveMarker(marker, steps, delta[lng], delta[lat]), 100);
        }
    }

    var onPositionChange = () => {
        var origin = { lng: 43.234, lat: 33.734 }
        var destination = { lng: 43.234, lat: 37.734 }

        var delta = {};
        delta[lng] = (destination.lng - origin.lng) / numDeltas;
        delta[lat] = (destination.lat - origin.lat) / numDeltas;
        moveMarker(null, delta, 0)
    }

    useEffect(() => {
        if (map.current) return; // initialize map only once
        map.current = new mapboxgl.Map({
            container: mapContainer.current,
            style: 'mapbox://styles/ahmedafifi/ckr0krxez6p641ao9vl2p71vf',
            center: [lng, lat],
            zoom: zoom
        });

        map.current.addControl(new mapboxgl.FullscreenControl());
        map.current.addControl(new mapboxgl.NavigationControl());

        var coordinates = []
        units.forEach(function (unit, index, units) {
            var el = document.createElement('div');
            el.className = 'map-unit' + unit.group;
            window.$(el).bind("click", () => {
                setSelectedUnit(() => {
                    console.log("hello")
                    return {name: unit.name, ip:unit.ip}
                })
            });
            coordinates.push([unit.lng, unit.lat])
            var marker = new mapboxgl.Marker(el)
                .setLngLat([unit.lng, unit.lat])
                .addTo(map.current);

            units[index] = { ...unit, marker: marker };
        });

        map.current.fitBounds(getBounds(coordinates));
    });

    return (
        <>
            <MapPopup selectedUnit={selectedUnit} />
            <div>
                <div ref={mapContainer} className="map-container" />
            </div>
        </>
    );
}

export default Home;