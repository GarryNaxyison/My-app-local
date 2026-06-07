import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";
import { motion, useMotionTemplate, useMotionValue } from "motion/react";

import { cn } from "@/lib/utils";

const buttonVariants = cva(
  "group relative inline-flex items-center justify-center gap-2 overflow-hidden whitespace-nowrap rounded-xl text-sm font-semibold transition-all duration-300 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0",
  {
    variants: {
      variant: {
        default:
          "bg-gradient-to-br from-blue-600 via-blue-600 to-blue-500 text-white shadow-lg shadow-blue-600/25 before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/20 before:to-transparent before:opacity-0 before:transition-opacity hover:shadow-xl hover:shadow-blue-600/30 hover:before:opacity-100",
        secondary:
          "bg-gradient-to-br from-slate-100 via-white to-slate-100 text-slate-950 shadow-lg shadow-slate-900/10 before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/30 before:to-transparent before:opacity-0 before:transition-opacity hover:before:opacity-100 dark:from-white/12 dark:via-white/10 dark:to-white/5 dark:text-white dark:shadow-black/25",
        outline:
          "border-2 border-slate-300 bg-white/70 text-slate-950 shadow-sm shadow-black/5 backdrop-blur-sm hover:border-blue-400 hover:bg-white dark:border-white/15 dark:bg-white/5 dark:text-white dark:hover:bg-white/10",
        ghost:
          "backdrop-blur-sm hover:bg-black/5 hover:text-current dark:hover:bg-white/10",
        destructive:
          "bg-gradient-to-br from-rose-600 via-rose-600 to-rose-500 text-white shadow-lg shadow-rose-600/25 before:absolute before:inset-0 before:bg-gradient-to-br before:from-white/20 before:to-transparent before:opacity-0 before:transition-opacity hover:shadow-xl hover:shadow-rose-600/30 hover:before:opacity-100",
        approve:
          "bg-green-600/10 text-green-700 hover:bg-green-600/20 focus-visible:ring-green-600/20 dark:bg-green-400/10 dark:text-green-300 dark:hover:bg-green-400/20 dark:focus-visible:ring-green-400/40",
        reject:
          "bg-rose-600/10 text-rose-700 hover:bg-rose-600/20 focus-visible:ring-rose-600/20 dark:bg-rose-400/10 dark:text-rose-300 dark:hover:bg-rose-400/20 dark:focus-visible:ring-rose-400/40",
        shimmer:
          "animate-shimmer border border-slate-800/50 bg-[length:200%_100%] bg-[linear-gradient(110deg,#000103,45%,#1e2631,55%,#000103)] text-white shadow-2xl shadow-blue-600/30 before:absolute before:inset-0 before:translate-x-[-200%] before:animate-shimmer-slide before:bg-gradient-to-r before:from-transparent before:via-white/10 before:to-transparent hover:shadow-blue-500/50",
        glow:
          "animate-gradient-x bg-gradient-to-r from-violet-600 via-purple-600 to-fuchsia-600 text-white shadow-lg shadow-purple-500/50 before:absolute before:inset-[-2px] before:-z-10 before:rounded-xl before:bg-gradient-to-r before:from-violet-600 before:via-purple-600 before:to-fuchsia-600 before:opacity-75 before:blur-md hover:shadow-2xl hover:shadow-purple-500/60",
      },
      size: {
        default: "h-11 px-4",
        sm: "h-9 px-3 text-xs",
        lg: "h-12 px-6",
        icon: "size-11",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  },
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, children, ...props }, ref) => {
    const Comp = asChild ? Slot : "button";
    const mouseX = useMotionValue(0);
    const mouseY = useMotionValue(0);
    const background = useMotionTemplate`radial-gradient(circle at ${mouseX}px ${mouseY}px, rgba(255,255,255,0.18), transparent 80%)`;

    if (asChild) {
      return (
        <Comp className={cn(buttonVariants({ variant, size, className }))} ref={ref} {...props}>
          {children}
        </Comp>
      );
    }

    return (
      <motion.button
        className={cn(buttonVariants({ variant, size, className }))}
        ref={ref}
        onMouseMove={(event) => {
          const { left, top } = event.currentTarget.getBoundingClientRect();
          mouseX.set(event.clientX - left);
          mouseY.set(event.clientY - top);
        }}
        whileHover={{ scale: variant === "ghost" ? 1.01 : 1.02 }}
        whileTap={{ scale: 0.98 }}
        transition={{ type: "spring", stiffness: 400, damping: 17 }}
        {...(props as any)}
      >
        {variant !== "ghost" ? (
          <motion.span
            aria-hidden="true"
            className="pointer-events-none absolute inset-0 opacity-0 transition-opacity duration-300 group-hover:opacity-100"
            style={{ background }}
          />
        ) : null}
        <span className="relative z-10 flex items-center justify-center gap-2">{children}</span>
      </motion.button>
    );
  },
);
Button.displayName = "Button";

export { Button, buttonVariants };
