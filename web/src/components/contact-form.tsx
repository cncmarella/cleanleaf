"use client";

import { useEffect, useId, useRef, useState } from "react";
import { submitContact, type ContactResult } from "@/lib/api";

const fieldClass =
  "w-full rounded-lg border border-border bg-background px-4 py-3 text-sm text-foreground placeholder:text-muted/70 focus:border-leaf-500 focus:outline-none";

type Status = "idle" | "submitting" | "done";

export function ContactForm() {
  const formId = useId();
  const [status, setStatus] = useState<Status>("idle");
  const [result, setResult] = useState<ContactResult | null>(null);
  // Lets a second submit cancel an in-flight first one.
  const inFlight = useRef<AbortController | null>(null);

  // Records which CTA/page led here (from a ?source= query param) for enquiry
  // attribution. Read after mount so the /contact page stays statically
  // prerendered; falls back to "direct" when there is no param.
  const sourceRef = useRef<HTMLInputElement>(null);
  useEffect(() => {
    const source = new URLSearchParams(window.location.search).get("source");
    if (source && sourceRef.current) {
      sourceRef.current.value = source;
    }
  }, []);

  async function handleSubmit(event: React.SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    const data = new FormData(form);

    inFlight.current?.abort();
    const controller = new AbortController();
    inFlight.current = controller;

    setStatus("submitting");
    setResult(null);

    const outcome = await submitContact(
      {
        name: String(data.get("name") ?? ""),
        email: String(data.get("email") ?? ""),
        phone: String(data.get("phone") ?? ""),
        subject: String(data.get("subject") ?? ""),
        message: String(data.get("message") ?? ""),
        source: String(data.get("source") ?? ""),
        website: String(data.get("website") ?? ""),
      },
      controller.signal,
    );

    // A superseded submit must not overwrite the newer one's result.
    if (controller.signal.aborted) return;

    setResult(outcome);
    setStatus(outcome.ok ? "done" : "idle");
    if (outcome.ok) {
      form.reset();
    }
  }

  const fieldErrors = result && !result.ok ? (result.fields ?? {}) : {};

  function errorFor(field: string) {
    const message = fieldErrors[field];
    if (!message) return null;
    return (
      <p id={`${formId}-${field}-error`} className="mt-1.5 text-xs text-red-600 dark:text-red-400">
        {message}
      </p>
    );
  }

  function ariaFor(field: string) {
    return fieldErrors[field]
      ? { "aria-invalid": true, "aria-describedby": `${formId}-${field}-error` }
      : {};
  }

  return (
    <form onSubmit={handleSubmit} noValidate className="space-y-5">
      {/* Attribution: which page/CTA led here. Populated from ?source= on mount;
          defaults to "direct" when the visitor came straight to /contact. */}
      <input type="hidden" name="source" defaultValue="direct" ref={sourceRef} />

      {/* Honeypot: hidden from people, tempting to bots. */}
      <div aria-hidden="true" className="absolute left-[-9999px]">
        <label htmlFor={`${formId}-website`}>Website</label>
        <input
          id={`${formId}-website`}
          name="website"
          type="text"
          tabIndex={-1}
          autoComplete="off"
        />
      </div>

      <div className="grid gap-5 sm:grid-cols-2">
        <div>
          <label htmlFor={`${formId}-name`} className="mb-1.5 block text-sm font-medium">
            Name <span className="text-red-600">*</span>
          </label>
          <input
            id={`${formId}-name`}
            name="name"
            type="text"
            required
            autoComplete="name"
            className={fieldClass}
            placeholder="Your name"
            {...ariaFor("name")}
          />
          {errorFor("name")}
        </div>

        <div>
          <label htmlFor={`${formId}-email`} className="mb-1.5 block text-sm font-medium">
            Email <span className="text-red-600">*</span>
          </label>
          <input
            id={`${formId}-email`}
            name="email"
            type="email"
            required
            autoComplete="email"
            className={fieldClass}
            placeholder="you@example.com"
            {...ariaFor("email")}
          />
          {errorFor("email")}
        </div>

        <div>
          <label htmlFor={`${formId}-phone`} className="mb-1.5 block text-sm font-medium">
            Phone <span className="font-normal text-muted">(optional)</span>
          </label>
          <input
            id={`${formId}-phone`}
            name="phone"
            type="tel"
            autoComplete="tel"
            className={fieldClass}
            placeholder="10-digit mobile number"
            {...ariaFor("phone")}
          />
          {errorFor("phone")}
        </div>

        <div>
          <label htmlFor={`${formId}-subject`} className="mb-1.5 block text-sm font-medium">
            Subject
          </label>
          <input
            id={`${formId}-subject`}
            name="subject"
            type="text"
            className={fieldClass}
            placeholder="Bulk order, dealership, product advice…"
            {...ariaFor("subject")}
          />
          {errorFor("subject")}
        </div>
      </div>

      <div>
        <label htmlFor={`${formId}-message`} className="mb-1.5 block text-sm font-medium">
          Message <span className="text-red-600">*</span>
        </label>
        <textarea
          id={`${formId}-message`}
          name="message"
          required
          rows={6}
          className={`${fieldClass} resize-y`}
          placeholder="Tell us your crop, the pest or disease you're seeing, and the acreage."
          {...ariaFor("message")}
        />
        {errorFor("message")}
      </div>

      <button
        type="submit"
        disabled={status === "submitting"}
        className="rounded-full bg-leaf-700 px-8 py-3.5 text-sm font-semibold text-white transition-colors hover:bg-leaf-800 disabled:cursor-not-allowed disabled:opacity-60"
      >
        {status === "submitting" ? "Sending…" : "Send enquiry"}
      </button>

      {/* aria-live so screen readers announce the outcome without a focus jump. */}
      <div aria-live="polite" className="min-h-6">
        {result && (
          <p
            className={`text-sm ${
              result.ok
                ? "text-leaf-700 dark:text-leaf-400"
                : "text-red-600 dark:text-red-400"
            }`}
          >
            {result.message}
          </p>
        )}
      </div>
    </form>
  );
}
