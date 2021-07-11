import "ol/ol.css"
import Graticule from "ol/layer/Graticule"
import Map from "ol/Map"
import Stroke from "ol/style/Stroke"
import Text from "ol/style/Text"
import View from "ol/View"
import { fromLonLat, transform } from "ol/proj"
import Feature from "ol/Feature"
import { Circle } from "ol/geom"
import { OSM, Vector as VectorSource } from "ol/source"
import { Fill, Style } from "ol/style"
import { Draw, Translate, Select } from "ol/interaction"
import { getCenter as getExtentCenter } from 'ol/extent'
import { Tile as TileLayer, Vector as VectorLayer } from "ol/layer"
import { METERS_PER_UNIT } from 'ol/proj/Units'
import Overlay from 'ol/Overlay'
import { shiftKeyOnly } from "ol/events/condition"
import ExtentInteraction from "ol/interaction/Extent";

let prefix = ''
let mode = ''
let numUnits = 1
let numCmds = 1
let halfRange = true

const START_CENTER = fromLonLat([6.7318473939, 0.3320770836])

function getRange() {
  return parseFloat(document.getElementById('range').value)
}
function setRange(range) {
  document.getElementById('range').value = range
}

function calcInterval(zlen) {
  let d = 1.0 / (1 << zlen) // [0, 1]
  d *= 2                    // [0, 2]
  d = d > 1 ? (d - 2) : d   // [-1, 1]
  d *= 180                  // [-180, 180]
  return d
}

function getZLen() {
  return parseInt(document.getElementById('zlen').value)
}

function bestZoom(zlen) {
  return zlen - 1
}
function getMaxZoom(zlen) {
  return bestZoom(zlen) * (1 + 0.113)
}
function getMinZoom(zlen) {
  return bestZoom(zlen) * (1 - 0.070)
}

function indexToDegrees(i) {
  /*i*/                    // [0, 0xFFFF]
  let d = i / 0xFFFF       // [0, 1]
  d *= 2                   // [0, 2]
  d = d > 1 ? d - 2 : d    // [-1, 1]
  d *= 180                 // [-180, 180]
  return d
}

function degreesToIndex(d) {
  /*d*/                  // [-180, 180]
  d /= 180               // [-1, 1]
  d = d < 0 ? d + 2 : d  // [0, 2]
  d /= 2                 // [0, 1]
  d *= 0xFFFF            // [0, 0xFFFF]
  return d
}

function zlenMask(zlen) {
  return ~(0xFFFF >> zlen)
}

function getZoneID(lon, lat, zlen) {
  const x = degreesToIndex(lon) & zlenMask(zlen)
  const y = degreesToIndex(lat) & zlenMask(zlen)
  return x.toString(16) + '.' + y.toString(16) + '/' + zlen
}

const raster = new TileLayer({
  source: new OSM({
    wrapX: false,
    visible: true
  })
})

const view = new View({ center: START_CENTER })
function updateViewLabels() {
  document.getElementById('zoomLbl').textContent = view.getZoom().toFixed(5)
}

view.on('change', updateViewLabels)
function updateView() {
  view.setZoom(bestZoom(getZLen()))
  if (document.getElementById('zoomBoundry').checked) {
    view.setMaxZoom(getMaxZoom(getZLen()))
    view.setMinZoom(getMinZoom(getZLen()))
  }
  updateViewLabels()
}
updateView()

document.getElementById('zoomBoundry').addEventListener('change', () => {
  if (document.getElementById('zoomBoundry').checked) {
    view.setMaxZoom(getMaxZoom(getZLen()))
    view.setMinZoom(getMinZoom(getZLen()))
  } else {
    view.setMaxZoom(17)
    view.setMinZoom(.5)
  }
})

document.getElementById('lonlatBtn').addEventListener('click', () => {
  const lon = document.getElementById('lonInput').value
  const lat = document.getElementById('latInput').value
  view.animate({ center: fromLonLat([lon, lat]) })
})
document.getElementById('zoneidBtn').addEventListener('click', () => {
  const zlen = getZLen()

  let x = document.getElementById('zoneidInput1').value
  x = parseInt(x, 16) & zlenMask(zlen)

  let y = document.getElementById('zoneidInput2').value
  y = parseInt(y, 16) & zlenMask(zlen)

  view.animate({ center: fromLonLat([indexToDegrees(x), indexToDegrees(y)]), zoom: bestZoom(zlen) })
  updateViewLabels()
})

