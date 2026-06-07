"use client";

import { MeshGradient } from "@paper-design/shaders-react";
import type React from "react";

export function ShaderBackground({ children }: { children: React.ReactNode }) {
  return (
    <div className="relative min-h-screen w-full overflow-hidden bg-[#f7fbff]">
      <MeshGradient
        className="absolute inset-0 h-full w-full opacity-70"
        colors={["#f8fbff", "#bcd7ff", "#ffffff", "#dff5ff", "#f4e7ff"]}
        speed={0.18}
      />
      <MeshGradient
        className="absolute inset-0 h-full w-full opacity-30"
        colors={["#ffffff", "#79a8ff", "#dbfff4", "#ffffff"]}
        speed={0.12}
      />
      {children}
    </div>
  );
}
