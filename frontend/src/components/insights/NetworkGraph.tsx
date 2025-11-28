import React, { useEffect, useRef, useState, useCallback, useMemo } from 'react';
import ForceGraph2D, { ForceGraphMethods } from 'react-force-graph-2d';
import { api } from '@/lib/api';
import { useLanguage } from '@/lib/i18n/LanguageContext';
import { useTheme } from 'next-themes';

interface Node {
  id: string;
  label: string;
  interactionCount: number;
  avgSentiment: number;
  x?: number;
  y?: number;
  vx?: number;
  vy?: number;
}

interface Link {
  source: string | Node;
  target: string | Node;
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
  const [hoverNode, setHoverNode] = useState<Node | null>(null);
  const [activeFilter, setActiveFilter] = useState<'positive' | 'neutral' | 'negative' | null>(null);

  const { t } = useLanguage();
  const { theme } = useTheme();
  const isDark = theme === 'dark';

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const fgRef = useRef<any>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  // Memoize highlight sets to avoid recalculating on every render unless hover changes
  const { highlightNodes, highlightLinks } = useMemo(() => {
    const hNodes = new Set<string>();
    const hLinks = new Set<string>();

    if (hoverNode && graphData) {
      hNodes.add(hoverNode.id);
      graphData.links.forEach(link => {
        const sourceId = typeof link.source === 'object' ? (link.source as Node).id : link.source;
        const targetId = typeof link.target === 'object' ? (link.target as Node).id : link.target;

        if (sourceId === hoverNode.id || targetId === hoverNode.id) {
          hLinks.add(`${sourceId}-${targetId}`); // Store unique link identifier if needed, or just check source/target
          hNodes.add(sourceId);
          hNodes.add(targetId);
        }
      });
    }
    return { highlightNodes: hNodes, highlightLinks: hLinks };
  }, [hoverNode, graphData]);

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
        const response = await api.get<NetworkGraphData>('/insights/network');
        setGraphData(response.data);
      } catch (err: any) {
        console.error("Failed to fetch network graph:", err);
        setError(err.response?.data?.error || err.message || t('insights.loadingGraphError'));
      } finally {
        setLoading(false);
        setTimeout(() => {
          if (fgRef.current) {
            fgRef.current.d3Force('charge')?.strength(-100); // Adjust repulsion
            fgRef.current.d3Force('link')?.distance(50); // Adjust link distance
            fgRef.current.zoomToFit(400, 50);
          }
        }, 200);
      }
    }

    fetchGraphData();
  }, [t]);

  const getNodeColor = useCallback((node: Node) => {
    if (node.avgSentiment > 0.3) return '#10b981'; // Emerald 500
    if (node.avgSentiment < -0.3) return '#ef4444'; // Red 500
    return '#64748b'; // Slate 500
  }, []);

  const getNodeCategory = useCallback((node: Node) => {
    if (node.avgSentiment > 0.3) return 'positive';
    if (node.avgSentiment < -0.3) return 'negative';
    return 'neutral';
  }, []);

  const paintNode = useCallback((node: Node, ctx: CanvasRenderingContext2D, globalScale: number) => {
    const isHovered = node === hoverNode;
    const isHighlighted = highlightNodes.has(node.id);

    // Determine if node should be dimmed based on filter or hover
    let isDimmed = false;

    if (activeFilter) {
      // If filter is active, dim nodes that don't match the filter
      if (getNodeCategory(node) !== activeFilter) {
        isDimmed = true;
      }
    } else if (hoverNode) {
      // If no filter but hovering, dim non-highlighted nodes
      if (!isHighlighted) {
        isDimmed = true;
      }
    }

    // Dynamic size based on interactions (log scale for better distribution)
    const baseRadius = 4;
    const sizeMultiplier = Math.log(node.interactionCount + 1) * 2;
    const radius = baseRadius + sizeMultiplier;

    const x = node.x || 0;
    const y = node.y || 0;

    // Draw Node
    ctx.beginPath();
    ctx.arc(x, y, radius, 0, 2 * Math.PI, false);

    // Fill
    const color = getNodeColor(node);
    if (isDimmed) {
      ctx.fillStyle = 'rgba(203, 213, 225, 0.3)'; // Very transparent slate
    } else {
      ctx.fillStyle = color;
    }
    ctx.fill();

    // Border (Highlight)
    if ((isHovered || isHighlighted) && !isDimmed) {
      ctx.lineWidth = isHovered ? 2 : 1;
      ctx.strokeStyle = isDark ? '#fff' : '#334155';
      ctx.stroke();
    }

    // Label (Only show on hover or if highlighted, and not dimmed)
    if (!isDimmed && (isHovered || (isHighlighted && globalScale > 1.5))) {
      const label = node.label.split('@')[0];
      const fontSize = 12 / globalScale;
      ctx.font = `${isHovered ? 'bold' : ''} ${fontSize}px Sans-Serif`;
      ctx.textAlign = 'center';
      ctx.textBaseline = 'bottom';

      // Background for text readability
      const textWidth = ctx.measureText(label).width;
      const bckgDimensions = [textWidth, fontSize].map(n => n + fontSize * 0.2);

      ctx.fillStyle = isDark ? 'rgba(0, 0, 0, 0.8)' : 'rgba(255, 255, 255, 0.8)';
      ctx.fillRect(x - bckgDimensions[0] / 2, y - radius - bckgDimensions[1] - 2, bckgDimensions[0], bckgDimensions[1]);

      ctx.fillStyle = isDark ? '#f1f5f9' : '#1e293b';
      ctx.fillText(label, x, y - radius - 4);
    }
  }, [hoverNode, highlightNodes, isDark, getNodeColor, activeFilter, getNodeCategory]);

  const toggleFilter = (filter: 'positive' | 'neutral' | 'negative') => {
    setActiveFilter(prev => prev === filter ? null : filter);
  };

  if (loading) {
    return (
      <div className="flex flex-col items-center justify-center h-96 text-slate-400 gap-3">
        <div className="animate-spin w-8 h-8 border-2 border-indigo-500 border-t-transparent rounded-full"></div>
        <p className="text-sm font-medium">{t('insights.loadingGraph')}</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center h-96 text-red-500 bg-red-50/50 rounded-xl border border-red-100">
        <p>{t('common.error')}: {error}</p>
      </div>
    );
  }

  if (!graphData || graphData.nodes.length === 0) {
    return (
      <div className="flex items-center justify-center h-96 text-slate-400 bg-slate-50/50 rounded-xl border border-slate-100 border-dashed">
        <p>{t('insights.noGraphData')}</p>
      </div>
    );
  }

  return (
    <div ref={containerRef} className="w-full h-full relative group">
      <div className="absolute top-4 right-4 z-10 flex gap-2">
        <button
          onClick={() => toggleFilter('positive')}
          className={`flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-medium shadow-sm transition-all border ${activeFilter === 'positive'
              ? 'bg-emerald-50 border-emerald-200 text-emerald-700 ring-1 ring-emerald-200'
              : 'bg-white/90 backdrop-blur border-slate-200 text-slate-600 hover:bg-slate-50'
            }`}
        >
          <span className={`w-2 h-2 rounded-full ${activeFilter === 'positive' ? 'bg-emerald-600' : 'bg-emerald-500'}`}></span>
          <span>{t('insights.positive')}</span>
        </button>

        <button
          onClick={() => toggleFilter('neutral')}
          className={`flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-medium shadow-sm transition-all border ${activeFilter === 'neutral'
              ? 'bg-slate-100 border-slate-300 text-slate-700 ring-1 ring-slate-300'
              : 'bg-white/90 backdrop-blur border-slate-200 text-slate-600 hover:bg-slate-50'
            }`}
        >
          <span className={`w-2 h-2 rounded-full ${activeFilter === 'neutral' ? 'bg-slate-600' : 'bg-slate-500'}`}></span>
          <span>{t('insights.neutral')}</span>
        </button>

        <button
          onClick={() => toggleFilter('negative')}
          className={`flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs font-medium shadow-sm transition-all border ${activeFilter === 'negative'
              ? 'bg-red-50 border-red-200 text-red-700 ring-1 ring-red-200'
              : 'bg-white/90 backdrop-blur border-slate-200 text-slate-600 hover:bg-slate-50'
            }`}
        >
          <span className={`w-2 h-2 rounded-full ${activeFilter === 'negative' ? 'bg-red-600' : 'bg-red-500'}`}></span>
          <span>{t('insights.negative')}</span>
        </button>
      </div>

      <ForceGraph2D
        ref={fgRef}
        width={dimensions.width}
        height={dimensions.height}
        graphData={graphData}
        nodeId="id"
        nodeLabel={(node) => {
          // HTML Tooltip
          return `
             <div class="px-3 py-2 bg-slate-800 text-white text-xs rounded shadow-lg">
               <div class="font-bold mb-1">${node.label}</div>
               <div>${t('insights.interactions')}: ${node.interactionCount}</div>
               <div>${t('insights.sentiment')}: ${node.avgSentiment.toFixed(2)}</div>
             </div>
           `;
        }}
        nodeCanvasObject={paintNode}
        onNodeHover={setHoverNode}
        linkColor={(link) => {
          const sourceId = typeof link.source === 'object' ? (link.source as Node).id : link.source;
          const targetId = typeof link.target === 'object' ? (link.target as Node).id : link.target;

          let isDimmed = false;

          if (activeFilter) {
            // If filter is active, link is dimmed if either end is dimmed (i.e. not matching filter)
            // Actually, if we filter by node type, we probably only want to show links between visible nodes
            const sourceNode = graphData.nodes.find(n => n.id === sourceId);
            const targetNode = graphData.nodes.find(n => n.id === targetId);

            if (sourceNode && getNodeCategory(sourceNode) !== activeFilter) isDimmed = true;
            if (targetNode && getNodeCategory(targetNode) !== activeFilter) isDimmed = true;
          } else if (hoverNode) {
            const isConnected = sourceId === hoverNode.id || targetId === hoverNode.id;
            if (!isConnected) isDimmed = true;
          }

          return isDimmed ? 'rgba(203, 213, 225, 0.1)' : '#e2e8f0';
        }}
        linkWidth={(link) => {
          const sourceId = typeof link.source === 'object' ? (link.source as Node).id : link.source;
          const targetId = typeof link.target === 'object' ? (link.target as Node).id : link.target;
          if (hoverNode && (sourceId === hoverNode.id || targetId === hoverNode.id)) return 2;
          return 1;
        }}
        enableNodeDrag={true}
        d3AlphaDecay={0.05}
        d3VelocityDecay={0.3}
        warmupTicks={100}
        cooldownTicks={100}
      />
    </div>
  );
}
