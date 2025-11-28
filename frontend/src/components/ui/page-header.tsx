import React from "react";
import { cn } from "@/lib/utils";

interface PageHeaderProps {
    title: string;
    children?: React.ReactNode;
    className?: string;
}

export function PageHeader({ title, children, className }: PageHeaderProps) {
    return (
        <div className={cn("flex items-center justify-between mb-4", className)}>
            <h1 className="text-2xl font-bold text-slate-800 tracking-tight">{title}</h1>
            {children && <div className="flex items-center gap-2">{children}</div>}
        </div>
    );
}
