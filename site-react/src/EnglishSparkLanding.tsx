import { useState, type ComponentPropsWithoutRef, type KeyboardEvent as ReactKeyboardEvent, type ReactNode } from "react";
import { motion, useReducedMotion } from "framer-motion";
import {
  ArrowRight,
  BookOpen,
  BrainCircuit,
  Camera,
  CheckCircle,
  Download,
  MessageCircle,
  Mic,
  Repeat2,
  Star,
} from "lucide-react";
import { SocialIconLinks } from "./components/SocialLinks";
import { PremiumHeroAnimation } from "@/components/ui/premium-hero-animation";
import { getLandingContent, type LandingLocale } from "./landingContent";

import { HeroLesson, ProductWindow } from "./ProductWindows";

const WEB_APP_HREF = "/app/";
const TELEGRAM_HREF = "https://t.me/NERIVAapp_bot";

const featureIcons = [BrainCircuit, Mic, Camera, Repeat2] as const;
const mobileIcons = [Star, BookOpen, Repeat2, Mic] as const;

type EnglishSparkLandingProps = {
  language?: LandingLocale;
};

type RevealSectionProps = Omit<ComponentPropsWithoutRef<typeof motion.section>, "children"> & { children: ReactNode };
type RevealArticleProps = Omit<ComponentPropsWithoutRef<typeof motion.article>, "children"> & { children: ReactNode };
type RevealBlockProps = Omit<ComponentPropsWithoutRef<typeof motion.div>, "children"> & { children: ReactNode };

const revealViewport = { once: true, amount: 0.12, margin: "0px 0px -10% 0px" };

function handleTabNavigation(event: ReactKeyboardEvent<HTMLDivElement>, active: number, select: (index: number) => void) {
  const tabs = event.currentTarget.querySelectorAll<HTMLButtonElement>('[role="tab"]');
  const count = tabs.length;
  const rtl = getComputedStyle(event.currentTarget).direction === "rtl";
  const direction = rtl ? -1 : 1;
  let next: number;
  switch (event.key) {
    case "ArrowRight": next = (active + direction + count) % count; break;
    case "ArrowLeft": next = (active - direction + count) % count; break;
    case "Home": next = 0; break;
    case "End": next = count - 1; break;
    default: return;
  }
  event.preventDefault();
  select(next);
  tabs[next]?.focus();
}

