import { useEffect, useState } from "react";
import { Download, X } from "lucide-react";

type BrowserInstallPrompt = Event & { prompt: () => Promise<unknown> };

export function pwaUsesRussian() {
  const language = new URLSearchParams(location.search).get("lang") || navigator.language;
  return language?.toLowerCase().startsWith("ru");
}

function isInstalled() {
  return window.matchMedia("(display-mode: standalone)").matches
    || (navigator as Navigator & { standalone?: boolean }).standalone === true;
}

export function PwaInstall() {
  const [open, setOpen] = useState(() => new URLSearchParams(location.search).get("install") === "1" && !isInstalled());
  const [prompt, setPrompt] = useState<BrowserInstallPrompt | null>(null);
  const russian = pwaUsesRussian();
  const title = russian ? "Установить NERIVA" : "Install NERIVA";
  const appleMobile = /iPhone|iPad|iPod/i.test(`${navigator.userAgent} ${navigator.platform}`)
    || (navigator.platform === "MacIntel" && navigator.maxTouchPoints > 1);

  const dismiss = () => {
    const url = new URL(location.href);
    url.searchParams.delete("install");
    history.replaceState(history.state, "", url);
    setOpen(false);
  };

  useEffect(() => {
    if (!open) return;
    const capture = (event: Event) => {
      if (typeof (event as BrowserInstallPrompt).prompt !== "function") return;
      event.preventDefault();
      setPrompt(event as BrowserInstallPrompt);
    };
    window.addEventListener("beforeinstallprompt", capture);
    window.addEventListener("appinstalled", dismiss);
    return () => {
      window.removeEventListener("beforeinstallprompt", capture);
      window.removeEventListener("appinstalled", dismiss);
    };
  }, [open]);

  const install = async () => {
    if (!prompt) return;
    const currentPrompt = prompt;
    setPrompt(null);
    await currentPrompt.prompt().catch(() => undefined);
  };

  if (!open) return null;

  const instructions = appleMobile
    ? (russian ? 'Откройте эту страницу в Safari. Нажмите «Поделиться», затем выберите На экран «Домой».' : 'Open this page in Safari. Tap Share, then Add to Home Screen.')
    : (russian ? 'В меню браузера выберите «Установить приложение» или «Добавить на главный экран». В Safari на Mac выберите «Файл» → «Добавить в Dock».' : 'In your browser menu, choose Install app or Add to Home Screen. On Mac in Safari, choose File, then Add to Dock.');

  return (
    <aside className="pwa-install-panel" aria-label={title} role="region" onKeyDown={(event) => {
      if (event.key === "Escape") dismiss();
    }}>
      <Download size={20} aria-hidden="true" />
      <div>
        <h2>{title}</h2>
        <p role="status">{prompt ? (russian ? "Открывайте приложение с иконки на своём устройстве." : "Open the app from its icon on your device.") : instructions}</p>
        {prompt ? <button className="pwa-install-panel__action" type="button" onClick={() => void install()}>{title}</button> : null}
      </div>
      <button className="pwa-install-panel__close" type="button" aria-label={russian ? "Закрыть" : "Close"} onClick={dismiss}><X size={18} /></button>
    </aside>
  );
}
