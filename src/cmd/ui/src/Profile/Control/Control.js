import './Control.css';
import React, {useRef} from 'react';
import RecordAudio from '../../RecordAudio/RecordAudio'


function Control({ unit }) {
    const recordAudioRef = useRef(null);
    var onSendAudio = (audio) => {
        console.log("audio", audio);
    }
    return (
        <div className="control-container">
            <RecordAudio onSend={onSendAudio} ref={recordAudioRef}></RecordAudio>
            <div data-augmented-ui="tr-clip br-clip bl-clip-y border" class="control-item"></div>
            <img className="attack-icon" src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZlcnNpb249IjEuMSIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHhtbG5zOnN2Z2pzPSJodHRwOi8vc3ZnanMuY29tL3N2Z2pzIiB3aWR0aD0iNTEyIiBoZWlnaHQ9IjUxMiIgeD0iMCIgeT0iMCIgdmlld0JveD0iMCAwIDUxMS45NzggNTExLjk3OCIgc3R5bGU9ImVuYWJsZS1iYWNrZ3JvdW5kOm5ldyAwIDAgNTEyIDUxMiIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSIgY2xhc3M9IiI+PGc+PGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cGF0aCBkPSJtNzQuNzk1IDI4Ny4xODZjMTIuMDEyIDEyLjAyNiAxNi45MzQgMjkuNzUxIDEyLjgzMiA0Ni4yNi02LjMxMyAyNS40IDEuMjYgNTIuNjc2IDE5Ljc0NiA3MS4xNjJzNDUuNzMyIDI2LjAwMSA3MS4xNDcgMTkuNzQ2YzE2LjU4Mi00LjEzMSAzNC4yNjMuODIgNDYuMjc0IDEyLjgzMmwyMS4yMTEtMjEuMjExLTE1MC0xNTB6IiBmaWxsPSIjMTk5ZTlhIiBkYXRhLW9yaWdpbmFsPSIjMDAwMDAwIiBzdHlsZT0iIiBjbGFzcz0iIj48L3BhdGg+PHBhdGggZD0ibTc2LjgwOSA0OTguODA1IDQ4LjcxNy00OC43MjZjLTE0LjY4NS01LjE5MS0yOC4zMS0xMy4yMDYtMzkuMzY0LTI0LjI2LTExLjE2Ni0xMS4xNi0xOS4yOS0yNC42NzItMjQuNDM1LTM5LjIwN2wtNDguNTY1IDQ4LjU3NGMtMTcuNTQ5IDE3LjU0OS0xNy41NDkgNDYuMDg0IDAgNjMuNjMzIDE3LjU0OCAxNy41NSA0Ni4wOTYgMTcuNTM3IDYzLjY0Ny0uMDE0eiIgZmlsbD0iIzE5OWU5YSIgZGF0YS1vcmlnaW5hbD0iIzAwMDAwMCIgc3R5bGU9IiIgY2xhc3M9IiI+PC9wYXRoPjxwYXRoIGQ9Im0yOTYuMjg2IDMxLjAwNWg4NC44NTN2MjkuOTk3aC04NC44NTN6IiB0cmFuc2Zvcm09Im1hdHJpeCguNzA3IC0uNzA3IC43MDcgLjcwNyA2Ni42NzcgMjUyLjk4KSIgZmlsbD0iIzE5OWU5YSIgZGF0YS1vcmlnaW5hbD0iIzAwMDAwMCIgc3R5bGU9IiIgY2xhc3M9IiI+PC9wYXRoPjxwYXRoIGQ9Im0xNTguMjQ1IDMuNTc3aDI5Ljk5N3Y4NC44NTNoLTI5Ljk5N3oiIHRyYW5zZm9ybT0ibWF0cml4KC43MDcgLS43MDcgLjcwNyAuNzA3IDE4LjIxMyAxMzUuOTc2KSIgZmlsbD0iIzE5OWU5YSIgZGF0YS1vcmlnaW5hbD0iIzAwMDAwMCIgc3R5bGU9IiIgY2xhc3M9IiI+PC9wYXRoPjxwYXRoIGQ9Im0yNDAuOTc4LjAwM2gzMHY5MWgtMzB6IiBmaWxsPSIjMTk5ZTlhIiBkYXRhLW9yaWdpbmFsPSIjMDAwMDAwIiBzdHlsZT0iIiBjbGFzcz0iIj48L3BhdGg+PHBhdGggZD0ibTI4Ny4xODMgNDM3LjE4NmMxMi4wMTItMTIuMDEyIDI5LjY5Mi0xNi45NjMgNDYuMjc0LTEyLjgzMiAyNS40MTUgNi4yNTUgNTIuNjYxLTEuMjYgNzEuMTQ3LTE5Ljc0NnMyNi4wNi00NS43NjIgMTkuNzQ2LTcxLjE2MmMtNC4xMDItMTYuNTA5LjgyLTM0LjIzMyAxMi44MzItNDYuMjZsLTIxLjIxMS0yMS4yMTEtMTUwIDE1MHoiIGZpbGw9IiMxOTllOWEiIGRhdGEtb3JpZ2luYWw9IiMwMDAwMDAiIHN0eWxlPSIiIGNsYXNzPSIiPjwvcGF0aD48cGF0aCBkPSJtNDUwLjI1MSAzODYuNjEyYy01LjE0NSAxNC41MzUtMTMuMjcgMjguMDQ2LTI0LjQzNSAzOS4yMDctMTEuMDU0IDExLjA1NC0yNC42NzkgMTkuMDY5LTM5LjM2NCAyNC4yNmw0OC43MTcgNDguNzI2YzE3LjU1MSAxNy41NTEgNDYuMDk5IDE3LjU2MyA2My42NDcuMDE1IDE3LjU0OS0xNy41NDkgMTcuNTQ5LTQ2LjA4NCAwLTYzLjYzM3oiIGZpbGw9IiMxOTllOWEiIGRhdGEtb3JpZ2luYWw9IiMwMDAwMDAiIHN0eWxlPSIiIGNsYXNzPSIiPjwvcGF0aD48cGF0aCBkPSJtMTkyLjM1NiAyMzQuNzgxIDQyLjQyMi00Mi40MjItMTYyLjM0OS0xNjIuMzU2Yy0xOS4wNzItMTkuMDcyLTQ1LjQ2OS0zMC03Mi40MjItMzAgMCAyNi45NTMgMTAuOTI4IDUzLjM1IDMwIDcyLjQyMnoiIGZpbGw9IiMxOTllOWEiIGRhdGEtb3JpZ2luYWw9IiMwMDAwMDAiIHN0eWxlPSIiIGNsYXNzPSIiPjwvcGF0aD48cGF0aCBkPSJtMjc5LjA4NyAyOTMuOTk0aDU5Ljk5NHYzMC4xOTRoLTU5Ljk5NHoiIHRyYW5zZm9ybT0ibWF0cml4KC43MDcgLS43MDcgLjcwNyAuNzA3IC0xMjguMDMyIDMwOS4wODYpIiBmaWxsPSIjMTk5ZTlhIiBkYXRhLW9yaWdpbmFsPSIjMDAwMDAwIiBzdHlsZT0iIiBjbGFzcz0iIj48L3BhdGg+PHBhdGggZD0ibTE3MS4wMDIgMjk4LjU1NyA0Mi40MjIgNDIuNDIyIDI2OC41NTQtMjY4LjU1NGMxOS4wNzItMTkuMDcyIDMwLTQ1LjQ2OSAzMC03Mi40MjItMjYuOTUzIDAtNTMuMzUgMTAuOTI4LTcyLjQyMiAzMHoiIGZpbGw9IiMxOTllOWEiIGRhdGEtb3JpZ2luYWw9IiMwMDAwMDAiIHN0eWxlPSIiIGNsYXNzPSIiPjwvcGF0aD48L2c+PC9nPjwvc3ZnPg==" />
            <img className="defense-icon" src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZlcnNpb249IjEuMSIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHhtbG5zOnN2Z2pzPSJodHRwOi8vc3ZnanMuY29tL3N2Z2pzIiB3aWR0aD0iNTEyIiBoZWlnaHQ9IjUxMiIgeD0iMCIgeT0iMCIgdmlld0JveD0iMCAwIDQzMS45NDIgNDMxLjk0MiIgc3R5bGU9ImVuYWJsZS1iYWNrZ3JvdW5kOm5ldyAwIDAgNTEyIDUxMiIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSIgY2xhc3M9IiI+PGc+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+Cgk8Zz4KCQk8cGF0aCBkPSJNMzg1LjczMSwyNy41MWMwLjAzNS00LjQxOC0zLjUxOC04LjAyOC03LjkzNi04LjA2M2MtMS4xNTItMC4wMDktMi4yOTMsMC4yMzEtMy4zNDQsMC43MDMgICAgYy01MC4zODUsMjMuNDUtMTA5LjY5NCwxNi40MTktMTUzLjItMTguMTZjLTMuMDItMi42NTMtNy41NC0yLjY1My0xMC41NiwwYy00My42ODksMzQuMTM3LTEwMi43NDEsNDEuMTM3LTE1My4yLDE4LjE2ICAgIGMtNC4wMy0xLjgxMi04Ljc2NS0wLjAxMy0xMC41NzcsNC4wMTZjLTAuNDcyLDEuMDUxLTAuNzEyLDIuMTkxLTAuNzAzLDMuMzQ0YzAsMS42OCwwLjU2LDE2Ni42NCwwLDIyMS4yOCAgICBjLTAuMzIsMjkuMiw2Ny44NCwxMjQsMTY1LjY4LDE4Mi4wOGMyLjQ3NSwxLjQyOSw1LjUyNSwxLjQyOSw4LDBjOTgtNTguMDgsMTY2LjE2LTE1Mi44OCwxNjUuODQtMTgyLjA4ICAgIEMzODUuMTcxLDE5NC4xNSwzODUuNzMxLDI5LjE5LDM4NS43MzEsMjcuNTF6IE0zNjkuNzMxLDI0OC44N2MwLDE3LjYtNTYsMTA2LTE1My43NiwxNjUuNzYgICAgYy05OC4wOC01OS42OC0xNTMuOTItMTQ4LjA4LTE1My43Ni0xNjUuNjhjMC40OC00Ni4xNiwwLjE2LTE3MC40LDAtMjA5LjZjNTEuNjA3LDE4Ljg4NCwxMDkuMjIyLDEwLjkxLDE1My43Ni0yMS4yOCAgICBjNDQuNTM4LDMyLjE5LDEwMi4xNTMsNDAuMTY0LDE1My43NiwyMS4yOEMzNjkuNzMxLDc4LjQ3LDM2OS4zMzEsMjAyLjc5LDM2OS43MzEsMjQ4Ljg3eiIgZmlsbD0iIzE5OWU5YSIgZGF0YS1vcmlnaW5hbD0iIzAwMDAwMCIgc3R5bGU9IiI+PC9wYXRoPgoJPC9nPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+Cgk8Zz4KCQk8cGF0aCBkPSJNMzQzLjk3MSw3Ni4xNWMwLjAzNS00LjQxOC0zLjUxOC04LjAyOC03LjkzNi04LjA2M2MtMS4xNTItMC4wMDktMi4yOTMsMC4yMzEtMy4zNDQsMC43MDMgICAgYy0zNi41NzEsMTcuMTkzLTc5LjcxOSwxMi4yMzctMTExLjQ0LTEyLjhjLTMuMDItMi42NTMtNy41NC0yLjY1My0xMC41NiwwYy0zMS43MjgsMjQuNzQ3LTc0LjU4MywyOS44MDMtMTExLjIsMTMuMTIgICAgYy00LjAzLTEuODEyLTguNzY1LTAuMDEzLTEwLjU3Nyw0LjAxNmMtMC40NzIsMS4wNTEtMC43MTIsMi4xOTItMC43MDMsMy4zNDRjMCwxLjIsMC40LDEyMy4zNiwwLDE2My44NCAgICBjLTAuMjQsMjIsNTAuNjQsOTMuMiwxMjMuNjgsMTM2LjU2aDBjMi40NzUsMS40MjksNS41MjUsMS40MjksOCwwYzczLjItNDMuMzYsMTI0LjA4LTExNC41NiwxMjQuMDgtMTM2Ljg4ICAgIEMzNDMuNTcxLDE5OS41MSwzNDMuOTcxLDc3LjQzLDM0My45NzEsNzYuMTV6IE0zMjcuOTcxLDIzOS45OWMwLDExLjI4LTQwLDc2LjMyLTExMiwxMjBjLTcxLjQ0LTQzLjkyLTExMi0xMDguOTYtMTEyLTEyMCAgICBjMC41Ni0zMi42NCwwLTEyMCwwLTE1MmMzNy43MDgsMTMuMjYyLDc5LjUxNCw3LjI5LDExMi0xNmMzMi40ODYsMjMuMjksNzQuMjkyLDI5LjI2MiwxMTIsMTZWMjM5Ljk5eiIgZmlsbD0iIzE5OWU5YSIgZGF0YS1vcmlnaW5hbD0iIzAwMDAwMCIgc3R5bGU9IiI+PC9wYXRoPgoJPC9nPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjwvZz48L3N2Zz4=" />
            <img onClick = {() => {recordAudioRef.current.openModal()}} className="audio-icon" src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZlcnNpb249IjEuMSIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHhtbG5zOnN2Z2pzPSJodHRwOi8vc3ZnanMuY29tL3N2Z2pzIiB3aWR0aD0iNTEyIiBoZWlnaHQ9IjUxMiIgeD0iMCIgeT0iMCIgdmlld0JveD0iMCAwIDM4OS4xMiAzODkuMTIiIHN0eWxlPSJlbmFibGUtYmFja2dyb3VuZDpuZXcgMCAwIDUxMiA1MTIiIHhtbDpzcGFjZT0icHJlc2VydmUiPjxnPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgoJPGc+CgkJPHBhdGggZD0iTTIxNS4wNzIsMGgtNDEuMDI3Yy0zMC4wOTIsMC01NC41NzksMjQuNTE2LTU0LjU3OSw1NC42NTN2MTM2LjQ5NGMwLDMuNzcsMy4wNTMsNi44MjcsNi44MjcsNi44MjdoMTM2LjUzMyAgICBjMy43NzMsMCw2LjgyNy0zLjA1Nyw2LjgyNy02LjgyN1Y1NC42NTNDMjY5LjY1MywyNC41MTYsMjQ1LjE2NiwwLDIxNS4wNzIsMHogTTI1NiwxODQuMzJIMTMzLjEyVjU0LjY1MyAgICBjMC0yMi42MDcsMTguMzYtNDAuOTk5LDQwLjkyNi00MC45OTloNDEuMDI4QzIzNy42NCwxMy42NTMsMjU2LDMyLjA0NiwyNTYsNTQuNjUzVjE4NC4zMnoiIGZpbGw9IiMxOTllOWEiIGRhdGEtb3JpZ2luYWw9IiMwMDAwMDAiIHN0eWxlPSIiPjwvcGF0aD4KCTwvZz4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgoJPGc+CgkJPHBhdGggZD0iTTExOS40NjcsMTg0LjMydjQ0LjM3M2MwLDI0LjQ2NywxOC4zNTksNDQuMzczLDQwLjkyNiw0NC4zNzNoNjguMzM1YzIyLjU2NiwwLDQwLjkyNi0xOS45MDcsNDAuOTI2LTQ0LjM3M1YxODQuMzIgICAgSDExOS40Njd6IE0yNTUuOTk4LDIyOC42OTNjMCwxNi45NC0xMi4yMzMsMzAuNzItMjcuMjczLDMwLjcyaC02OC4zMzNjLTE1LjAzOSwwLTI3LjI3My0xMy43OC0yNy4yNzMtMzAuNzJ2LTMwLjcyaDEyMi44NzggICAgVjIyOC42OTN6IiBmaWxsPSIjMTk5ZTlhIiBkYXRhLW9yaWdpbmFsPSIjMDAwMDAwIiBzdHlsZT0iIj48L3BhdGg+Cgk8L2c+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KCTxnPgoJCTxwYXRoIGQ9Ik0xNjcuMjUzLDYxLjQ0SDEzMy4xMmMtMy43NzMsMC02LjgyNywzLjA1Ny02LjgyNyw2LjgyN3MzLjA1Myw2LjgyNyw2LjgyNyw2LjgyN2gzNC4xMzNjMy43NzMsMCw2LjgyNy0zLjA1Nyw2LjgyNy02LjgyNyAgICBTMTcxLjAyNyw2MS40NCwxNjcuMjUzLDYxLjQ0eiIgZmlsbD0iIzE5OWU5YSIgZGF0YS1vcmlnaW5hbD0iIzAwMDAwMCIgc3R5bGU9IiI+PC9wYXRoPgoJPC9nPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+Cgk8Zz4KCQk8cGF0aCBkPSJNMTY3LjI1Myw4OC43NDdIMTMzLjEyYy0zLjc3MywwLTYuODI3LDMuMDU3LTYuODI3LDYuODI3YzAsMy43NywzLjA1Myw2LjgyNyw2LjgyNyw2LjgyN2gzNC4xMzMgICAgYzMuNzczLDAsNi44MjctMy4wNTcsNi44MjctNi44MjdDMTc0LjA4LDkxLjgwMywxNzEuMDI3LDg4Ljc0NywxNjcuMjUzLDg4Ljc0N3oiIGZpbGw9IiMxOTllOWEiIGRhdGEtb3JpZ2luYWw9IiMwMDAwMDAiIHN0eWxlPSIiPjwvcGF0aD4KCTwvZz4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgoJPGc+CgkJPHBhdGggZD0iTTI1Niw2MS40NGgtMzQuMTMzYy0zLjc3MywwLTYuODI3LDMuMDU3LTYuODI3LDYuODI3czMuMDUzLDYuODI3LDYuODI3LDYuODI3SDI1NmMzLjc3MywwLDYuODI3LTMuMDU3LDYuODI3LTYuODI3ICAgIFMyNTkuNzczLDYxLjQ0LDI1Niw2MS40NHoiIGZpbGw9IiMxOTllOWEiIGRhdGEtb3JpZ2luYWw9IiMwMDAwMDAiIHN0eWxlPSIiPjwvcGF0aD4KCTwvZz4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgoJPGc+CgkJPHBhdGggZD0iTTI1Niw4OC43NDdoLTM0LjEzM2MtMy43NzMsMC02LjgyNywzLjA1Ny02LjgyNyw2LjgyN2MwLDMuNzcsMy4wNTMsNi44MjcsNi44MjcsNi44MjdIMjU2ICAgIGMzLjc3MywwLDYuODI3LTMuMDU3LDYuODI3LTYuODI3QzI2Mi44MjcsOTEuODAzLDI1OS43NzMsODguNzQ3LDI1Niw4OC43NDd6IiBmaWxsPSIjMTk5ZTlhIiBkYXRhLW9yaWdpbmFsPSIjMDAwMDAwIiBzdHlsZT0iIj48L3BhdGg+Cgk8L2c+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KCTxnPgoJCTxwYXRoIGQ9Ik0xNjcuMjUzLDExNi4wNTNIMTMzLjEyYy0zLjc3MywwLTYuODI3LDMuMDU3LTYuODI3LDYuODI3czMuMDUzLDYuODI3LDYuODI3LDYuODI3aDM0LjEzMyAgICBjMy43NzMsMCw2LjgyNy0zLjA1Nyw2LjgyNy02LjgyN1MxNzEuMDI3LDExNi4wNTMsMTY3LjI1MywxMTYuMDUzeiIgZmlsbD0iIzE5OWU5YSIgZGF0YS1vcmlnaW5hbD0iIzAwMDAwMCIgc3R5bGU9IiI+PC9wYXRoPgoJPC9nPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+Cgk8Zz4KCQk8cGF0aCBkPSJNMjU2LDExNi4wNTNoLTM0LjEzM2MtMy43NzMsMC02LjgyNywzLjA1Ny02LjgyNyw2LjgyN3MzLjA1Myw2LjgyNyw2LjgyNyw2LjgyN0gyNTZjMy43NzMsMCw2LjgyNy0zLjA1Nyw2LjgyNy02LjgyNyAgICBTMjU5Ljc3MywxMTYuMDUzLDI1NiwxMTYuMDUzeiIgZmlsbD0iIzE5OWU5YSIgZGF0YS1vcmlnaW5hbD0iIzAwMDAwMCIgc3R5bGU9IiI+PC9wYXRoPgoJPC9nPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+Cgk8Zz4KCQk8cGF0aCBkPSJNMTY3LjI1MywxNDMuMzZIMTMzLjEyYy0zLjc3MywwLTYuODI3LDMuMDU3LTYuODI3LDYuODI3czMuMDUzLDYuODI3LDYuODI3LDYuODI3aDM0LjEzMyAgICBjMy43NzMsMCw2LjgyNy0zLjA1Nyw2LjgyNy02LjgyN1MxNzEuMDI3LDE0My4zNiwxNjcuMjUzLDE0My4zNnoiIGZpbGw9IiMxOTllOWEiIGRhdGEtb3JpZ2luYWw9IiMwMDAwMDAiIHN0eWxlPSIiPjwvcGF0aD4KCTwvZz4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgoJPGc+CgkJPHBhdGggZD0iTTI1NiwxNDMuMzZoLTM0LjEzM2MtMy43NzMsMC02LjgyNywzLjA1Ny02LjgyNyw2LjgyN3MzLjA1Myw2LjgyNyw2LjgyNyw2LjgyN0gyNTZjMy43NzMsMCw2LjgyNy0zLjA1Nyw2LjgyNy02LjgyNyAgICBTMjU5Ljc3MywxNDMuMzYsMjU2LDE0My4zNnoiIGZpbGw9IiMxOTllOWEiIGRhdGEtb3JpZ2luYWw9IiMwMDAwMDAiIHN0eWxlPSIiPjwvcGF0aD4KCTwvZz4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgoJPGc+CgkJPHBhdGggZD0iTTI5Ni45NiwxMzguNzE4Yy0zLjc3MywwLTYuODI3LDMuMDU3LTYuODI3LDYuODI3djczLjk5OWMwLDM3LjA0LTMwLjYyNiw2Ny4xNzYtNjguMjczLDY3LjE3NmgtNTQuNiAgICBjLTM3LjY0NywwLTY4LjI3My0zMC4xMzYtNjguMjczLTY3LjE3NnYtNzMuOTk5YzAtMy43Ny0zLjA1My02LjgyNy02LjgyNy02LjgyN3MtNi44MjcsMy4wNTctNi44MjcsNi44Mjd2NzMuOTk5ICAgIGMwLDQ0LjU3LDM2Ljc1Myw4MC44MjksODEuOTI3LDgwLjgyOWg1NC42YzQ1LjE3NCwwLDgxLjkyNy0zNi4yNiw4MS45MjctODAuODI5di03My45OTkgICAgQzMwMy43ODcsMTQxLjc3NSwzMDAuNzMzLDEzOC43MTgsMjk2Ljk2LDEzOC43MTh6IiBmaWxsPSIjMTk5ZTlhIiBkYXRhLW9yaWdpbmFsPSIjMDAwMDAwIiBzdHlsZT0iIj48L3BhdGg+Cgk8L2c+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KCTxnPgoJCTxwYXRoIGQ9Ik0xOTQuNTYsMjk1LjcyOWMtMy43NzMsMC02LjgyNywzLjA1Ny02LjgyNyw2LjgyN3Y3NS4wOTNjMCwzLjc3LDMuMDUzLDYuODI3LDYuODI3LDYuODI3czYuODI3LTMuMDU3LDYuODI3LTYuODI3ICAgIHYtNzUuMDkzQzIwMS4zODcsMjk4Ljc4NiwxOTguMzMzLDI5NS43MjksMTk0LjU2LDI5NS43Mjl6IiBmaWxsPSIjMTk5ZTlhIiBkYXRhLW9yaWdpbmFsPSIjMDAwMDAwIiBzdHlsZT0iIj48L3BhdGg+Cgk8L2c+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KCTxnPgoJCTxwYXRoIGQ9Ik0yNDkuMTczLDM3NS40NjdIMTM5Ljk0N2MtMy43NzMsMC02LjgyNywzLjA1Ny02LjgyNyw2LjgyN2MwLDMuNzcsMy4wNTMsNi44MjcsNi44MjcsNi44MjdoMTA5LjIyNyAgICBjMy43NzMsMCw2LjgyNy0zLjA1Nyw2LjgyNy02LjgyN0MyNTYsMzc4LjUyMywyNTIuOTQ3LDM3NS40NjcsMjQ5LjE3MywzNzUuNDY3eiIgZmlsbD0iIzE5OWU5YSIgZGF0YS1vcmlnaW5hbD0iIzAwMDAwMCIgc3R5bGU9IiI+PC9wYXRoPgoJPC9nPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjwvZz48L3N2Zz4=" />            <img className="video-icon" src="data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHZlcnNpb249IjEuMSIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHhtbG5zOnN2Z2pzPSJodHRwOi8vc3ZnanMuY29tL3N2Z2pzIiB3aWR0aD0iNTEyIiBoZWlnaHQ9IjUxMiIgeD0iMCIgeT0iMCIgdmlld0JveD0iMCAwIDUxMi4wMDIgNTEyLjAwMiIgc3R5bGU9ImVuYWJsZS1iYWNrZ3JvdW5kOm5ldyAwIDAgNTEyIDUxMiIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSIgY2xhc3M9IiI+PGc+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+Cgk8Zz4KCQk8cGF0aCBkPSJNNDYyLjAwMiw5Mi4wMDJoLTQyLjAwMVY1MGMwLTI3LjU3LTIyLjQzLTUwLTUwLTUwaC0zMjBjLTI3LjU3LDAtNTAsMjIuNDMtNTAsNTB2MzIwYzAsMjcuNTcsMjIuNDMsNTAsNTAsNTBoNDIuMDAxICAgIHY0Mi4wMDJjMCwyNy41NywyMi40Myw1MCw1MCw1MGgzMjBjMjcuNTcsMCw1MC0yMi40Myw1MC01MFYxNDIuMDAxQzUxMi4wMDMsMTE0LjQzMSw0ODkuNTczLDkyLjAwMiw0NjIuMDAyLDkyLjAwMnogTTUwLjAwMSw0MDAgICAgYy0xNi41NDIsMC0zMC0xMy40NTctMzAtMzBWNTBjMC0xNi41NDIsMTMuNDU4LTMwLDMwLTMwaDMyMGMxNi41NDIsMCwzMCwxMy40NTgsMzAsMzB2MzIwYzAsMTYuNTQyLTEzLjQ1OCwzMC0zMCwzMEg1MC4wMDF6ICAgICBNNDkyLjAwMiw0NjIuMDAyYzAsMTYuNTQyLTEzLjQ1OCwzMC0zMCwzMGgtMzIwYy0xNi41NDIsMC0zMC0xMy40NTgtMzAtMzBWNDIwaDI1Ny45OTljMjcuNTcsMCw1MC0yMi40Myw1MC01MFYxMTIuMDAyaDQyLjAwMSAgICBjMTYuNTQyLDAsMzAsMTMuNDU4LDMwLDMwVjQ2Mi4wMDJ6IiBmaWxsPSIjMTk5ZTlhIiBkYXRhLW9yaWdpbmFsPSIjMDAwMDAwIiBzdHlsZT0iIj48L3BhdGg+Cgk8L2c+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KCTxnPgoJCTxwYXRoIGQ9Ik0xNDIuMjUsNDU3LjAwMmgtMC4yN2MtNS41MjMsMC0xMCw0LjQ3Ny0xMCwxMHM0LjQ3NywxMCwxMCwxMGgwLjI3YzUuNTIzLDAsMTAtNC40NzcsMTAtMTAgICAgUzE0Ny43NzQsNDU3LjAwMiwxNDIuMjUsNDU3LjAwMnoiIGZpbGw9IiMxOTllOWEiIGRhdGEtb3JpZ2luYWw9IiMwMDAwMDAiIHN0eWxlPSIiPjwvcGF0aD4KCTwvZz4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgoJPGc+CgkJPHBhdGggZD0iTTQ2Mi4wMjUsNDU3LjAwMkgxNzAuOThjLTUuNTIzLDAtMTAsNC40NzctMTAsMTBzNC40NzcsMTAsMTAsMTBoMjkxLjA0NGM1LjUyMiwwLDEwLTQuNDc3LDEwLTEwICAgIFM0NjcuNTQ3LDQ1Ny4wMDIsNDYyLjAyNSw0NTcuMDAyeiIgZmlsbD0iIzE5OWU5YSIgZGF0YS1vcmlnaW5hbD0iIzAwMDAwMCIgc3R5bGU9IiI+PC9wYXRoPgoJPC9nPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+Cgk8Zz4KCQk8cGF0aCBkPSJNMTEwLjAzNSwzNWgtMC4yN2MtNS41MjMsMC0xMCw0LjQ3Ny0xMCwxMHM0LjQ3NywxMCwxMCwxMGgwLjI3YzUuNTIzLDAsMTAtNC40NzcsMTAtMTBTMTE1LjU1OSwzNSwxMTAuMDM1LDM1eiIgZmlsbD0iIzE5OWU5YSIgZGF0YS1vcmlnaW5hbD0iIzAwMDAwMCIgc3R5bGU9IiI+PC9wYXRoPgoJPC9nPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+Cgk8Zz4KCQk8cGF0aCBkPSJNODEuMDM2LDM1SDUwLjAwMWMtNS41MjMsMC0xMCw0LjQ3Ny0xMCwxMHM0LjQ3NywxMCwxMCwxMGgzMS4wMzRjNS41MjMsMCwxMC00LjQ3NywxMC0xMCAgICBDOTEuMDM2LDM5LjQ3Nyw4Ni41NTksMzUsODEuMDM2LDM1eiIgZmlsbD0iIzE5OWU5YSIgZGF0YS1vcmlnaW5hbD0iIzAwMDAwMCIgc3R5bGU9IiI+PC9wYXRoPgoJPC9nPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+Cgk8Zz4KCQk8cGF0aCBkPSJNMjkzLjQ4NSwyMDEuMTI5Yy0yLjM3Ny04Ljc5OS04LjAzNy0xNi4xNDYtMTUuOTIyLTIwLjY3NmwtNDguOTE4LTI4LjI0MmMtMC40OTEtMC4zMTktMC45ODYtMC42MjUtMS40ODQtMC45MTEgICAgbC00OC43OTItMjguMTcyYy01LjU4OS0zLjY4NS0xMi4wODYtNS42MzEtMTguODEzLTUuNjMxYy0xOC44NzEsMC0zNC4yMjQsMTUuMzUyLTM0LjIyNCwzNC4yMjR2MTE2LjYxMSAgICBjMCwwLjQwNywwLjAyNSwwLjgwOSwwLjA3MiwxLjIwM2MwLjE5Nyw1LjU4LDEuNzcyLDExLjA2LDQuNTg2LDE1LjkyYzYuMDksMTAuNTE2LDE3LjQyOCwxNy4wNDgsMjkuNTksMTcuMDQ4ICAgIGM1Ljk4NiwwLDExLjktMS41OTIsMTcuMDg5LTQuNTk4bDUwLjQ5Mi0yOS4xNTNjMC4yODYtMC4xNjYsMC41NjEtMC4zNDMsMC44MjUtMC41MzJsNDkuMjE5LTI4LjQxNCAgICBjNS4zNzMtMy4wMDEsOS44NC03LjQxNywxMi45MjMtMTIuNzc5QzI5NC42NjksMjE5LjEyNSwyOTUuODYxLDIwOS45MjcsMjkzLjQ4NSwyMDEuMTI5eiBNMjcyLjc5LDIxNy4wNjEgICAgYy0xLjI4NCwyLjIzMy0zLjE0MSw0LjA2Ny01LjM2OSw1LjMwNGMtMC4wNDksMC4wMjctMC4wOTksMC4wNTUtMC4xNDgsMC4wODNMMjE3LjE0LDI1MS4zOWMtMC4yODcsMC4xNjUtMC41NjIsMC4zNDMtMC44MjYsMC41MzIgICAgbC00OS42NTYsMjguNjdjLTIuMTU5LDEuMjUtNC42MDcsMS45MTItNy4wNzgsMS45MTJjLTUuMDUsMC05Ljc1Ny0yLjcwOS0xMi4yODMtNy4wNzFjLTEuMjUzLTIuMTY0LTEuOTE0LTQuNjE2LTEuOTEyLTcuMDkxICAgIGMwLTAuMzQtMC4wMTctMC42NzktMC4wNTItMS4wMTVWMTUxLjcyMWgwLjAwMWMwLTcuODQzLDYuMzgxLTE0LjIyNCwxNC4yMjQtMTQuMjI0YzIuODUxLDAsNS41OTUsMC44MzYsNy45MzYsMi40MTcgICAgYzAuMTk1LDAuMTMxLDAuMzk0LDAuMjU2LDAuNTk3LDAuMzc0bDQ5LjA5MSwyOC4zNDRjMC4yNDYsMC4xNDIsMC40ODIsMC4yOTIsMC43MTYsMC40NDdjMC4xNjYsMC4xMDksMC4zMzUsMC4yMTQsMC41MDcsMC4zMTQgICAgbDQ5LjE3NywyOC4zOTFjMy4yNywxLjg3OSw1LjYxMiw0LjkxOSw2LjU5Niw4LjU2MUMyNzUuMTYxLDIwOS45ODcsMjc0LjY2OCwyMTMuNzkzLDI3Mi43OSwyMTcuMDYxeiIgZmlsbD0iIzE5OWU5YSIgZGF0YS1vcmlnaW5hbD0iIzAwMDAwMCIgc3R5bGU9IiI+PC9wYXRoPgoJPC9nPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjxnIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyI+CjwvZz4KPGcgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KPC9nPgo8ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciPgo8L2c+CjwvZz48L3N2Zz4=" />
            <svg className="ctrl-svg" version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" x="0px" y="0px"
                viewBox="0 0 1000 1000" styles="enable-background:new 0 0 1000 1000;" space="preserve">
                <circle class="st0" cx="500" cy="500" r="302.8">
                    <animateTransform attributeType="xml"
                        attributeName="transform"
                        type="rotate"
                        from="0 500 500"
                        to="360 500 500"
                        dur="100s"
                        repeatCount="indefinite" />
                </circle>
                <circle class="st1" cx="500" cy="500" r="237.7">
                    <animateTransform attributeType="xml"
                        attributeName="transform"
                        type="rotate"
                        from="0 500 500"
                        to="360 500 500"
                        dur="40s"
                        repeatCount="indefinite" />
                </circle>
                <circle class="st2" cx="500" cy="500" r="366.8" transform="rotate(0 500 500)">
                    <animateTransform attributeType="xml"
                        attributeName="transform"
                        type="rotate"
                        from="0 500 500"
                        to="-360 500 500"
                        dur="50s"
                        repeatCount="indefinite" />
                </circle>
                <circle class="st3" cx="500" cy="500" r="385.1" />
            </svg>
        </div>
    );
}
export default Control;
