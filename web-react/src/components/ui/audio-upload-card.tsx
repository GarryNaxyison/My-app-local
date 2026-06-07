"use client";

import React, { useCallback, useEffect, useMemo, useState } from "react";
import { AnimatePresence, motion } from "framer-motion";
import { FileAudio, Upload, X } from "lucide-react";
import { cn } from "@/lib/utils";

interface AudioUploadCardProps {
  className?: string;
  triggerAnimation?: boolean;
  onAnimationComplete?: () => void;
  title?: string;
  description?: string;
  file?: File | null;
  onFileChange?: (file: File | null) => void;
  compact?: boolean;
}

interface WaveformProps {
  width?: number;
  height?: number;
  bars?: number;
}

interface UploadCardBaseProps {
  children?: React.ReactNode;
  className?: string;
  isDragOver?: boolean;
  isUploading?: boolean;
}

const Waveform = ({ width = 300, height = 40, bars = 60 }: WaveformProps) => {
  const barsArray = useMemo(() => {
    const barWidth = (width / bars) * 0.5;
    const spacing = (width / bars) * 0.4;
    const centerY = height / 2;
    const cornerRadius = barWidth * 0.7;
    const totalBarsWidth = (barWidth + spacing) * bars - spacing;
    const startX = (width - totalBarsWidth) / 2;

    return Array.from({ length: bars }, (_, i) => {
      const x = startX + i * (barWidth + spacing);
      const barHeight = Math.random() * (height * 0.6) + height * 0.1;
      const topY = centerY - barHeight / 2;
      return (
        <rect
          key={i}
          x={x}
          y={topY}
          width={barWidth}
          height={barHeight}
          rx={cornerRadius}
          ry={cornerRadius}
          fill="currentColor"
          className="text-muted-foreground/60"
        />
      );
    });
  }, [bars, height, width]);

  return (
    <svg width={width} height={height} viewBox={`0 0 ${width} ${height}`}>
      {barsArray}
    </svg>
  );
};

const UploadCardBase = ({ children, className, isDragOver = false, isUploading = false }: UploadCardBaseProps) => {
  const hasChildren = React.Children.count(children) > 0;

  return (
    <div
      className={cn(
        "rounded-xl border-2 border-dashed border-border/60 p-6 backdrop-blur-sm min-h-[120px] flex items-center justify-center relative transition-colors duration-200",
        !isUploading && "cursor-pointer hover:bg-accent/20",
        isUploading
          ? "bg-primary/20 border-primary/60"
          : isDragOver
            ? "bg-accent/40 border-accent/80 shadow-inner"
            : "bg-card",
        className,
      )}
    >
      {!hasChildren && (
        <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
          <Upload
            size={48}
            className={cn(
              "transition-colors duration-200",
              isDragOver ? "text-primary" : "text-muted",
              isUploading && "text-primary",
            )}
          />
        </div>
      )}
      <div className="relative z-10 w-full">{children}</div>
    </div>
  );
};

interface AudioComponentProps {
  isAnimating: boolean;
  onAnimationComplete?: () => void;
  filename?: string;
  onRemove?: () => void;
}

const truncateFilename = (filename: string, maxLength = 20) => {
  if (filename.length <= maxLength) return filename;
  const extension = filename.split(".").pop();
  const nameWithoutExt = filename.replace(`.${extension}`, "");
  const truncatedName = nameWithoutExt.substring(0, Math.max(1, maxLength - 3 - (extension?.length || 0)));
  return `${truncatedName}...${extension}`;
};

