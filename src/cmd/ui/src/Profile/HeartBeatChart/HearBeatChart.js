import React, { useRef, useEffect, useState } from 'react';
import Chart from 'react-apexcharts'
import './HeartBeatChart.css'

class HeartBeatChart extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            series: [{
                data: [10, 41, 35, 51, 49, 62, 69, 91, 148, 100, 150 ,170, 122,155 ,177 ,122 ,133 ,144 ,174]
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
            },
        };
    }

    render() {
        return (
            <div id="chart">
                <Chart options={this.state.options} series={this.state.series} type="line" height={400} className="hearbeat-chart"/>
            </div>
        )
    }
} export default HeartBeatChart;