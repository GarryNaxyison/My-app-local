import { useState, type ComponentPropsWithoutRef, type ReactNode } from "react";
import { motion, useReducedMotion } from "framer-motion";
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
  X,
} from "lucide-react";
import { PremiumHeroAnimation } from "@/components/ui/premium-hero-animation";
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
  dashboard: { src: "/assets/product/dashboard-progress-dark.png", altKey: "dashboard" },
  mobileHome: { src: "/assets/product/mobile-home-progress-dark.png", altKey: "mobileHome" },
  aiTutor: { src: "/assets/product/ai-tutor-lesson-correction-dark.png", altKey: "aiTutor" },
  mistakes: { src: "/assets/product/mistakes-review-dark.png", altKey: "mistakes" },
  voice: { src: "/assets/product/voice-pronunciation-score-dark.png", altKey: "voice" },
  photo: { src: "/assets/product/photo-translation-dark.png", altKey: "photo" },
  notes: { src: "/assets/product/notes-phrasebook-dark.png", altKey: "notes" },
  telegram: { src: "/assets/product/telegram-app-dark.png", altKey: "telegram" },
  mobileLesson: { src: "/assets/product/mobile-lesson-correction-dark.png", altKey: "mobileLesson" },
  mobileMistakes: { src: "/assets/product/mobile-mistakes-dark.png", altKey: "mobileMistakes" },
  mobileVoice: { src: "/assets/product/mobile-voice-pronunciation-dark.png", altKey: "mobileVoice" },
} satisfies Record<string, ProductImage>;

const featureIcons = [BrainCircuit, Mic, Camera, Repeat2] as const;
const mobileIcons = [Star, BookOpen, Repeat2, Mic] as const;

type EnglishSparkLandingProps = {
  language?: LandingLocale;
};

type RevealSectionProps = Omit<ComponentPropsWithoutRef<typeof motion.section>, "children"> & { children: ReactNode };
type RevealArticleProps = Omit<ComponentPropsWithoutRef<typeof motion.article>, "children"> & { children: ReactNode };

const revealViewport = { once: true, amount: 0.18 };

function RevealSection({ children, transition, ...props }: RevealSectionProps) {
  const reduceMotion = useReducedMotion();

  if (reduceMotion) {
    return <section {...(props as ComponentPropsWithoutRef<"section">)}>{children}</section>;
  }

  return (
    <motion.section
      initial={{ opacity: 0, y: 34, scale: 0.985 }}
      whileInView={{ opacity: 1, y: 0, scale: 1 }}
      viewport={revealViewport}
      transition={{ duration: 0.72, ease: [0.22, 1, 0.36, 1], ...transition }}
      {...props}
    >
      {children}
    </motion.section>
  );
}

function RevealArticle({ children, transition, ...props }: RevealArticleProps) {
  const reduceMotion = useReducedMotion();

  if (reduceMotion) {
    return <article {...(props as ComponentPropsWithoutRef<"article">)}>{children}</article>;
  }

  return (
    <motion.article
      initial={{ opacity: 0, y: 26, filter: "blur(10px)" }}
      whileInView={{ opacity: 1, y: 0, filter: "blur(0px)" }}
      viewport={revealViewport}
      transition={{ duration: 0.62, ease: [0.22, 1, 0.36, 1], ...transition }}
      {...props}
    >
      {children}
    </motion.article>
  );
}

function ProductShot({
  image,
  alt,
  className = "",
  onOpen,
}: {
  image: ProductImage;
  alt: string;
  className?: string;
  onOpen: (image: ProductImage, alt: string) => void;
}) {
  return (
    <figure className={`product-shot ${className}`.trim()}>
      <button className="product-shot__button" type="button" onClick={() => onOpen(image, alt)} aria-label={alt}>
        <img src={image.src} alt={alt} loading="lazy" />
      </button>
    </figure>
  );
}

