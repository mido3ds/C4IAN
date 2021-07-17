import React from 'react';
import Chart from 'react-apexcharts'
import { getSensorsData } from '../../Api/Api'
import moment from 'moment'
import './HeartBeatChart.css'

class HeartBeatChart extends React.Component {

    getData() {
        if (this.props.unit) {
            getSensorsData(this.props.unit.ip).then(sensorData => {
                if (!sensorData || !sensorData.length) return null;
                var data = []
                var time = []
                console.log(sensorData)
                sensorData.forEach((item, index) => {
                    data.push(item.heartbeat)
                    time.push(moment.unix(item.time).format('hh:mm'))
                })
                return { data: data, time: time }
            })
            return null;
        }
    }

    componentDidMount() {
        var data = []
        var time = []
        if (!this.props.unit) return

        getSensorsData(this.props.unit.ip).then(sensorData => {
            if (!sensorData || !sensorData.length) return null;
            sensorData.forEach((item, index) => {
                data.push(item.heartbeat)
                time.push(item.time)
            })
        })
        if (!data || !time) return
        console.log(time)
        this.setState({
            series: [{
                data: data
            }],
            options: {
                chart: {
                    height: 500,
                    type: 'line',
                    zoom: {
                        enabled: false
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
                },
                xaxis: {
                    show: false,
                    categories: time,
                    labels: {
                        show: true,
                        formatter: function(val) {
                            return moment.unix(time[val]).format('hh:mm:ss') // formats to hours:minutes
                        } 
                    }
                }
            },
        });
    }

    constructor(props) {
        super(props);

        this.state = {
            series: [{
                data: []
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
                        tools: {
                            download: false // <== line to add
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
                    categories: [],
                }
            },
        }
    }

    render() {
        return (
            <>{!this.state ?
                <div className="no-data-heartbeat-msg">
                    <p> No data to be previewed </p>
                </div> :
                <div id="chart">
                    <Chart options={this.state.options} series={this.state.series} type="line" height={400} className="hearbeat-chart" />
                </div>
            }
            </>
        )
    }
} export default HeartBeatChart;