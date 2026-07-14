"use client";

import { LogOut, Undo2 } from "lucide-react";
import { motion, AnimatePresence } from "motion/react";
import React, { useEffect, useState } from "react";

type LogoutButtonProps = {
  onConfirm: () => void;
  disabled?: boolean;
  label: string;
  cancelLabel: string;
  iconOnly?: boolean;
};

export default function LogoutButton({ onConfirm, disabled, label, cancelLabel, iconOnly = false }: LogoutButtonProps) {
  const [armed, setArmed] = useState(false);
  const [count, setCount] = useState(4);
  const [isAnimating, setIsAnimating] = useState(false);

  useEffect(() => {
    if (!armed) return;
    if (count === 0) {
      onConfirm();
      return;
    }
    const timer = setTimeout(() => setCount((c) => c - 1), 1000);
    return () => clearTimeout(timer);
  }, [armed, count, onConfirm]);

  const handleClick = (newState: boolean) => {
    if (isAnimating || disabled) return;
    setIsAnimating(true);
    setArmed(newState);
    if (newState) setCount(4);
    setTimeout(() => setIsAnimating(false), 400);
  };

  return (
    <div className="flex items-center justify-center">
      <AnimatePresence mode="popLayout" initial={false}>
        {!armed ? (
          <motion.button
            key="logout"
            layoutId="logoutButton"
            onClick={() => handleClick(true)}
            whileTap={{ scale: 0.95 }}
            style={{ pointerEvents: isAnimating || disabled ? "none" : "auto", minWidth: iconOnly ? 44 : 104, minHeight: 44, visibility: "visible" }}
            initial={{ backgroundColor: "#FE322A", filter: "blur(0px)", opacity: 1 }}
            animate={{ backgroundColor: "#FE322A", filter: "blur(0px)", opacity: 1 }}
            exit={{ backgroundColor: "#FFEDF1", filter: "blur(1px)", opacity: 0 }}
            className="logout-button-v2 flex items-center justify-center gap-2 overflow-hidden rounded-full px-5 py-3 text-white"
            transition={{
              layout: { duration: 0.4, ease: [0.77, 0, 0.175, 1] },
              backgroundColor: { duration: 0.4, ease: "easeInOut" },
              filter: { duration: 0.1, ease: "easeInOut" },
              opacity: { duration: 0.2, ease: "easeOut" },
            }}
            aria-label={label}
          >
            <LogOut size={15} />
            {!iconOnly ? (
              <motion.span layoutId="logoutText" className="flex">
                {label.split("").map((char, index) => (
                  <motion.span
                    key={`logout-${index}`}
                    initial={{ y: 20, opacity: 0, scale: 0.3 }}
                    animate={{ y: 0, opacity: 1, scale: 1 }}
                    exit={{ y: -20, opacity: 0, scale: 0.3 }}
                    transition={{ duration: 0.3, delay: index * 0.005, ease: [0.785, 0.135, 0.15, 0.86] }}
                    style={{ display: "inline-block", whiteSpace: "pre" }}
                  >
                    {char}
                  </motion.span>
                ))}
              </motion.span>
            ) : null}
          </motion.button>
        ) : iconOnly ? (
          <motion.button
            key="cancel-icon"
            layoutId="logoutButton"
            onClick={() => handleClick(false)}
            whileTap={{ scale: 0.95 }}
            style={{ pointerEvents: isAnimating || disabled ? "none" : "auto", minWidth: 44, minHeight: 44, visibility: "visible" }}
            initial={{ backgroundColor: "#FE322A", filter: "blur(1px)", opacity: 0 }}
            animate={{ backgroundColor: "#FFEDF1", filter: "blur(0px)", opacity: 1 }}
            exit={{ backgroundColor: "#FE322A", filter: "blur(1px)", opacity: 0 }}
            className="logout-button-v2 logout-button-v2--icon-confirm flex items-center justify-center gap-1 overflow-hidden rounded-full px-3 py-3"
            transition={{
              layout: { duration: 0.4, ease: [0.77, 0, 0.175, 1] },
              backgroundColor: { duration: 0.4, ease: "easeInOut" },
              filter: { duration: 0.2, ease: "easeInOut" },
              opacity: { duration: 0.2, ease: "easeIn" },
            }}
            aria-label={cancelLabel}
          >
            <motion.div className="flex shrink-0 items-center justify-center rounded-full bg-[#FE322A] p-1.5" initial={{ opacity: 0, scale: 0.5 }} animate={{ opacity: 1, scale: 1 }} exit={{ opacity: 0, scale: 0.5 }}>
              <Undo2 className="h-4 w-4 text-white" />
            </motion.div>
            <motion.div className="relative flex min-w-[24px] shrink-0 items-center justify-center overflow-hidden rounded-full bg-[#FE322A] px-3 py-2 text-xs font-semibold text-white" initial={{ opacity: 0, scale: 0.5 }} animate={{ opacity: 1, scale: 1 }} exit={{ opacity: 0, scale: 0.5 }}>
              <AnimatePresence mode="popLayout">
                <motion.span key={count} initial={{ opacity: 0, y: 10, scale: 0.8 }} animate={{ opacity: 1, y: 0, scale: 1 }} exit={{ opacity: 0, y: -10, scale: 0.8 }} transition={{ duration: 0.2, ease: [0.33, 1, 0.68, 1] }} className="absolute">
                  {count}
                </motion.span>
              </AnimatePresence>
            </motion.div>
          </motion.button>
        ) : (
          <motion.button
            key="cancel"
            layoutId="logoutButton"
            onClick={() => handleClick(false)}
            whileTap={{ scale: 0.95 }}
            style={{ pointerEvents: isAnimating || disabled ? "none" : "auto", minWidth: 104, minHeight: 44, visibility: "visible" }}
            initial={{ backgroundColor: "#FE322A", filter: "blur(1px)", opacity: 0 }}
            animate={{ backgroundColor: "#FFEDF1", filter: "blur(0px)", opacity: 1 }}
            exit={{ backgroundColor: "#FE322A", filter: "blur(1px)", opacity: 0 }}
            className="logout-button-v2 flex items-center gap-2 overflow-hidden rounded-full px-3 py-3"
            transition={{
              layout: { duration: 0.4, ease: [0.77, 0, 0.175, 1] },
              backgroundColor: { duration: 0.4, ease: "easeInOut" },
              filter: { duration: 0.2, ease: "easeInOut" },
              opacity: { duration: 0.2, ease: "easeIn" },
            }}
            aria-label={cancelLabel}
          >
            <motion.div className="flex shrink-0 items-center justify-center rounded-full bg-[#FE322A] p-1.5" initial={{ opacity: 0, scale: 0.5 }} animate={{ opacity: 1, scale: 1 }} exit={{ opacity: 0, scale: 0.5 }}>
              <Undo2 className="h-4 w-4 text-white" />
            </motion.div>
            <motion.span layoutId="logoutText" className="flex font-medium text-[#FE322A]">
              {cancelLabel.split("").map((char, index) => (
                <motion.span
                  key={`cancel-${index}`}
                  initial={{ y: 20, opacity: 0, scale: 0.3 }}
                  animate={{ y: 0, opacity: 1, scale: 1 }}
                  exit={{ y: -20, opacity: 0, scale: 0.3 }}
                  transition={{ duration: 0.3, delay: index * 0.006, ease: [0.785, 0.135, 0.15, 0.86] }}
                  style={{ display: "inline-block", whiteSpace: "pre" }}
                >
                  {char}
                </motion.span>
              ))}
            </motion.span>
            <motion.div className="relative flex min-w-[32px] shrink-0 items-center justify-center overflow-hidden rounded-full bg-[#FE322A] px-4 py-3 text-sm font-semibold text-white" initial={{ opacity: 0, scale: 0.5 }} animate={{ opacity: 1, scale: 1 }} exit={{ opacity: 0, scale: 0.5 }}>
              <AnimatePresence mode="popLayout">
                <motion.span key={count} initial={{ opacity: 0, y: 10, scale: 0.8 }} animate={{ opacity: 1, y: 0, scale: 1 }} exit={{ opacity: 0, y: -10, scale: 0.8 }} transition={{ duration: 0.2, ease: [0.33, 1, 0.68, 1] }} className="absolute">
                  {count}
                </motion.span>
              </AnimatePresence>
            </motion.div>
          </motion.button>
        )}
      </AnimatePresence>
    </div>
  );
}
