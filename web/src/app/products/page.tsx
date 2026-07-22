import type { Metadata } from "next";
import Link from "next/link";
import { Container } from "@/components/container";
import { products } from "@/lib/site";

export const metadata: Metadata = {
  title: "Products",
  description:
    "Insecticides, fungicides, herbicides and plant nutrition products from CleanLeaf Technologies, with the crops each range is suited to.",
};

export default function ProductsPage() {
  return (
    <>
      <section className="border-b border-border bg-surface py-16 sm:py-20">
        <Container>
          <h1 className="text-4xl font-bold tracking-tight sm:text-5xl">
            Products
          </h1>
          <p className="mt-5 max-w-2xl text-lg leading-relaxed text-muted">
            Four ranges covering the season end to end. Pack sizes, active
            ingredients and registration details are available on request.
          </p>
        </Container>
      </section>

      <section className="py-16 sm:py-20">
        <Container>
          <div className="grid gap-6 lg:grid-cols-2">
            {products.map((product) => (
              <article
                key={product.slug}
                className="flex flex-col rounded-2xl border border-border p-7"
              >
                <h2 className="text-xl font-semibold">{product.name}</h2>
                <p className="mt-3 text-sm leading-relaxed text-muted">
                  {product.summary}
                </p>
                <div className="mt-6 border-t border-border pt-5">
                  <h3 className="text-xs font-semibold uppercase tracking-wide text-muted">
                    Suited to
                  </h3>
                  <ul className="mt-3 flex flex-wrap gap-2">
                    {product.crops.map((crop) => (
                      <li
                        key={crop}
                        className="rounded-full bg-leaf-50 px-3 py-1 text-xs font-medium text-leaf-800 dark:bg-leaf-950 dark:text-leaf-300"
                      >
                        {crop}
                      </li>
                    ))}
                  </ul>
                </div>
              </article>
            ))}
          </div>

          <div className="mt-14 rounded-2xl border border-border bg-surface p-8 text-center">
            <h2 className="text-xl font-semibold">
              Looking for something specific?
            </h2>
            <p className="mx-auto mt-3 max-w-xl text-sm leading-relaxed text-muted">
              Tell us your crop and the pest or disease you are dealing with,
              and we will point you to the right product and dosage.
            </p>
            <Link
              href="/contact"
              className="mt-7 inline-block rounded-full bg-leaf-700 px-7 py-3.5 text-sm font-semibold text-white transition-colors hover:bg-leaf-800"
            >
              Ask us
            </Link>
          </div>
        </Container>
      </section>
    </>
  );
}
