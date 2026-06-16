import { motion } from "framer-motion";
import {
  ArrowRight,
  BookOpen,
  BrainCircuit,
  Camera,
  CheckCircle,
  ChevronRight,
  Headphones,
  Laptop,
  MessageCircle,
  Mic,
  Repeat2,
  ScanText,
  Smartphone,
  Zap,
} from "lucide-react";
import { SparklesCore } from "@/components/ui/sparkles";

const WEB_APP_HREF = "/app/";
const TELEGRAM_HREF = "https://t.me/poliglot_ai_bot";

const proofChips = ["A1-C2 routes", "35 languages later", "Web + Telegram profile", "Voice and photo practice"] as const;

const dailySteps = [
  ["01", "Learn the phrase", "Get the meaning, grammar hint, natural example, and one focused prompt."],
  ["02", "Use it in context", "Practice through a short roleplay, voice answer, or photo-based task."],
  ["03", "Review the weak spot", "Mistakes, notes, weak words, XP, and streak point to the next repetition."],
] as const;

const modules = [
  {
    icon: BrainCircuit,
    title: "AI Tutor",
    body: "Guided lessons turn a phrase into an answer, a correction, and a next action.",
  },
  {
    icon: MessageCircle,
    title: "Roleplay",
    body: "Travel, work, exam, or casual speaking scenarios stay on topic and correct your reply.",
  },
  {
    icon: Mic,
    title: "Voice Coach",
    body: "Pronunciation practice highlights weak words, rhythm, and the next sentence to repeat.",
  },
  {
    icon: Camera,
    title: "Photo Practice",
    body: "A menu or sign becomes translation, context, notes, and a short practice loop.",
  },
] as const;

const memoryNodes = ["Mistakes", "Notes", "Weak words", "XP", "Streak", "Offline decks"] as const;

type LandingPlan = {
  name: string;
  label: string;
  oldPrice?: string;
  price: string;
  body: string;
  limits: readonly string[];
  featured?: boolean;
};

const plans: readonly LandingPlan[] = [
  {
    name: "Free",
    label: "Try the loop",
    price: "0 ₽",
    body: "Basic text practice for trying lessons, word training, notes, and progress without payment.",
    limits: ["5 lessons per day", "15 practice messages", "Basic notes and phrasebook"],
  },
  {
    name: "Premium",
    label: "Daily speaking plan",
    oldPrice: "1000 ₽",
    price: "300 ₽",
    body: "The main plan for daily progress with AI Tutor, roleplay, pronunciation, voice checks, and photo tools.",
    limits: ["50 lessons per day", "200 practice messages", "20 voice checks up to 30 seconds"],
    featured: true,
  },
  {
    name: "Platinum",
    label: "Intensive preparation",
    oldPrice: "2000 ₽",
    price: "590 ₽",
    body: "Higher limits for travel, work, exam preparation, and longer AI-dialogue sessions.",
    limits: ["100 lessons per day", "500 practice messages", "60 voice checks up to 30 seconds"],
  },
] as const;

const reviews = [
  ["Anna", "Travel practice", "I rehearsed check-in, cafe orders, and transport before the trip. The phrases stayed in notes for quick review."],
  ["Marat", "English for work", "I use the web app for longer lessons and Telegram for weak words before calls. The same profile keeps it simple."],
  ["Sofia", "Pronunciation", "Voice practice shows which words sound weak and gives a better sentence to repeat right away."],
] as const;