function RevealSection({ children, transition, ...props }: RevealSectionProps) {
  const reduceMotion = useReducedMotion();

  if (reduceMotion) {
    return <section {...(props as ComponentPropsWithoutRef<"section">)}>{children}</section>;
  }

  return (
    <motion.section
      initial={{ opacity: 0, y: 18 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={revealViewport}
      transition={{ duration: 0.48, ease: [0.22, 1, 0.36, 1], ...transition }}
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
      initial={{ opacity: 0, y: 12 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={revealViewport}
      transition={{ duration: 0.42, ease: [0.22, 1, 0.36, 1], ...transition }}
      {...props}
    >
      {children}
    </motion.article>
  );
}

function RevealBlock({ children, transition, ...props }: RevealBlockProps) {
  const reduceMotion = useReducedMotion();

  if (reduceMotion) {
    return <div {...(props as ComponentPropsWithoutRef<"div">)}>{children}</div>;
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 30, filter: "blur(10px)" }}
      whileInView={{ opacity: 1, y: 0, filter: "blur(0px)" }}
      viewport={revealViewport}
      transition={{ duration: 0.66, ease: [0.22, 1, 0.36, 1], ...transition }}
      {...props}
    >
      {children}
    </motion.div>
  );
}

export function EnglishSparkLanding({ language = "en" }: EnglishSparkLandingProps) {
  const reduceMotion = useReducedMotion();
  const copy = getLandingContent(language);
  const [activeFeature, setActiveFeature] = useState(0);
  const [activeMobile, setActiveMobile] = useState(0);
  const active = copy.features.items[activeFeature] ?? copy.features.items[0];
  const ActiveIcon = featureIcons[activeFeature] ?? BrainCircuit;
  const activeMobileItem = copy.mobile.items[activeMobile] ?? copy.mobile.items[0];
  const ActiveMobileIcon = mobileIcons[activeMobile] ?? Star;
  return (
    <main className="english-spark-landing" data-visual-anchor="official-neriva-brand">
      <section className="landing-hero" data-hero-preset="dark">
        {!reduceMotion && <div className="landing-hero__matter" aria-hidden="true"><PremiumHeroAnimation /></div>}
        <div className="landing-hero__inner">
          <div className="landing-hero__copy">
            <div className="hero-brand-caption">
              <div className="hero-brand-emblem" aria-hidden="true">
                <img src="/assets/brand/neriva-symbol.png" alt="" width="1024" height="569" />
              </div>
              <span className="hero-announcement">{copy.hero.announcement}</span>
            </div>
            <span className="eyebrow">{copy.hero.eyebrow}</span>
            <h1>{copy.hero.title}</h1>
            <p className="hero-lead">{copy.hero.lead}</p>
            <div className="hero-actions">
              <a className="entry-cta hero-action hero-action--primary" data-entry="web-app" href={WEB_APP_HREF}>
                {copy.hero.start} <ArrowRight size={18} />
              </a>
              <a className="entry-cta hero-action hero-action--secondary" data-entry="telegram" href={TELEGRAM_HREF}>
                {copy.nav.telegram} <MessageCircle size={18} />
              </a>
            </div>
            <a className="hero-install-link" href={`${WEB_APP_HREF}?install=1&lang=${language}`}>
              <Download size={18} aria-hidden="true" />
              <span>
                <strong>{language === "ru" ? "Установить NERIVA" : "Install NERIVA"}</strong>
                <small>{language === "ru" ? "На телефон, планшет или компьютер" : "On your phone, tablet or computer"}</small>
              </span>
            </a>
            <div className="hero-proof" aria-label="NERIVA product status">
              {copy.hero.proof.map(([title, body]) => (
                <span key={title}>
                  <strong>{title}</strong>
                  {body}
                </span>
              ))}
            </div>
          </div>

          <HeroLesson key={language} language={language} />
        </div>
      </section>

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
          <ProductWindow kind="aiTutor" language={language} />
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
          <div className="feature-carousel__tabs" role="tablist" aria-label={copy.features.title} onKeyDown={(event) => handleTabNavigation(event, activeFeature, setActiveFeature)}>
            {copy.features.items.map((feature, index) => {
              const Icon = featureIcons[index] ?? BrainCircuit;
              return (
                <button
                  className={index === activeFeature ? "feature-carousel__tab is-active" : "feature-carousel__tab"}
                  key={feature.title}
                  type="button"
                  role="tab"
                  id={`site-feature-tab-${index}`}
                  aria-controls="site-feature-panel"
                  aria-selected={index === activeFeature}
                  tabIndex={index === activeFeature ? 0 : -1}
                  onClick={() => setActiveFeature(index)}
                >
                  <Icon size={18} />
                  <span>{feature.title}</span>
                  <small>{feature.short}</small>
                </button>
              );
            })}
          </div>
          <RevealBlock className="feature-carousel__viewport">
            <RevealArticle className="feature-panel" id="site-feature-panel" role="tabpanel" aria-labelledby={`site-feature-tab-${activeFeature}`}>
              <div className="feature-panel__copy">
                <ActiveIcon size={26} />
                <span>{active.stat}</span>
                <h3>{active.title}</h3>
                <p>{active.body}</p>
              </div>
              <ProductWindow key={active.imageKey} kind={active.imageKey} language={language} />
            </RevealArticle>
          </RevealBlock>
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
      </RevealSection>

      <RevealSection className="stitch-metrics" aria-label="NERIVA operating metrics">
        {copy.metrics.map((metric, index) => (
          <RevealArticle className="stitch-metric-card" key={metric.label} transition={{ delay: index * 0.08 }}>
            <span>{metric.label}</span>
            <strong>{metric.value}</strong>
            <p>{metric.body}</p>
          </RevealArticle>
        ))}
      </RevealSection>

      <RevealSection className="spark-section outcome-section" aria-label={copy.outcomes.title}>
        <div className="section-copy">
          <span className="eyebrow">{copy.outcomes.eyebrow}</span>
          <h2>{copy.outcomes.title}</h2>
          <p>{copy.outcomes.body}</p>
        </div>
        <div className="outcome-grid">
          {copy.outcomes.items.map((item, index) => (
            <RevealArticle className="outcome-card" key={item.title} transition={{ delay: index * 0.08 }}>
              <span>{String(index + 1).padStart(2, "0")}</span>
              <h3>{item.title}</h3>
              <p>{item.body}</p>
            </RevealArticle>
          ))}
        </div>
      </RevealSection>

      <RevealSection className="spark-section use-case-section" aria-label={copy.useCases.title}>
        <div className="section-copy">
          <span className="eyebrow">{copy.useCases.eyebrow}</span>
          <h2>{copy.useCases.title}</h2>
          <p>{copy.useCases.body}</p>
        </div>
        <div className="use-case-grid">
          {copy.useCases.items.map((item, index) => {
            return (
              <RevealArticle className="use-case-card" key={item.title} transition={{ delay: index * 0.07 }}>
                <ProductWindow kind={item.imageKey} language={language} />
                <div>
                  <h3>{item.title}</h3>
                  <p>{item.body}</p>
                  <a href={WEB_APP_HREF}>{item.cta}</a>
                </div>
              </RevealArticle>
            );
          })}
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

      <RevealSection className="spark-section proof-section" aria-label={copy.proof.title}>
        <div className="proof-section__layout">
          <div className="section-copy">
            <span className="eyebrow">{copy.proof.eyebrow}</span>
            <h2>{copy.proof.title}</h2>
            <p>{copy.proof.body}</p>
          </div>
          <div className="proof-section__shots" aria-label={copy.proof.title}>
            <ProductWindow kind="aiTutor" language={language} />
            <ProductWindow kind="mistakes" language={language} />
          </div>
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
            return (
              <RevealArticle className="ecosystem-card" key={item.title}>
                <ProductWindow kind={item.imageKey} language={language} />
                <div>
                  <strong>{item.title}</strong>
                  <p>{item.body}</p>
                </div>
              </RevealArticle>
            );
          })}
        </div>
      </RevealSection>

      <RevealSection className="spark-section mobile-section">
        <div className="section-copy">
          <span className="eyebrow">{copy.mobile.eyebrow}</span>
          <h2>{copy.mobile.title}</h2>
          <p>{copy.mobile.body}</p>
        </div>
        <div className="mobile-carousel">
          <div className="mobile-carousel__tabs" role="tablist" aria-label={copy.mobile.title} onKeyDown={(event) => handleTabNavigation(event, activeMobile, setActiveMobile)}>
            {copy.mobile.items.map((item, index) => {
              const Icon = mobileIcons[index] ?? Star;
              return (
                <button
                  className={index === activeMobile ? "mobile-carousel__tab is-active" : "mobile-carousel__tab"}
                  key={item.title}
                  type="button"
                  role="tab"
                  id={`site-mobile-tab-${index}`}
                  aria-controls="site-mobile-panel"
                  aria-selected={index === activeMobile}
                  tabIndex={index === activeMobile ? 0 : -1}
                  onClick={() => setActiveMobile(index)}
                >
                  <Icon size={18} />
                  <span>{item.title}</span>
                  <small>{item.short}</small>
                </button>
              );
            })}
          </div>
          <RevealBlock className="mobile-carousel__viewport">
            <RevealArticle className="mobile-panel" id="site-mobile-panel" role="tabpanel" aria-labelledby={`site-mobile-tab-${activeMobile}`}>
              <div className="mobile-panel__copy">
                <ActiveMobileIcon size={26} />
                <span>{activeMobileItem.stat}</span>
                <h3>{activeMobileItem.title}</h3>
                <p>{activeMobileItem.body}</p>
              </div>
              <ProductWindow key={activeMobileItem.imageKey} kind={activeMobileItem.imageKey} language={language} />
            </RevealArticle>
          </RevealBlock>
        </div>
      </RevealSection>

      <RevealSection className="spark-section community-section">
        <div className="section-copy">
          <span className="eyebrow">{copy.community.eyebrow}</span>
          <h2>{copy.community.title}</h2>
          <p>{copy.community.body}</p>
        </div>
        <RevealBlock className="community-panel">
          <div>
            <Star size={22} />
            <strong>{copy.community.follow}</strong>
            <span>{copy.community.note}</span>
            <div className="community-panel__chips">
              {copy.community.chips.map((chip) => <small key={chip}>{chip}</small>)}
            </div>
          </div>
          <SocialIconLinks />
        </RevealBlock>
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
        <RevealBlock className="final-cta-visual" aria-label={copy.finalCta.visualLabel}>
          <strong>{copy.finalCta.visualLabel}</strong>
          <span>{copy.finalCta.visualBody}</span>
          <ol>
            {copy.finalCta.flow.map((step) => <li key={step}>{step}</li>)}
          </ol>
        </RevealBlock>
      </RevealSection>

    </main>
  );
}
