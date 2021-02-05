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
This document lists details of our graduation project requirements and specifications.

# Project Description
A `mobile ad-hoc network` communication system for military, for operations in areas with no internet infrastructure.
The system connects the command center(s) with deployed units in two-way communications.

# System Architecture
The system is composed of devices (nodes) running linux-based operating systems and have certain programs running in them.

## Nodes
All nodes are provided with wireless communication modules that follow `IEEE 802.11` standards.

There are 2 types of nodes: units and command centers.

### Units
Devices with deployed units in the operation field, connected with: 

- LCD screen with resolution of `48x84` pixels.
- Helmet video camera. 
- Audio input. 
- Keybad. 
- GPS (or any other position detection system.)
- Heartbeat sensor.

Features:

- Low power consumption.
- Running on battery.
- Low wireless range.
- High mobility.
- Operated by one person.

### Command Centers
High-end computers at the command and control centers, accessed by units leaders.

Features:

- Capabale of high power consumption.
- Powerfull CPUs.
- Big storage and RAM.
- Operated by multiple peeople with multiple wide screens.
- Wide wireless range.
- Installed nearby the operation field, and has a connection to devices in the field.
- Low (or zero) mobility.

## Programs {#subsec:programs}
Programs running in devices are running as daemons, started at the startup of the system and are always running and restarted on failure.

### Units
Each unit has a public and private key, and a map of command centers `IP`s and their corresponding public keys.

> Extension: units announce their `IP`s to command centers and share their keys dynamically.

A unit device has 2 programs:
- `Router`: implements routing protocol.
- `Unit Client Daemon`: Connected to device hardware and network interface and provide all unit features.

### Command Centers
Each command center has a public and private key, and a map of units `IP`s and their corresponding public keys.

> Extension: command centers announce their `IP`s to units and share their keys dynamically.

A command center computer has 3 programs:
- `Router`: implements routing protocol, same router as in unit devices.
- `Command Client Daemon`: Exposes an interface to `UI` program, connects to units clients and handles all communication with units.
- `Command Client UI`: Connects to `Command Client Daemon`, shows all data in the daemon and controls it.

# Functional Requirements
## Units
- Stream video from combat-cameras to command center(s) only if the latter requested them. Video streaming terminates if the unit received end stream request, or the start request wasn't refreshed after 1 minute.
- Stream the heartbeat of the device owner and their position every 10 seconds.
- Store all the recorded video and sensors data locally, in a rolling db, where new data override old data when there is no left space. This feature could be controlled in a configuration file.
- If the owner requested:
    + Send audio messages from microphone.
    + Send code messages (every code has its predefined meaning.)
- Receive audio messages from command centers into a queue.
- Play received audio messages from the queue instantly.
- Receive and show code messages.
- Temporarily store audio and code messages.
- Access stored audio and code messages and delete them.

## Command Centers
- Can be accessed from multiple computers.
- Send audio commands and command codes (every code has its predefined meaning) to one (unicast), some (multicast) or all (boradcast) of the deployed units devices.
- Store all received messages, video streams and other data (position and heartbeat).
- Show old archieved data.
- Play any audio message and video stream. Received messages don't autoplay, but the UI shows if a received message was already heard or not.
- Group units.
- Show map of all units, their group color.
- If didn't receive heartbeat and position from some unit for 2 minutes, mark it as inactive, show list of inactive units and show their latest position on map with inactive color.
- Show no-heartbeat warning, when some unit gives near-zero heartbeat for some time period. Color units with no heartbeat and show them clearly on map.
- Show statistics about selected unit or group of units. Statistics include:
    + their movement speed over time, 
    + average/max/min/median speed,
    + and connection outages with command center over time.

# Non-functional Requirements
## Reliability
The following must be delivered reliably (with gurantee of delivery):

- Code messages.
- Audio messages.

The following can be delivered unreliably (*no* gurantee of delivery):

- Video streams.
- Position and heartbeat messages (minimum 80% delivery success rate).

## Speed
The system allows nodes to communicate with low latency and high throughput.
Video streams must be viewable at minimum of 20 fps.

## Routing
> TODO: define routing algorithm

The system uses a complex routing protocol that utilizes redundancy in the topology to increase communication reliability.

## Security
- All transmitted data are encrypted.
- Authentiction is required for accessing command center by its UI.
- All stored data in command centers and units are encrypted.
- Units don't persist any data, messages self destruct after a 3 minutes of receiving them.

# Testbeds
The system will be tested in 2 different environments: virtual and actual hardware.

## Virtual
Using virtualization/emulation, each node (unit/command center) will be deployed in a virtual machine.
Each node will have a static ip equivalent to that stored in nodes databases.

Ther should be `UI` for units' clients that:
- connects with them over forwarded ports,
- receives their screens and audio, 
- and sends them button actions and fake audio/video/position/heartbeat inputs,

Mininet-wifi will be used to simulate the wiereless connections and create topologies.

The following mobility models should be tested:

- Random Walk
- Truncated Levy Walk
- Random Direction
- Random Way Point
- Gauss Markov
- Reference Point
- Time Variant Community

Different toplogies with up-to 25 nodes should be tested.

## Hardware
1. Install clients on our laptops.
2. Create a minimum topology with at least one command center.
3. Use the virtual `UI` for unit client.
4. Test streaming video/audio in a small ad-hoc network.
5. Remove one unit and test that the rest of the units noticed that and changed their routes.

> Extension: create actual hardware for the unit, modify the client to read actual inputs instead of taking them virtually.

# Deliverables
- Source code of all programs listed in Subsection \ref{subsec:programs}.
- Instructions on how to:
    + Configure devices, connect inputs.
    + Install all dependencies.
    + Configure all software.
    + Install and run all software.
- A paper that describes the modification(s) to the routing protocol, if any.
- Experiments' results about latency and throughput using different mobility models.
