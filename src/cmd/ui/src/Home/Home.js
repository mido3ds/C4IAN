import React, { useRef, useEffect, useState, forwardRef, useImperativeHandle } from 'react';
import mapboxgl from '!mapbox-gl'; // eslint-disable-line import/no-webpack-loader-syntax
import MapPopup from './MapPopUp/MapPopUp'
import GroupSelect from '../GroupSelect/GroupSelect'
import ChatWindow from './ChatWindow/ChatWindow';
import { NotificationManager } from 'react-notifications';
import { getUnits, getMembers, getGroups } from '../Api/Api'
import './Home.css'

const HeartBeatThreshold = 50

mapboxgl.accessToken = 'pk.eyJ1IjoiYWhtZWRhZmlmaSIsImEiOiJja3F6YzJibGUwNXEyMnNsZ2U2N2lod2xqIn0.U2YYTWHCYqkCUBaAFd9MfA';

const numDeltas = 50;

const Home = forwardRef(({ selectedTab, port }, ref) => {
    const mapContainer = useRef(null);
    const msgsRef = useRef(null);
    const map = useRef(null);
    const [units, setUnits] = useState({})
    const [selectedUnit, setSelectedUnit] = useState(null);

    useImperativeHandle(ref, () => ({
        onNewMessage(message) {
            msgsRef.current.onNewMessage(message)
        },
        onChangePort(port) {
            onGetPort(port)
        }
    }));

    var onTimeout = (unitIP) => {
        setUnits(units => {
            units[unitIP].active = false;
            if (selectedTab !== "Log Out")
                NotificationManager.error(units[unitIP].name + " is inactive for 2 minutes!", '', 1000, () => {}, true);
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

        if (units[unitIP].marker) units[unitIP].marker.remove()
        units[unitIP].marker = new mapboxgl.Marker(el)
            .setLngLat([units[unitIP].lng, units[unitIP].lat])
            .addTo(map.current);
    }

    var onDataChange = (newData) => {
        setUnits(units => {
            if (newData.heartbeat <= HeartBeatThreshold && !units[newData.src].InDanger) {
                units[newData.src].InDanger = true;
                if (selectedTab !== "Log Out")
                    NotificationManager.error(units[newData.src].name + " is in danger!!", '', 1000, () => {}, true);
            } else if (newData.heartbeat > HeartBeatThreshold && units[newData.src].InDanger) {
                units[newData.src].InDanger = false;
                if (selectedTab !== "Log Out")
                    NotificationManager.info(units[newData.src].name + " is no more in danger", '', 1000, () => {}, true);
            }

            if (!units[newData.src].active) {
                if (selectedTab !== "Log Out")
                    NotificationManager.info(units[newData.src].name + " is active now", '', 1000, () => {}, true);
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
                map.current.fitBounds(getBounds(units));
            }


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

    var loadData = (receivedPort) => {
        var eventSource = new EventSource("http://localhost:" + receivedPort + "/events")
        eventSource.addEventListener("sensors-data", ev => {
            onDataChange(JSON.parse(ev.data))
        })
        getUnits(receivedPort).then(initialData => {
            getMembers(receivedPort).then(members => {
                getGroups(receivedPort).then(groups => {
                    setUnits(units => {
                        var groupsIPs = groups.map(group => group.ip).sort();
                        var groupIDs = {}
                        groupsIPs.forEach((ip, index) => {
                            groupIDs[ip] = index
                        });

                        initialData.forEach(unit => {
                            units[unit.ip] = { name: unit.name, ip: unit.ip, active: unit.active, lng: unit.lon, lat: unit.lat, heartbeat: unit.heartbeat, InDanger: unit.heartbeat < HeartBeatThreshold }
                        });

                        members.forEach(membership => {
                            units[membership.unitIP] = { ...units[membership.unitIP], groupID: groupIDs[membership.groupIP] }
                        });

                        let activeUnits = false
                        for (const unitIP in units) {
                            if (units[unitIP].lng !== 1000 && units[unitIP].lat !== 1000) {
                                drawMarker(unitIP, units)
                                activeUnits = true
                            }
                        }

                        if(initialData?.length && activeUnits) {
                            map.current.setCenter([initialData[Math.floor(initialData.length/2)].lon, initialData[Math.floor(initialData.length/2)].lat]);
                            map.current.setZoom(15);
                        }

                        return units
                    })
                })
            })
        })

    }

    var onGetPort = (receivedPort) => {
        loadData(receivedPort)
    }

    useEffect(() => {
        if (Object.keys(units).length) return
        if (map.current) return;
        map.current = new mapboxgl.Map({
            container: mapContainer.current,
            style: 'mapbox://styles/ahmedafifi/ckr3eqazg5ndn18p3pgmuffc1',
            center: [0, 0],
            zoom: 0
        });
        
        map.current.addControl(new mapboxgl.FullscreenControl());
        map.current.addControl(new mapboxgl.NavigationControl());
    }, []);

    return (
        <>
            {!port ? <> </> :
                <>
                    <img alt="focus" onClick={() => loadData(port)} className="focus-icon" src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZlcnNpb249IjEuMSIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHhtbG5zOnN2Z2pzPSJodHRwOi8vc3ZnanMuY29tL3N2Z2pzIiB3aWR0aD0iNTEyIiBoZWlnaHQ9IjUxMiIgeD0iMCIgeT0iMCIgdmlld0JveD0iMCAwIDUxMiA1MTIiIHN0eWxlPSJlbmFibGUtYmFja2dyb3VuZDpuZXcgMCAwIDUxMiA1MTIiIHhtbDpzcGFjZT0icHJlc2VydmUiIGNsYXNzPSIiPjxnPjxwYXRoIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgZD0ibTI1NiAzNjIuNjY3OTY5Yy01OC44MTY0MDYgMC0xMDYuNjY3OTY5LTQ3Ljg1MTU2My0xMDYuNjY3OTY5LTEwNi42Njc5NjlzNDcuODUxNTYzLTEwNi42Njc5NjkgMTA2LjY2Nzk2OS0xMDYuNjY3OTY5IDEwNi42Njc5NjkgNDcuODUxNTYzIDEwNi42Njc5NjkgMTA2LjY2Nzk2OS00Ny44NTE1NjMgMTA2LjY2Nzk2OS0xMDYuNjY3OTY5IDEwNi42Njc5Njl6bTAtMTgxLjMzNTkzOGMtNDEuMTcxODc1IDAtNzQuNjY3OTY5IDMzLjQ5NjA5NC03NC42Njc5NjkgNzQuNjY3OTY5czMzLjQ5NjA5NCA3NC42Njc5NjkgNzQuNjY3OTY5IDc0LjY2Nzk2OSA3NC42Njc5NjktMzMuNDk2MDk0IDc0LjY2Nzk2OS03NC42Njc5NjktMzMuNDk2MDk0LTc0LjY2Nzk2OS03NC42Njc5NjktNzQuNjY3OTY5em0wIDAiIGZpbGw9IiNmZmZmZmYiIGRhdGEtb3JpZ2luYWw9IiMwMDAwMDAiIHN0eWxlPSIiIGNsYXNzPSIiPjwvcGF0aD48cGF0aCB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGQ9Im0yNTYgNDY5LjMzMjAzMWMtMTE3LjYzMjgxMiAwLTIxMy4zMzIwMzEtOTUuNjk5MjE5LTIxMy4zMzIwMzEtMjEzLjMzMjAzMXM5NS42OTkyMTktMjEzLjMzMjAzMSAyMTMuMzMyMDMxLTIxMy4zMzIwMzEgMjEzLjMzMjAzMSA5NS42OTkyMTkgMjEzLjMzMjAzMSAyMTMuMzMyMDMxLTk1LjY5OTIxOSAyMTMuMzMyMDMxLTIxMy4zMzIwMzEgMjEzLjMzMjAzMXptMC0zOTQuNjY0MDYyYy05OS45ODgyODEgMC0xODEuMzMyMDMxIDgxLjM0Mzc1LTE4MS4zMzIwMzEgMTgxLjMzMjAzMXM4MS4zNDM3NSAxODEuMzMyMDMxIDE4MS4zMzIwMzEgMTgxLjMzMjAzMSAxODEuMzMyMDMxLTgxLjM0Mzc1IDE4MS4zMzIwMzEtMTgxLjMzMjAzMS04MS4zNDM3NS0xODEuMzMyMDMxLTE4MS4zMzIwMzEtMTgxLjMzMjAzMXptMCAwIiBmaWxsPSIjZmZmZmZmIiBkYXRhLW9yaWdpbmFsPSIjMDAwMDAwIiBzdHlsZT0iIiBjbGFzcz0iIj48L3BhdGg+PHBhdGggeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiBkPSJtMjU2IDEwNi42Njc5NjljLTguODMyMDMxIDAtMTYtNy4xNjc5NjktMTYtMTZ2LTc0LjY2Nzk2OWMwLTguODMyMDMxIDcuMTY3OTY5LTE2IDE2LTE2czE2IDcuMTY3OTY5IDE2IDE2djc0LjY2Nzk2OWMwIDguODMyMDMxLTcuMTY3OTY5IDE2LTE2IDE2em0wIDAiIGZpbGw9IiNmZmZmZmYiIGRhdGEtb3JpZ2luYWw9IiMwMDAwMDAiIHN0eWxlPSIiIGNsYXNzPSIiPjwvcGF0aD48cGF0aCB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIGQ9Im0yNTYgNTEyYy04LjgzMjAzMSAwLTE2LTcuMTY3OTY5LTE2LTE2di03NC42Njc5NjljMC04LjgzMjAzMSA3LjE2Nzk2OS0xNiAxNi0xNnMxNiA3LjE2Nzk2OSAxNiAxNnY3NC42Njc5NjljMCA4LjgzMjAzMS03LjE2Nzk2OSAxNi0xNiAxNnptMCAwIiBmaWxsPSIjZmZmZmZmIiBkYXRhLW9yaWdpbmFsPSIjMDAwMDAwIiBzdHlsZT0iIiBjbGFzcz0iIj48L3BhdGg+PHBhdGggeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiBkPSJtOTAuNjY3OTY5IDI3MmgtNzQuNjY3OTY5Yy04LjgzMjAzMSAwLTE2LTcuMTY3OTY5LTE2LTE2czcuMTY3OTY5LTE2IDE2LTE2aDc0LjY2Nzk2OWM4LjgzMjAzMSAwIDE2IDcuMTY3OTY5IDE2IDE2cy03LjE2Nzk2OSAxNi0xNiAxNnptMCAwIiBmaWxsPSIjZmZmZmZmIiBkYXRhLW9yaWdpbmFsPSIjMDAwMDAwIiBzdHlsZT0iIiBjbGFzcz0iIj48L3BhdGg+PHBhdGggeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIiBkPSJtNDk2IDI3MmgtNzQuNjY3OTY5Yy04LjgzMjAzMSAwLTE2LTcuMTY3OTY5LTE2LTE2czcuMTY3OTY5LTE2IDE2LTE2aDc0LjY2Nzk2OWM4LjgzMjAzMSAwIDE2IDcuMTY3OTY5IDE2IDE2cy03LjE2Nzk2OSAxNi0xNiAxNnptMCAwIiBmaWxsPSIjZmZmZmZmIiBkYXRhLW9yaWdpbmFsPSIjMDAwMDAwIiBzdHlsZT0iIiBjbGFzcz0iIj48L3BhdGg+PC9nPjwvc3ZnPg==" />
                    <GroupSelect port={port}></GroupSelect>
                    <ChatWindow port={port} ref={msgsRef}></ChatWindow>
                    <MapPopup selectedUnit={units[selectedUnit]} />
                </>
            }
            <div>
                <div ref={mapContainer} className="map-container" />
            </div>

        </>
    );
});
export default Home;
