import { site } from "@/lib/site";

/**
 * Wordmark with an inline leaf glyph. Inline SVG rather than an image file so
 * it stays crisp and inherits the current text colour — swap for the real
 * brand mark once we have it as an SVG.
 */
export function Logo({ className = "" }: { className?: string }) {
  return (
    <span className={`inline-flex items-center gap-2.5 ${className}`}>
      <svg
        viewBox="0 0 32 32"
        aria-hidden="true"
        className="h-8 w-8 shrink-0 text-leaf-600 dark:text-leaf-400"
      >
        <circle cx="16" cy="16" r="15" className="fill-current opacity-15" />
        <path
          d="M23.5 8.5c0 7.2-4.4 12-9.6 12-1.5 0-2.8-.4-3.8-1.1 1.1-5.9 5.9-9.5 13.4-10.9Z"
          className="fill-current"
        />
        <path
          d="M8.5 24c1-5.2 4-9.4 8.6-11.7"
          fill="none"
          stroke="currentColor"
          strokeWidth="1.8"
          strokeLinecap="round"
        />
      </svg>
      <span className="flex flex-col leading-none">
        <span className="text-lg font-bold tracking-tight text-leaf-800 dark:text-leaf-200">
          {site.shortName.toUpperCase()}
        </span>
        <span className="text-[0.6rem] font-medium uppercase tracking-[0.2em] text-muted">
          Technologies
        </span>
      </span>
    </span>
  );
}
