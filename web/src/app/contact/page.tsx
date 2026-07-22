import type { Metadata } from "next";
import { Container } from "@/components/container";
import { ContactForm } from "@/components/contact-form";
import { site } from "@/lib/site";

export const metadata: Metadata = {
  title: "Contact",
  description: `Get in touch with ${site.name} in Ghatkesar, Hyderabad for product advice, bulk orders and dealership enquiries.`,
};

export default function ContactPage() {
  return (
    <>
      <section className="border-b border-border bg-surface py-16 sm:py-20">
        <Container>
          <h1 className="text-4xl font-bold tracking-tight sm:text-5xl">
            Contact us
          </h1>
          <p className="mt-5 max-w-2xl text-lg leading-relaxed text-muted">
            Product advice, bulk orders, dealership enquiries — send us a note
            and we will get back to you.
          </p>
        </Container>
      </section>

      <section className="py-16 sm:py-20">
        <Container>
          <div className="grid gap-14 lg:grid-cols-3">
            <div className="lg:col-span-2">
              <ContactForm />
            </div>

            <aside className="space-y-8">
              <div>
                <h2 className="text-sm font-semibold uppercase tracking-wide text-muted">
                  Visit
                </h2>
                <address className="mt-3 text-sm not-italic leading-relaxed">
                  {site.address.line1}
                  <br />
                  {site.address.line2}
                </address>
              </div>

              <div>
                <h2 className="text-sm font-semibold uppercase tracking-wide text-muted">
                  Call
                </h2>
                <p className="mt-3 text-sm">
                  <a
                    href={site.phoneHref}
                    className="text-leaf-700 underline-offset-4 hover:underline dark:text-leaf-400"
                  >
                    {site.phone}
                  </a>
                </p>
              </div>

              <div>
                <h2 className="text-sm font-semibold uppercase tracking-wide text-muted">
                  Email
                </h2>
                <p className="mt-3 text-sm">
                  <a
                    href={site.emailHref}
                    className="break-all text-leaf-700 underline-offset-4 hover:underline dark:text-leaf-400"
                  >
                    {site.email}
                  </a>
                </p>
              </div>
            </aside>
          </div>
        </Container>
      </section>
    </>
  );
}
