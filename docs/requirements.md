---
title: \huge{Tactical MANET Project Requirements}
author:
- Mahmoud Adas
- Yosry Mohammad
- Ahmed Mahmoud 
- Abdulrahman Khalid
date: \today
classoption:
- twocolumn
---

# Abstract
This document lists details of the graduation project requirements and specifications.

# Project Description
A communication system for military, used in areas with no internet infrastructure. 
The system connects the command centre(s) with deployed units in two-way communications.

## Nodes
Nodes include:
- Fixed number of pre-known command centres computers.
- Devices with soldiers, connected with sensors, dashcam and audio input.

All nodes are provided with wireless communication modules that follow `IEEE 802.11` standards.

## Functional Features
The system should let the soldiers devices:

- stream video from dash cams,
- stream audio from microphones,
- stream raw data from various sensors (e.g GPS, thermal sensors, health sensors, etc \dots),
- and send message codes (every code has its predefined meaning)

to all the command centres.

The system should let the command centres:

- send audio commands,
- and send command codes (every code has its predefined meaning)

to one (unicast), some (multicast) or all (boradcast) of the soldiers devices.

## Non-functional Features
The system should allow the units to communicate securely, with low latency and high throughput.

The system have to use a complex routing protocol that utilizes redundancy in the topology to increase communication reliability.

The system should be ready to deploy to devices with low-power microprocessors running linux. 

# Deliverables
- Application source code.
- Routing protocol implementation.
- Instructions on how to:
    + Attach inputs.
    + Configure devices.
    + Install and run all software
- A paper that describes the modification(s) to the routing protocol, if any.
- Experiments' results about latency and throughput using different mobility models.