const source = new VectorSource()
const vector = new VectorLayer({
  source: source,
  style: (f) => [
    new Style({
      stroke: new Stroke({
        width: 2,
        color: [255, 0, 0, 1],
      }),
      fill: new Fill({
        color: [255, 0, 0, .1],
      }),
      text: new Text({
        font: '17px Calibri,sans-serif',
        fill: new Fill({ color: '#000' }),
        stroke: new Stroke({
          color: '#fff', width: 2
        }),
        text: f.getId()
      })
    })
  ]
})

let select = new Select({
  filter: (f) => {
    return f.getId() !== undefined
  },
  style: (f) => [
    new Style({
      stroke: new Stroke({
        width: 2,
        color: [0, 0, 255, 1],
      }),
      fill: new Fill({
        color: [255, 0, 0, .1],
      }),
      text: new Text({
        font: '17px Calibri,sans-serif',
        fill: new Fill({ color: '#000' }),
        stroke: new Stroke({
          color: '#fff', width: 2
        }),
        text: f.getId()
      })
    })
  ]
})

let selectedFeature = null
select.on('select', (e) => {
  if (mode === 'delete') {
    e.selected.forEach((f) => { source.removeFeature(f) })
  } else {
    if (e.selected.length == 1) {
      selectedFeature = e.selected[0]
    } else {
      selectedFeature = null
    }
  }
})

var translate = new Translate({
  features: select.getFeatures(),
})

const draw = new Draw({
  source: source,
  type: "Circle"
})

const map = new Map({
  layers: [raster, vector],
  target: "map",
  view: view
})


var cross = new Overlay({
  element: document.getElementById('overlay'),
  stopEvent: false,
  positioning: 'center-center'
});
cross.setPosition(START_CENTER);
map.addOverlay(cross);
function updateCross() {
  const center = view.getCenter()
  coord = transform(center, 'EPSG:3857', 'EPSG:4326')
  document.getElementById('lonLbl').textContent = coord[0].toFixed(10)
  document.getElementById('latLbl').textContent = coord[1].toFixed(10)
  cross.setPosition(center)
  document.getElementById('zoneidLbl').textContent = getZoneID(coord[0], coord[1], getZLen()).toUpperCase()
}

map.on('pointermove', updateCross)
map.on('moveend', updateCross)

let grid
function updateGrid() {
  if (grid) {
    map.removeLayer(grid)
  }
  grid = new Graticule({
    intervals: [calcInterval(getZLen())],
    strokeStyle: new Stroke({
      color: "rgba(255,120,0,0.9)",
      width: 2,
      lineDash: [0.5, 4]
    }),
    showLabels: true,
    wrapX: false
  })
  map.addLayer(grid)
}
updateGrid()

function onZlenChanged() {
  const zlen = getZLen()
  if (zlen > 16) {
    document.getElementById('zlen').value = 16
  } else if (zlen <= 0) {
    document.getElementById('zlen').value = 1
  }

  updateGrid()
  updateView()
  document.getElementById('zoneidLbl').textContent = getZoneID(coord[0], coord[1], getZLen()).toUpperCase()
}

document.getElementById('zlen').addEventListener('change', onZlenChanged)

function newCirlceFeature(center, range, name) {
  if (halfRange) {
    range /= 2
  }

  const f = new Feature(new Circle(center, range))
  f.setId(name)
  return f;
}

function drawCircleInMeter(center, range, name) {
  source.addFeature(newCirlceFeature(center, range, name));
}

map.on('singleclick', (e) => {
  if (mode === 'add') {
    let name = ''
    if (prefix === 'u') {
      name = `${prefix}${numUnits++}`
    } else if (prefix === 'c') {
      name = `${prefix}${numCmds++}`
    }

    drawCircleInMeter(e.coordinate, getRange(), name)
  }
})

function onActionChange() {
  if (document.getElementById('add-unit').checked) {
    mode = 'add'
    prefix = 'u'
    select.setMap(null)
    translate.setMap(null)
    draw.setMap(null)
  } else if (document.getElementById('add-cmd').checked) {
    mode = 'add'
    prefix = 'c'
    select.setMap(null)
    translate.setMap(null)
    draw.setMap(null)
  } else if (document.getElementById('move').checked) {
    mode = 'move'
    map.addInteraction(select)
    map.addInteraction(translate)
    draw.setMap(null)
  } else if (document.getElementById('delete').checked) {
    mode = 'delete'
    map.addInteraction(select)
    translate.setMap(null)
    draw.setMap(null)
  }
}

document.getElementsByName('action').forEach((i) => i.addEventListener('change', onActionChange))
onActionChange()

