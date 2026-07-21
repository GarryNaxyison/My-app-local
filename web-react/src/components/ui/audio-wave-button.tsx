"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { Loader2, PauseCircle, PlayCircle, Volume2 } from "lucide-react";
import { motion, type Variants } from "framer-motion";
import { apiBlob } from "@/lib/api";
import { cn } from "@/lib/utils";

type AudioWaveButtonProps = {
  label: string;
  text?: string;
  wordId?: string;
  targetLanguage?: string;
  className?: string;
  compact?: boolean;
  showText?: boolean;
};

function buildWaveVariants(): Variants[] {
  return Array.from({ length: 30 }, () => ({
    initial: {
      scaleY: 1.5,
      transition: { duration: 0.5 },
    },
    animate: {
      scaleY: [1, Math.random() * 1.2 + 1, 1],
      transition: {
        duration: Math.random() * 0.5 + 0.5,
        repeat: Infinity,
        ease: "easeInOut",
        delay: Math.random() * 0.5,
      },
    },
  }));
}

export function AudioWaveButton({ label, text = "", wordId, targetLanguage, className, compact, showText = true }: AudioWaveButtonProps) {
  const variants = useMemo(() => buildWaveVariants(), []);
  const waveHeights = useMemo(() => Array.from({ length: 30 }, () => Math.random() * 20 + 5), []);
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const objectUrlRef = useRef<string | null>(null);
  const [playing, setPlaying] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    audioRef.current?.pause();
    audioRef.current = null;
    setPlaying(false);
    setLoading(false);
    setError("");
    if (objectUrlRef.current) {
      URL.revokeObjectURL(objectUrlRef.current);
      objectUrlRef.current = null;
    }
  }, [text, wordId, targetLanguage]);

  const stop = () => {
    audioRef.current?.pause();
    if (audioRef.current) audioRef.current.currentTime = 0;
    setPlaying(false);
  };

  const resolveAudio = async () => {
    if (objectUrlRef.current) return objectUrlRef.current;
    let blob: Blob;
    if (wordId) {
      try {
        blob = await apiBlob("/api/words/pronunciation", { word_id: wordId });
      } catch (reason) {
        if (!text.trim()) throw reason;
        blob = await apiBlob("/api/tools/translator-speech", { text, target_language: targetLanguage || "" });
      }
    } else {
      blob = await apiBlob("/api/tools/translator-speech", { text, target_language: targetLanguage || "" });
    }
    const url = URL.createObjectURL(blob);
    objectUrlRef.current = url;
    return url;
  };

  const play = async () => {
    if (playing) {
      stop();
      return;
    }
    if (!wordId && !text.trim()) return;
    setLoading(true);
    setError("");
    try {
      const url = await resolveAudio();
      const audio = audioRef.current || new Audio();
      audioRef.current = audio;
      audio.src = url;
      audio.onended = () => setPlaying(false);
      await audio.play();
      setPlaying(true);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Audio is unavailable.");
    } finally {
      setLoading(false);
    }
  };
  const visibleLabel = error ? "Retry audio" : loading ? "Loading..." : label;
  const visibleText = showText ? text.trim() : "";

  return (
    <button
      className={cn("audio-wave-button-v2", compact && "is-compact", playing && "is-playing", loading && "is-loading", className)}
      type="button"
      onClick={() => void play()}
      disabled={loading || (!wordId && !text.trim())}
      aria-busy={loading}
      title={error || label}
    >
      <span className="audio-wave-button-v2__icon">
        {loading ? <Loader2 className="audio-wave-button-v2__loader" size={compact ? 22 : 30} /> : playing ? <PauseCircle size={compact ? 24 : 32} fill="currentColor" /> : <PlayCircle size={compact ? 24 : 32} fill="currentColor" />}
      </span>
      <span className="audio-wave-button-v2__meta">
        <small><Volume2 size={12} />{visibleLabel}</small>
        {error ? <em className="audio-wave-button-v2__error">{error}</em> : null}
        {visibleText ? <span className="audio-wave-button-v2__text">{visibleText}</span> : null}
        <span className="audio-wave-button-v2__waves" aria-hidden="true">
          {variants.map((variant, index) => (
            <motion.i
              key={index}
              style={{ height: `${waveHeights[index]}px` }}
              variants={variant}
              initial="initial"
              animate={playing || loading ? "animate" : "initial"}
            />
          ))}
        </span>
      </span>
    </button>
  );
}
