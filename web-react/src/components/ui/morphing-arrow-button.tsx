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
};

const MorphingArrowButton = ({ direction, onClick, disabled, className, label }: MorphingArrowButtonProps) => {
  const [isHovered, setIsHovered] = useState(false);
  const isLeft = direction === "left";

  const containerVariants = {
    initial: { width: "64px", x: 0 },
    hover: { width: "112px", x: 0 },
  };
  const buttonVariants = {
    initial: { borderRadius: "50%", height: "64px", padding: "0" },
    hover: { borderRadius: isLeft ? "50px 14px 14px 50px" : "14px 50px 50px 14px", height: "64px", padding: "0 10px" },
  };
  const lineVariants = { initial: { width: 0 }, hover: { width: "calc(100% - 50px)" } };
  const arrowVariants = { initial: { x: "-50%" }, hover: { x: isLeft ? "-120%" : "20%" } };
  const Icon = isLeft ? ChevronLeft : ChevronRight;

  return (
    <div className={cn("morph-arrow-wrap inline-block w-[120px] overflow-visible", className)}>
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
              className={cn("absolute top-1/2 h-0.5 -translate-y-1/2 bg-current", isLeft ? "right-5" : "left-5")}
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
              <Icon className="h-6 w-6" />
            </motion.div>
          </div>
        </motion.button>
      </motion.div>
    </div>
  );
};

export default MorphingArrowButton;
