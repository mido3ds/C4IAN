---
title: \huge{Tactical MANET}
---

<!-- use this template before submitting https://drive.google.com/file/d/13Ot9Atu3ej9Qkhx067NDUw7WJr63n8wF/view?fbclid=IwAR2cSzbeAUklcX4H2SILcNGqgvFeBzn9GKeCGNjVsJfFKL69gAVUFI6j4T8 -->

# Team Members

| Name               | Email                                            |
|--------------------|--------------------------------------------------|
| Mahmoud Adas       | \texttt{mahmoud.ibrahim97@eng-st.cu.edu.eg}      |
| Yosry Mohammad     | \texttt{yosry.mohammad99@eng-st.cu.edu.eg}       |
| Ahmed Mahmoud      | \texttt{Ahmed.Afifi98@eng-st.cu.edu.eg}          |
| Abdulrahman Khalid | \texttt{abdulrahman.elshafie98@eng-st.cu.edu.eg} |

# 1. Problem Statement
<!-- Introduction to the problem (max 30 words) -->
A `mobile ad-hoc network` communication system for military, for operations in areas with no internet infrastructure.
Deployed units can stream audio, video and sensors readings to command Centers.
Command Centers can stream audio and message codes to some/all unit(s).

# 2. Motivation
<!-- Why are you motivated to work on this problem? (max 30 words) -->
We are interested in decentralized/distributed algorithms and designing/building complex systems.

# 3. System Architecture
<!-- In this section, draw the block diagram of your system showing the flow between
different modules. -->
Figure \ref{fig:modules} shows the modules diagram.

![Modules Diagram \label{fig:modules}](figures/modules_diagram.png)

# 4. List of Deliverables
<!-- State the main modules of your system with its function, inputs and expected outputs
- Number of modules must be at least equal to number of team members
- Max number of modules including the integration of whole project must not exceed 6
modules -->
| Module Name        | Function                                                                                                     | Input                                                                                    | Expected Output                                | % of used Libraries |
|--------------------|--------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------|------------------------------------------------|---------------------|
| Unit Client        | Stream and receive streams to/from command Centers                                                           | Device audio, video, sensors and message codes. Streams and messages from command center | Send streams and show play audio/messages      | TODO                |
| Cmd. Center Client | Stream and receive streams to/from deployed units. Shows a map of all units with their statistics            | Audio and message codes. Streams and messages from deplyed units                         | Send streams and show play audio/messages      | TODO                |
| Router             | Determine how a certain ip-packet should be forwarded. Implements some `MANET` ad-hoc protocol               | IP packet (with final destination) to forward                                            | Path from this node to final destination       | TODO                |
| Testbed            | Build, configure and monitor the simulation/emulation of the `MANEt`. Define the topology and mobility model | User commands and arguments or configuration file                                        | Commands to emulation, simulation or actual-HW | TODO                |