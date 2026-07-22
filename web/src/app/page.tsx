import Link from "next/link";
import { Container } from "@/components/container";
import { products, site, strengths } from "@/lib/site";

export default function HomePage() {
  return (
    <>
      {/* Hero */}
      <section className="relative overflow-hidden border-b border-border bg-surface">
        {/* Decorative wash; kept behind content and hidden from assistive tech. */}
        <div
          aria-hidden="true"
          className="pointer-events-none absolute -right-24 -top-32 h-96 w-96 rounded-full bg-leaf-300/25 blur-3xl dark:bg-leaf-800/25"
        />
        <Container className="relative py-20 sm:py-28">
          <div className="max-w-2xl">
            <p className="inline-flex items-center rounded-full border border-leaf-200 bg-leaf-50 px-3 py-1 text-xs font-semibold uppercase tracking-wide text-leaf-800 dark:border-leaf-900 dark:bg-leaf-950 dark:text-leaf-300">
              Made in Ghatkesar, Hyderabad
            </p>
            <h1 className="mt-6 text-4xl font-bold leading-[1.1] tracking-tight sm:text-5xl lg:text-6xl">
              Crop protection you can{" "}
              <span className="text-leaf-700 dark:text-leaf-400">trust</span>
            </h1>
            <p className="mt-6 text-lg leading-relaxed text-muted">
              {site.name} formulates and supplies insecticides, fungicides,
              herbicides and plant nutrition for farms across Telangana — backed
              by batch testing and practical, crop-specific advice.
            </p>
            <div className="mt-9 flex flex-col gap-3 sm:flex-row">
              <Link
                href="/contact"
                className="rounded-full bg-leaf-700 px-7 py-3.5 text-center text-sm font-semibold text-white transition-colors hover:bg-leaf-800"
              >
                Request a quote
              </Link>
              <Link
                href="/products"
                className="rounded-full border border-border px-7 py-3.5 text-center text-sm font-semibold text-foreground transition-colors hover:bg-background"
              >
                Browse products
              </Link>
            </div>
          </div>
        </Container>
      </section>

      {/* Why us */}
      <section className="py-20 sm:py-24">
        <Container>
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
            Why growers work with us
          </h2>
          <div className="mt-12 grid gap-8 sm:grid-cols-2 lg:grid-cols-4">
            {strengths.map((item) => (
              <div key={item.title}>
                <div
                  aria-hidden="true"
                  className="mb-4 h-1 w-10 rounded-full bg-leaf-500"
                />
                <h3 className="text-base font-semibold">{item.title}</h3>
                <p className="mt-2.5 text-sm leading-relaxed text-muted">
                  {item.body}
                </p>
              </div>
            ))}
          </div>
        </Container>
      </section>

      {/* Product preview */}
      <section className="border-y border-border bg-surface py-20 sm:py-24">
        <Container>
          <div className="flex flex-wrap items-end justify-between gap-4">
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">
              Our range
            </h2>
            <Link
              href="/products"
              className="text-sm font-semibold text-leaf-700 underline-offset-4 hover:underline dark:text-leaf-400"
            >
              See all products →
            </Link>
          </div>
          <div className="mt-12 grid gap-5 sm:grid-cols-2">
            {products.map((product) => (
              <article
                key={product.slug}
                className="rounded-2xl border border-border bg-background p-6"
              >
                <h3 className="text-lg font-semibold">{product.name}</h3>
                <p className="mt-2.5 text-sm leading-relaxed text-muted">
                  {product.summary}
                </p>
              </article>
            ))}
          </div>
        </Container>
      </section>

      {/* Closing CTA */}
      <section className="py-20 sm:py-24">
        <Container>
          <div className="rounded-3xl bg-leaf-800 px-8 py-14 text-center dark:bg-leaf-900 sm:px-14">
            <h2 className="text-3xl font-bold tracking-tight text-white sm:text-4xl">
              Need a recommendation for your crop?
            </h2>
            <p className="mx-auto mt-4 max-w-xl text-leaf-100">
              Tell us the crop and the problem you are seeing. We will suggest a
              product, a dosage and the right time to spray.
            </p>
            <div className="mt-9 flex flex-col justify-center gap-3 sm:flex-row">
              <Link
                href="/contact"
                className="rounded-full bg-white px-7 py-3.5 text-sm font-semibold text-leaf-900 transition-colors hover:bg-leaf-50"
              >
                Send an enquiry
              </Link>
              <a
                href={site.phoneHref}
                className="rounded-full border border-leaf-500 px-7 py-3.5 text-sm font-semibold text-white transition-colors hover:bg-leaf-700"
              >
                Call {site.phone}
              </a>
            </div>
          </div>
        </Container>
      </section>
    </>
  );
}
