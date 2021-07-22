import React from 'react';
import ReactHlsPlayer from 'react-hls-player';
import './Stream.css'

const M3U8Name   = "index.m3u8";

const HLSConfig = {
    maxLoadingDelay: 4,
    minAutoBitrate: 0,
    lowLatencyMode: true,
    liveMaxLatencyDurationCount: 10
}

function Stream({ onEndStream, stream, port }) {
    return (
        <>
            {stream ?
                <div className="stream-wrapper">
                    <img alt="stop-stream" onClick={() => {onEndStream(stream)}} className="stop-stream" src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZlcnNpb249IjEuMSIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHhtbG5zOnN2Z2pzPSJodHRwOi8vc3ZnanMuY29tL3N2Z2pzIiB3aWR0aD0iNTEyIiBoZWlnaHQ9IjUxMiIgeD0iMCIgeT0iMCIgdmlld0JveD0iMCAwIDMwLjA1IDMwLjA1IiBzdHlsZT0iZW5hYmxlLWJhY2tncm91bmQ6bmV3IDAgMCA1MTIgNTEyIiB4bWw6c3BhY2U9InByZXNlcnZlIiBjbGFzcz0iIj48Zz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KCTxwYXRoIHN0eWxlPSIiIGQ9Ik0xOC45OTMsMTAuNjg4aC03LjkzNmMtMC4xOSwwLTAuMzQ2LDAuMTQ5LTAuMzQ2LDAuMzQydjguMDIyYzAsMC4xODksMC4xNTUsMC4zNDQsMC4zNDYsMC4zNDQgICBoNy45MzZjMC4xOSwwLDAuMzQ0LTAuMTU0LDAuMzQ0LTAuMzQ0VjExLjAzQzE5LjMzNiwxMC44MzgsMTkuMTgzLDEwLjY4OCwxOC45OTMsMTAuNjg4eiIgZmlsbD0iIzE5OWU5YSIgZGF0YS1vcmlnaW5hbD0iIzAzMDEwNCIgY2xhc3M9IiI+PC9wYXRoPgoJPHBhdGggc3R5bGU9IiIgZD0iTTE1LjAyNiwwQzYuNzI5LDAsMC4wMDEsNi43MjYsMC4wMDEsMTUuMDI1UzYuNzI5LDMwLjA1LDE1LjAyNiwzMC4wNSAgIGM4LjI5OCwwLDE1LjAyMy02LjcyNiwxNS4wMjMtMTUuMDI1UzIzLjMyNCwwLDE1LjAyNiwweiBNMTUuMDI2LDI3LjU0Yy02LjkxMiwwLTEyLjUxNi01LjYwNC0xMi41MTYtMTIuNTE1ICAgYzAtNi45MTQsNS42MDQtMTIuNTE3LDEyLjUxNi0xMi41MTdjNi45MTMsMCwxMi41MTQsNS42MDMsMTIuNTE0LDEyLjUxN0MyNy41NCwyMS45MzYsMjEuOTM5LDI3LjU0LDE1LjAyNiwyNy41NHoiIGZpbGw9IiMxOTllOWEiIGRhdGEtb3JpZ2luYWw9IiMwMzAxMDQiIGNsYXNzPSIiPjwvcGF0aD4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8L2c+PC9zdmc+" />
                    <ReactHlsPlayer
                        src={"http://localhost:" + port + "/api/stream/" + stream.src + "/" + stream.id + "/" + M3U8Name}
                        autoPlay={true}
                        controls={true}
                        hlsConfig={HLSConfig}  
                    />
                </div> : <> </>
            }
        </>
    );
} export default Stream;
