import { useState } from "react";
import { motion } from "framer-motion";
import {
  Activity,
  ArrowRight,
  BrainCircuit,
  Camera,
  CheckCircle,
  Languages,
  LineChart,
  MessageCircle,
  Mic,
  Zap,
} from "lucide-react";
import { GenerativeArtScene } from "@/components/ui/anomalous-matter-hero";
import { LandingSocialProof, SocialIconLinks } from "./components/SocialLinks";
import { landingSeoCopy } from "./landingSeoContent";

const WEB_APP_HREF = "/app/";
const TELEGRAM_HREF = "https://t.me/NERIVAapp_bot";

const demoScreens = [
  {
    tab: "Lesson",
    avatar: "AI",
    intro: "NERIVA gives the phrase, meaning, example, and one short next step.",
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

const shortFaq = landingSeoCopy.en.questions.slice(0, 4);

type EnglishSparkLandingProps = {
  siteTheme?: "dark" | "light";
};

export function EnglishSparkLanding({ siteTheme = "dark" }: EnglishSparkLandingProps) {
  const [activeTab, setActiveTab] = useState(0);
  const activeDemo = demoScreens[activeTab] ?? demoScreens[0];
  const activeLabel = activeDemo.tab;

  return (
    <main className="english-spark-landing">
      <section className="landing-hero" data-hero-preset={siteTheme}>
        <div className="landing-hero__matter">
          <GenerativeArtScene animate variant={siteTheme} color="#7bdcff" particleColor="#f5d27a" />
        </div>
        <div className="landing-hero__veil" aria-hidden="true" />

        <div className="landing-hero__inner">
          <motion.div initial={{ opacity: 0, y: 24 }} animate={{ opacity: 1, y: 0 }} transition={{ duration: 0.66 }} className="landing-hero__copy">
            <span className="eyebrow">Premium AI tutor for daily practice</span>
            <h1>NERIVA: Practice speaking before the moment matters</h1>
            <p className="hero-lead">Lessons, dialogue, pronunciation, photo translation, mistake review, and progress live in one profile across web, PWA, and Telegram.</p>
            <div className="hero-actions">
              <a className="entry-cta hero-action hero-action--primary" data-entry="web-app" href={WEB_APP_HREF}>
                Start free <ArrowRight size={18} />
              </a>
              <a className="entry-cta hero-action hero-action--secondary" data-entry="telegram" href={TELEGRAM_HREF}>
                Open Telegram
              </a>
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
            <LandingSocialProof />
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
                  <strong>NERIVA</strong>
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

      <section className="lesson-scenario" aria-label="How NERIVA runs a short lesson">
        <div className="lesson-scenario__copy">
          <span className="eyebrow">Live product scenario</span>
          <h2>Every session ends with a visible next step</h2>
          <p>NERIVA explains a phrase, asks for an answer, checks it, saves the mistake, and gives a repeatable prompt for the next review.</p>
          <div className="lesson-scenario__tags" aria-label="Lesson flow">
            <span className="lesson-scenario__tag">phrase</span>
            <span className="lesson-scenario__tag">answer</span>
            <span className="lesson-scenario__tag">correction</span>
            <span className="lesson-scenario__tag">review</span>
          </div>
        </div>

        <div className="lesson-scenario__console" aria-label="AI Tutor lesson example">
          <div className="cockpit-console__top">
            <span>Live learning loop</span>
            <strong>84/100</strong>
          </div>
          <div className="cockpit-console__screen">
            <div className="cockpit-status">
              <Activity size={18} />
              <span>AI Tutor is checking the answer</span>
            </div>
            <div className="cockpit-dialogue">
              <p>Say it naturally: "Could you help me check in?"</p>
              <p>Better: "I have a reservation under my name."</p>
            </div>
            <div className="cockpit-metric-grid">
              <span>
                <Zap size={17} />
                XP +12
              </span>
              <span>
                <LineChart size={17} />
                weak word fixed
              </span>
              <span>
                <Languages size={17} />
                translation saved
              </span>
            </div>
          </div>
        </div>

        <figure className="lesson-scenario__media">
          <img src="/assets/scenarios/speaking-ai-tutor.jpg" alt="Speaking practice lesson preview" loading="lazy" />
          <figcaption>Speaking practice, voice correction, and saved weak words in one loop.</figcaption>
        </figure>
      </section>

      <section id="features" className="spark-section modules-section">
        <div className="section-copy">
          <span className="eyebrow">Product modules</span>
          <h2>Four tools, one learning profile</h2>
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

      <section id="pricing" className="pricing-section">
        <div className="section-copy">
          <span className="eyebrow">Pricing</span>
          <h2>Start free, then unlock the daily practice limits you actually need</h2>
          <p className="pricing-lead">Payment is available through Telegram Stars and YooKassa/SBP.</p>
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
        </div>
      </section>

      <section className="spark-section community-section">
        <div className="section-copy">
          <span className="eyebrow">Community</span>
          <h2>Premium key giveaways and product updates in the NERIVA channels</h2>
          <p>Follow the community for short lessons, product updates, and monthly Premium key giveaways.</p>
        </div>
        <div className="community-panel">
          <div>
            <strong>Follow NERIVA</strong>
            <span>Short lessons, product updates, and learning tips.</span>
          </div>
          <SocialIconLinks />
        </div>
      </section>

      <section id="faq" className="spark-section landing-faq">
        <div className="section-copy">
          <span className="eyebrow">FAQ</span>
          <h2>Questions before starting</h2>
        </div>
        <div className="landing-faq__grid">
          {shortFaq.map((item) => (
            <article className="landing-faq__item" key={item.question}>
              <h3>{item.question}</h3>
              <p>{item.answer}</p>
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
