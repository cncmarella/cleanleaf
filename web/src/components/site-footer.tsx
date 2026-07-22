import Link from "next/link";
import { Container } from "@/components/container";
import { Logo } from "@/components/logo";
import { nav, products, site } from "@/lib/site";

export function SiteFooter() {
  return (
    <footer className="border-t border-border bg-surface">
      <Container className="py-12">
        <div className="grid gap-10 sm:grid-cols-2 lg:grid-cols-4">
          <div className="lg:col-span-2">
            <Logo />
            <p className="mt-4 max-w-sm text-sm leading-relaxed text-muted">
              {site.description}
            </p>
          </div>

          <div>
            <h2 className="text-sm font-semibold text-foreground">Explore</h2>
            <ul className="mt-4 space-y-2.5">
              {nav.map((item) => (
                <li key={item.href}>
                  <Link
                    href={item.href}
                    className="text-sm text-muted transition-colors hover:text-foreground"
                  >
                    {item.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div>
            <h2 className="text-sm font-semibold text-foreground">Reach us</h2>
            <address className="mt-4 space-y-2.5 text-sm not-italic text-muted">
              <p>
                {site.address.line1}
                <br />
                {site.address.line2}
              </p>
              <p>
                <a
                  href={site.phoneHref}
                  className="transition-colors hover:text-foreground"
                >
                  {site.phone}
                </a>
              </p>
              <p>
                <a
                  href={site.emailHref}
                  className="break-all transition-colors hover:text-foreground"
                >
                  {site.email}
                </a>
              </p>
            </address>
          </div>
        </div>

        <div className="mt-10 flex flex-col gap-3 border-t border-border pt-6 text-xs text-muted sm:flex-row sm:items-center sm:justify-between">
          <p>
            © {new Date().getFullYear()} {site.name}. All rights reserved.
          </p>
          <p>
            Always read the product label before use.{" "}
            <Link href="/products" className="underline underline-offset-2">
              See our range ({products.length} categories)
            </Link>
          </p>
        </div>
      </Container>
    </footer>
  );
}
