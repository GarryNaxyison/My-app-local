"use client";

import type React from "react";
import { motion } from "framer-motion";
import { cn } from "@/lib/utils";

type FallingPatternProps = React.ComponentProps<"div"> & {
  color?: string;
  backgroundColor?: string;
  duration?: number;
  blurIntensity?: string;
  density?: number;
};

export function FallingPattern({
  color = "rgba(122, 162, 255, 0.62)",
  backgroundColor = "#050914",
  duration = 150,
  blurIntensity = "1em",
  density = 1,
  className,
}: FallingPatternProps) {
  const patterns = [
    "235", "252", "150", "253", "204", "134", "179", "299", "215", "281", "158", "210",
  ].flatMap((height, index) => {
    const y = Number(height);
    const x = index * 25;
    return [
      `radial-gradient(4px 100px at ${x}px ${y}px, ${color}, transparent)`,
      `radial-gradient(4px 100px at ${x + 300}px ${y}px, ${color}, transparent)`,
      `radial-gradient(1.5px 1.5px at ${x + 150}px ${y / 2}px, ${color} 100%, transparent 150%)`,
    ];
  });

  const backgroundSizes = Array.from({ length: patterns.length }, (_, index) => {
    const heights = [235, 252, 150, 253, 204, 134, 179, 299, 215, 281, 158, 210];
    return `300px ${heights[Math.floor(index / 3)]}px`;
  }).join(", ");

  const startPositions = Array.from({ length: patterns.length }, (_, index) => `${index * 7}px ${120 + index * 11}px`).join(", ");
  const endPositions = Array.from({ length: patterns.length }, (_, index) => `${index * 7}px ${5200 + index * 379}px`).join(", ");

  return (
    <div className={cn("relative h-full w-full overflow-hidden", className)}>
      <motion.div
        className="relative size-full"
        initial={{ opacity: 0 }}
        animate={{ opacity: 1 }}
        transition={{ duration: 0.24 }}
      >
        <motion.div
          className="relative z-0 size-full"
          style={{
            backgroundColor,
            backgroundImage: patterns.join(", "),
            backgroundSize: backgroundSizes,
          }}
          variants={{
            initial: { backgroundPosition: startPositions },
            animate: {
              backgroundPosition: [startPositions, endPositions],
              transition: { duration, ease: "linear", repeat: Number.POSITIVE_INFINITY },
            },
          }}
          initial="initial"
          animate="animate"
        />
      </motion.div>
      <div
        className="absolute inset-0 z-[1]"
        style={{
          backdropFilter: `blur(${blurIntensity})`,
          backgroundImage: `radial-gradient(circle at 50% 50%, transparent 0, transparent 2px, ${backgroundColor} 2px)`,
          backgroundSize: `${8 * density}px ${8 * density}px`,
        }}
      />
    </div>
  );
}
