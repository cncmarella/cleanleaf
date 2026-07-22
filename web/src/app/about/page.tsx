import type { Metadata } from "next";
import { Container } from "@/components/container";
import { site, strengths } from "@/lib/site";

export const metadata: Metadata = {
  title: "About",
  description: `${site.name} is a crop protection manufacturer based in Ghatkesar, Hyderabad, supplying insecticides, fungicides, herbicides and plant nutrition products.`,
};

export default function AboutPage() {
  return (
    <>
      <section className="border-b border-border bg-surface py-16 sm:py-20">
        <Container>
          <h1 className="text-4xl font-bold tracking-tight sm:text-5xl">
            About {site.shortName}
          </h1>
          <p className="mt-5 max-w-2xl text-lg leading-relaxed text-muted">
            A crop protection manufacturer working out of Ghatkesar, Hyderabad —
            close enough to the fields we serve that our advice comes from the
            same season our customers are having.
          </p>
        </Container>
      </section>

      <section className="py-16 sm:py-20">
        <Container>
          <div className="grid gap-12 lg:grid-cols-3">
            <div className="space-y-5 leading-relaxed text-muted lg:col-span-2">
              <h2 className="text-2xl font-bold tracking-tight text-foreground">
                What we do
              </h2>
              <p>
                {site.name} formulates and supplies agricultural inputs across
                four ranges: insecticides, fungicides, herbicides and plant
                nutrition. Our products are made for the crops grown around us —
                cotton, paddy, chilli, groundnut and vegetables — and for the
                pest pressure those crops actually face.
              </p>
              <p>
                Every batch is checked for active-ingredient concentration and
                stability before dispatch, and labelled with the dosage,
                safety guidance and pre-harvest interval a grower needs to stay
                within residue limits.
              </p>
              <p>
                We work directly with dealers and with farmers. If you are not
                sure which product fits your situation, tell us the crop and
                what you are seeing in the field, and we will make a specific
                recommendation rather than sell you the biggest pack.
              </p>
            </div>

            <aside className="rounded-2xl border border-border bg-surface p-7">
              <h2 className="text-sm font-semibold uppercase tracking-wide text-muted">
                Our facility
              </h2>
              <address className="mt-4 text-sm not-italic leading-relaxed text-foreground">
                {site.address.line1}
                <br />
                {site.address.line2}
              </address>
              <dl className="mt-6 space-y-4 text-sm">
                <div>
                  <dt className="font-semibold text-muted">Phone</dt>
                  <dd className="mt-1">
                    <a
                      href={site.phoneHref}
                      className="text-leaf-700 underline-offset-4 hover:underline dark:text-leaf-400"
                    >
                      {site.phone}
                    </a>
                  </dd>
                </div>
                <div>
                  <dt className="font-semibold text-muted">Email</dt>
                  <dd className="mt-1">
                    <a
                      href={site.emailHref}
                      className="break-all text-leaf-700 underline-offset-4 hover:underline dark:text-leaf-400"
                    >
                      {site.email}
                    </a>
                  </dd>
                </div>
              </dl>
            </aside>
          </div>

          <div className="mt-16 border-t border-border pt-14">
            <h2 className="text-2xl font-bold tracking-tight">How we work</h2>
            <div className="mt-9 grid gap-8 sm:grid-cols-2">
              {strengths.map((item) => (
                <div key={item.title}>
                  <h3 className="text-base font-semibold">{item.title}</h3>
                  <p className="mt-2.5 text-sm leading-relaxed text-muted">
                    {item.body}
                  </p>
                </div>
              ))}
            </div>
          </div>
        </Container>
      </section>
    </>
  );
}
