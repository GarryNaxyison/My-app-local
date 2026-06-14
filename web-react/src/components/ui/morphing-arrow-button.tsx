"use client";

import { useState } from "react";
import { motion } from "framer-motion";
import { ChevronLeft, ChevronRight } from "lucide-react";

import { cn } from "@/lib/utils";

type MorphingArrowButtonProps = {
  direction: "left" | "right";
  onClick?: () => void;
  disabled?: boolean;
  className?: string;
  label?: string;
  size?: "default" | "compact";
};

const MorphingArrowButton = ({ direction, onClick, disabled, className, label, size = "default" }: MorphingArrowButtonProps) => {
  const [isHovered, setIsHovered] = useState(false);
  const isLeft = direction === "left";
  const isCompact = size === "compact";
  const metrics = isCompact
    ? { initialWidth: "40px", hoverWidth: "56px", height: "40px", padding: "0 6px", lineWidth: "calc(100% - 30px)", lineInset: isLeft ? "right-3" : "left-3", icon: "h-5 w-5", wrap: "w-[58px]" }
    : { initialWidth: "64px", hoverWidth: "112px", height: "64px", padding: "0 10px", lineWidth: "calc(100% - 50px)", lineInset: isLeft ? "right-5" : "left-5", icon: "h-6 w-6", wrap: "w-[120px]" };

  const containerVariants = {
    initial: { width: metrics.initialWidth, x: 0 },
    hover: { width: metrics.hoverWidth, x: 0 },
  };
  const buttonVariants = {
    initial: { borderRadius: "50%", height: metrics.height, padding: "0" },
    hover: { borderRadius: isLeft ? "50px 14px 14px 50px" : "14px 50px 50px 14px", height: metrics.height, padding: metrics.padding },
  };
  const lineVariants = { initial: { width: 0 }, hover: { width: metrics.lineWidth } };
  const arrowVariants = { initial: { x: "-50%" }, hover: { x: isLeft ? "-120%" : "20%" } };
  const Icon = isLeft ? ChevronLeft : ChevronRight;

  return (
    <div className={cn("morph-arrow-wrap inline-block overflow-visible", metrics.wrap, className)}>
      <motion.div
        className={cn("flex items-center", isLeft ? "justify-end" : "justify-start")}
        variants={containerVariants}
        initial="initial"
        animate={isHovered ? "hover" : "initial"}
        transition={{ duration: 0.42, ease: "easeInOut" }}
        onHoverStart={() => setIsHovered(true)}
        onHoverEnd={() => setIsHovered(false)}
      >
        <motion.button
          className="relative flex w-full cursor-pointer items-center justify-center overflow-hidden border border-current bg-transparent disabled:pointer-events-none disabled:opacity-35"
          variants={buttonVariants}
          initial="initial"
          animate={isHovered ? "hover" : "initial"}
          transition={{ duration: 0.42, ease: "easeInOut" }}
          onClick={onClick}
          disabled={disabled}
          type="button"
          aria-label={label || direction}
        >
          <div className="relative flex h-full w-full items-center">
            <motion.div
              className={cn("absolute top-1/2 h-0.5 -translate-y-1/2 bg-current", metrics.lineInset)}
              variants={lineVariants}
              initial="initial"
              animate={isHovered ? "hover" : "initial"}
              transition={{ duration: 0.42, ease: "easeInOut" }}
            />
            <motion.div
              className="absolute left-1/2 top-1/2 -translate-y-1/2"
              variants={arrowVariants}
              initial="initial"
              animate={isHovered ? "hover" : "initial"}
              transition={{ duration: 0.42, ease: "easeInOut" }}
            >
              <Icon className={metrics.icon} />
            </motion.div>
          </div>
        </motion.button>
      </motion.div>
    </div>
  );
};

export default MorphingArrowButton;
