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
import { getLandingContent, type LandingContent, type LandingLocale } from "./landingContent";

const WEB_APP_HREF = "/app/";
const TELEGRAM_HREF = "https://t.me/NERIVAapp_bot";

type ProductImageKey = keyof LandingContent["images"];
type ProductImage = {
  src: string;
  altKey: ProductImageKey;
};

const productImages = {
  dashboard: { src: "/assets/product/dashboard-progress.png", altKey: "dashboard" },
  mobileHome: { src: "/assets/product/mobile-home-progress.png", altKey: "mobileHome" },
  aiTutor: { src: "/assets/product/ai-tutor-lesson-correction.png", altKey: "aiTutor" },
  mistakes: { src: "/assets/product/mistakes-review.png", altKey: "mistakes" },
  voice: { src: "/assets/product/voice-pronunciation-score.png", altKey: "voice" },
  photo: { src: "/assets/product/photo-translation.png", altKey: "photo" },
  notes: { src: "/assets/product/notes-phrasebook.png", altKey: "notes" },
  premium: { src: "/assets/product/premium-plans.png", altKey: "premium" },
  telegram: { src: "/assets/product/telegram-app-light.png", altKey: "telegram" },
  mobileLesson: { src: "/assets/product/mobile-lesson-correction.png", altKey: "mobileLesson" },
  mobileMistakes: { src: "/assets/product/mobile-mistakes.png", altKey: "mobileMistakes" },
  mobileVoice: { src: "/assets/product/mobile-voice-pronunciation.png", altKey: "mobileVoice" },
} satisfies Record<string, ProductImage>;

const featureIcons = [BrainCircuit, Mic, Camera, Repeat2] as const;
const featureImages = [productImages.aiTutor, productImages.voice, productImages.photo, productImages.mistakes] as const;

type EnglishSparkLandingProps = {
  siteTheme?: "dark" | "light";
  language?: LandingLocale;
};

function ProductShot({ image, alt, className = "" }: { image: ProductImage; alt: string; className?: string }) {
  return (
    <figure className={`product-shot ${className}`.trim()}>
      <img src={image.src} alt={alt} loading="lazy" />
    </figure>
  );
}

