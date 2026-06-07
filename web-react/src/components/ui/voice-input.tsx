"use client";

import React from "react";
import { Mic } from "lucide-react";
import { AnimatePresence, motion } from "motion/react";

import { cn } from "@/lib/utils";

interface VoiceInputProps {
  isListening?: boolean;
  onStart?: () => void;
  onStop?: () => void;
  onToggle?: () => void;
}

export function VoiceInput({
  className,
  isListening,
  onStart,
  onStop,
  onToggle,
}: React.ComponentProps<"div"> & VoiceInputProps) {
  const [internalListening, setInternalListening] = React.useState(false);
  const [time, setTime] = React.useState(0);
  const listening = typeof isListening === "boolean" ? isListening : internalListening;
  const previousRef = React.useRef(listening);

  React.useEffect(() => {
    let intervalId: ReturnType<typeof setInterval> | undefined;

    if (listening) {
      if (!previousRef.current) onStart?.();
      intervalId = setInterval(() => {
        setTime((t) => t + 1);
      }, 1000);
    } else {
      if (previousRef.current) onStop?.();
      setTime(0);
    }

    previousRef.current = listening;
    return () => {
      if (intervalId) clearInterval(intervalId);
    };
  }, [listening, onStart, onStop]);

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins.toString().padStart(2, "0")}:${secs.toString().padStart(2, "0")}`;
  };

  const onClickHandler = () => {
    if (onToggle) {
      onToggle();
      return;
    }
    setInternalListening((current) => !current);
  };

  return (
    <div className={cn("flex flex-col items-center justify-center", className)}>
      <motion.div
        className="voice-input-shell flex cursor-pointer items-center justify-center rounded-full border p-2"
        layout
        transition={{ layout: { duration: 0.4 } }}
        onClick={onClickHandler}
        role="button"
        tabIndex={0}
        onKeyDown={(event) => {
          if (event.key === "Enter" || event.key === " ") onClickHandler();
        }}
      >
        <div className="flex h-6 w-6 items-center justify-center">
          {listening ? (
            <motion.div
              className="h-4 w-4 rounded-sm bg-primary"
              animate={{ rotate: [0, 180, 360] }}
              transition={{ duration: 2, repeat: Number.POSITIVE_INFINITY, ease: "easeInOut" }}
            />
          ) : (
            <Mic />
          )}
        </div>
        <AnimatePresence mode="wait">
          {listening && (
            <motion.div
              initial={{ opacity: 0, width: 0, marginLeft: 0 }}
              animate={{ opacity: 1, width: "auto", marginLeft: 8 }}
              exit={{ opacity: 0, width: 0, marginLeft: 0 }}
              transition={{ duration: 0.4 }}
              className="flex items-center justify-center gap-2 overflow-hidden"
            >
              <div className="flex items-center justify-center gap-0.5">
                {[...Array(12)].map((_, i) => (
                  <motion.div
                    key={i}
                    className="w-0.5 rounded-full bg-primary"
                    initial={{ height: 2 }}
                    animate={{ height: listening ? [2, 3 + Math.random() * 10, 3 + Math.random() * 5, 2] : 2 }}
                    transition={{
                      duration: listening ? 1 : 0.3,
                      repeat: listening ? Infinity : 0,
                      delay: listening ? i * 0.05 : 0,
                      ease: "easeInOut",
                    }}
                  />
                ))}
              </div>
              <div className="w-10 text-center text-xs text-muted-foreground">{formatTime(time)}</div>
            </motion.div>
          )}
        </AnimatePresence>
      </motion.div>
    </div>
  );
}
