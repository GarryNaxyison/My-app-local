"use client";

import { Home, SearchX } from "lucide-react";
import { Button } from "@/components/ui/button";

type NotFoundPageProps = {
  kicker?: string;
  title?: string;
  description?: string;
  actionLabel?: string;
  onHome?: () => void;
};

export function NotFoundPage({
  kicker = "404",
  title = "Page not found",
  description = "The page you are looking for is not available.",
  actionLabel = "Go to home",
  onHome,
}: NotFoundPageProps) {
  return (
    <section className="not-found-page-v2" role="alert">
      <div className="not-found-page-v2__scene" aria-hidden="true">
        <div className="not-found-page-v2__orb">
          <SearchX size={56} />
        </div>
        <span>404</span>
      </div>
      <div className="not-found-page-v2__copy">
        <span className="eyebrow">{kicker}</span>
        <h1>{title}</h1>
        <p>{description}</p>
        <Button type="button" onClick={onHome}>
          <Home size={18} />
          {actionLabel}
        </Button>
      </div>
    </section>
  );
}

export default NotFoundPage;
