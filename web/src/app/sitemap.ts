import type { MetadataRoute } from "next";
import { nav, site } from "@/lib/site";

// Derived from the same nav list the header renders, so a new page shows up in
// the sitemap the moment it is linked — no second list to keep in sync.
export default function sitemap(): MetadataRoute.Sitemap {
  const lastModified = new Date();
  return nav.map(({ href }) => ({
    url: new URL(href, site.url).toString(),
    lastModified,
    changeFrequency: "monthly",
    priority: href === "/" ? 1 : 0.8,
  }));
}
