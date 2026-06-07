"use client";

import { useMemo, useRef, useState } from "react";
import { PauseCircle, PlayCircle, Volume2 } from "lucide-react";
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

export function AudioWaveButton({ label, text = "", wordId, targetLanguage, className, compact }: AudioWaveButtonProps) {
  const variants = useMemo(() => buildWaveVariants(), []);
  const waveHeights = useMemo(() => Array.from({ length: 30 }, () => Math.random() * 20 + 5), []);
  const audioRef = useRef<HTMLAudioElement | null>(null);
  const objectUrlRef = useRef<string | null>(null);
  const [playing, setPlaying] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const stop = () => {
    audioRef.current?.pause();
    if (audioRef.current) audioRef.current.currentTime = 0;
    setPlaying(false);
  };

  const resolveAudio = async () => {
    if (objectUrlRef.current) return objectUrlRef.current;
    const blob = wordId
      ? await apiBlob("/api/words/pronunciation", { word_id: wordId })
      : await apiBlob("/api/tools/translator-speech", { text, target_language: targetLanguage || "" });
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

  return (
    <button
      className={cn("audio-wave-button-v2", compact && "is-compact", playing && "is-playing", className)}
      type="button"
      onClick={() => void play()}
      disabled={loading || (!wordId && !text.trim())}
    >
      <span className="audio-wave-button-v2__icon">
        {playing ? <PauseCircle size={compact ? 24 : 32} fill="currentColor" /> : <PlayCircle size={compact ? 24 : 32} fill="currentColor" />}
      </span>
      <span className="audio-wave-button-v2__meta">
        <small><Volume2 size={12} />{error || label}</small>
        <span className="audio-wave-button-v2__waves" aria-hidden="true">
          {variants.map((variant, index) => (
            <motion.i
              key={index}
              style={{ height: `${waveHeights[index]}px` }}
              variants={variant}
              initial="initial"
              animate={playing ? "animate" : "initial"}
            />
          ))}
        </span>
      </span>
    </button>
  );
}