function importFile(fileContent) {
  const json = JSON.parse(fileContent)

  let avgCenter = [0, 0]
  let total = 0
  setRange(json.range)

  let features = json.nodes.map(n => {
    const center = fromLonLat([n.lon, n.lat])
    avgCenter[0] += center[0]
    avgCenter[1] += center[1]
    total++

    return newCirlceFeature(center, json.range, n.name);
  })

  if (total === 0) {
    avgCenter = START_CENTER
  } else {
    avgCenter[0] /= total
    avgCenter[1] /= total
  }

  // set zlen
  document.getElementById('zlen').value = json.zlen
  onZlenChanged()

  source.clear()
  source.addFeatures(features)

  view.animate({ center: avgCenter, zoom: bestZoom(getZLen()) })
  updateViewLabels()
}

function replaceFeatures(features) {
  // remove existing
  features.forEach(f => {
    f2 = source.getFeatureById(f.getId())
    if (f2) {
      source.removeFeature(f2)
    }
  })

  source.addFeatures(features)
}

function importFromMininet(data, change) {
  const json = JSON.parse(data)

  let avgCenter = [0, 0]
  let total = 0
  setRange(json.range)

  let features = json.nodes.filter(n => {
    if (selectedFeature && n.name === selectedFeature.getId()) {
      return false
    }
    return true
  }).map(n => {
    const center = fromLonLat([n.lon, n.lat])
    avgCenter[0] += center[0]
    avgCenter[1] += center[1]
    total++

    return newCirlceFeature(center, json.range, n.name);
  })

  replaceFeatures(features)

  if (change) {
    if (total === 0) {
      avgCenter = START_CENTER
    } else {
      avgCenter[0] /= total
      avgCenter[1] /= total
    }

    // set zlen
    document.getElementById('zlen').value = json.zlen
    onZlenChanged()

    view.animate({ center: avgCenter, zoom: bestZoom(getZLen()) })
    updateViewLabels()
  }
}

function downloadToFile(content, filename) {
  const a = document.createElement('a')
  const file = new Blob([content], { type: 'text/plain' })

  a.href = URL.createObjectURL(file)
  a.download = filename
  a.click()

  URL.revokeObjectURL(a.href)
}

function featureToJSON(f) {
  const center = getExtentCenter(f.getGeometry().getExtent())
  const lonlat = transform(center, 'EPSG:3857', 'EPSG:4326')
  const lon = lonlat[0]
  const lat = lonlat[1]
  const name = f.getId()

  return { name, lon, lat }
}

function exportFile() {
  const nodes = source.getFeatures().map(featureToJSON)
  const range = getRange()
  const zlen = getZLen()

  downloadToFile(JSON.stringify({ nodes, range, zlen }), 'exported.topo')
}

document.getElementById('exportBtn').addEventListener('click', exportFile)
document.getElementById('importBtn').addEventListener('click', () => {
  console.log('import..')
  document.getElementById('file-input').value = null
  document.getElementById('file-input').click()
})
document.getElementById('file-input').addEventListener('change', () => {
  var file = document.getElementById("file-input").files[0]
  if (file) {
    var reader = new FileReader()
    reader.readAsText(file, "UTF-8")
    reader.onload = (evt) => importFile(evt.target.result)
    reader.onerror = function (evt) {
      alert("couldn't read file")
      console.error("couldn't read file")
    }
  }
})

function onRangeChanged() {
  let range = getRange()
  if (range < 50) {
    setRange(50)
  } else if (range > 50000) {
    setRange(50000)
  }
  range = getRange()
  
  source.getFeatures().forEach(f => {
    source.removeFeature(f)
  
    const center = getExtentCenter(f.getGeometry().getExtent())
    const name = f.getId()
    drawCircleInMeter(center, range, name)
  })
}

document.getElementById('range').addEventListener('change', onRangeChanged)
document.getElementById('halfRange').addEventListener('change', () => {
  halfRange = document.getElementById('halfRange').checked
  onRangeChanged()
})

function onSendMsg(socket) {
  if (selectedFeature) {
    const json = { type: 'setNodePosition', ...featureToJSON(selectedFeature) }
    socket.send(JSON.stringify(json))
  }
}

let firstReceive = true
function onReceiveMsg(msg) {
  if (firstReceive) {
    // clear
    numUnits = 1
    numCmds = 1

    // import
    importFromMininet(msg.data, true)

    firstReceive = false
  } else {
    importFromMininet(msg.data, false)
  }
}

let minVel = 10
let maxVel = 30
let aggregation = 0.9

let socket = null
let mobStarted = false
let extent = null

