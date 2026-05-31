import { useSyncExternalStore } from "react";

function subscribe(query: string, onChange: () => void): () => void {
  const media = window.matchMedia(query);
  media.addEventListener("change", onChange);
  return () => media.removeEventListener("change", onChange);
}

function getSnapshot(query: string): boolean {
  return window.matchMedia(query).matches;
}

function getServerSnapshot(): boolean {
  return false;
}

export function useMediaQuery(query: string): boolean {
  return useSyncExternalStore(
    (onChange) => subscribe(query, onChange),
    () => getSnapshot(query),
    getServerSnapshot,
  );
}

export function useIsMobile(): boolean {
  return useMediaQuery("(max-width: 767px)");
}
