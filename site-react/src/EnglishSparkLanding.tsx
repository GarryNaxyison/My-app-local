import { useState } from "react";
import { motion } from "framer-motion";
import {
  ArrowRight,
  BrainCircuit,
  Camera,
  CheckCircle,
  ChevronRight,
  Laptop,
  MessageCircle,
  Mic,
  Repeat2,
  Smartphone,
  Zap,
} from "lucide-react";
import { GenerativeArtScene } from "@/components/ui/anomalous-matter-hero";

const WEB_APP_HREF = "/app/";
const TELEGRAM_HREF = "https://t.me/poliglot_ai_bot";

const heroGoals = [
  ["Travel", "hotel, cafe, airport"],
  ["Work", "calls, emails, small talk"],
  ["Exam", "vocabulary, grammar, speaking"],
  ["Conversation", "live answers without memorizing"],
] as const;

const learningLanguageShowcase = [
  ["English", "A1-C2", "speaking routes"],
  ["Espanol", "A1-C2", "travel dialogues"],
  ["Deutsch", "A1-C2", "work and exams"],
  ["Francais", "A1-C2", "listening and speech"],
  ["Italiano", "A1-C2", "real-life scenes"],
  ["Chinese", "A1-C2", "words and practice"],
  ["Japanese", "A1-C2", "lessons and review"],
  ["Korean", "A1-C2", "vocabulary and speech"],
] as const;

const demoScreens = [
  {
    tab: "Lesson",
    avatar: "AI",
    intro: "Poliglot AI gives the phrase, meaning, example, and one short next step.",
    label: "Lesson",
    output: "I would like to book a table for tonight.",
    hint: "The phrase moves straight into your route and review loop.",
  },
  {
    tab: "Dialogue",
    avatar: "DL",
    intro: "The role keeps context and asks you to answer like you would in the real situation.",
    label: "Hotel check-in",
    output: "Could you help me check in? I have a reservation.",
    hint: "AI corrects the answer and gives a more natural version.",
  },
  {
    tab: "Voice",
    avatar: "VO",
    intro: "Say the phrase, get a score, weak words, and a repeat prompt for shadowing.",
    label: "Voice",
    output: "Please speak a little slower.",
    hint: "Voice history shows what became easier after repetition.",
  },
  {
    tab: "Photo",
    avatar: "PH",
    intro: "Photograph a menu, sign, or task and turn it into a learning scenario.",
    label: "Photo",
    output: "No peanuts, please. How spicy is this dish?",
    hint: "The photo becomes translation, a note, and a practice prompt.",
  },
] as const;

const demoTabs = demoScreens.map((screen) => screen.tab);

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
  {
    name: "Anna",
    role: "Travel practice",
    avatar: "/assets/testimonials/anna.jpg",
    text: "I rehearsed check-in, cafe orders, and transport before the trip. The phrases stayed in notes for quick review.",
  },
  {
    name: "Marat",
    role: "English for work",
    avatar: "/assets/testimonials/marat.jpg",
    text: "I use the web app for longer lessons and Telegram for weak words before calls. The same profile keeps it simple.",
  },
  {
    name: "Sofia",
    role: "Pronunciation",
    avatar: "/assets/testimonials/sofia.jpg",
    text: "Voice practice shows which words sound weak and gives a better sentence to repeat right away.",
  },
] as const;

