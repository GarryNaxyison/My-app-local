"use client";

import { Check, ShieldCheck } from "lucide-react";
import { motion } from "framer-motion";
import { Badge } from "./badge";
import { Button } from "./button";
import { BorderTrail } from "./border-trail";
import { cn } from "@/lib/utils";

const plans = [
  {
    name: "Free",
    badge: "Для старта",
    description: "Познакомиться с ботом и выработать привычку.",
    price: "0 ₽",
    period: "навсегда",
    cta: "Начать бесплатно",
    features: ["5 уроков в день", "15 сообщений практики", "текстовый переводчик"],
  },
  {
    name: "Premium",
    badge: "Лучший выбор",
    description: "Больше уроков, голос, фото и расширенная практика.",
    oldPrice: "1000 ₽",
    price: "300 ₽",
    period: "в месяц",
    cta: "Подключить Premium",
    features: ["50 уроков в день", "200 сообщений практики", "20 голосовых в день"],
    highlight: true,
  },
  {
    name: "Platinum",
    badge: "Максимум",
    description: "Для поездок, экзаменов и учебных рывков.",
    oldPrice: "2000 ₽",
    price: "590 ₽",
    period: "в месяц",
    cta: "Выбрать Platinum",
    features: ["100 уроков в день", "500 сообщений практики", "60 голосовых в день"],
  },
];

export function Pricing() {
  return (
    <section id="pricing" className="relative overflow-hidden bg-slate-950 py-24 text-white">
      <div className="pointer-events-none absolute inset-0 bg-[linear-gradient(to_right,rgba(255,255,255,.08)_1px,transparent_1px),linear-gradient(to_bottom,rgba(255,255,255,.08)_1px,transparent_1px)] bg-[size:32px_32px] [mask-image:radial-gradient(ellipse_at_center,black_12%,transparent_70%)]" />
      <div className="relative mx-auto w-full max-w-6xl space-y-10 px-6">
        <motion.div
          initial={{ opacity: 0, y: 18 }}
          whileInView={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.65 }}
          viewport={{ once: true }}
          className="mx-auto max-w-2xl text-center"
        >
          <div className="mx-auto inline-flex rounded-lg border border-white/15 px-4 py-1 font-mono text-sm text-white/80">
            Тарифы
          </div>
          <h2 className="mt-5 text-3xl font-bold md:text-5xl">
            Цены видны сразу, лимиты без мелкого шрифта
          </h2>
          <p className="mt-5 text-base leading-7 text-white/65">
            Free подходит для знакомства. Premium — для регулярной практики. Platinum — для интенсивного режима.
          </p>
        </motion.div>

        <div className="grid gap-4 md:grid-cols-3">
          {plans.map((plan, index) => (
            <motion.article
              key={plan.name}
              initial={{ opacity: 0, y: 24 }}
              whileInView={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.65, delay: index * 0.08 }}
              viewport={{ once: true }}
              className={cn(
                "relative overflow-hidden rounded-lg border bg-slate-900 p-5 shadow-2xl",
                plan.highlight ? "border-blue-400/60" : "border-white/10",
              )}
            >
              {plan.highlight && (
                <BorderTrail
                  className="bg-white"
                  size={120}
                  style={{
                    boxShadow:
                      "0 0 46px 18px rgb(255 255 255 / 32%), 0 0 90px 42px rgb(59 130 246 / 28%)",
                  }}
                />
              )}
              <div className="relative">
                <div className="flex items-center justify-between gap-3">
                  <h3 className="text-xl font-semibold">{plan.name}</h3>
                  <Badge variant={plan.highlight ? "default" : "secondary"}>{plan.badge}</Badge>
                </div>
                <p className="mt-3 min-h-12 text-sm leading-6 text-white/62">{plan.description}</p>
                <div className="mt-8 flex items-end gap-2">
                  {plan.oldPrice && <span className="mb-1 text-sm text-white/38 line-through">{plan.oldPrice}</span>}
                  <span className="text-4xl font-extrabold">{plan.price}</span>
                  <span className="mb-1 text-sm text-white/55">{plan.period}</span>
                </div>
                <Button asChild variant={plan.highlight ? "default" : "outline"} className="mt-6 w-full">
                  <a href="/app">{plan.cta}</a>
                </Button>
                <ul className="mt-6 space-y-3 border-t border-white/10 pt-5">
                  {plan.features.map((feature) => (
                    <li key={feature} className="flex items-center gap-2 text-sm text-white/72">
                      <Check className="h-4 w-4 text-blue-300" />
                      {feature}
                    </li>
                  ))}
                </ul>
              </div>
            </motion.article>
          ))}
        </div>

        <div className="flex items-center justify-center gap-2 text-sm text-white/62">
          <ShieldCheck className="h-4 w-4" />
          <span>Лимиты: Free, Premium и Platinum синхронизируются с web app и Telegram.</span>
        </div>

        <div className="overflow-hidden rounded-lg border border-white/10 bg-white/[0.03]">
          <div className="grid grid-cols-4 border-b border-white/10 bg-white/[0.04] text-sm font-bold text-white">
            <span className="p-3">Лимит</span>
            <span className="p-3">Free</span>
            <span className="p-3">Premium</span>
            <span className="p-3">Platinum</span>
          </div>
          {[
            ["Уроки в день", "5", "50", "100"],
            ["Сообщения практики", "15", "200", "500"],
            ["Голосовые задания", "нет", "20", "60"],
          ].map((row) => (
            <div key={row[0]} className="grid grid-cols-4 border-b border-white/10 text-sm text-white/70 last:border-b-0">
              {row.map((cell, index) => (
                <span key={`${row[0]}-${index}`} className="p-3">
                  {cell}
                </span>
              ))}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
