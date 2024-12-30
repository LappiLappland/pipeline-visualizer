'use client';
import { useContext, useEffect, useMemo, useRef, useState } from "react";
import VisualizerNode from "./VisualizerNode";
import VisualizerNodesGroup from "./VisualizerNodesGroup";
import NodesGroup from "./VisualizerNodesGroup";
import getVisualizerDataObject from "@/app/helpers/getVisualizerDataObject";
import { PipelineContext } from "../contexts/PipelineContext";
import { useParams } from "next/navigation";

export default function Visualizer({

}) {

    const {pipelineDataByLevels: dataByLevel, pipelineData: data, permissions, stats} = useContext(PipelineContext)
    const params = useParams();

    const mainRef = useRef();
    const mainContainerRef = useRef();
    //const [draggingPoint, setDraggingPoint] = useState(null);
    const [hoveredNode, setHoveredNode] = useState(null);
    const draggingPoint = useRef(null);
    const touchDistance = useRef(null);
    const currentCoords = useRef({x: 0, y: 0});
    const currentScale = useRef(0.9);

    // Setup events
    useEffect(() => {
        let ref = mainRef.current;
        if (mainRef.current) {
            mainRef.current.addEventListener('wheel', handleWheel, { passive: false });

            mainRef.current.addEventListener('touchmove', handleMove, { passive: false });
            mainRef.current.addEventListener('touchstart', handleStart, { passive: false });
            mainRef.current.addEventListener('touchend', handleEnd, { passive: false });

            document.addEventListener('mousemove', handleMove);
            document.addEventListener('mouseup', handleEnd);
            mainRef.current.addEventListener('mousedown', handleStart);
        }

        return () => {
            if (ref) {
                ref.removeEventListener('wheel', handleWheel);

                ref.removeEventListener('touchmove', handleMove);
                ref.removeEventListener('touchstart', handleStart);
                ref.removeEventListener('touchend', handleEnd);

                document.removeEventListener('mousemove', handleMove);
                document.removeEventListener('mouseup', handleEnd);
                ref.removeEventListener('mousedown', handleStart);
            }
        };
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    // Init scale and position.
    useEffect(() => {
        updateScreen(currentCoords.current, currentScale.current);
    }, []);

    function getNodesGroups(dataLevel) {
        const nodesGroups = [];
        let i = 0;
        while (dataLevel[i]) {
            const nodes = dataLevel[i].map((nodeId) => {
                const node = data[nodeId];
                let isHovered = false;
                if (hoveredNode) {
                    const hoveredNodeReal = data[hoveredNode];
                    const isHoveredCurrent = node.id === hoveredNode;
                    //const isHoveredParent = node.parents.some(parent => parent.id === hoveredNode);
                    //const isHoveredChild = node.dependencies.some(child => child.id === hoveredNode);
                    const isHoveredBloodline = hoveredNodeReal.bloodline.some(node => nodeId === node.id);
                    //isHovered = isHoveredCurrent || isHoveredParent || isHoveredChild;
                    isHovered = isHoveredCurrent || isHoveredBloodline;
                }
                
                return (
                    <VisualizerNode
                        key={[node.id, node.status].join('@%-')}
                        id={node.id}
                        pipelineId={params.pipelineId}
                        title={node.title}
                        canExecute={permissions.canExecute && stats.pipelineStatus !== 'pending'}
                        status={node.status}
                        progress={node.progress > 99 ? 0 : node.progress}
                        startedAt={node.startedAt}
                        finishedAt={node.finishedAt}
                        startedBy={node.startedBy}
                        isUnhovered={hoveredNode && !isHovered}
                        isHovered={hoveredNode && isHovered}
                        isMainHovered={hoveredNode === node.id}
                        onMouseEnter={() => setHoveredNode(node.id)}
                        onMouseLeave={() => {
                            if (hoveredNode === node.id) {
                                setHoveredNode(null);
                            }
                        }}
                    />
                )
            });
            nodesGroups.push(
                <VisualizerNodesGroup
                    key={i}
                >
                    {nodes}
                </VisualizerNodesGroup>
            )
            i++;
        }
        return nodesGroups;
    }

    function getNodesLinks(dataLevel) {
        const nodesLinks = [];
        const nodesLinksFront = [];
        
        let i = 0;
        while (dataLevel[i]) {
            dataLevel[i].forEach((nodeId) => {
                const node = data[nodeId]
                node.parents.forEach(parent => {
                    let isHovered = false;
                    if (hoveredNode) {
                        const isHoveredJoin = parent.id === hoveredNode;
                        const isHoveredExit = node.id === hoveredNode;
                        isHovered = isHoveredJoin || isHoveredExit;

                        //const hoveredNodeReal = data[hoveredNode];
                        //const isHoveredCurrent = node.id === hoveredNode;
                        //const isHoveredBloodline = hoveredNodeReal.bloodline.some(node => nodeId === node.id);
                        //isHovered = isHoveredCurrent || isHoveredBloodline
                    }

                    const levelDiff = parent.dependencyDepth - node.dependencyDepth - 1;

                    const xStart = 348 * i + 300;
                    const xEnd = xStart + 48;
                    const yStart = 75 * (node.row) + 25;
                    const yEnd = 75 * (parent.row) + 25

                    //let draw = `M ${xStart} ${yStart} L ${xEnd} ${yEnd}`;
                    let draw = `M ${xStart} ${yStart} H ${levelDiff === 0 ? xEnd : xEnd + 348 * levelDiff}`;
                    if (yStart !== yEnd) {
                        const minus = yStart > yEnd ? '-' : '';
                        draw = `M ${xStart} ${yStart}
                        ${levelDiff === 0 ? '' : 'h ' + 348 * levelDiff}
                        c 12 0 24 0 24 ${minus}12
                        V ${minus ? yEnd + 12 : yEnd - 12}
                        c 0 0 0 ${minus}12 24 ${minus}12
                        `
                    }

                    nodesLinks.push(
                        <path
                        key={`${node.id}_to_${parent.id}`}
                        className={`transition-colors ease-out delay-50 relative stroke-gray-85  ${hoveredNode ? 'opacity-50' : ''}`}
                        fill="none"
                        d={draw}
                        />
                    );
                    nodesLinksFront.push(
                        <path
                        key={`${node.id}_to_${parent.id}`}
                        className={`transition-colors ease-out delay-50 relative ${isHovered ? 'stroke-blue-400' : 'stroke-transparent'}`}
                        fill="none"
                        d={draw}
                        />
                    )
                });
            });
            i++;
        }

        return [nodesLinks, nodesLinksFront];
    }

    function updateScreen(coords, scale) {
        if (mainContainerRef.current) {
            currentCoords.current = {x: coords.x, y: coords.y};
            currentScale.current = Math.min(Math.max(0.125, scale), 4);
            mainContainerRef.current.style.transform = `translate(${coords.x}px, ${coords.y}px) scale(${scale})`;
        }
    }

    function handleStart(e) {
        //e.preventDefault();
        
        const clientX = e.touches ? e.touches[0].clientX : e.clientX;
        const clientY = e.touches ? e.touches[0].clientY : e.clientY;

        const coordsNew = {
            x: clientX - currentCoords.current.x,
            y: clientY - currentCoords.current.y,
        }

        draggingPoint.current = coordsNew

        if (mainRef.current) {
            mainRef.current.style.cursor = 'grabbing';
        }
    }
    
    function handleMove(e) {
        if (e.touches && e.touches.length === 2 && touchDistance.current) {
            e.preventDefault();
    
            const touch1 = e.touches[0];
            const touch2 = e.touches[1];
            const distance = Math.sqrt(
                Math.pow(touch2.clientX - touch1.clientX, 2) +
                Math.pow(touch2.clientY - touch1.clientY, 2)
            );
    
            const delta = distance - touchDistance.current;

            const zoomFactor = delta > 0 ? 1.05 : 0.95;
            currentScale.current *= zoomFactor;

            updateScreen(currentCoords.current, currentScale.current);
    
            touchDistance.current = distance;
        } 
        else if (draggingPoint.current) {
            e.preventDefault();

            const clientX = e.touches ? e.touches[0].clientX : e.clientX;
            const clientY = e.touches ? e.touches[0].clientY : e.clientY;
    
            const newX = clientX - draggingPoint.current.x;
            const newY = clientY - draggingPoint.current.y;

            updateScreen({ x: newX, y: newY }, currentScale.current);
        }
    }
    
    function handleEnd() {
        draggingPoint.current = null;
        touchDistance.current = null;

        if (mainRef.current) {
            mainRef.current.style.cursor = 'grab';
        }
    }

    function handleWheel(e) {
        e.stopPropagation();
        e.preventDefault();
        let scale = currentScale.current + e.deltaY * -0.001;
        updateScreen(currentCoords.current, scale);
    }

    return (
        <div className={`w-full h-full p-8 overflow-hidden 
            select-none
            min-h-64 max-h-screen
            ${draggingPoint ? 'cursor-grabbing' : 'cursor-grab'}
            `}
            ref={mainRef}
        >
            <div className="flex justify-center items-center"
            ref={mainContainerRef}
            >
                <div className="inline-flex flex-row gap-12 relative"
                >
                    <svg className="absolute w-full h-full stroke-[2px] stroke-white">
                        {getNodesLinks(dataByLevel)[0]}
                    </svg>
                    <svg className="absolute w-full h-full stroke-[3px]">
                        {getNodesLinks(dataByLevel)[1]}
                    </svg>
                    {getNodesGroups(dataByLevel)}
                </div>
            </div>
        </div>
    )
}