export function EnglishSparkLanding({ siteTheme = "light", language = "en" }: EnglishSparkLandingProps) {
  const copy = getLandingContent(language);
  const heroVariant = siteTheme === "dark" ? "dark" : "light";
  const heroColor = siteTheme === "dark" ? "#7bdcff" : "#002fa7";
  const heroParticleColor = siteTheme === "dark" ? "#f5d27a" : "#002fa7";
  const imgAlt = (image: ProductImage) => copy.images[image.altKey];

  return (
    <main className="english-spark-landing" data-visual-anchor="swiss-editorial">
      <section className="landing-hero" data-hero-preset={siteTheme}>
        <div className="landing-hero__matter" aria-hidden="true">
          <GenerativeArtScene
            animate
            transparentBackdrop
            variant={heroVariant}
            color={heroColor}
            particleColor={heroParticleColor}
          />
        </div>
        <div className="landing-hero__inner">
          <div className="landing-hero__copy">
            <span className="eyebrow">{copy.hero.eyebrow}</span>
            <h1>{copy.hero.title}</h1>
            <p className="hero-lead">{copy.hero.lead}</p>
            <div className="hero-actions">
              <a className="entry-cta hero-action hero-action--primary" data-entry="web-app" href={WEB_APP_HREF}>
                {copy.nav.webApp} <ArrowRight size={18} />
              </a>
              <a className="entry-cta hero-action hero-action--secondary" data-entry="telegram" href={TELEGRAM_HREF}>
                {copy.nav.telegram} <MessageCircle size={18} />
              </a>
            </div>
            <div className="hero-proof" aria-label="NERIVA product status">
              {copy.hero.proof.map(([title, body]) => (
                <span key={title}>
                  <strong>{title}</strong>
                  {body}
                </span>
              ))}
            </div>
          </div>

          <div className="mobile-landing-rail" aria-label="NERIVA mobile quick view">
            {copy.mobile.rail.map((item) => (
              <span key={item}>{item}</span>
            ))}
          </div>

          <div className="hero-product-frame" aria-label="NERIVA product screenshots">
            <ProductShot image={productImages.dashboard} alt={imgAlt(productImages.dashboard)} className="product-shot--desktop" />
            <ProductShot image={productImages.mobileHome} alt={imgAlt(productImages.mobileHome)} className="product-shot--phone" />
          </div>
        </div>
      </section>

      <section className="lesson-scenario" aria-label="How NERIVA runs a short lesson">
        <div className="lesson-scenario__copy">
          <span className="eyebrow">{copy.scenario.eyebrow}</span>
          <h2>{copy.scenario.title}</h2>
          <p>{copy.scenario.body}</p>
          <div className="lesson-scenario__tags" aria-label="Lesson flow">
            {copy.scenario.tags.map((tag) => (
              <span className="lesson-scenario__tag" key={tag}>
                {tag}
              </span>
            ))}
          </div>
        </div>
        <div className="scenario-shot" aria-label="AI Tutor lesson example">
          <ProductShot image={productImages.aiTutor} alt={imgAlt(productImages.aiTutor)} />
          <div className="scenario-shot__caption">
            <span>{copy.scenario.captionLabel}</span>
            <strong>{copy.scenario.caption}</strong>
          </div>
        </div>
      </section>

      <section id="features" className="spark-section modules-section">
        <div className="section-copy">
          <span className="eyebrow">{copy.features.eyebrow}</span>
          <h2>{copy.features.title}</h2>
        </div>
        <div className="feature-grid">
          {copy.features.items.map((feature, index) => {
            const Icon = featureIcons[index] ?? BrainCircuit;
            const image = featureImages[index] ?? productImages.aiTutor;
            return (
              <article className="feature-panel" key={feature.title}>
                <div className="feature-panel__copy">
                  <Icon size={22} />
                  <h3>{feature.title}</h3>
                  <p>{feature.body}</p>
                </div>
                <ProductShot image={image} alt={imgAlt(image)} />
              </article>
            );
          })}
        </div>
      </section>

      <section className="spark-section notes-section" aria-label="Notes and phrasebook">
        <div className="section-copy">
          <span className="eyebrow">{copy.notes.eyebrow}</span>
          <h2>{copy.notes.title}</h2>
          <p>{copy.notes.body}</p>
        </div>
        <ProductShot image={productImages.notes} alt={imgAlt(productImages.notes)} />
      </section>

      <section id="pricing" className="pricing-section">
        <div className="section-copy">
          <span className="eyebrow">{copy.pricing.eyebrow}</span>
          <h2>{copy.pricing.title}</h2>
          <p className="pricing-lead">{copy.pricing.lead}</p>
        </div>
        <div className="pricing-layout">
          <div className="pricing-grid">
            {copy.pricing.plans.map((plan, index) => (
              <article className={index === 1 ? "plan-card is-featured" : "plan-card"} key={plan.name}>
                <span className="plan-label">{plan.label}</span>
                <h3>{plan.name}</h3>
                <div className="price-line">
                  {plan.oldPrice ? <s>{plan.oldPrice}</s> : null}
                  <strong>{plan.price}</strong>
                  <small>{plan.period}</small>
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
                <a href={WEB_APP_HREF}>{copy.pricing.choose}</a>
              </article>
            ))}
          </div>
          <ProductShot image={productImages.premium} alt={imgAlt(productImages.premium)} className="product-shot--pricing" />
        </div>
        <div className="payment-methods" aria-label="Payment methods">
          {copy.pricing.methods.map((method) => (
            <span key={method}>{method}</span>
          ))}
        </div>
      </section>

      <section className="spark-section telegram-section">
        <div className="section-copy">
          <span className="eyebrow">{copy.telegram.eyebrow}</span>
          <h2>{copy.telegram.title}</h2>
          <p>{copy.telegram.body}</p>
        </div>
        <div className="telegram-panel">
          <ProductShot image={productImages.telegram} alt={imgAlt(productImages.telegram)} />
          <div className="telegram-panel__copy">
            <BookOpen size={22} />
            <strong>{copy.telegram.note}</strong>
            <a className="entry-cta" data-entry="telegram" href={TELEGRAM_HREF}>
              {copy.telegram.cta} <ArrowRight size={18} />
            </a>
          </div>
        </div>
      </section>

      <section className="spark-section mobile-section">
        <div className="section-copy">
          <span className="eyebrow">{copy.mobile.eyebrow}</span>
          <h2>{copy.mobile.title}</h2>
          <p>{copy.mobile.body}</p>
        </div>
        <div className="mobile-strip" aria-label="Mobile web app screenshots">
          <ProductShot image={productImages.mobileLesson} alt={imgAlt(productImages.mobileLesson)} />
          <ProductShot image={productImages.mobileMistakes} alt={imgAlt(productImages.mobileMistakes)} />
          <ProductShot image={productImages.mobileVoice} alt={imgAlt(productImages.mobileVoice)} />
        </div>
      </section>

      <section className="spark-section community-section">
        <div className="section-copy">
          <span className="eyebrow">{copy.community.eyebrow}</span>
          <h2>{copy.community.title}</h2>
          <p>{copy.community.body}</p>
        </div>
        <div className="community-panel">
          <div>
            <Star size={22} />
            <strong>{copy.community.follow}</strong>
            <span>{copy.community.note}</span>
          </div>
          <SocialIconLinks />
        </div>
      </section>

      <section id="faq" className="spark-section landing-faq">
        <div className="section-copy">
          <span className="eyebrow">{copy.faq.eyebrow}</span>
          <h2>{copy.faq.title}</h2>
        </div>
        <div className="landing-faq__grid">
          {copy.faq.items.map((item) => (
            <article className="landing-faq__item" key={item.question}>
              <h3>{item.question}</h3>
              <p>{item.answer}</p>
            </article>
          ))}
        </div>
      </section>

      <section className="final-cta-section">
        <div className="final-cta-section__copy">
          <Repeat2 size={32} />
          <h2>{copy.finalCta.title}</h2>
          <p>{copy.finalCta.body}</p>
          <div>
            <a className="entry-cta" data-entry="web-app" href={WEB_APP_HREF}>
              {copy.nav.webApp} <ArrowRight size={18} />
            </a>
            <a className="entry-cta" data-entry="telegram" href={TELEGRAM_HREF}>
              {copy.nav.telegram} <ArrowRight size={18} />
            </a>
          </div>
        </div>
        <div className="final-cta-visual" aria-label={copy.finalCta.visualLabel}>
          <strong>{copy.finalCta.visualLabel}</strong>
          <span>{copy.finalCta.visualBody}</span>
        </div>
      </section>
    </main>
  );
}