export function EnglishSparkLanding() {
  const [activeTab, setActiveTab] = useState(0);
  const activeDemo = demoScreens[activeTab] ?? demoScreens[0];
  const activeLabel = activeDemo.tab;

  return (
    <main className="english-spark-landing" data-no-translate="">
      <section className="landing-hero">
        <div className="landing-hero__matter">
          <GenerativeArtScene animate color="#7bdcff" particleColor="#f5d27a" />
        </div>
        <div className="landing-hero__veil" aria-hidden="true" />

        <div className="landing-hero__inner">
          <motion.div initial={{ opacity: 0, y: 24 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.66 }} className="landing-hero__copy">
            <span className="eyebrow">Premium AI tutor for daily practice</span>
            <h1>Practice speaking before the moment matters</h1>
            <p className="hero-lead">Lessons, dialogue, pronunciation, photo translation, mistake review, and progress live in one profile across web, PWA, and Telegram.</p>
            <div className="hero-actions">
              <a className="entry-cta hero-action hero-action--primary" data-entry="web-app" href={WEB_APP_HREF}>
                Start free <ArrowRight size={18} />
              </a>
              <a className="entry-cta hero-action hero-action--secondary" data-entry="telegram" href={TELEGRAM_HREF}>
                Open Telegram
              </a>
            </div>
            <div className="hero-goals" aria-label="Learning goals">
              {heroGoals.map(([goal, note]) => (
                <a key={goal} href={WEB_APP_HREF}>
                  <strong>{goal}</strong>
                  <span>{note}</span>
                </a>
              ))}
            </div>
            <div className="hero-language-picker" aria-label="Popular learning languages">
              <h3>Popular learning languages:</h3>
              <div>
                {learningLanguageShowcase.map(([language, level, note]) => (
                  <a key={language} href={WEB_APP_HREF}>
                    <strong>{language}</strong>
                    <span>{level}</span>
                    <small>{note}</small>
                  </a>
                ))}
              </div>
            </div>
            <div className="hero-path" aria-label="Learning route">
              <span>
                <b>1</b> Goal
              </span>
              <span>
                <b>2</b> Lesson
              </span>
              <span>
                <b>3</b> Speech
              </span>
              <span>
                <b>4</b> Review
              </span>
            </div>
            <div className="hero-proof">
              <span>
                <strong>35 languages</strong>
                interface coverage
              </span>
              <span>
                <strong>A1-C2</strong>
                learning levels
              </span>
              <span>
                <strong>Free</strong>
                start without payment
              </span>
            </div>
          </motion.div>

          <motion.div initial={{ opacity: 0, x: 38 }} animate={{ opacity: 1, x: 0 }} transition={{ duration: 0.72, delay: 0.08 }} className="hero-demo">
            <div className="hero-demo__top">
              <div>
                <span>Today</span>
                <strong>AI Tutor Cockpit</strong>
              </div>
              <span className="hero-demo__score">84/100</span>
            </div>
            <div className="hero-demo__tabs" role="tablist" aria-label="Feature demo">
              {demoTabs.map((tab, index) => (
                <button key={tab} type="button" className={activeTab === index ? "is-active" : undefined} onClick={() => setActiveTab(index)}>
                  {tab}
                </button>
              ))}
            </div>
            <motion.div key={activeLabel} initial={{ opacity: 0, y: 12 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.22 }} className="hero-demo__panel">
              <div className="demo-message">
                <div className="demo-avatar">{activeDemo.avatar}</div>
                <div>
                  <strong>Poliglot AI</strong>
                  <p>{activeDemo.intro}</p>
                </div>
              </div>
              <div className="demo-output">
                <span>{activeDemo.label}</span>
                <p>{activeDemo.output}</p>
              </div>
              <div className="demo-wave" aria-hidden="true">
                {Array.from({ length: 22 }, (_, index) => (
                  <i key={index} style={{ animationDelay: `${index * 0.035}s` }} />
                ))}
              </div>
              <div className="demo-hint">
                <CheckCircle size={17} />
                <span>{activeDemo.hint}</span>
              </div>
            </motion.div>
            <div className="hero-demo__skeletons" aria-label="Stable loading states">
              <div className="skeleton-card" aria-label="Stable route preview">
                <span className="skeleton-line" />
                <span className="skeleton-line" />
                <span className="skeleton-line" />
              </div>
              <div className="skeleton-card" aria-label="Stable review preview">
                <span className="skeleton-line" />
                <span className="skeleton-line" />
                <span className="skeleton-line" />
              </div>
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
          {reviews.map((review) => (
            <article className="review-card" key={review.name}>
              <div>
                <img src={review.avatar} alt={review.name} loading="lazy" />
                <div>
                  <strong>{review.name}</strong>
                  <small>{review.role}</small>
                </div>
              </div>
              <p>{review.text}</p>
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
