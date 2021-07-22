import React, { useRef, useEffect, useState, forwardRef, useImperativeHandle } from 'react';
import mapboxgl from '!mapbox-gl'; // eslint-disable-line import/no-webpack-loader-syntax
import MapPopup from './MapPopUp/MapPopUp'
import GroupSelect from '../GroupSelect/GroupSelect'
import ChatWindow from './ChatWindow/ChatWindow';
import { NotificationManager } from 'react-notifications';
import { getUnits, getMembers } from '../Api/Api'
import { groupIDs } from '../groupIDs'
import './Home.css'

const HeartBeatThreshold = 60

mapboxgl.accessToken = 'pk.eyJ1IjoiYWhtZWRhZmlmaSIsImEiOiJja3F6YzJibGUwNXEyMnNsZ2U2N2lod2xqIn0.U2YYTWHCYqkCUBaAFd9MfA';

const numDeltas = 50;

const Home = forwardRef(({selectedTab, port}, ref) => {
    const mapContainer = useRef(null);
    const msgsRef = useRef(null);
    const map = useRef(null);
    const [units, setUnits] = useState({})
    const [selectedUnit, setSelectedUnit] = useState(null);

    useImperativeHandle(ref, () => ({
        onNewMessage(message) {
            msgsRef.current.onNewMessage(message)
        }
    }));

    var onTimeout = (unitIP) => {
        setUnits(units => {
            units[unitIP].active = false;
            if(selectedTab !== "Log Out")
                NotificationManager.error(units[unitIP].name + " is inactive for 2 minutes!");
            drawMarker(unitIP, units)
            return units
        })
    }

    var drawMarker = (unitIP, units) => {
        var el = document.createElement('div');
        el.className = units[unitIP].InDanger ? "map-unit-danger" :
                        !units[unitIP].active ? "map-unit-inactive" :
                        units[unitIP].hasOwnProperty("groupID") ? 'map-unit' + units[unitIP].groupID :
                        "map-unit";
        window.$(el).bind("click", () => {
            setSelectedUnit(unitIP)
        });

        if(units[unitIP].marker) units[unitIP].marker.remove()
        units[unitIP].marker = new mapboxgl.Marker(el)
            .setLngLat([units[unitIP].lng, units[unitIP].lat])
            .addTo(map.current);
    } 

    var onDataChange = (newData) => {
        setUnits(units => {
            if (newData.heartbeat <= HeartBeatThreshold && !units[newData.src].InDanger) {
                units[newData.src].InDanger = true;
                if(selectedTab !== "Log Out")
                    NotificationManager.error(units[newData.src].name + " is in danger!!");
            } else if (newData.heartbeat > HeartBeatThreshold && units[newData.src].InDanger) {
                units[newData.src].InDanger = false;
                if(selectedTab !== "Log Out")
                    NotificationManager.info(units[newData.src].name + " is no more in danger");
            } 

            if (!units[newData.src].active) {
                if(selectedTab !== "Log Out")
                    NotificationManager.info(units[newData.src].name + " is active now");
                units[newData.src].active = true;
            }
            

            if (units[newData.src].hasOwnProperty("marker")) {
                drawMarker(newData.src, units)
                onPositionChange(newData, units, units[newData.src].marker)
                clearTimeout(units[newData.src].timerID);
                units[newData.src].lng = newData.lon
                units[newData.src].lat = newData.lat
            } else {
                units[newData.src].lng = newData.lon
                units[newData.src].lat = newData.lat
                drawMarker(newData.src, units)
            }

            //map.current.fitBounds(getBounds(units));

            units[newData.src].heartbeat = newData.heartbeat
            units[newData.src].timerID = setTimeout(() => { onTimeout(newData.src) }, 2 * 60 * 1000)

            return units;
        })
    }

    var getBounds = unitsCopy => {
        var coordinates = []
        for (const value of Object.values(unitsCopy)) {
            if (value.lng !== 1000 && value.lat !== 1000) {
                coordinates.push([value.lng, value.lat])
            }
        }
       
        if (!coordinates.length)
            return [[0, 0], [0, 0]]

        var lngB = [Number.MAX_SAFE_INTEGER, Number.MIN_SAFE_INTEGER]
        var latB = [Number.MAX_SAFE_INTEGER, Number.MIN_SAFE_INTEGER]
        coordinates.map((coordinate) => {
            if (coordinate[0] < lngB[0]) lngB[0] = coordinate[0];
            if (coordinate[0] > lngB[1]) lngB[1] = coordinate[0];
            if (coordinate[1] < latB[0]) latB[0] = coordinate[1];
            if (coordinate[1] > latB[1]) latB[1] = coordinate[1];
            return coordinate
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

    var onPositionChange = (newData, units, marker) => {
        var origin = [units[newData.src].lng, units[newData.src].lat]
        var destination = [newData.lon, newData.lat]

        var delta = [];
        delta[0] = (destination[0] - origin[0]) / numDeltas;
        delta[1] = (destination[1] - origin[1]) / numDeltas;
        moveMarker(marker, 0, delta)
    }

    useEffect(() => {
        if(!port) return
        if (Object.keys(units).length ) return
        if (map.current) return; 
        map.current = new mapboxgl.Map({
            container: mapContainer.current,
            style: 'mapbox://styles/ahmedafifi/ckr3eqazg5ndn18p3pgmuffc1',
            center: [0,0],
            zoom: 10
        });
        map.current.addControl(new mapboxgl.FullscreenControl());
        map.current.addControl(new mapboxgl.NavigationControl());

        var eventSource = new EventSource("http://localhost:" + port + "/events")
        eventSource.addEventListener("sensors-data", ev => {
            onDataChange(JSON.parse(ev.data))
        })

        getUnits(port).then(initialData => {
            getMembers(port).then(members => {
                setUnits(units => {
                    initialData.forEach(unit => {
                        units[unit.ip] = { name: unit.name, ip: unit.ip, active: unit.active, lng: unit.lon, lat: unit.lat, heartbeat: unit.heartbeat, InDanger: unit.heartbeat < HeartBeatThreshold }
                    });

                    members.forEach(membership => {
                        units[membership.unitIP] = { ...units[membership.unitIP], groupID: groupIDs[membership.groupIP] }
                    });

                    for (const unitIP in units) {
                        if (units[unitIP].lng !== 1000 && units[unitIP].lat !== 1000) {
                            drawMarker(unitIP, units)
                        }
                    }

                    map.current.fitBounds(getBounds(units));
                    
                    return units
                })
            })
        })
    }, [port]);

    return (
        <>
            <GroupSelect port={port}></GroupSelect>
            <ChatWindow port={port} ref={msgsRef}></ChatWindow>
            <MapPopup selectedUnit={units[selectedUnit]} />
            <div>
                <div ref={mapContainer} className="map-container" />
            </div>
        </>
    );
});
export default Home;
