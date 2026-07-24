import { site } from "@/lib/site";

/**
 * schema.org Organization data, emitted once site-wide from the root layout so
 * search engines and AI crawlers can read the business name, address and
 * contact details as structured facts rather than scraping them from markup.
 * Sourced from site.ts — the same single source the visible pages use.
 */
const organization = {
  "@context": "https://schema.org",
  "@type": "Organization",
  name: site.name,
  alternateName: site.shortName,
  url: site.url,
  description: site.description,
  telephone: site.phoneHref.replace("tel:", ""),
  email: site.email,
  address: {
    "@type": "PostalAddress",
    streetAddress: site.address.line1,
    addressLocality: site.address.locality,
    addressRegion: site.address.region,
    postalCode: site.address.postalCode,
    addressCountry: site.address.country,
  },
  contactPoint: {
    "@type": "ContactPoint",
    contactType: "sales",
    telephone: site.phoneHref.replace("tel:", ""),
    email: site.email,
    areaServed: site.address.region,
    availableLanguage: ["en", "te", "hi"],
  },
} as const;

export function StructuredData() {
  return (
    <script
      type="application/ld+json"
      // Escape "<" per Next's JSON-LD guidance to avoid XSS via injected markup.
      dangerouslySetInnerHTML={{
        __html: JSON.stringify(organization).replace(/</g, "\\u003c"),
      }}
    />
  );
}
