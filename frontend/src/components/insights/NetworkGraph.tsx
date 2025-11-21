import React, { useEffect, useRef, useState } from 'react';
import ForceGraph2D from 'react-force-graph-2d';
import apiClient from '@/lib/api';

interface Node {
  id: string;
  label: string;
  interactionCount: number;
  avgSentiment: number;
  __bckgDimensions?: [number, number];
  x?: number;
  y?: number;
}

interface Link {
  source: string;
  target: string;
  weight: number;
}

interface NetworkGraphData {
  nodes: Node[];
  links: Link[];
}

export default function NetworkGraph() {
  const [graphData, setGraphData] = useState<NetworkGraphData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [dimensions, setDimensions] = useState({ width: 800, height: 600 });
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const fgRef = useRef<any>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!containerRef.current) return;

    const resizeObserver = new ResizeObserver((entries) => {
      for (const entry of entries) {
        const { width, height } = entry.contentRect;
        setDimensions({ width, height });
      }
    });

    resizeObserver.observe(containerRef.current);

    return () => resizeObserver.disconnect();
  }, []);

  useEffect(() => {
    async function fetchGraphData() {
      try {
        const response = await apiClient.get<NetworkGraphData>('/insights/network');
        setGraphData(response.data);
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      } catch (err: any) {
        console.error("Failed to fetch network graph:", err);
        setError(err.response?.data?.error || err.message || 'Failed to fetch network graph.');
      } finally {
        setLoading(false);
        // Wait for next tick to ensure graph is rendered with data
        setTimeout(() => {
          if (fgRef.current) {
            fgRef.current.zoomToFit(400);
            fgRef.current.d3ReheatSimulation();
          }
        }, 100);
      }
    }

    fetchGraphData();
  }, []);

  if (loading) {
    return <div className="flex items-center justify-center h-64 text-gray-500">Loading network graph...</div>;
  }

  if (error) {
    return <div className="flex items-center justify-center h-64 text-red-500">Error: {error}</div>;
  }

  if (!graphData || graphData.nodes.length === 0) {
    return <div className="flex items-center justify-center h-64 text-gray-500">No network data available.</div>;
  }

  const nodeColor = (node: Node) => {
    // Color nodes based on sentiment (e.g., positive green, negative red, neutral grey)
    if (node.avgSentiment > 0.3) return 'green';
    if (node.avgSentiment < -0.3) return 'red';
    return 'grey';
  };

  const nodeLabel = (node: Node) => `${node.label} (Interactions: ${node.interactionCount}, Sentiment: ${node.avgSentiment.toFixed(2)})`;

  return (
    <div ref={containerRef} className="w-full h-full bg-slate-50/50 overflow-hidden">
      <ForceGraph2D
        ref={fgRef}
        width={dimensions.width}
        height={dimensions.height}
        graphData={graphData}
        nodeId="id"
        nodeLabel={nodeLabel}
        nodeAutoColorBy="label"
        nodeCanvasObject={(node: Node, ctx, globalScale) => {
          const label = nodeLabel(node as Node);
          const fontSize = 12 / globalScale;
          ctx.font = `${fontSize}px "Inter", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif`;
          
          // 1. Draw Node Circle
          const radius = 5;
          ctx.beginPath();
          ctx.arc((node.x || 0), (node.y || 0), radius, 0, 2 * Math.PI, false);
          ctx.fillStyle = nodeColor(node as Node);
          ctx.fill();
          
          // 2. Draw Text Label
          ctx.textAlign = 'center';
          ctx.textBaseline = 'top';
          ctx.fillStyle = '#475569'; // Slate-600 for text
          ctx.fillText(label, (node.x || 0), (node.y || 0) + radius + 2);

          // Store dimensions for pointer interaction (approximate for circle)
          node.__bckgDimensions = [radius * 2, radius * 2]; 
        }}
        nodePointerAreaPaint={(node: Node, color, ctx) => {
          ctx.fillStyle = color;
          const radius = 5;
          ctx.beginPath();
          ctx.arc((node.x || 0), (node.y || 0), radius + 2, 0, 2 * Math.PI, false); // Slightly larger hit area
          ctx.fill();
        }}
        linkColor={() => 'rgba(0,0,0,0.2)'}
        linkWidth={link => link.weight * 0.5} // Example: Thicker links for stronger relationships
        enableNodeDrag={true}
        d3AlphaDecay={0.04} // Slower decay for more stable graph after interactions
        d3VelocityDecay={0.2} // Slower velocity decay for smoother movements
      />
    </div>
  );
}