export function EnglishSparkLanding() {
  return (
    <main className="english-spark-landing">
      <section className="spark-hero">
        <div className="spark-hero__sparkles" aria-hidden="true">
          <SparklesCore background="transparent" minSize={0.35} maxSize={1.1} particleDensity={115} particleColor="#ffffff" speed={0.75} className="h-full w-full" />
        </div>
        <div className="spark-hero__veil" aria-hidden="true" />
        <div className="spark-hero__content">
          <motion.div className="spark-hero__copy" initial={{ opacity: 0, y: 28 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.64 }}>
            <span className="eyebrow">AI Tutor for real situations</span>
            <h1>Practice speaking before the moment matters</h1>
            <p className="hero-lead">Lessons, roleplay, voice, photos, and review loops work in one learning profile across web and Telegram.</p>
            <div className="hero-actions">
              <a className="entry-cta hero-action hero-action--primary" data-entry="web-app" href={WEB_APP_HREF}>
                Open Web app <ArrowRight size={18} />
              </a>
              <a className="entry-cta hero-action hero-action--secondary" data-entry="telegram" href={TELEGRAM_HREF}>
                Open Telegram <MessageCircle size={18} />
              </a>
            </div>
            <div className="spark-hero__proof" aria-label="Product proof">
              {proofChips.map((chip) => (
                <span key={chip}>{chip}</span>
              ))}
            </div>
          </motion.div>

          <motion.div className="hero-product__card" initial={{ opacity: 0, x: 32 }} animate={{ opacity: 1, x: 0 }} transition={{ duration: 0.66, delay: 0.08 }}>
            <div className="hero-product__top">
              <span>Live lesson loop</span>
              <strong>84%</strong>
            </div>
            <div className="hero-product__lesson">
              <p>
                Could you help me check in?
                <span>meaning, grammar, next reply</span>
              </p>
              <div>
                <strong>+12 XP</strong>
                <span>mistake saved</span>
              </div>
            </div>
            <div className="skeleton-card" aria-label="Stable loading preview">
              <span className="skeleton-line" />
              <span className="skeleton-line" />
              <span className="skeleton-line" />
            </div>
            <div className="skeleton-card" aria-label="Stable review preview">
              <span className="skeleton-line" />
              <span className="skeleton-line" />
              <span className="skeleton-line" />
            </div>
            <div className="hero-product__voice" aria-hidden="true">
              {Array.from({ length: 12 }, (_, index) => (
                <i key={index} style={{ height: `${18 + (index % 5) * 8}px` }} />
              ))}
            </div>
          </motion.div>
        </div>
      </section>

      <section className="spark-section daily-loop-section">
        <div className="section-copy">
          <span className="eyebrow">Daily loop</span>
          <h2>One short session always ends with the next useful step</h2>
          <p>Poliglot AI is built around a simple loop: learn, use, review. Every module feeds the same progress profile.</p>
        </div>
        <div className="daily-steps">
          {dailySteps.map(([num, title, body]) => (
            <article className="daily-step" key={title}>
              <strong>{num}</strong>
              <h3>{title}</h3>
              <p>{body}</p>
            </article>
          ))}
        </div>
      </section>

      <section id="features" className="spark-section modules-section">
        <div className="section-copy">
          <span className="eyebrow">Product modules</span>
          <h2>Lessons, conversation, voice, and photos work as one system</h2>
        </div>
        <div className="module-grid">
          {modules.map((module) => {
            const Icon = module.icon;
            return (
              <article className="module-card" key={module.title}>
                <Icon size={24} />
                <h3>{module.title}</h3>
                <p>{module.body}</p>
              </article>
            );
          })}
        </div>
      </section>

      <section className="memory-loop-section">
        <div className="memory-loop-section__copy">
          <span className="eyebrow">Review memory</span>
          <h2>Your weak spots stay visible until they become easy</h2>
          <p>Review is not a separate folder. It is the connective tissue between lessons, voice, notes, and Telegram practice.</p>
        </div>
        <div className="memory-grid">
          {memoryNodes.map((node) => (
            <span className="memory-node" key={node}>
              <Repeat2 size={16} />
              {node}
            </span>
          ))}
        </div>
      </section>

      <section className="spark-section entry-section">
        <div className="section-copy">
          <span className="eyebrow">Two entry points</span>
          <h2>Use the same profile for deep sessions and quick practice</h2>
        </div>
        <div className="entry-grid">
          <article className="entry-panel">
            <Laptop size={26} />
            <h3>Web app</h3>
            <p>Open the dashboard for deep sessions, progress, pricing, mistakes, and longer AI Tutor work.</p>
            <a className="entry-cta" data-entry="web-app" href={WEB_APP_HREF}>
              Open Web app <ChevronRight size={18} />
            </a>
          </article>
          <article className="entry-panel">
            <Smartphone size={26} />
            <h3>Telegram</h3>
            <p>Start quick practice, voice checks, reminders, and weak-word review from the same profile.</p>
            <a className="entry-cta" data-entry="telegram" href={TELEGRAM_HREF}>
              Open Telegram <ChevronRight size={18} />
            </a>
          </article>
        </div>
      </section>

      <section id="pricing" className="pricing-section">
        <div className="section-copy">
          <span className="eyebrow">Pricing</span>
          <h2>Start free, then unlock the daily practice limits you actually need</h2>
          <p className="pricing-lead">Payment is available through Telegram Stars, YooKassa/SBP, TON, and USDT.</p>
        </div>
        <div className="pricing-grid">
          {plans.map((plan) => (
            <article className={plan.featured ? "plan-card is-featured" : "plan-card"} key={plan.name}>
              <span className="plan-label">{plan.label}</span>
              <h3>{plan.name}</h3>
              <div className="price-line">
                {plan.oldPrice ? <s>{plan.oldPrice}</s> : null}
                <strong>{plan.price}</strong>
                {plan.oldPrice ? <small>per month</small> : <small>starter access</small>}
              </div>
              <p>{plan.body}</p>
              <ul>
                {plan.limits.map((item) => (
                  <li key={item}>
                    <CheckCircle size={16} />
                    {item}
                  </li>
                ))}
              </ul>
              <a href={WEB_APP_HREF}>Choose plan</a>
            </article>
          ))}
        </div>
        <div className="payment-methods" aria-label="Payment methods">
          <span>Telegram Stars</span>
          <span>YooKassa/SBP</span>
          <span>TON</span>
          <span>USDT</span>
        </div>
      </section>

      <section id="reviews" className="spark-section reviews-section">
        <div className="section-copy">
          <span className="eyebrow">User stories</span>
          <h2>Built for concrete speaking moments</h2>
        </div>
        <div className="reviews-grid">
          {reviews.map(([name, role, text]) => (
            <article className="review-card" key={name}>
              <div>
                <span>{name.slice(0, 1)}</span>
                <div>
                  <strong>{name}</strong>
                  <small>{role}</small>
                </div>
              </div>
              <p>{text}</p>
            </article>
          ))}
        </div>
      </section>

      <section className="final-cta-section">
        <Zap size={34} />
        <h2>Open the loop and run the first lesson today</h2>
        <p>Use the web app for a full session or Telegram for a fast practice check.</p>
        <div>
          <a className="entry-cta" data-entry="web-app" href={WEB_APP_HREF}>
            Web app <ArrowRight size={18} />
          </a>
          <a className="entry-cta" data-entry="telegram" href={TELEGRAM_HREF}>
            Telegram <ArrowRight size={18} />
          </a>
        </div>
      </section>
    </main>
  );
}