const mobBtn = document.getElementById('startMobility')
mobBtn.addEventListener('click', () => {
  const type = 'setMobility'
  if (mobStarted) {
    mobBtn.textContent = 'Start'
    mobStarted = false
    socket.send(JSON.stringify({ type, action: 'stop' }))
  } else {
    mobBtn.textContent = 'Stop'
    socket.send(JSON.stringify({ type, action: 'start', aggregation, minVel, maxVel, extent }))
    mobStarted = true
  }
})

function isOpen(ws) { return ws && ws.readyState === ws.OPEN }

const cnctBtn = document.getElementById('connect')

function disconnect() {
  console.log("disconnecting")
  cnctBtn.innerHTML = 'Connect'
  cnctBtn.className = 'toConnect'

  if (isOpen(socket)) {
    socket.close(1000, "disconnect")
  }
  socket = null

  document.getElementById('add-unit').disabled = false
  document.getElementById('add-cmd').disabled = false
  document.getElementById('delete').disabled = false
  firstReceive = true

  mobBtn.textContent = 'Start'
  mobStarted = false
  document.getElementById('mob-fieldset').disabled = true
  document.getElementById('range').disabled = false

  map.removeInteraction(extentInter)
  document.getElementById('extentTip').className = 'hidden'
}

let inter = null
function connect() {
  const UPDATE_RATE = 1 / 15.0 * 1000 // 15fps

  console.log("connecting")
  socket = new WebSocket('ws://localhost:' + document.getElementById('portInput').value)
  socket.onopen = () => {
    console.log('connected')
    firstReceive = true
    cnctBtn.innerHTML = 'Disconnect'
    cnctBtn.className = 'toDisconnect'

    inter = setInterval(() => {
      if (isOpen(socket)) {
        onSendMsg(socket)
      } else {
        clearInterval(inter)
        disconnect()
      }
    }, UPDATE_RATE)

    document.getElementById('add-unit').disabled = true
    document.getElementById('add-cmd').disabled = true
    document.getElementById('delete').disabled = true
    document.getElementById('move').checked = true
    onActionChange()
    document.getElementById('mob-fieldset').disabled = false
    document.getElementById('range').disabled = true

    map.addInteraction(extentInter)
    document.getElementById('extentTip').className = ''
  }
  socket.onerror = disconnect
  socket.onmessage = onReceiveMsg
}

cnctBtn.addEventListener('click', () => {
  if (socket) {
    disconnect()
  } else {
    connect()
  }
})

document.getElementById('reset').addEventListener('click', () => {
  disconnect()
  source.clear()
  numUnits = 1
  numCmds = 1
  extent = null
  document.getElementById('zlen').value = 16
  onZlenChanged()
  view.setZoom(bestZoom(16))
  view.animate({ center: START_CENTER })
  updateViewLabels()
})

function onMobChange() {
  if (mobStarted) {
    socket.send(JSON.stringify({ type: 'setMobility', action: 'change', aggregation, minVel, maxVel, extent }))
  }
}

const upInput = document.getElementById('minVelocity')
const downInput = document.getElementById('maxVelocity')
function onVelocityChanged() {
  const v1 = parseFloat(upInput.value)
  const v2 = parseFloat(downInput.value)

  maxVel = v1 > v2 ? v1 : v2
  document.getElementById('maxVelocityText').textContent = maxVel

  minVel = v1 < v2 ? v1 : v2
  document.getElementById('minVelocityText').textContent = minVel

  onMobChange()
}

upInput.addEventListener('input', onVelocityChanged)
downInput.addEventListener('input', onVelocityChanged)

addEventListener('input', e => {
  let _t = e.target;
  _t.parentNode.style.setProperty(`--${_t.id}`, +_t.value)
}, false);

document.getElementById('aggregInput').addEventListener('input', () => {
  aggregation = document.getElementById('aggregInput').value / 100
  document.getElementById('aggregText').textContent = aggregation.toFixed(2)
  onMobChange()
})

const extentInter = new ExtentInteraction({
  condition: shiftKeyOnly,
  boxStyle: (f) => [
    new Style({
      stroke: new Stroke({
        width: 2,
        color: [0, 0, 0, 1],
      }),
      fill: new Fill({
        color: [255, 255, 255, .1],
      }),
    })
  ]
});

extentInter.on("extentchanged", (e) => {
  if (e.extent) {
    const [x1, y2, x2, y1] = e.extent
    const [lon1, lat1] = transform([x1, y1], 'EPSG:3857', 'EPSG:4326')
    const [lon2, lat2] = transform([x2, y2], 'EPSG:3857', 'EPSG:4326')
    extent = {
      lon1, lat1,
      lon2, lat2,
    }
    onMobChange()
  } else {
    extent = null
  }
});