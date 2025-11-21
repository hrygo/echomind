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
  const [dimensions, setDimensions] = useState({ width: 1, height: 1 });
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
    <div ref={containerRef} className="w-full h-[600px] bg-gray-50 rounded-lg shadow-md overflow-hidden">
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
          const textWidth = ctx.measureText(label).width;
          const bckgDimensions = [textWidth, fontSize].map(n => n + fontSize * 0.2); // some padding

          ctx.fillStyle = 'rgba(255, 255, 255, 0.8)';
          ctx.fillRect((node.x || 0) - bckgDimensions[0] / 2, (node.y || 0) - bckgDimensions[1] / 2, bckgDimensions[0], bckgDimensions[1]);
          ctx.textAlign = 'center';
          ctx.textBaseline = 'middle';
          ctx.fillStyle = nodeColor(node as Node);
          ctx.fillText(label, (node.x || 0), (node.y || 0));

          node.__bckgDimensions = bckgDimensions as [number, number];
        }}
        nodePointerAreaPaint={(node: Node, color, ctx) => {
          ctx.fillStyle = color;
          const bckgDimensions = node.__bckgDimensions;
          if (bckgDimensions) {
            ctx.fillRect((node.x || 0) - bckgDimensions[0] / 2, (node.y || 0) - bckgDimensions[1] / 2, bckgDimensions[0], bckgDimensions[1]);
          }
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
