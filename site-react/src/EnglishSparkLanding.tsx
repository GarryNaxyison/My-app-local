import {
  ArrowRight,
  BookOpen,
  BrainCircuit,
  Camera,
  CheckCircle,
  MessageCircle,
  Mic,
  Repeat2,
  Star,
} from "lucide-react";
import { GenerativeArtScene } from "@/components/ui/anomalous-matter-hero";
import { SocialIconLinks } from "./components/SocialLinks";
import { landingSeoCopy } from "./landingSeoContent";

const WEB_APP_HREF = "/app/";
const TELEGRAM_HREF = "https://t.me/NERIVAapp_bot";

type ProductImage = {
  src: string;
  alt: string;
};

const productImages = {
  dashboard: {
    src: "/assets/product/dashboard-progress.png",
    alt: "NERIVA web app dashboard with level, progress, XP, and streak status",
  },
  mobileHome: {
    src: "/assets/product/mobile-home-progress.png",
    alt: "NERIVA mobile dashboard with today progress",
  },
  aiTutor: {
    src: "/assets/product/ai-tutor-lesson-correction.png",
    alt: "AI tutor lesson with a corrected hotel check-in answer",
  },
  mistakes: {
    src: "/assets/product/mistakes-review.png",
    alt: "Mistakes list with corrected English phrases and audio review",
  },
  voice: {
    src: "/assets/product/voice-pronunciation-score.png",
    alt: "Pronunciation practice with score and weak words",
  },
  photo: {
    src: "/assets/product/photo-translation.png",
    alt: "Photo translation screen with OCR result and practice prompt",
  },
  notes: {
    src: "/assets/product/notes-phrasebook.png",
    alt: "Phrasebook and notes with saved useful phrases",
  },
  premium: {
    src: "/assets/product/premium-plans.png",
    alt: "Free Premium and Platinum plan limits in the web app",
  },
  telegram: {
    src: "/assets/product/telegram-app-light.png",
    alt: "NERIVA Telegram bot menu in light theme",
  },
  mobileLesson: {
    src: "/assets/product/mobile-lesson-correction.png",
    alt: "Mobile AI tutor lesson correction",
  },
  mobileMistakes: {
    src: "/assets/product/mobile-mistakes.png",
    alt: "Mobile mistakes review screen",
  },
  mobileVoice: {
    src: "/assets/product/mobile-voice-pronunciation.png",
    alt: "Mobile voice pronunciation practice",
  },
} satisfies Record<string, ProductImage>;

