import React, { useRef, useEffect, useState } from 'react';
import mapboxgl from '!mapbox-gl'; // eslint-disable-line import/no-webpack-loader-syntax
import { getMsgs, getMsgData } from './Api'
import './App.css';

mapboxgl.accessToken = 'pk.eyJ1IjoiYWhtZWRhZmlmaSIsImEiOiJja3F6YzJibGUwNXEyMnNsZ2U2N2lod2xqIn0.U2YYTWHCYqkCUBaAFd9MfA';

function App() {
  const MapContainer = useRef(null);
  const map = useRef(null);
  const [markers, setMarkers] = useState({})
  const [msgs, setMsgs] = useState([])
  const [selectedMsg, setSelectedMsg] = useState(null)


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

  var addMarkers = (locations) => {
    setMarkers(markers => {
      var coordinates = []
      Object.keys(locations).forEach(ip => {
        if (markers[ip]) markers[ip].remove()

        markers[ip] = new mapboxgl.Marker({ color: 'black' })
          .setLngLat([locations[ip].lon, locations[ip].lat])
          .addTo(map.current);

        coordinates.push([locations[ip].lon, locations[ip].lat])
      })
      return markers;
    })
  }

  var drawPath = (pathIps) => {
    setMarkers(markers => {
      var coordinates = []
      pathIps.forEach((ip, index) => {
        var color = 'green'
        var lng = markers[ip].getLngLat().lng
        var lat = markers[ip].getLngLat().lat
        coordinates.push([lng, lat])

        if (index === 0) color = 'red' 
        if (index === pathIps.length - 1) color = 'blue'
        if (markers[ip]) markers[ip].remove()

        markers[ip] = new mapboxgl.Marker({ color: color })
          .setLngLat([lng, lat])
          .addTo(map.current);
      })

      if(typeof(map.current.getSource('route')) !== 'undefined') {
        map.current.removeLayer('route')
        map.current.removeSource('route')
      }

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

      map.current.setZoom(15)
      map.current.setCenter([markers[pathIps[0]].getLngLat().lng, markers[pathIps[0]].getLngLat().lat])
      return markers;
    })
  }

  useEffect(() => {
    map.current = new mapboxgl.Map({
      container: MapContainer.current,
      style: 'mapbox://styles/ahmedafifi/ckr3eqazg5ndn18p3pgmuffc1',
      center: [-112, 70],
      zoom: 7
    });

    map.current.addControl(new mapboxgl.FullscreenControl());
    map.current.addControl(new mapboxgl.NavigationControl());
    
    getMsgs().then(msgs => {
      setMsgs(msgs)
    })
  }, [])

  useEffect(() => {
    setSelectedMsg(selectedMsg => {
      if(!selectedMsg) return selectedMsg
      getMsgData(selectedMsg).then(msgData => {
        addMarkers(msgData["locations"])
        drawPath(msgData["path"])
      })
    })
  }, [selectedMsg])

  return (
    <>
      {!msgs || !msgs.length ? <> </> :
        <div className="msgs-container">
          {
            msgs.map((msg, index) => {
              return <div className="msg" onClick={() => setSelectedMsg(msg.hash)}> 
                      <img alt="msg" src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZlcnNpb249IjEuMSIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHhtbG5zOnN2Z2pzPSJodHRwOi8vc3ZnanMuY29tL3N2Z2pzIiB3aWR0aD0iNTEyIiBoZWlnaHQ9IjUxMiIgeD0iMCIgeT0iMCIgdmlld0JveD0iMCAwIDUxMiA1MTIiIHN0eWxlPSJlbmFibGUtYmFja2dyb3VuZDpuZXcgMCAwIDUxMiA1MTIiIHhtbDpzcGFjZT0icHJlc2VydmUiIGNsYXNzPSIiPjxnPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgoJPGc+CgkJPHBhdGggZD0iTTQ2Nyw2MUg0NUMyMC4yMTgsNjEsMCw4MS4xOTYsMCwxMDZ2MzAwYzAsMjQuNzIsMjAuMTI4LDQ1LDQ1LDQ1aDQyMmMyNC43MiwwLDQ1LTIwLjEyOCw0NS00NVYxMDYgICAgQzUxMiw4MS4yOCw0OTEuODcyLDYxLDQ2Nyw2MXogTTQ2MC43ODYsOTFMMjU2Ljk1NCwyOTQuODMzTDUxLjM1OSw5MUg0NjAuNzg2eiBNMzAsMzk5Ljc4OFYxMTIuMDY5bDE0NC40NzksMTQzLjI0TDMwLDM5OS43ODh6ICAgICBNNTEuMjEzLDQyMWwxNDQuNTctMTQ0LjU3bDUwLjY1Nyw1MC4yMjJjNS44NjQsNS44MTQsMTUuMzI3LDUuNzk1LDIxLjE2Ny0wLjA0NkwzMTcsMjc3LjIxM0w0NjAuNzg3LDQyMUg1MS4yMTN6IE00ODIsMzk5Ljc4NyAgICBMMzM4LjIxMywyNTZMNDgyLDExMi4yMTJWMzk5Ljc4N3oiIGZpbGw9IiNmZmZmZmYiIGRhdGEtb3JpZ2luYWw9IiMwMDAwMDAiIHN0eWxlPSIiIGNsYXNzPSIiPjwvcGF0aD4KCTwvZz4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8L2c+PC9zdmc+" />               
                       <span> {"Message " + index} </span> 
                    </div>
            })
          }
        </div>
      }
      <div>
        <div ref={MapContainer} className="map-container" />
      </div>
    </>
  );
}

export default App;
