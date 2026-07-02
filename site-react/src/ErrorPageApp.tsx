import { AlertTriangle, ArrowRight, Home, MessageCircle, RefreshCw } from "lucide-react";

export type ErrorPageKind = "not-found" | "maintenance";

type ErrorPageCopy = {
  eyebrow: string;
  title: string;
  status: string;
  body: string;
  english: string;
  primaryLabel: string;
  primaryHref?: string;
  secondaryLabel: string;
  secondaryHref: string;
  detail: string;
};

const copyByKind: Record<ErrorPageKind, ErrorPageCopy> = {
  "not-found": {
    eyebrow: "NERIVA / 404",
    title: "Страница не найдена",
    status: "NOT_FOUND",
    body: "Адрес не совпал с публичной страницей NERIVA. Вернитесь на главную или откройте приложение.",
    english: "Page not found. The public landing may have moved.",
    primaryLabel: "На главную",
    primaryHref: "/poliglot-ai.html",
    secondaryLabel: "Открыть web app",
    secondaryHref: "/app/",
    detail: "Если ссылка пришла из Telegram или письма, проверьте, что она скопирована полностью.",
  },
  maintenance: {
    eyebrow: "NERIVA / status",
    title: "Технические работы",
    status: "MAINTENANCE",
    body: "Мы обновляем NERIVA или восстанавливаем соединение с приложением. Попробуйте обновить страницу через пару минут.",
    english: "Maintenance is in progress. Please retry in a few minutes.",
    primaryLabel: "Обновить",
    secondaryLabel: "На главную",
    secondaryHref: "/poliglot-ai.html",
    detail: "Лендинг и справочные страницы могут открываться отдельно от приложения, пока backend возвращается в работу.",
  },
};

function PrimaryAction({ copy }: { copy: ErrorPageCopy }) {
  if (copy.primaryHref) {
    return (
      <a className="error-page__button error-page__button--primary" href={copy.primaryHref}>
        <Home size={18} aria-hidden="true" />
        {copy.primaryLabel}
      </a>
    );
  }

  return (
    <button className="error-page__button error-page__button--primary" type="button" onClick={() => window.location.reload()}>
      <RefreshCw size={18} aria-hidden="true" />
      {copy.primaryLabel}
    </button>
  );
}

export function ErrorPageApp({ kind }: { kind: ErrorPageKind }) {
  const copy = copyByKind[kind];

  return (
    <main className="error-page" data-error-kind={kind}>
      <section className="error-page__shell" aria-labelledby="error-page-title">
        <div className="error-page__status" aria-hidden="true">
          <span>{copy.status}</span>
          <span>EDGE_READY</span>
          <span>POLIGLOT_AI</span>
        </div>

        <div className="error-page__content">
          <p className="error-page__eyebrow">
            <AlertTriangle size={18} aria-hidden="true" />
            {copy.eyebrow}
          </p>
          <h1 id="error-page-title">{copy.title}</h1>
          <p className="error-page__body">{copy.body}</p>
          <p className="error-page__english">{copy.english}</p>

          <div className="error-page__actions">
            <PrimaryAction copy={copy} />
            <a className="error-page__button" href={copy.secondaryHref}>
              <ArrowRight size={18} aria-hidden="true" />
              {copy.secondaryLabel}
            </a>
            <a className="error-page__button" href="https://t.me/NERIVAapp_bot">
              <MessageCircle size={18} aria-hidden="true" />
              Telegram
            </a>
          </div>
        </div>

        <p className="error-page__detail">{copy.detail}</p>
      </section>
    </main>
  );
}
