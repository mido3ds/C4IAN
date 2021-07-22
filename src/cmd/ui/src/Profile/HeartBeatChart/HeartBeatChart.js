import React, { useEffect, useState } from 'react';
import Chart from 'react-apexcharts'
import { getSensorsData } from '../../Api/Api'
import moment from 'moment'
import './HeartBeatChart.css'

function HeartBeatChart({ unit, port }) {
    const [series, setSeries] = useState([{ data: [] }])
    const [options, setOptions] = useState({
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
    })

    
    useEffect(() => {
        if (!unit || !port) return

        var data = []
        var time = []

        getSensorsData(unit.ip, port).then(sensorData => {
            if (!sensorData || !sensorData.length) {
                setSeries(series => {
                    return [{ data: [] }]
                })

                setOptions(options => {
                    return {
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
                        },
                    }
                })
                return;
            }

            sensorData.forEach((item, index) => {
                data.push(item.heartbeat)
                time.push(item.time)
            })

            if (!data || !time) return

            setSeries(() => {
                return [{ data: data }]
            })

            setOptions(() => {
                return {
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
                        categories: time,
                        labels: {
                            show: true,
                            formatter: function (val) {
                                return moment.unix(val/ (1000*1000)).format('hh:mm:ss') 
                            }
                        }
                    }
                }
            })
        })
    }, [unit, port])

    return (
        <>{!series[0].data.length ?
            <div className="no-data-heartbeat-msg">
                <p> No data to be previewed </p>
            </div> :
            <div id="chart">
                <Chart options={options} series={series} type="line" height={400} className="hearbeat-chart" />
            </div>
        }
        </>
    )

} export default HeartBeatChart;