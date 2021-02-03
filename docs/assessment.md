---
title: \huge{Tactical MANET Assessment}
author:
- Mahmoud Adas
- Yosry Mohammad
- Ahmed Mahmoud 
- Abdulrahman Khalid
date: Jan 26, 2021
classoption:
- twocolumn
---

# Abstract
A proposal and assessment for a graudation project idea about an application for tactical mobile ad-hoc networks.

# Project Description
A communication application for tactical teams, used in areas where there is no internet infrastructure. It should provide:

- Video streaming from soldiers' dash cams.
- Audio streaming from soldiers' microphones. 
- Transmitting raw data from various sensors (e.g GPS, thermal sensors, health sensors, etc \dots).
- Sending message codes (every code has its predefined meaning), then transmitting them live to the command center(s).

The application should allow the units to communicate securely, with low latency and high throughput.

The application should use a complex routing protocol that utilizes redundancy in the topology to increase communication reliability. We have to implement the routing protocol and build the application on top of it.

The application should be ready to deploy to devices with low-power microprocessors running linux. Unfortunately, we have to emulate those devices for budget reasons, and to be able to test various stressful mobility models, which would be nearly impossible using actual hardware. 

# Use Cases
- Nation wide emergency situations, when most/all of the internet infrastructure is gone, and emergency tactical units need to communicate with the command center with no preconfigured setup.
- Plot a frequently updated map of deployed units' positions and their current activities in places with no internet coverage.
- Rapidly notify nearby units of any danger.
- Use complex multi-layer digital encryption that is not feasible in analog communications via RF.
- Collecting live statistics from deployed units to quickly examine their behaviour and changes in the environment to be able to adjust strategy and commands quickly.
- During training and manoeuvres, to experiment with different hypothesis about soldiers' optimal behaviour. 

# Challenges
- Understand ad-hoc networks, and different families of routing protocols.
- Choose, study and implement a suitable routing protocol.
- Understand and work with linux routing interface (Netlink and RTNL).
- We may need to adjust the TCP protocol to be suitable for an ad-hoc environment (use ATCP).
- Video streaming is already complex enough using internet infrastructure, it is harder and more complex in ad-hoc networks because they are more dynamic.
- We need to understand and configure the emulators.

# Deliverables
- Application source code with instructions on how to build a linux image and how to attach sensors.
- Routing protocol implementation.
- A paper that describes the modification(s) to the routing protocol, if any.
- Experiments' results about latency and throughput using different mobility models.