export function EnglishSparkLanding({ language = "en" }: EnglishSparkLandingProps) {
  const copy = getLandingContent(language);
  const [activeFeature, setActiveFeature] = useState(0);
  const [activeMobile, setActiveMobile] = useState(0);
  const [preview, setPreview] = useState<{ image: ProductImage; alt: string } | null>(null);
  const active = copy.features.items[activeFeature] ?? copy.features.items[0];
  const ActiveIcon = featureIcons[activeFeature] ?? BrainCircuit;
  const activeImage = productImages[active.imageKey] ?? productImages.aiTutor;
  const activeMobileItem = copy.mobile.items[activeMobile] ?? copy.mobile.items[0];
  const ActiveMobileIcon = mobileIcons[activeMobile] ?? Star;
  const activeMobileImage = productImages[activeMobileItem.imageKey] ?? productImages.mobileHome;
  const imgAlt = (image: ProductImage) => copy.images[image.altKey];
  const openPreview = (image: ProductImage, alt: string) => setPreview({ image, alt });

  return (
    <main className="english-spark-landing" data-visual-anchor="premium-tech-narrative">
      <section className="landing-hero" data-hero-preset="dark">
        <div className="landing-hero__matter" aria-hidden="true">
          <PremiumHeroAnimation />
        </div>
        <div className="landing-hero__inner">
          <div className="landing-hero__copy">
            <span className="hero-announcement">{copy.hero.announcement}</span>
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

          <div className="hero-product-frame" aria-label="NERIVA product screenshots">
            <ProductShot image={productImages.dashboard} alt={imgAlt(productImages.dashboard)} className="product-shot--desktop" onOpen={openPreview} />
            <ProductShot image={productImages.mobileHome} alt={imgAlt(productImages.mobileHome)} className="product-shot--phone" onOpen={openPreview} />
            <div className="hero-product-frame__badges" aria-hidden="true">
              {copy.hero.badges.map((badge) => <span key={badge}>{badge}</span>)}
            </div>
          </div>
        </div>
      </section>

      <RevealSection className="stitch-metrics" aria-label="NERIVA operating metrics">
        {copy.metrics.map((metric, index) => (
          <RevealArticle className="stitch-metric-card" key={metric.label} transition={{ delay: index * 0.08 }}>
            <span>{metric.label}</span>
            <strong>{metric.value}</strong>
            <p>{metric.body}</p>
          </RevealArticle>
        ))}
      </RevealSection>

      <RevealSection className="before-after-bridge" aria-label="Before and after NERIVA">
        <div className="section-copy">
          <span className="eyebrow">{copy.beforeAfter.eyebrow}</span>
          <h2>{copy.beforeAfter.title}</h2>
        </div>
        <div className="before-after-bridge__grid">
          <RevealArticle>
            <span>01</span>
            <h3>{copy.beforeAfter.beforeTitle}</h3>
            <ul>
              {copy.beforeAfter.before.map((item) => <li key={item}>{item}</li>)}
            </ul>
          </RevealArticle>
          <RevealArticle className="is-after" transition={{ delay: 0.08 }}>
            <span>02</span>
            <h3>{copy.beforeAfter.afterTitle}</h3>
            <ul>
              {copy.beforeAfter.after.map((item) => <li key={item}>{item}</li>)}
            </ul>
          </RevealArticle>
        </div>
      </RevealSection>

      <RevealSection className="lesson-scenario" aria-label="How NERIVA runs a short lesson">
        <div className="lesson-scenario__copy">
          <span className="eyebrow">{copy.scenario.eyebrow}</span>
          <h2>{copy.scenario.title}</h2>
          <p>{copy.scenario.body}</p>
          <div className="lesson-scenario__tags" aria-label="Lesson flow">
            {copy.scenario.tags.map((tag) => <span className="lesson-scenario__tag" key={tag}>{tag}</span>)}
          </div>
        </div>
        <div className="scenario-shot" aria-label="AI Tutor lesson example">
          <ProductShot image={productImages.aiTutor} alt={imgAlt(productImages.aiTutor)} onOpen={openPreview} />
          <div className="scenario-shot__caption">
            <span>{copy.scenario.captionLabel}</span>
            <strong>{copy.scenario.caption}</strong>
          </div>
        </div>
      </RevealSection>

      <RevealSection id="features" className="spark-section modules-section">
        <div className="section-copy">
          <span className="eyebrow">{copy.features.eyebrow}</span>
          <h2>{copy.features.title}</h2>
        </div>
        <div className="feature-carousel">
          <div className="feature-carousel__tabs" role="tablist" aria-label={copy.features.title}>
            {copy.features.items.map((feature, index) => {
              const Icon = featureIcons[index] ?? BrainCircuit;
              return (
                <button
                  className={index === activeFeature ? "feature-carousel__tab is-active" : "feature-carousel__tab"}
                  key={feature.title}
                  type="button"
                  role="tab"
                  aria-selected={index === activeFeature}
                  onClick={() => setActiveFeature(index)}
                >
                  <Icon size={18} />
                  <span>{feature.title}</span>
                  <small>{feature.short}</small>
                </button>
              );
            })}
          </div>
          <div className="feature-carousel__viewport">
            <article className="feature-panel" role="tabpanel">
              <div className="feature-panel__copy">
                <ActiveIcon size={26} />
                <span>{active.stat}</span>
                <h3>{active.title}</h3>
                <p>{active.body}</p>
              </div>
              <ProductShot image={activeImage} alt={imgAlt(activeImage)} onOpen={openPreview} />
            </article>
          </div>
        </div>
      </RevealSection>

      <RevealSection className="spark-section methodology-timeline" aria-label={copy.methodology.title}>
        <div className="section-copy">
          <span className="eyebrow">{copy.methodology.eyebrow}</span>
          <h2>{copy.methodology.title}</h2>
          <p>{copy.methodology.body}</p>
        </div>
        <div className="methodology-timeline__steps">
          {copy.methodology.steps.map((step) => (
            <RevealArticle className="methodology-step" key={step.label}>
              <span>{step.label}</span>
              <h3>{step.title}</h3>
              <p>{step.body}</p>
            </RevealArticle>
          ))}
        </div>
      </RevealSection>

      <RevealSection id="pricing" className="pricing-section">
        <div className="section-copy">
          <span className="eyebrow">{copy.pricing.eyebrow}</span>
          <h2>{copy.pricing.title}</h2>
          <p className="pricing-lead">{copy.pricing.lead}</p>
        </div>
        <div className="pricing-grid">
          {copy.pricing.plans.map((plan, index) => (
            <RevealArticle className={index === 1 ? "plan-card is-featured" : "plan-card"} key={plan.name} transition={{ delay: index * 0.08 }}>
              {index === 1 ? <span className="plan-recommendation">Recommended</span> : null}
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
              <div className="plan-access" aria-label={`${plan.name} access`}>
                {plan.access.map((item) => <span key={item}>{item}</span>)}
              </div>
              <a href={WEB_APP_HREF}>{copy.pricing.choose}</a>
            </RevealArticle>
          ))}
        </div>
        <div className="payment-methods" aria-label="Payment methods">
          {copy.pricing.methods.map((method) => <span key={method}>{method}</span>)}
        </div>
      </RevealSection>

      <RevealSection className="spark-section ecosystem-showcase" aria-label={copy.ecosystem.title}>
        <div className="section-copy">
          <span className="eyebrow">{copy.ecosystem.eyebrow}</span>
          <h2>{copy.ecosystem.title}</h2>
          <p>{copy.ecosystem.body}</p>
        </div>
        <div className="ecosystem-showcase__grid">
          {copy.ecosystem.items.map((item) => {
            const image = productImages[item.imageKey] ?? productImages.dashboard;
            return (
              <RevealArticle className="ecosystem-card" key={item.title}>
                <ProductShot image={image} alt={imgAlt(image)} onOpen={openPreview} />
                <div>
                  <strong>{item.title}</strong>
                  <p>{item.body}</p>
                </div>
              </RevealArticle>
            );
          })}
        </div>
      </RevealSection>

      <RevealSection className="spark-section telegram-section">
        <div className="section-copy">
          <span className="eyebrow">{copy.telegram.eyebrow}</span>
          <h2>{copy.telegram.title}</h2>
          <p>{copy.telegram.body}</p>
        </div>
        <div className="telegram-panel">
          <ProductShot image={productImages.telegram} alt={imgAlt(productImages.telegram)} onOpen={openPreview} />
          <div className="telegram-panel__copy">
            <BookOpen size={22} />
            <strong>{copy.telegram.note}</strong>
            <div className="telegram-panel__chips">
              {copy.telegram.chips.map((chip) => <span key={chip}>{chip}</span>)}
            </div>
            <a className="entry-cta" data-entry="telegram" href={TELEGRAM_HREF}>
              {copy.telegram.cta} <ArrowRight size={18} />
            </a>
          </div>
        </div>
      </RevealSection>

      <RevealSection className="spark-section mobile-section">
        <div className="section-copy">
          <span className="eyebrow">{copy.mobile.eyebrow}</span>
          <h2>{copy.mobile.title}</h2>
          <p>{copy.mobile.body}</p>
        </div>
        <div className="mobile-carousel">
          <div className="mobile-carousel__tabs" role="tablist" aria-label={copy.mobile.title}>
            {copy.mobile.items.map((item, index) => {
              const Icon = mobileIcons[index] ?? Star;
              return (
                <button
                  className={index === activeMobile ? "mobile-carousel__tab is-active" : "mobile-carousel__tab"}
                  key={item.title}
                  type="button"
                  role="tab"
                  aria-selected={index === activeMobile}
                  onClick={() => setActiveMobile(index)}
                >
                  <Icon size={18} />
                  <span>{item.title}</span>
                  <small>{item.short}</small>
                </button>
              );
            })}
          </div>
          <div className="mobile-carousel__viewport">
            <article className="mobile-panel" role="tabpanel">
              <div className="mobile-panel__copy">
                <ActiveMobileIcon size={26} />
                <span>{activeMobileItem.stat}</span>
                <h3>{activeMobileItem.title}</h3>
                <p>{activeMobileItem.body}</p>
              </div>
              <ProductShot image={activeMobileImage} alt={imgAlt(activeMobileImage)} onOpen={openPreview} />
            </article>
          </div>
        </div>
      </RevealSection>

      <RevealSection className="spark-section community-section">
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
            <div className="community-panel__chips">
              {copy.community.chips.map((chip) => <small key={chip}>{chip}</small>)}
            </div>
          </div>
          <SocialIconLinks />
        </div>
      </RevealSection>

      <RevealSection id="faq" className="spark-section landing-faq">
        <div className="section-copy">
          <span className="eyebrow">{copy.faq.eyebrow}</span>
          <h2>{copy.faq.title}</h2>
        </div>
        <div className="landing-faq__grid">
          {copy.faq.items.map((item) => (
            <RevealArticle className="landing-faq__item" key={item.question}>
              <h3>{item.question}</h3>
              <p>{item.answer}</p>
            </RevealArticle>
          ))}
        </div>
      </RevealSection>

      <RevealSection className="final-cta-section">
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
          <ol>
            {copy.finalCta.flow.map((step) => <li key={step}>{step}</li>)}
          </ol>
        </div>
      </RevealSection>

      {preview ? (
        <div className="image-preview" role="dialog" aria-modal="true" aria-label={preview.alt} onClick={() => setPreview(null)}>
          <button className="image-preview__close" type="button" onClick={() => setPreview(null)} aria-label="Close image preview">
            <X size={24} />
          </button>
          <figure className="image-preview__frame" onClick={(event) => event.stopPropagation()}>
            <img src={preview.image.src} alt={preview.alt} />
            <figcaption>{preview.alt}</figcaption>
          </figure>
        </div>
      ) : null}
    </main>
  );
}