const features = [
  {
    icon: BrainCircuit,
    title: "AI Tutor",
    body: "Guided lessons turn a phrase into an answer, a correction, and a next action.",
    image: productImages.aiTutor,
  },
  {
    icon: Mic,
    title: "Voice Coach",
    body: "Pronunciation practice highlights weak words, rhythm, and the next sentence to repeat.",
    image: productImages.voice,
  },
  {
    icon: Camera,
    title: "Photo Practice",
    body: "A menu or sign becomes translation, context, notes, and a short practice loop.",
    image: productImages.photo,
  },
  {
    icon: Repeat2,
    title: "Mistake Review",
    body: "Saved mistakes, weak words, and corrected phrases return when they are useful.",
    image: productImages.mistakes,
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

function ProductShot({ image, className = "" }: { image: ProductImage; className?: string }) {
  return (
    <figure className={`product-shot ${className}`.trim()}>
      <img src={image.src} alt={image.alt} loading="lazy" />
    </figure>
  );
}

export function EnglishSparkLanding({ siteTheme = "light" }: EnglishSparkLandingProps) {
  return (
    <main className="english-spark-landing" data-visual-anchor="swiss-editorial">
      <section className="landing-hero" data-hero-preset={siteTheme}>
        <div className="landing-hero__matter" aria-hidden="true">
          <GenerativeArtScene animate variant="light" color="#002fa7" particleColor="#002fa7" />
        </div>
        <div className="landing-hero__inner">
          <div className="landing-hero__copy">
            <span className="eyebrow">AI tutor in web app and Telegram</span>
            <h1>NERIVA: Practice speaking before the moment matters</h1>
            <p className="hero-lead">
              NERIVA combines lessons, practice, listening, vocabulary, translator tools and pronunciation scoring. Learn in Telegram or the web app, get structured tasks, voice feedback and visible progress.
            </p>
            <div className="hero-actions">
              <a className="entry-cta hero-action hero-action--primary" data-entry="web-app" href={WEB_APP_HREF}>
                Web app <ArrowRight size={18} />
              </a>
              <a className="entry-cta hero-action hero-action--secondary" data-entry="telegram" href={TELEGRAM_HREF}>
                Telegram <MessageCircle size={18} />
              </a>
            </div>
            <div className="hero-proof" aria-label="NERIVA product status">
              <span>
                <strong>35 interface languages</strong>
                curated landing layer
              </span>
              <span>
                <strong>A1-C2</strong>
                lesson levels
              </span>
              <span>
                <strong>Free, Premium, Platinum</strong>
                clear daily limits
              </span>
            </div>
          </div>

          <div className="hero-product-frame" aria-label="NERIVA product screenshots">
            <ProductShot image={productImages.dashboard} className="product-shot--desktop" />
            <ProductShot image={productImages.mobileHome} className="product-shot--phone" />
          </div>
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
        <div className="scenario-shot" aria-label="AI Tutor lesson example">
          <ProductShot image={productImages.aiTutor} />
          <div className="scenario-shot__caption">
            <span>Hotel check-in</span>
            <strong>Correction, explanation, audio sample, and the user's answer in one lesson.</strong>
          </div>
        </div>
      </section>

      <section id="features" className="spark-section modules-section">
        <div className="section-copy">
          <span className="eyebrow">Product modules</span>
          <h2>Four tools, one learning profile</h2>
        </div>
        <div className="feature-grid">
          {features.map((feature) => {
            const Icon = feature.icon;
            return (
              <article className="feature-panel" key={feature.title}>
                <div className="feature-panel__copy">
                  <Icon size={22} />
                  <h3>{feature.title}</h3>
                  <p>{feature.body}</p>
                </div>
                <ProductShot image={feature.image} />
              </article>
            );
          })}
        </div>
      </section>

      <section className="spark-section notes-section" aria-label="Notes and phrasebook">
        <div className="section-copy">
          <span className="eyebrow">Notes / Phrasebook</span>
          <h2>Saved phrases stay useful after the lesson</h2>
          <p>Notes, saved translations, weak words, and phrasebook examples keep the practice from disappearing after one chat.</p>
        </div>
        <ProductShot image={productImages.notes} />
      </section>

      <section id="pricing" className="pricing-section">
        <div className="section-copy">
          <span className="eyebrow">Pricing</span>
          <h2>Start free, then unlock the daily practice limits you actually need</h2>
          <p className="pricing-lead">Payment is available through Telegram Stars and YooKassa/SBP.</p>
        </div>
        <div className="pricing-layout">
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
          <ProductShot image={productImages.premium} className="product-shot--pricing" />
        </div>
        <div className="payment-methods" aria-label="Payment methods">
          <span>Telegram Stars</span>
          <span>YooKassa/SBP</span>
        </div>
      </section>

      <section className="spark-section telegram-section">
        <div className="section-copy">
          <span className="eyebrow">Telegram</span>
          <h2>Telegram is the quick entry, not a second product</h2>
          <p>Use the bot for fast starts, voice, photos, reminders, and the same learning profile that continues in the web app.</p>
        </div>
        <div className="telegram-panel">
          <ProductShot image={productImages.telegram} />
          <div className="telegram-panel__copy">
            <BookOpen size={22} />
            <strong>Lesson, practice, pronunciation, mistakes, progress, and Premium are one tap away.</strong>
            <a className="entry-cta" data-entry="telegram" href={TELEGRAM_HREF}>
              Open Telegram <ArrowRight size={18} />
            </a>
          </div>
        </div>
      </section>

      <section className="spark-section mobile-section">
        <div className="section-copy">
          <span className="eyebrow">Mobile web app</span>
          <h2>Short practice fits the phone screen</h2>
          <p>Dashboard, lesson correction, mistakes, and pronunciation keep the same loop on mobile.</p>
        </div>
        <div className="mobile-strip" aria-label="Mobile web app screenshots">
          <ProductShot image={productImages.mobileLesson} />
          <ProductShot image={productImages.mobileMistakes} />
          <ProductShot image={productImages.mobileVoice} />
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
            <Star size={22} />
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
        <Repeat2 size={32} />
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
