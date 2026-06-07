"use client";

import { Star } from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";

const testimonials = [
  {
    name: "Мария",
    role: "A2 English · 18 дней подряд",
    stars: 5,
    initials: "М",
    content:
      "После голосового ответа видно, где именно речь съехала. Не просто “молодец”, а нормальный совет, который можно повторить вслух.",
  },
  {
    name: "Тимур",
    role: "B1 Deutsch · 126 фраз",
    stars: 5,
    initials: "Т",
    content:
      "Диалоги стали привычкой: открыл Telegram, ответил на пару вопросов и сразу получил исправления.",
  },
  {
    name: "Анна",
    role: "Spanish A1 · 420 слов",
    stars: 5,
    initials: "А",
    content:
      "Фото-перевод выручает в поездке, а потом бот делает из этой ситуации практику.",
  },
];

export default function TestimonialSection() {
  return (
    <section id="reviews" className="bg-white py-24 text-slate-950">
      <div className="mx-auto w-full max-w-6xl px-6">
        <div className="max-w-3xl">
          <span className="text-sm font-bold uppercase tracking-[0.24em] text-blue-700">Отзывы</span>
          <h2 className="mt-4 text-3xl font-bold md:text-5xl">
            Пользователи видят прогресс, а не просто чат
          </h2>
          <p className="mt-5 text-lg leading-8 text-slate-600">
            Poliglot AI помогает заниматься коротко, но каждый день: уроки, голос, словарь и награды собираются в понятный маршрут.
          </p>
        </div>

        <div className="mt-10 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
          {testimonials.map((item) => (
            <article key={item.name} className="rounded-lg border border-slate-200 bg-white p-5 shadow-sm">
              <div className="flex gap-1" aria-label={`${item.stars} out of 5 stars`}>
                {Array.from({ length: 5 }).map((_, i) => (
                  <Star
                    key={i}
                    className={cn(
                      "h-4 w-4",
                      i < item.stars ? "fill-blue-600 stroke-blue-600" : "fill-slate-200 stroke-transparent",
                    )}
                  />
                ))}
              </div>
              <p className="my-5 leading-7 text-slate-700">{item.content}</p>
              <footer className="flex items-center gap-2">
                <Avatar className="h-8 w-8 border border-slate-200 shadow">
                  <AvatarFallback className="bg-slate-950 text-white">{item.initials}</AvatarFallback>
                </Avatar>
                <cite className="text-sm font-semibold not-italic text-slate-950">{item.name}</cite>
                <span aria-hidden className="h-1 w-1 rounded-full bg-slate-300" />
                <span className="text-sm text-slate-500">{item.role}</span>
              </footer>
            </article>
          ))}
        </div>
      </div>
    </section>
  );
}
