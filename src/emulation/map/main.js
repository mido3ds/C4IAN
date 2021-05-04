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
import Overlay from 'ol/Overlay';

let prefix = ''
let mode = ''
let numUnits = 1
let numCmds = 1

const START_CENTER = fromLonLat([6.7318473939, 0.3320770836])

function getRange() {
  return parseFloat(document.getElementById('range').value)
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
  return zlen-1
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
  view.setMaxZoom(getMaxZoom(getZLen()))
  view.setMinZoom(getMinZoom(getZLen()))
  updateViewLabels()
}
updateView()

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

select.on('select', (e) => {
  if (mode === 'delete') {
    e.selected.forEach((f) => { source.removeFeature(f) })
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

function drawCircleInMeter(center, range, name) {
  let view = map.getView()
  let projection = view.getProjection()
  let resolutionAtEquator = view.getResolution()
  let pointResolution = projection.getPointResolutionFunc()(resolutionAtEquator, center)
  let resolutionFactor = resolutionAtEquator / pointResolution
  range = (range / METERS_PER_UNIT.m) * resolutionFactor

  const f = new Feature(new Circle(center, getRange()))
  f.setId(name)
  source.addFeature(f)
}

map.on('singleclick', (e) => {
  if (mode === 'add') {
    let name = ''
    if (prefix === 'u') {
      name = `${prefix}${numUnits++}`
    } else if (prefix === 'c') {
      name = `${prefix}${numCmds++}`
    }

    drawCircleInMeter(e.coordinate, range, name)
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

  const features = json.nodes.map(n => {
    const center = fromLonLat([n.lon, n.lat])
    avgCenter[0] += center[0]
    avgCenter[1] += center[1]
    total++

    const f = new Feature(new Circle(center, json.range))
    f.setId(n.name)
    return f
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

function downloadToFile(content, filename) {
  const a = document.createElement('a')
  const file = new Blob([content], { type: 'text/plain' })

  a.href = URL.createObjectURL(file)
  a.download = filename
  a.click()

  URL.revokeObjectURL(a.href)
}

function exportFile() {
  const nodes = source.getFeatures().map(f => {
    const center = getExtentCenter(f.getGeometry().getExtent())
    const lonlat = transform(center, 'EPSG:3857', 'EPSG:4326')
    const lon = lonlat[0]
    const lat = lonlat[1]
    const name = f.getId()

    return { name, lon, lat }
  })
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

document.getElementById('range').addEventListener('change', () => {
  let range = getRange()
  if (range < 50) {
    document.getElementById('range').value = 50
  } else if (range > 50000) {
    document.getElementById('range').value = 50000
  }
  range = getRange()

  source.getFeatures().forEach(f => {
    source.removeFeature(f)

    const center = getExtentCenter(f.getGeometry().getExtent())
    const name = f.getId()
    drawCircleInMeter(center, range, name)
  })
})

document.getElementById('reset').addEventListener('click', () => {
  source.clear()
  numUnits = 1
  numCmds = 1
  document.getElementById('zlen').value = 16
  onZlenChanged()
  view.setZoom(bestZoom(16))
  view.animate({ center: START_CENTER })
  updateViewLabels()
})