const AudioComponent = ({ isAnimating, onAnimationComplete, filename = "audio.mp3", onRemove }: AudioComponentProps) => {
  const [isRemoving, setIsRemoving] = useState(false);
  const [shouldShow, setShouldShow] = useState(false);

  useEffect(() => {
    if (isAnimating) setShouldShow(true);
  }, [isAnimating]);

  if (!shouldShow && !isRemoving) return null;

  const displayName = truncateFilename(filename);

  const handleRemoveComplete = () => {
    setShouldShow(false);
    setIsRemoving(false);
    onRemove?.();
  };

  return (
    <AnimatePresence>
      {shouldShow && (
        <motion.div
          className="absolute z-20"
          initial={{ right: 20, bottom: 20, opacity: 0 }}
          animate={
            isRemoving
              ? {
                  scale: 0,
                  opacity: 0,
                  filter: "blur(8px)",
                  transition: { duration: 0.4, ease: [0.23, 1, 0.32, 1] },
                }
              : {
                  left: "50%",
                  top: "calc(50% - 15px)",
                  x: "-50%",
                  y: "-50%",
                  opacity: 1,
                  transition: { duration: 0.6, ease: [0.23, 1, 0.32, 1] },
                }
          }
          exit={{ scale: 0, opacity: 0, filter: "blur(8px)", transition: { duration: 0.4, ease: [0.23, 1, 0.32, 1] } }}
          style={{ transformOrigin: "center" }}
          onAnimationComplete={isRemoving ? handleRemoveComplete : onAnimationComplete}
        >
          <motion.div
            initial={{ scale: 1.5 }}
            animate={
              isRemoving
                ? { scale: 0, transition: { duration: 0.4 } }
                : { scale: [1.5, 1.1, 1], transition: { duration: 0.8, ease: [0.68, -0.55, 0.265, 1.55] } }
            }
            className="rounded-lg border border-border/30 bg-muted px-2 py-1.5 backdrop-blur-sm shadow-lg relative group"
          >
            <button
              onClick={() => setIsRemoving(true)}
              className="absolute -top-2 -right-2 w-5 h-5 bg-destructive text-destructive-foreground rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-200 hover:scale-110 z-30"
              type="button"
            >
              <X size={12} />
            </button>
            <div className="w-full flex items-center justify-center">
              <Waveform width={180} height={32} bars={40} />
            </div>
            <div>
              <span className="text-xs text-foreground/60 font-medium text-left">{displayName}</span>
            </div>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
};

export function AudioUploadCard({
  className,
  triggerAnimation = false,
  onAnimationComplete,
  title = "Upload Your Audio",
  description = "Drop in your recordings and start transcribing instantly.",
  file,
  onFileChange,
  compact = false,
}: AudioUploadCardProps) {
  const [isAnimating, setIsAnimating] = useState(false);
  const [isDragOver, setIsDragOver] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [uploadedFile, setUploadedFile] = useState<File | null>(file || null);
  const fileInputRef = React.useRef<HTMLInputElement>(null);

  useEffect(() => {
    setUploadedFile(file || null);
    if (file) setIsAnimating(true);
  }, [file]);

  useEffect(() => {
    if (triggerAnimation) setIsAnimating(true);
  }, [triggerAnimation]);

  const commitFile = useCallback(
    (audioFile: File) => {
      setUploadedFile(audioFile);
      onFileChange?.(audioFile);
      setIsUploading(true);
      setTimeout(() => {
        setIsUploading(false);
        setIsAnimating(true);
      }, 200);
    },
    [onFileChange],
  );

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragOver(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragOver(false);
  }, []);

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      setIsDragOver(false);
      const files = Array.from(e.dataTransfer.files);
      const audioFile = files.find((candidate) => candidate.type.startsWith("audio/"));
      if (audioFile) commitFile(audioFile);
    },
    [commitFile],
  );

  const handleFileSelect = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const selectedFile = e.target.files?.[0];
      if (selectedFile && selectedFile.type.startsWith("audio/")) commitFile(selectedFile);
    },
    [commitFile],
  );

  const handleRemoveFile = useCallback(() => {
    setUploadedFile(null);
    setIsAnimating(false);
    onFileChange?.(null);
    if (fileInputRef.current) fileInputRef.current.value = "";
  }, [onFileChange]);

  const handleBaseClick = useCallback(() => {
    if (!isUploading) {
      fileInputRef.current?.click();
    }
  }, [isUploading]);

  if (compact) {
    return (
      <motion.div
        className={cn("audio-upload-card-v2 audio-upload-card-v2--button", isDragOver && "is-dragging", className)}
        initial={{ opacity: 0, y: 6 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.22, ease: [0.23, 1, 0.32, 1] }}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
      >
        <button
          type="button"
          className={cn("attachment-button-v2 audio-attachment-button-v2", uploadedFile && "has-file", isUploading && "is-uploading")}
          onClick={handleBaseClick}
          aria-label={title}
        >
          {isUploading ? <Upload size={18} className="spin" /> : <FileAudio size={19} />}
        </button>
        <input ref={fileInputRef} type="file" accept="audio/*" onChange={handleFileSelect} className="sr-only" />
        <AnimatePresence>
          {uploadedFile ? (
            <motion.div
              className="file-chip-v2"
              initial={{ opacity: 0, x: -8, scale: 0.96 }}
              animate={{ opacity: 1, x: 0, scale: 1 }}
              exit={{ opacity: 0, x: -8, scale: 0.96 }}
              transition={{ duration: 0.24, ease: [0.23, 1, 0.32, 1] }}
            >
              <FileAudio size={15} />
              <span>{truncateFilename(uploadedFile.name, 24)}</span>
              <button type="button" onClick={handleRemoveFile} aria-label="Remove audio">
                <X size={13} />
              </button>
            </motion.div>
          ) : null}
        </AnimatePresence>
      </motion.div>
    );
  }

  return (
    <motion.div
      className={cn("relative w-full max-w-md mx-auto", className)}
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ type: "spring", stiffness: 300, damping: 30 }}
    >
      <div className="relative overflow-hidden rounded-xl border border-border/50 bg-card p-6 text-center">
        <div className="flex flex-col justify-center space-y-8">
          <div className="relative w-full mx-auto">
            <div
              className="relative"
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
              onClick={handleBaseClick}
            >
              <UploadCardBase isDragOver={isDragOver} isUploading={isUploading} />
              <input ref={fileInputRef} type="file" accept="audio/*" onChange={handleFileSelect} className="sr-only" />
            </div>
          </div>
          <div className="flex flex-col items-start">
            <h2 className="text-lg font-semibold text-foreground text-left">{title}</h2>
            <p className="text-sm text-muted-foreground text-left">{description}</p>
          </div>
        </div>
        <AudioComponent
          isAnimating={isAnimating}
          onAnimationComplete={onAnimationComplete}
          filename={uploadedFile?.name}
          onRemove={handleRemoveFile}
        />
      </div>
    </motion.div>
  );
}
