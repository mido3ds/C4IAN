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

let prefix = ''
let mode = ''
let numUnits = 1
let numCmds = 1

const START_CENTER = fromLonLat([31.2108959, 30.0272552])

function getRange() {
  return parseFloat(document.getElementById('range').value)
}

function calcInterval(zlen) {
  return 825.3644039596112 / (2 ** zlen)
}

function getZLen() {
  return parseInt(document.getElementById('zlen').value)
}

const bestZoom = [
  1.9106593850919042, // 0
  1.9106593850919042, // 1
  1.9106593850919042, // 2
  2.1592278779030680, // 3
  3.4465588756583085, // 4
  3.4465588756583085, // 5
  5.2093235813695510, // 6
  5.6673399339085660, // 7
  6.7044641422046390, // 8
  7.6488708143325940, // 9
  8.5171023882495000, // 10
  9.6466901351308640, // 11
  10.478585300924003, // 12
  11.321827975883714, // 13
  12.889894303560485, // 14
  13.561995542821041, // 15
  14.800000000000000, // 16
]
function getMaxZoom(zlen) {
  return bestZoom[zlen] * (1 + 0.113)
}
function getMinZoom(zlen) {
  return bestZoom[zlen] * (1 - 0.113)
}

const raster = new TileLayer({
  source: new OSM({
    wrapX: false,
    visible: true
  })
})

const view = new View({ center: START_CENTER })
function updateView() {
  view.setZoom(bestZoom[getZLen()])
  view.setMaxZoom(getMaxZoom(getZLen()))
  view.setMinZoom(getMinZoom(getZLen()))
}
updateView()

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
  updateGrid()
  updateView()
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
  let n = 0

  const features = json.nodes.map(n => {
    const center = fromLonLat([n.lon, n.lat])
    avgCenter[0] += center[0]
    avgCenter[1] += center[1]

    const f = new Feature(new Circle(center, json.range))
    f.setId(n.name)
    return f
  })

  if (n === 0) {
    avgCenter = START_CENTER
  } else {
    avgCenter[0] /= n
    avgCenter[1] /= n
  }

  // set zlen
  document.getElementById('zlen').value = json.zlen
  onZlenChanged()

  source.clear()
  source.addFeatures(features)

  view.animate({ center: avgCenter, zoom: bestZoom[getZLen()] })
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
  const range = getRange()
  source.getFeatures().forEach(f => {
    source.removeFeature(f)

    const center = getExtentCenter(f.getGeometry().getExtent())
    const name = f.getId()
    drawCircleInMeter(center, range, name)
  })
})

document.getElementById('clear').addEventListener('click', () => source.clear())