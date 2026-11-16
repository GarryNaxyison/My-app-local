/** The supplied wordmark is an artwork asset, not a typeset replacement. */
export function OfficialBrand({ className = "", decorative = false }: { className?: string; decorative?: boolean }) {
  return (
    <span className={`official-brand ${className}`.trim()}>
      <img
        src="/assets/brand/neriva-wordmark.png"
        alt={decorative ? "" : "NERIVA"}
        aria-hidden={decorative || undefined}
        width="1800"
        height="329"
        decoding="async"
      />
    </span>
  );
}
