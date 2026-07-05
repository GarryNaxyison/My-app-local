"use client";

import { useEffect, useMemo, useRef } from "react";
import { motion } from "motion/react";
import { Link2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { AudioWaveButton } from "@/components/ui/audio-wave-button";

export interface ChatPerson {
  name: string;
  avatar?: string;
}

export interface ChatMessageItem {
  id: string | number;
  sender: "left" | "right";
  type?: "text" | "text-with-links";
  content: string;
  attachments?: Array<{ type: "image"; url: string; name: string }>;
  links?: Array<{ text: string }>;
  audio?: Array<{ label: string; text: string; wordId?: string; targetLanguage?: string }>;
}

export interface ChatConfig {
  leftPerson: ChatPerson;
  rightPerson: ChatPerson;
  messages: ChatMessageItem[];
  targetLanguage?: string;
}

export interface UiConfig {
  className?: string;
  backgroundColor?: string;
  leftBubble?: string;
  rightBubble?: string;
}

export default function ChatComponent({ config, uiConfig = {} }: { config: ChatConfig; uiConfig?: UiConfig }) {
  const messages = useMemo(() => config.messages.slice(0, 12), [config.messages]);
  const scrollRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    const node = scrollRef.current;
    if (!node) return;
    node.scrollTo({ top: 0, behavior: "auto" });
  }, [messages[0]?.id]);

  return (
    <div
      className={cn("chat-interface rounded-3xl border border-black/10 bg-white/80 p-4 shadow-2xl shadow-slate-900/10 backdrop-blur-xl dark:border-white/10 dark:bg-slate-950/72", uiConfig.className)}
      style={{ backgroundColor: uiConfig.backgroundColor }}
    >
      <div className="chat-interface__scroll" ref={scrollRef}>
        {messages.map((message, index) => {
          const isLeft = message.sender === "left";
          const person = isLeft ? config.leftPerson : config.rightPerson;
          return (
            <motion.div
              className={cn("chat-interface__row", !isLeft && "chat-interface__row--right")}
              key={message.id}
              initial={{ opacity: 0, y: 12, scale: 0.98 }}
              animate={{ opacity: 1, y: 0, scale: 1 }}
              transition={{ delay: Math.min(index * 0.025, 0.16), duration: 0.22 }}
            >
              <div className="chat-interface__group">
                <div className={cn("chat-interface__meta", !isLeft && "chat-interface__meta--right")}>
                  <span>{person.name}</span>
                </div>
                <div className={cn("chat-interface__bubble", !isLeft && "chat-interface__bubble--right")}>
                  <p>{message.content}</p>
                  {message.attachments?.length ? (
                    <div className="chat-interface__attachments">
                      {message.attachments.map((attachment) => (
                        attachment.type === "image" ? (
                          <figure className="chat-interface__image" key={`${attachment.url}-${attachment.name}`}>
                            <img src={attachment.url} alt={attachment.name} />
                            <figcaption>{attachment.name}</figcaption>
                          </figure>
                        ) : null
                      ))}
                    </div>
                  ) : null}
                  {message.links?.length ? (
                    <div className="chat-interface__links">
                      {message.links.map((link) => (
                        <small key={link.text}><Link2 size={12} />{link.text}</small>
                      ))}
                    </div>
                  ) : null}
                  {message.audio?.length ? (
                    <div className="chat-interface__audio">
                      {message.audio.map((clip, clipIndex) => (
                        <AudioWaveButton
                          key={`${clip.label}-${clipIndex}`}
                          label={clip.label}
                          text={clip.text}
                          wordId={clip.wordId}
                          targetLanguage={clip.targetLanguage || config.targetLanguage}
                          compact
                        />
                      ))}
                    </div>
                  ) : null}
                </div>
              </div>
            </motion.div>
          );
        })}
      </div>
    </div>
  );
}
