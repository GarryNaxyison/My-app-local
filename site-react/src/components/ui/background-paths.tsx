"use client";

import type { ReactNode } from "react";
import { motion } from "framer-motion";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

function FloatingPaths({ position }: { position: number }) {
  const paths = Array.from({ length: 28 }, (_, i) => ({
    id: i,
    d: `M-${380 - i * 5 * position} -${189 + i * 6}C-${380 - i * 5 * position} -${189 + i * 6} -${
      312 - i * 5 * position
    } ${216 - i * 6} ${152 - i * 5 * position} ${343 - i * 6}C${616 - i * 5 * position} ${470 - i * 6} ${
      684 - i * 5 * position
    } ${875 - i * 6} ${684 - i * 5 * position} ${875 - i * 6}`,
    width: 0.5 + i * 0.03,
  }));

  return (
    <div className="absolute inset-0 pointer-events-none">
      <svg className="h-full w-full text-slate-950 dark:text-white" viewBox="0 0 696 316" fill="none">
        <title>Background Paths</title>
        {paths.map((path) => (
          <motion.path
            key={path.id}
            d={path.d}
            stroke="currentColor"
            strokeWidth={path.width}
            strokeOpacity={0.08 + path.id * 0.02}
            initial={{ pathLength: 0.3, opacity: 0.42 }}
            animate={{
              pathLength: 1,
              opacity: [0.18, 0.46, 0.18],
              pathOffset: [0, 1, 0],
            }}
            transition={{
              duration: 24 + (path.id % 8),
              repeat: Number.POSITIVE_INFINITY,
              ease: "linear",
            }}
          />
        ))}
      </svg>
    </div>
  );
}

export function BackgroundPaths({
  title = "Background Paths",
  lead,
  cta = "Discover Excellence",
  href = "#",
  children,
  className,
}: {
  title?: string;
  lead?: string;
  cta?: string;
  href?: string;
  children?: ReactNode;
  className?: string;
}) {
  const words = title.split(" ");

  return (
    <section className={cn("relative overflow-hidden bg-white text-slate-950", className)}>
      <div className="absolute inset-0">
        <FloatingPaths position={1} />
        <FloatingPaths position={-1} />
      </div>

      <div className="relative z-10 mx-auto grid min-h-[620px] max-w-7xl items-center gap-10 px-6 py-24 lg:grid-cols-[1fr_0.9fr]">
        <motion.div initial={{ opacity: 0 }} whileInView={{ opacity: 1 }} viewport={{ once: true }} transition={{ duration: 1 }}>
          <h2 className="max-w-3xl text-4xl font-bold sm:text-6xl">
            {words.map((word, wordIndex) => (
              <span key={wordIndex} className="mr-3 inline-block last:mr-0">
                {word.split("").map((letter, letterIndex) => (
                  <motion.span
                    key={`${wordIndex}-${letterIndex}`}
                    initial={{ y: 72, opacity: 0 }}
                    whileInView={{ y: 0, opacity: 1 }}
                    viewport={{ once: true }}
                    transition={{
                      delay: wordIndex * 0.08 + letterIndex * 0.02,
                      type: "spring",
                      stiffness: 150,
                      damping: 24,
                    }}
                    className="inline-block bg-gradient-to-r from-slate-950 to-slate-600 bg-clip-text text-transparent"
                  >
                    {letter}
                  </motion.span>
                ))}
              </span>
            ))}
          </h2>
          {lead && <p className="mt-6 max-w-2xl text-lg leading-8 text-slate-600">{lead}</p>}
          <div className="mt-8 inline-block rounded-lg bg-gradient-to-b from-black/10 to-white/10 p-px shadow-lg backdrop-blur-lg">
            <Button
              asChild
              variant="ghost"
              className="rounded-lg border border-black/10 bg-white/95 px-7 py-6 text-base font-semibold text-black hover:bg-white"
            >
              <a href={href}>
                <span className="opacity-90">{cta}</span>
                <span className="ml-3 opacity-70">→</span>
              </a>
            </Button>
          </div>
        </motion.div>
        {children && <div className="relative">{children}</div>}
      </div>
    </section>
  );
}
