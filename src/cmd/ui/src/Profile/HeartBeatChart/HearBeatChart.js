import React, { useRef, useEffect, useState } from 'react';
import Chart from 'react-apexcharts'
import { getSensorsData } from '../../Api/Api'
import moment from 'moment'
import './HeartBeatChart.css'

class HeartBeatChart extends React.Component {

    getData() {
        var sensorData = getSensorsData(this.props.unit.ip)
        var data = []
        var time = []
        sensorData.forEach((item, index) => {
            data.push(item.heartbeat)
            time.push(moment.unix(item.time).format('hh:mm:ss'))
        })
        console.log(time)
        return {data:data, time:time}
    }

    constructor(props) {
        super(props);
        var graphData = this.getData()
        
        this.state = {
            series: [{
                data: graphData.data
            }],
            options: {
                chart: {
                    height: 500,
                    type: 'line',
                    zoom: {
                        enabled: false
                    },
                    toolbar: {
                        show: false,
                        tools:{
                          download:false // <== line to add
                        }
                      },
                      
                },
                dataLabels: {
                    enabled: false
                },
                stroke: {
                    curve: 'straight'
                },
                grid: {
                    show: true,
                    borderColor: 'rgb(25, 158, 154)',
                    position: 'back',
                    xaxis: {
                        lines: {
                            show: false
                        }
                    },   
                },
                xaxis: {
                    categories: graphData.time,
                }
            },
        };
    }

    render() {
        return (
            <>
            <div id="chart">
                <Chart options={this.state.options} series={this.state.series} type="line" height={400} className="hearbeat-chart"/>
            </div>
            </>
        )
    }
} export default HeartBeatChart;