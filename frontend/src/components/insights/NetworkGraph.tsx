import React, { useEffect, useRef, useState } from 'react';
import ForceGraph2D from 'react-force-graph-2d';
import apiClient from '@/lib/api';

interface Node {
  id: string;
  label: string;
  interactionCount: number;
  avgSentiment: number;
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
  const fgRef = useRef();

  useEffect(() => {
    async function fetchGraphData() {
      try {
        const response = await apiClient.get<NetworkGraphData>('/insights/network');
        setGraphData(response.data);
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
    <div className="w-full h-[600px] bg-gray-50 rounded-lg shadow-md">
      <ForceGraph2D
        ref={fgRef}
        graphData={graphData}
        nodeId="id"
        nodeLabel={nodeLabel}
        nodeAutoColorBy="label"
        nodeCanvasObject={(node, ctx, globalScale) => {
          const label = nodeLabel(node as Node);
          const fontSize = 12 / globalScale;
          ctx.font = `${fontSize}px Sans-Serif`;
          const textWidth = ctx.measureText(label).width;
          const bckgDimensions = [textWidth, fontSize].map(n => n + fontSize * 0.2); // some padding

          ctx.fillStyle = 'rgba(255, 255, 255, 0.8)';
          ctx.fillRect(node.x - bckgDimensions[0] / 2, node.y - bckgDimensions[1] / 2, bckgDimensions[0], bckgDimensions[1]);
          ctx.textAlign = 'center';
          ctx.textBaseline = 'middle';
          ctx.fillStyle = nodeColor(node as Node);
          ctx.fillText(label, node.x, node.y);

          (node as any).__bckgDimensions = bckgDimensions;
        }}
        nodePointerAreaPaint={(node, color, ctx) => {
          ctx.fillStyle = color;
          const bckgDimensions = (node as any).__bckgDimensions;
          bckgDimensions && ctx.fillRect(node.x - bckgDimensions[0] / 2, node.y - bckgDimensions[1] / 2, bckgDimensions[0], bckgDimensions[1]);
        }}
        linkColor={() => 'rgba(0,0,0,0.2)'}
        linkWidth={link => link.weight * 0.5} // Example: Thicker links for stronger relationships
        enableNodeDrag={true}
        enableZoomPan={true}
        d3AlphaDecay={0.04} // Slower decay for more stable graph after interactions
        d3VelocityDecay={0.2} // Slower velocity decay for smoother movements
      />
    </div>
  );
}
