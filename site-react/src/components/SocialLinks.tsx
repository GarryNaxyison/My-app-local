type SocialLink = {
  name: string;
  href: string;
  label: string;
  iconSrc: string;
};

export const socialLinks: SocialLink[] = [
  { name: "YouTube", href: "https://www.youtube.com/@neriva_app", label: "Open NERIVA on YouTube", iconSrc: "/assets/social/youtube-full-color.svg" },
  { name: "Instagram", href: "https://www.instagram.com/neriva.ru", label: "Open NERIVA on Instagram", iconSrc: "/assets/social/instagram-logo-2022.svg" },
  { name: "TikTok", href: "https://tiktok.com/@nerivaru", label: "Open NERIVA on TikTok", iconSrc: "/assets/social/tiktok-icon.svg" },
  { name: "Telegram", href: "https://t.me/NERIVAapp_bot", label: "Open NERIVA on Telegram", iconSrc: "/assets/social/telegram-logo.png" },
];

export function SocialIconLinks({ className = "" }: { className?: string }) {
  return (
    <div className={`social-icon-links ${className}`.trim()} aria-label="NERIVA social channels">
      {socialLinks.map(({ name, href, label, iconSrc }) => (
        <a key={name} className={`social-icon-link social-icon-link--${name.toLowerCase()}`} href={href} aria-label={label} title={label} target="_blank" rel="noreferrer">
          <img src={iconSrc} alt="" loading="lazy" aria-hidden="true" />
        </a>
      ))}
    </div>
  );
}

export function LandingSocialProof() {
  return (
    <div className="landing-social-proof">
      <div>
        <strong>Follow NERIVA</strong>
        <span>Short lessons, product updates, and learning tips.</span>
      </div>
      <SocialIconLinks />
    </div>
  );
}

export function SiteFooterSocial({ title }: { title: string }) {
  return (
    <nav className="site-footer-social" aria-label={title}>
      <strong>{title}</strong>
      <SocialIconLinks />
    </nav>
  );
}
